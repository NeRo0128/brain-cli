# Brain CLI - Makefile
# Última actualización: 2026-09-05

# Variables
BINARY_NAME=brain-cli
MAIN_PATH=cmd/tui-assistant/main.go
BUILD_DIR=build
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Colores para output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: help
help: ## Muestra esta ayuda
	@echo "$(GREEN)Brain CLI - Comandos disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

.PHONY: run
run: ## Ejecuta la aplicación en modo desarrollo
	@echo "$(GREEN)Ejecutando Brain CLI...$(NC)"
	go run $(MAIN_PATH)

.PHONY: run-debug
run-debug: ## Ejecuta con debug logging activado
	@echo "$(GREEN)Ejecutando Brain CLI (modo debug)...$(NC)"
	DEBUG=1 go run $(MAIN_PATH)

.PHONY: build
build: ## Compila el binario
	@echo "$(GREEN)Compilando $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Binario creado en $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

.PHONY: build-all
build-all: ## Compila para múltiples plataformas
	@echo "$(GREEN)Compilando para múltiples plataformas...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "$(GREEN)✓ Binarios creados en $(BUILD_DIR)/$(NC)"

.PHONY: install
install: build ## Compila e instala en $GOPATH/bin
	@echo "$(GREEN)Instalando $(BINARY_NAME)...$(NC)"
	go install $(LDFLAGS) $(MAIN_PATH)
	@echo "$(GREEN)✓ Instalado en $(shell go env GOPATH)/bin/$(BINARY_NAME)$(NC)"

.PHONY: test
test: ## Ejecuta todos los tests
	@echo "$(GREEN)Ejecutando tests...$(NC)"
	go test -v ./...

.PHONY: test-cover
test-cover: ## Ejecuta tests con reporte de cobertura
	@echo "$(GREEN)Ejecutando tests con cobertura...$(NC)"
	go test -cover ./...

.PHONY: test-cover-html
test-cover-html: ## Genera reporte HTML de cobertura
	@echo "$(GREEN)Generando reporte de cobertura...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Reporte generado en coverage.html$(NC)"

.PHONY: bench
bench: ## Ejecuta benchmarks
	@echo "$(GREEN)Ejecutando benchmarks...$(NC)"
	go test -bench=. -benchmem ./...

.PHONY: lint
lint: ## Ejecuta linters de código
	@echo "$(GREEN)Ejecutando linters...$(NC)"
	go vet ./...
	go fmt ./...
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "$(YELLOW)⚠ golangci-lint no instalado$(NC)"

.PHONY: fmt
fmt: ## Formatea el código
	@echo "$(GREEN)Formateando código...$(NC)"
	go fmt ./...
	@command -v goimports >/dev/null 2>&1 && goimports -w . || echo "$(YELLOW)⚠ goimports no instalado$(NC)"

.PHONY: tidy
tidy: ## Limpia dependencias no usadas
	@echo "$(GREEN)Limpiando dependencias...$(NC)"
	go mod tidy
	@echo "$(GREEN)✓ Dependencias limpiadas$(NC)"

.PHONY: deps
deps: ## Descarga dependencias
	@echo "$(GREEN)Descargando dependencias...$(NC)"
	go mod download
	@echo "$(GREEN)✓ Dependencias descargadas$(NC)"

.PHONY: clean
clean: ## Limpia archivos generados
	@echo "$(GREEN)Limpiando archivos generados...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	go clean -cache
	@echo "$(GREEN)✓ Limpieza completada$(NC)"

.PHONY: check
check: lint test ## Ejecuta linters y tests

.PHONY: setup
setup: ## Configura el entorno de desarrollo
	@echo "$(GREEN)Configurando entorno de desarrollo...$(NC)"
	@mkdir -p data logs scripts
	@touch .env
	@echo "OMNIROUTE_API_KEY=your-api-key-here" > .env.example
	@echo "$(GREEN)✓ Entorno configurado$(NC)"
	@echo "$(YELLOW)Recuerda configurar tu .env con las credenciales necesarias$(NC)"

.PHONY: doctor
doctor: ## Verifica que todo esté configurado correctamente
	@echo "$(GREEN)Verificando configuración...$(NC)"
	@echo ""
	@echo "$(YELLOW)Go version:$(NC)"
	@go version
	@echo ""
	@echo "$(YELLOW)Go environment:$(NC)"
	@echo "  GOPATH: $(shell go env GOPATH)"
	@echo "  GOROOT: $(shell go env GOROOT)"
	@echo ""
	@echo "$(YELLOW)Dependencias instaladas:$(NC)"
	@command -v golangci-lint >/dev/null 2>&1 && echo "  ✓ golangci-lint" || echo "  ✗ golangci-lint (opcional)"
	@command -v goimports >/dev/null 2>&1 && echo "  ✓ goimports" || echo "  ✗ goimports (opcional)"
	@echo ""
	@echo "$(YELLOW)Directorios:$(NC)"
	@test -d configs && echo "  ✓ configs/" || echo "  ✗ configs/"
	@test -d scripts && echo "  ✓ scripts/" || echo "  ✗ scripts/"
	@test -d internal && echo "  ✓ internal/" || echo "  ✗ internal/"
	@echo ""
	@echo "$(GREEN)✓ Verificación completada$(NC)"

.DEFAULT_GOAL := help
