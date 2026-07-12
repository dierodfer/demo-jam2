# ---- Build ----
FROM golang:1.24 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# modernc.org/sqlite es Go puro: podemos compilar sin CGO (binario portable).
RUN CGO_ENABLED=0 GOOS=linux go build -o /portal-go .

# ---- Runtime ----
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /portal-go /app/portal-go

# La base de datos SQLite vive en el volumen montado en /data; las migraciones
# se montan en /migrations desde shared/migrations (ver docker-compose).
ENV PORTAL_DB_PATH=/data/portal.db
ENV MIGRATIONS_PATH=/migrations
ENV SERVER_PORT=8081
EXPOSE 8081

ENTRYPOINT ["/app/portal-go"]
