# Guía para agentes de IA — Portal de Empleado

Este repo es una demo: un portal de empleado implementado **por duplicado** en
4 stacks para comparar cómo la IA genera la misma característica en cada uno.
Tu trabajo casi siempre será **añadir una característica nueva manteniendo la
paridad entre las 4 implementaciones**. Esta guía es la receta.

## Arquitectura en 30 segundos

| Componente | Stack | Puerto | Habla con |
|---|---|---|---|
| `backend-java/` | Spring Boot 4.1 (Java 25), JPA, sqlite-jdbc | 8080 | `data/portal.db` |
| `backend-go/` | Go, `net/http`, `modernc.org/sqlite` | 8081 | `data/portal.db` (el mismo fichero) |
| `frontend-react/` | React 19 + Vite | 5173 | backend Java por defecto |
| `frontend-vue/` | Vue 3 + Vite | 5174 | backend Go por defecto |

Fuentes de verdad **compartidas** (no hay código compartido entre backends,
solo estos ficheros):

- `shared/openapi.yaml` — el contrato de API. Ambos backends lo implementan
  idéntico (mismas rutas, mismo JSON).
- `shared/migrations/*.sql` — el esquema de la BBDD. Ambos backends aplican
  estas migraciones al arrancar (tabla `schema_migration` registra cuáles ya
  corrieron). **Ningún backend crea esquema por su cuenta**: Hibernate está en
  `ddl-auto=none` y el Go no tiene DDL propio.

Login simulado: username `admin` (env `SEED_USERNAME`), **cualquier**
contraseña. Sesión por cookie (`JSESSIONID` en Java, `session_id` en Go);
el frontend siempre llama con `credentials: 'include'`.

## La regla de oro: paridad

Una característica no está terminada hasta que existe en **los 5 sitios**:
contrato, backend Java, backend Go, frontend React, frontend Vue. Si te piden
solo una parte, deja constancia explícita de qué quedó pendiente.

### Receta para añadir una característica (en este orden)

1. **Contrato** — añade rutas y esquemas a `shared/openapi.yaml`.
2. **Esquema** — si hay tablas/columnas nuevas, crea `shared/migrations/NNN_descripcion.sql`
   con el siguiente número libre. **Nunca edites una migración ya aplicada**;
   los cambios van en una migración nueva. SQL simple, sin `;` dentro de
   literales (los runners parten por `;`).
3. **Backend Java** — entidad con `@Table`/`@Column` explícitos (tabla y
   columnas en `snake_case`, igual que la migración), repositorio, DTOs
   `record` (JSON en `camelCase`, sin exponer `username` ni `empleado_id`),
   controlador que lee `empleadoId` de la sesión (ver `CertificacionController`
   como plantilla). Siembra de datos demo en `DataSeeder` si aplica.
4. **Backend Go** — struct con tags JSON idénticos al DTO de Java, handlers
   registrados como `mux.HandleFunc("GET /api/...", ...)`, mismo scoping por
   sesión, misma siembra en `seed()` (ver los handlers de certificaciones como
   plantilla). Los datos sembrados deben ser **idénticos** a los del Java.
5. **Frontends** — para una sección nueva del portal: añade la entrada en
   `SECCIONES` de `frontend-react/src/App.jsx` **y** `frontend-vue/src/App.vue`
   (misma id y label), crea el componente en cada framework, añade las
   funciones de API y los estilos (ver "Ficheros espejo").

### Ficheros espejo (deben ser IDÉNTICOS byte a byte)

React y Vue no comparten componentes, pero sí estilos y lógica pura:

- `frontend-react/src/styles.css` ↔ `frontend-vue/src/styles.css`
- `frontend-react/src/lib/vacaciones.js` ↔ `frontend-vue/src/lib/vacaciones.js`
- `frontend-react/src/lib/api.js` ↔ `frontend-vue/src/lib/api.js`
  (única diferencia permitida: la URL por defecto — 8080 en React, 8081 en Vue)

Si tocas uno, replica el cambio en el otro en el mismo commit.

## Convenciones

- **Idioma**: UI, comentarios, commits y esta documentación en español.
- **JSON**: claves en `camelCase` (`empresaEmisora`). **Columnas**: `snake_case`
  (`empresa_emisora`). **Fechas**: string `YYYY-MM-DD`.
- **Errores**: cuerpo `{"error": "mensaje"}`; `401` sin sesión, `404` si el
  recurso no existe o no pertenece al empleado de la sesión.
- **Sin Spring Security**, sin librerías nuevas salvo necesidad real.
- **Sin subida/adjuntos de archivos** — decisión de producto de la demo.
- Secciones del portal sin implementar muestran el componente `NoDisponible`.
- Config por variables de entorno con defaults para desarrollo local:
  `SERVER_PORT`, `PORTAL_DB_PATH`, `MIGRATIONS_PATH`, `CORS_ALLOWED_ORIGIN`,
  `SEED_USERNAME` (backends); `VITE_API_BASE` (frontends).
- SQLite corre en modo WAL con `busy_timeout=5000` y **una** conexión por
  proceso; no subas el pool ni quites los PRAGMA: hay dos procesos escribiendo
  el mismo fichero.

## Cómo verificar tu trabajo

```bash
# Compilar todo
make install

# Tests de contrato (arranca antes el/los backend(s): make run-java / run-go)
make verify-java     # batería contra el Java (8080)
make verify-go       # batería contra el Go (8081)
make verify          # ambos + comparación de paridad de respuestas
```

`scripts/contract-test.mjs` ejecuta la misma batería contra cada backend y,
con los dos en marcha, comprueba que respondan con el mismo status y la misma
forma de JSON. **Amplíalo cuando añadas endpoints** — es la red de seguridad
que detecta cuándo los backends divergen. El test limpia lo que crea.

Para los frontends: `npm run build` en cada uno debe compilar sin errores, y
la prueba manual es login → sección nueva → operaciones CRUD → los datos
persisten tras recargar.

## Mapa del repo

```
shared/openapi.yaml           Contrato de API (fuente de verdad nº 1)
shared/migrations/            Esquema SQLite (fuente de verdad nº 2)
data/portal.db                BBDD compartida (generada; en .gitignore)
backend-java/src/main/java/com/empresa/portal/
  ├── config/                 MigrationRunner (Order 1), DataSeeder (Order 2), CORS
  ├── model/  repo/           Entidades JPA y repositorios
  └── web/    web/dto/        Controladores REST y records DTO
backend-go/main.go            Todo el backend Go en un fichero
frontend-react/src/           App.jsx (SECCIONES), components/, lib/, styles.css
frontend-vue/src/             App.vue (SECCIONES), components/, lib/, styles.css
scripts/contract-test.mjs     Tests de contrato (make verify)
docker-compose.*.yml          Stacks java+react y go+vue (mismo ./data)
Makefile                      make help lista todos los atajos
```
