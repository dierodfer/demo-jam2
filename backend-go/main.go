// Backend Go del Portal de Empleado.
//
// Implementa el mismo contrato (shared/openapi.yaml) que el backend Java:
//   POST /api/login   GET /api/me   PUT /api/me   POST /api/logout
//
// Usa el mismo fichero SQLite (./data/portal.db) y el mismo esquema de tabla,
// para que ambos backends sean intercambiables. La sesión se mantiene con una
// cookie propia (session_id) respaldada por un almacén en memoria.
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"crypto/rand"
	"encoding/hex"

	_ "modernc.org/sqlite"
)

const empleadoID int64 = 1

// Empleado replica exactamente el JSON del contrato (sin el campo username).
type Empleado struct {
	ID           int64  `json:"id"`
	Nombre       string `json:"nombre"`
	Email        string `json:"email"`
	Telefono     string `json:"telefono"`
	Puesto       string `json:"puesto"`
	Departamento string `json:"departamento"`
	Direccion    string `json:"direccion"`
	Foto         string `json:"foto"`
}

// Certificacion replica el JSON del contrato (sin empleadoId).
type Certificacion struct {
	ID             int64  `json:"id"`
	Conocimiento   string `json:"conocimiento"`
	EmpresaEmisora string `json:"empresaEmisora"`
	Fecha          string `json:"fecha"`
}

