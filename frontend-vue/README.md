# frontend-vue — ⛏ PENDIENTE

> **Estado:** a la espera de tu portal.
>
> Igual que `frontend-react`, este frontend se construirá cuando integremos
> tu portal. Reutilizará el **mismo diseño visual** (estilos/tokens, no
> componentes — React y Vue no comparten componentes de forma nativa).

## Qué se integrará aquí

- **Vue 3 + Vite 8**, puerto **5174**.
- Mismas pantallas: **login** y **perfil** con los 7 campos.
- Apunta por defecto al **backend Go** (`http://localhost:8081`), configurable
  con la variable de entorno `VITE_API_BASE`.

## Contrato

La API está documentada en [`../shared/openapi.yaml`](../shared/openapi.yaml) —
es el mismo contrato que consume el frontend React. Todas las peticiones deben
usar `credentials: 'include'` para enviar/recibir la cookie de sesión
(`session_id`).

## Cuando integremos el portal

1. Crear aquí el proyecto Vue 3 + Vite 8 con el mismo diseño visual que el React.
2. Añadir un `Dockerfile` (puerto 5174) y `.dockerignore`.
3. Descomentar el servicio `frontend-vue` en
   [`../docker-compose.go-vue.yml`](../docker-compose.go-vue.yml).
4. Ajustar el target `frontend-vue` del [`../Makefile`](../Makefile) si hace falta.
