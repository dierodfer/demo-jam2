# frontend-react — ⛏ PENDIENTE

> **Estado:** a la espera de tu portal.
>
> Decidimos parar aquí el frontend hasta que pases tu portal React existente
> para integrarlo (login + perfil con los 7 campos). El resto del stack
> (backends, OpenAPI, Docker y Makefile) ya está construido y funcionando.

## Qué se integrará aquí

- **React 19.2 + Vite 8**, puerto **5173**.
- Pantalla de **login** → `POST /api/login` (username `admin`, cualquier contraseña).
- Vista de **perfil** con los 7 campos, editable → `GET /api/me` y `PUT /api/me`.
- Apunta por defecto al **backend Java** (`http://localhost:8080`), configurable
  con la variable de entorno `VITE_API_BASE`.

## Contrato

La API está documentada en [`../shared/openapi.yaml`](../shared/openapi.yaml).
Todas las peticiones deben usar `credentials: 'include'` para enviar/recibir
la cookie de sesión (`JSESSIONID`).

## Cuando integremos el portal

1. Copiar aquí el código del portal React.
2. Añadir un `Dockerfile` (puerto 5173) y `.dockerignore`.
3. Descomentar el servicio `frontend-react` en
   [`../docker-compose.java-react.yml`](../docker-compose.java-react.yml).
4. Ajustar el target `frontend-react` del [`../Makefile`](../Makefile) si hace falta.
