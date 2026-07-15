# ---- Build ----
FROM golang:1.25 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# pgx es Go puro: podemos compilar sin CGO (binario portable).
RUN CGO_ENABLED=0 GOOS=linux go build -o /portal-go .

# ---- Runtime ----
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /portal-go /app/portal-go

# Postgres corre como servicio "postgres" en la red de Docker Compose; las
# migraciones se montan en /migrations desde shared/migrations (ver docker-compose).
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV MIGRATIONS_PATH=/migrations
ENV SERVER_PORT=8081
EXPOSE 8081

ENTRYPOINT ["/app/portal-go"]
