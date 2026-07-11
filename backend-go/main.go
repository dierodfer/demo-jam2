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
	"log"
	"net/http"
	"os"
	"sync"

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

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("no se pudo abrir la base de datos: %v", err)
	}
	defer db.Close()
	// SQLite es monousuario para escritura: un único conexión evita "database is locked".
	db.SetMaxOpenConns(1)

	if err := initSchema(db); err != nil {
		log.Fatalf("no se pudo inicializar el esquema: %v", err)
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

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS empleado (
			id           INTEGER PRIMARY KEY,
			username     TEXT,
			nombre       TEXT,
			email        TEXT,
			telefono     TEXT,
			puesto       TEXT,
			departamento TEXT,
			direccion    TEXT,
			foto         TEXT
		)`)
	return err
}

func seed(db *sql.DB, username string) error {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM empleado`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err := db.Exec(`
		INSERT INTO empleado (id, username, nombre, email, telefono, puesto, departamento, direccion, foto)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		empleadoID, username, "Ana García", "ana.garcia@empresa.com", "+34 600 123 456",
		"Desarrolladora Full Stack", "Tecnología", "Calle Mayor 1, 28013 Madrid",
		"https://i.pravatar.cc/300?u=ana.garcia")
	return err
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
