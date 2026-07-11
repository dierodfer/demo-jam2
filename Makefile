# Portal de Empleado — Makefile
#
# Nota: los frontends (React/Vue) están PENDIENTES de integrar el portal, así
# que sus targets avisan en lugar de fallar. El resto está operativo.

COMPOSE_JAVA_REACT = docker compose -f docker-compose.java-react.yml
COMPOSE_GO_VUE     = docker compose -f docker-compose.go-vue.yml

.PHONY: help install install-java install-go dev run-java run-go \
        up-java-react down-java-react up-go-vue down-go-vue \
        clean frontend-pending

help: ## Muestra esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

## ---- Instalación / compilación ----

install: install-java install-go ## Compila los backends (los frontends están pendientes)
	@$(MAKE) --no-print-directory frontend-pending

install-java: ## Compila el backend Java (mvn package)
	cd backend-java && mvn -B clean package -DskipTests

install-go: ## Compila el backend Go
	cd backend-go && go build ./...

## ---- Desarrollo local ----

dev: ## Arranca ambos backends en local (en segundo plano; Ctrl-C para parar)
	@echo "Arrancando backend Java (8080) y backend Go (8081)..."
	@mkdir -p data
	@( cd backend-java && mvn -q spring-boot:run ) & \
	 ( cd backend-go && SERVER_PORT=8081 go run . ) & \
	 wait

run-java: ## Arranca solo el backend Java (8080)
	@mkdir -p data
	cd backend-java && mvn -q spring-boot:run

run-go: ## Arranca solo el backend Go (8081)
	@mkdir -p data
	cd backend-go && SERVER_PORT=8081 go run .

## ---- Docker: stack Java + React ----

up-java-react: ## Levanta backend-java (+ frontend-react cuando exista)
	$(COMPOSE_JAVA_REACT) up --build -d

down-java-react: ## Para el stack Java + React
	$(COMPOSE_JAVA_REACT) down

## ---- Docker: stack Go + Vue ----

up-go-vue: ## Levanta backend-go (+ frontend-vue cuando exista)
	$(COMPOSE_GO_VUE) up --build -d

down-go-vue: ## Para el stack Go + Vue
	$(COMPOSE_GO_VUE) down

## ---- Limpieza ----

clean: ## Limpia artefactos de build de los backends
	cd backend-java && mvn -q clean || true
	rm -rf backend-go/portal-go
	@echo "Artefactos limpiados. (data/portal.db se conserva)"

frontend-pending:
	@echo ""
	@echo "  ⛏ frontend-react y frontend-vue están PENDIENTES de integrar tu portal."
	@echo "    Ver frontend-react/README.md y frontend-vue/README.md."