type certificacionInput struct {
	Conocimiento   string `json:"conocimiento"`
	EmpresaEmisora string `json:"empresaEmisora"`
	Fecha          string `json:"fecha"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type empleadoUpdate struct {
	Nombre       *string `json:"nombre"`
	Email        *string `json:"email"`
	Telefono     *string `json:"telefono"`
	Puesto       *string `json:"puesto"`
	Departamento *string `json:"departamento"`
	Direccion    *string `json:"direccion"`
	Foto         *string `json:"foto"`
}

// sessionStore es un almacén de sesiones en memoria: token -> empleadoId.
type sessionStore struct {
	mu       sync.RWMutex
	sessions map[string]int64
}

func newSessionStore() *sessionStore {
	return &sessionStore{sessions: make(map[string]int64)}
}

func (s *sessionStore) create(empID int64) string {
	token := randomToken()
	s.mu.Lock()
	s.sessions[token] = empID
	s.mu.Unlock()
	return token
}

func (s *sessionStore) get(token string) (int64, bool) {
	s.mu.RLock()
	id, ok := s.sessions[token]
	s.mu.RUnlock()
	return id, ok
}

func (s *sessionStore) delete(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func randomToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type server struct {
	db            *sql.DB
	sessions      *sessionStore
	allowedOrigin string
}

const cookieName = "session_id"

func main() {
	dbPath := env("PORTAL_DB_PATH", "../data/portal.db")
	port := env("SERVER_PORT", "8081")
	origin := env("CORS_ALLOWED_ORIGIN", "http://localhost:5174")
	seedUsername := env("SEED_USERNAME", "admin")
	migrationsDir := env("MIGRATIONS_PATH", "../shared/migrations")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("no se pudo abrir la base de datos: %v", err)
	}
	defer db.Close()
	// Una única conexión por proceso; WAL + busy_timeout permiten que los dos
	// backends (Java y Go) escriban el mismo fichero sin "database is locked".
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		log.Fatalf("no se pudo activar WAL: %v", err)
	}
	if _, err := db.Exec(`PRAGMA busy_timeout=5000`); err != nil {
		log.Fatalf("no se pudo fijar busy_timeout: %v", err)
	}

	if err := runMigrations(db, migrationsDir); err != nil {
		log.Fatalf("no se pudieron aplicar las migraciones: %v", err)
	}
	if err := seed(db, seedUsername); err != nil {
		log.Fatalf("no se pudo sembrar el empleado: %v", err)
	}

	s := &server{db: db, sessions: newSessionStore(), allowedOrigin: origin}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/login", s.handleLogin)
	mux.HandleFunc("GET /api/me", s.handleGetMe)
	mux.HandleFunc("PUT /api/me", s.handlePutMe)
	mux.HandleFunc("POST /api/logout", s.handleLogout)
	mux.HandleFunc("GET /api/certificaciones", s.handleListCertificaciones)
	mux.HandleFunc("POST /api/certificaciones", s.handleCreateCertificacion)
	mux.HandleFunc("PUT /api/certificaciones/{id}", s.handleUpdateCertificacion)
	mux.HandleFunc("DELETE /api/certificaciones/{id}", s.handleDeleteCertificacion)

	handler := s.withCORS(mux)

	log.Printf("Backend Go escuchando en :%s (db=%s, origin=%s)", port, dbPath, origin)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// runMigrations aplica en orden los ficheros .sql de shared/migrations que aún
// no consten en la tabla schema_migration. Es el MISMO mecanismo que usa el
// backend Java: el primero que arranca aplica la migración, el otro la ve ya
// registrada. Las migraciones aplicadas nunca se editan.
func runMigrations(db *sql.DB, dir string) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migration (
			version    TEXT PRIMARY KEY,
			applied_at TEXT NOT NULL
		)`); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir) // ReadDir devuelve ordenado por nombre
	if err != nil {
		return fmt.Errorf("no se pudo leer %s: %w", dir, err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".sql") {
			continue
		}

		var applied int
		if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migration WHERE version=?`, name).Scan(&applied); err != nil {
			return err
		}
		if applied > 0 {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}
		for _, stmt := range splitStatements(string(content)) {
			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				return fmt.Errorf("migración %s: %w", name, err)
			}
		}
		if _, err := tx.Exec(`INSERT INTO schema_migration (version, applied_at) VALUES (?, ?)`,
			name, time.Now().UTC().Format(time.RFC3339)); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		log.Printf("migración aplicada: %s", name)
	}
	return nil
}

// splitStatements separa un fichero .sql en sentencias por ";" y descarta los
// trozos que solo contienen comentarios o espacios. Las migraciones no deben
// usar ";" dentro de literales de texto.
func splitStatements(sqlText string) []string {
	var out []string
	for _, chunk := range strings.Split(sqlText, ";") {
		hasContent := false
		for _, line := range strings.Split(chunk, "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
				hasContent = true
				break
			}
		}
		if hasContent {
			out = append(out, chunk)
		}
	}
	return out
}

func seed(db *sql.DB, username string) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM empleado`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		if _, err := db.Exec(`
			INSERT INTO empleado (id, username, nombre, email, telefono, puesto, departamento, direccion, foto)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			empleadoID, username, "Ana García", "ana.garcia@empresa.com", "+34 600 123 456",
			"Desarrolladora Full Stack", "Tecnología", "Calle Mayor 1, 28013 Madrid",
			"https://i.pravatar.cc/300?u=ana.garcia"); err != nil {
			return err
		}
	}

	var certCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM certificacion`).Scan(&certCount); err != nil {
		return err
	}
	if certCount == 0 {
		certs := [][3]string{
			{"AWS Certified Developer – Associate", "AWS", "2025-11-17"},
			{"AWS Certified Solutions Architect – Associate", "AWS", "2025-11-17"},
			{"AWS Certified SysOps Administrator – Associate", "AWS", "2025-11-17"},
			{"Certificado PRL", "Avanta", "2026-05-29"},
		}
		for _, c := range certs {
			if _, err := db.Exec(`
				INSERT INTO certificacion (empleado_id, conocimiento, empresa_emisora, fecha)
				VALUES (?, ?, ?, ?)`, empleadoID, c[0], c[1], c[2]); err != nil {
				return err
			}
		}
	}
	return nil
}

// ---- CORS ----

func (s *server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.allowedOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Vary", "Origin")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ---- Handlers ----

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo no válido")
		return
	}

	var username string
	emp, err := s.loadEmpleado(&username)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Usuario o contraseña no válidos")
		return
	}
	// La contraseña se ignora; solo se comprueba el username sembrado.
	if body.Username == "" || body.Username != username {
		writeError(w, http.StatusUnauthorized, "Usuario o contraseña no válidos")
		return
	}

	token := s.sessions.create(emp.ID)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, http.StatusOK, emp)
}

func (s *server) handleGetMe(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.authed(r); !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}
	emp, err := s.loadEmpleado(nil)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}
	writeJSON(w, http.StatusOK, emp)
}

