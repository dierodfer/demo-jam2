# Portal de Empleado — Makefile

COMPOSE_BASE       = docker compose -f docker-compose.yml
COMPOSE_JAVA_REACT = docker compose -f docker-compose.yml -f docker-compose.java-react.yml
COMPOSE_GO_VUE     = docker compose -f docker-compose.yml -f docker-compose.go-vue.yml

.PHONY: help install install-java install-go install-react install-vue \
        dev run-java run-go run-react run-vue \
        db-up db-down verify verify-java verify-go \
        up-java-react down-java-react up-go-vue down-go-vue clean

help: ## Muestra esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

## ---- Instalación / compilación ----

install: install-java install-go install-react install-vue ## Instala/compila los 4 proyectos

install-java: ## Compila el backend Java (mvn package)
	cd backend-java && mvn -B clean package -DskipTests

install-go: ## Compila el backend Go
	cd backend-go && go build ./...

install-react: ## Instala dependencias del frontend React
	cd frontend-react && npm install

install-vue: ## Instala dependencias del frontend Vue
	cd frontend-vue && npm install

## ---- Base de datos (Postgres compartido) ----

db-up: ## Levanta solo el Postgres compartido (necesario para run-java/run-go/dev)
	$(COMPOSE_BASE) up -d --wait postgres

db-down: ## Para el Postgres compartido
	$(COMPOSE_BASE) down

## ---- Desarrollo local ----

dev: db-up ## Arranca los 4 proyectos en local contra el Postgres de Docker (Ctrl-C para parar)
	@echo "Backends: Java :8080, Go :8081 — Frontends: React :5173, Vue :5174"
	@( cd backend-java && mvn -q spring-boot:run ) & \
	 ( cd backend-go && SERVER_PORT=8081 go run . ) & \
	 ( cd frontend-react && npm run dev ) & \
	 ( cd frontend-vue && npm run dev ) & \
	 wait

run-java: ## Arranca solo el backend Java (8080; requiere make db-up)
	cd backend-java && mvn -q spring-boot:run

run-go: ## Arranca solo el backend Go (8081; requiere make db-up)
	cd backend-go && SERVER_PORT=8081 go run .

run-react: ## Arranca solo el frontend React (5173, dev server contra backend Java)
	cd frontend-react && npm run dev

run-vue: ## Arranca solo el frontend Vue (5174, dev server contra backend Go)
	cd frontend-vue && npm run dev

## ---- Tests de contrato (requieren el backend arrancado) ----

verify: ## Contrato contra ambos backends (8080 y 8081) + comparación de paridad
	node scripts/contract-test.mjs http://localhost:8080 http://localhost:8081

verify-java: ## Contrato contra el backend Java (8080)
	node scripts/contract-test.mjs http://localhost:8080

verify-go: ## Contrato contra el backend Go (8081)
	node scripts/contract-test.mjs http://localhost:8081

## ---- Docker: stack Java + React ----

up-java-react: ## Levanta postgres + backend-java (8080) + frontend-react estático (5173)
	$(COMPOSE_JAVA_REACT) up --build -d

down-java-react: ## Para el stack Java + React
	$(COMPOSE_JAVA_REACT) down

## ---- Docker: stack Go + Vue ----

up-go-vue: ## Levanta postgres + backend-go (8081) + frontend-vue estático (5174)
	$(COMPOSE_GO_VUE) up --build -d

down-go-vue: ## Para el stack Go + Vue
	$(COMPOSE_GO_VUE) down

## ---- Limpieza ----

clean: ## Limpia artefactos de build de los 4 proyectos
	cd backend-java && mvn -q clean || true
	rm -rf backend-go/portal-go
	rm -rf frontend-react/dist frontend-vue/dist
	@echo "Artefactos limpiados. (node_modules y el volumen de Postgres se conservan)"
