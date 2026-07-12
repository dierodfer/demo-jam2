# frontend-react

Frontend **React 19.2 + Vite 8** del Portal de Empleado. Puerto **5173**.
Réplica del portal original de Nunegal (login, cabecera y navegación por pestañas).

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
| `VITE_API_BASE` | `http://localhost:8080` (backend Java) | URL base de la API |

Las peticiones usan `credentials: 'include'` para la cookie de sesión.

## Comandos

```bash
npm install     # dependencias
npm run dev     # servidor de desarrollo en http://localhost:5173
npm run build   # build de producción en dist/
```

Con Docker se levanta desde la raíz del repo: `make up-java-react`.