func (s *server) handlePutMe(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.authed(r); !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}
	var body empleadoUpdate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo no válido")
		return
	}

	emp, err := s.loadEmpleado(nil)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}
	if body.Nombre != nil {
		emp.Nombre = *body.Nombre
	}
	if body.Email != nil {
		emp.Email = *body.Email
	}
	if body.Telefono != nil {
		emp.Telefono = *body.Telefono
	}
	if body.Puesto != nil {
		emp.Puesto = *body.Puesto
	}
	if body.Departamento != nil {
		emp.Departamento = *body.Departamento
	}
	if body.Direccion != nil {
		emp.Direccion = *body.Direccion
	}
	if body.Foto != nil {
		emp.Foto = *body.Foto
	}

	_, err = s.db.Exec(`
		UPDATE empleado
		SET nombre=?, email=?, telefono=?, puesto=?, departamento=?, direccion=?, foto=?
		WHERE id=?`,
		emp.Nombre, emp.Email, emp.Telefono, emp.Puesto, emp.Departamento,
		emp.Direccion, emp.Foto, emp.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo actualizar")
		return
	}
	writeJSON(w, http.StatusOK, emp)
}

func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(cookieName); err == nil {
		s.sessions.delete(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusNoContent)
}

// ---- Certificaciones ----

func (s *server) handleListCertificaciones(w http.ResponseWriter, r *http.Request) {
	empID, ok := s.authed(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}
	rows, err := s.db.Query(`
		SELECT id, conocimiento, empresa_emisora, fecha
		FROM certificacion WHERE empleado_id=? ORDER BY id ASC`, empID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo listar")
		return
	}
	defer rows.Close()

	list := []Certificacion{}
	for rows.Next() {
		var c Certificacion
		if err := rows.Scan(&c.ID, &c.Conocimiento, &c.EmpresaEmisora, &c.Fecha); err != nil {
			writeError(w, http.StatusInternalServerError, "no se pudo leer")
			return
		}
		list = append(list, c)
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *server) handleCreateCertificacion(w http.ResponseWriter, r *http.Request) {
	empID, ok := s.authed(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}
	var body certificacionInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo no válido")
		return
	}
	res, err := s.db.Exec(`
		INSERT INTO certificacion (empleado_id, conocimiento, empresa_emisora, fecha)
		VALUES (?, ?, ?, ?)`, empID, body.Conocimiento, body.EmpresaEmisora, body.Fecha)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo crear")
		return
	}
	id, _ := res.LastInsertId()
	writeJSON(w, http.StatusCreated, Certificacion{
		ID: id, Conocimiento: body.Conocimiento, EmpresaEmisora: body.EmpresaEmisora, Fecha: body.Fecha,
	})
}

func (s *server) handleUpdateCertificacion(w http.ResponseWriter, r *http.Request) {
	empID, ok := s.authed(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "No encontrado")
		return
	}
	var body certificacionInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo no válido")
		return
	}
	res, err := s.db.Exec(`
		UPDATE certificacion SET conocimiento=?, empresa_emisora=?, fecha=?
		WHERE id=? AND empleado_id=?`,
		body.Conocimiento, body.EmpresaEmisora, body.Fecha, id, empID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo actualizar")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "No encontrado")
		return
	}
	writeJSON(w, http.StatusOK, Certificacion{
		ID: id, Conocimiento: body.Conocimiento, EmpresaEmisora: body.EmpresaEmisora, Fecha: body.Fecha,
	})
}

func (s *server) handleDeleteCertificacion(w http.ResponseWriter, r *http.Request) {
	empID, ok := s.authed(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "No autenticado")
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, "No encontrado")
		return
	}
	res, err := s.db.Exec(`DELETE FROM certificacion WHERE id=? AND empleado_id=?`, id, empID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "no se pudo eliminar")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		writeError(w, http.StatusNotFound, "No encontrado")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---- Helpers ----

// authed devuelve el empleadoId asociado a la cookie de sesión, si es válida.
func (s *server) authed(r *http.Request) (int64, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return 0, false
	}
	return s.sessions.get(c.Value)
}

// loadEmpleado carga el empleado id=1. Si usernameOut no es nil, además
// escribe ahí el username almacenado (que no se expone en el JSON).
func (s *server) loadEmpleado(usernameOut *string) (*Empleado, error) {
	var e Empleado
	var username string
	err := s.db.QueryRow(`
		SELECT id, username, nombre, email, telefono, puesto, departamento, direccion, foto
		FROM empleado WHERE id=?`, empleadoID).
		Scan(&e.ID, &username, &e.Nombre, &e.Email, &e.Telefono, &e.Puesto,
			&e.Departamento, &e.Direccion, &e.Foto)
	if err != nil {
		return nil, err
	}
	if usernameOut != nil {
		*usernameOut = username
	}
	return &e, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
