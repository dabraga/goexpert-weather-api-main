# ==============================================================================
# Variáveis
# ==============================================================================
APP_NAME=weather-api-lab
PORT?=8080

# Cores
BLUE=\033[0;34m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m

# ==============================================================================
# Comandos de Desenvolvimento
# ==============================================================================
.PHONY: setup run test lint clean help

setup: ## Configura o ambiente
	@echo "$(BLUE)🔧 Configurando ambiente...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Ambiente configurado!$(NC)"
	@echo "$(YELLOW)📝 Configure WEATHER_API_KEY como variável de ambiente$(NC)"

run: ## Roda a aplicação
	@echo "$(BLUE)🚀 Iniciando aplicação na porta $(PORT)...$(NC)"
	@go run ./cmd/api

test: ## Roda os testes
	@echo "$(BLUE)🧪 Executando testes...$(NC)"
	@go test -v ./...

test-coverage: ## Roda os testes com cobertura
	@echo "$(BLUE)🧪 Executando testes com cobertura...$(NC)"
	@go test -v -cover ./...

lint: ## Executa o linter
	@echo "$(BLUE)🔍 Verificando código...$(NC)"
	@go vet ./...
	@go fmt ./...

clean: ## Limpa arquivos temporários
	@echo "$(BLUE)🧹 Limpando arquivos temporários...$(NC)"
	@go clean
	@echo "$(GREEN)✅ Limpeza concluída!$(NC)"

# ==============================================================================
# Comandos Docker
# ==============================================================================
docker-build: ## Build da imagem Docker
	@echo "$(BLUE)🐳 Construindo imagem Docker...$(NC)"
	@docker build -t $(APP_NAME) .

docker-run: ## Roda o container
	@echo "$(BLUE)🐳 Iniciando container...$(NC)"
	@docker run -p $(PORT):$(PORT) -e WEATHER_API_KEY=$(WEATHER_API_KEY) $(APP_NAME)

# ==============================================================================
# Ajuda
# ==============================================================================
help: ## Mostra essa ajuda
	@echo "$(BLUE)Comandos disponíveis:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
