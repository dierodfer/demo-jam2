# frontend-vue

Frontend **Vue 3 + Vite 8** del Portal de Empleado. Puerto **5174**.
Mismo diseño visual que `frontend-react` (comparte `styles.css` y la lógica de
`src/lib/`, con componentes reescritos en Vue — React y Vue no comparten
componentes de forma nativa).

## Pantallas

- **Login** — logo, «Acceso al Portal del Empleado», usuario/contraseña →
  `POST /api/login` (username `admin`, cualquier contraseña).
- **Datos del empleado** (pantalla inicial) — perfil con los 7 campos,
  editable → `GET /api/me` / `PUT /api/me`.
- **Vacaciones** — calendario anual con festivos, días disfrutados y
  seleccionados, resumen por barras y leyenda (datos estáticos de demo).
- **Resto de secciones** (Nóminas, Retenciones, …) — mensaje animado de
  «Sección no disponible» para la demo.

## Configuración

| Variable | Por defecto | Descripción |
|---|---|---|
| `VITE_API_BASE` | `http://localhost:8081` (backend Go) | URL base de la API |

Las peticiones usan `credentials: 'include'` para la cookie de sesión.

En **Docker** el frontend se sirve como **build estático con nginx** en el
mismo puerto 5174 (`VITE_API_BASE` se pasa como *build-arg* y queda horneada
en el build). En local, `npm run dev` usa el dev server de Vite con hot-reload
y lee la variable del entorno.

## Comandos

```bash
npm install     # dependencias
npm run dev     # servidor de desarrollo en http://localhost:5174
npm run build   # build de producción en dist/
```

Con Docker se levanta desde la raíz del repo: `make up-go-vue`.
