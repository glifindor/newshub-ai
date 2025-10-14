.PHONY: help build run stop clean proto migrate test docker-build docker-up docker-down

# Цвета для красивого вывода
BLUE=\033[0;34m
GREEN=\033[0;32m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Показать это сообщение помощи
	@echo "${BLUE}Доступные команды:${NC}"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-20s${NC} %s\n", $$1, $$2}'

# ========== Docker Commands ==========
docker-build: ## Собрать все Docker образы
	@echo "${BLUE}Building Docker images...${NC}"
	docker-compose build

docker-up: ## Запустить все сервисы в Docker
	@echo "${BLUE}Starting services...${NC}"
	docker-compose up -d
	@echo "${GREEN}Services started! Check status with 'make docker-ps'${NC}"

docker-down: ## Остановить все сервисы
	@echo "${BLUE}Stopping services...${NC}"
	docker-compose down

docker-restart: docker-down docker-up ## Перезапустить все сервисы

docker-ps: ## Показать статус контейнеров
	docker-compose ps

docker-logs: ## Показать логи всех сервисов
	docker-compose logs -f

docker-logs-gateway: ## Логи API Gateway
	docker-compose logs -f gateway

docker-logs-auth: ## Логи Auth Service
	docker-compose logs -f auth-service

docker-logs-news: ## Логи News Service
	docker-compose logs -f news-service

docker-clean: docker-down ## Удалить все контейнеры и volumes
	@echo "${RED}Removing containers and volumes...${NC}"
	docker-compose down -v
	docker system prune -f

# ========== Local Development ==========
run-auth: ## Запустить Auth Service локально
	cd auth-service && go run cmd/auth-service/main.go

run-news: ## Запустить News Service локально
	cd news-service && go run cmd/news-service/main.go

run-gateway: ## Запустить Gateway локально
	cd gateway && go run cmd/gateway/main.go

run-frontend: ## Запустить Frontend локально
	cd frontend && npm run dev

# ========== Protobuf Generation ==========
proto: proto-auth proto-news proto-seo proto-admin proto-media ## Сгенерировать все protobuf файлы

proto-auth: ## Сгенерировать protobuf для Auth Service
	@echo "${BLUE}Generating auth protobuf...${NC}"
	cd auth-service && protoc --go_out=. --go-grpc_out=. proto/auth.proto

proto-news: ## Сгенерировать protobuf для News Service
	@echo "${BLUE}Generating news protobuf...${NC}"
	cd news-service && protoc --go_out=. --go-grpc_out=. proto/news.proto

proto-seo: ## Сгенерировать protobuf для SEO Service
	@echo "${BLUE}Generating SEO protobuf...${NC}"
	cd seo-service && protoc --go_out=. --go-grpc_out=. proto/seo.proto

proto-admin: ## Сгенерировать protobuf для Admin Service
	@echo "${BLUE}Generating admin protobuf...${NC}"
	cd admin-service && protoc --go_out=. --go-grpc_out=. proto/admin.proto

proto-media: ## Сгенерировать protobuf для Media Service
	@echo "${BLUE}Generating media protobuf...${NC}"
	cd media-service && protoc --go_out=. --go-grpc_out=. proto/media.proto

# ========== Database Migrations ==========
migrate-up: ## Применить все миграции
	@echo "${BLUE}Running migrations...${NC}"
	migrate -path auth-service/migrations -database "postgresql://postgres:password@localhost:5432/auth_db?sslmode=disable" up
	migrate -path news-service/migrations -database "postgresql://postgres:password@localhost:5432/news_db?sslmode=disable" up

migrate-down: ## Откатить миграции
	@echo "${RED}Rolling back migrations...${NC}"
	migrate -path auth-service/migrations -database "postgresql://postgres:password@localhost:5432/auth_db?sslmode=disable" down
	migrate -path news-service/migrations -database "postgresql://postgres:password@localhost:5432/news_db?sslmode=disable" down

migrate-create: ## Создать новую миграцию (use: make migrate-create SERVICE=auth-service NAME=create_users)
	@if [ -z "$(SERVICE)" ] || [ -z "$(NAME)" ]; then \
		echo "${RED}Error: SERVICE and NAME are required${NC}"; \
		echo "Usage: make migrate-create SERVICE=auth-service NAME=create_users"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(SERVICE)/migrations -seq $(NAME)

# ========== Dependencies ==========
deps: deps-auth deps-news deps-gateway deps-frontend ## Установить все зависимости

deps-auth: ## Установить зависимости для Auth Service
	cd auth-service && go mod download

deps-news: ## Установить зависимости для News Service
	cd news-service && go mod download

deps-gateway: ## Установить зависимости для Gateway
	cd gateway && go mod download

deps-frontend: ## Установить зависимости для Frontend
	cd frontend && npm install

# ========== Testing ==========
test: test-auth test-news test-gateway ## Запустить все тесты

test-auth: ## Тесты Auth Service
	cd auth-service && go test -v ./...

test-news: ## Тесты News Service
	cd news-service && go test -v ./...

test-gateway: ## Тесты Gateway
	cd gateway && go test -v ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "${BLUE}Running tests with coverage...${NC}"
	cd auth-service && go test -coverprofile=coverage.out ./...
	cd news-service && go test -coverprofile=coverage.out ./...

# ========== Linting ==========
lint: ## Запустить линтеры для всех сервисов
	@echo "${BLUE}Running linters...${NC}"
	cd auth-service && golangci-lint run
	cd news-service && golangci-lint run
	cd gateway && golangci-lint run

# ========== Build ==========
build: build-auth build-news build-gateway ## Собрать все сервисы

build-auth: ## Собрать Auth Service
	cd auth-service && go build -o bin/auth-service cmd/auth-service/main.go

build-news: ## Собрать News Service
	cd news-service && go build -o bin/news-service cmd/news-service/main.go

build-gateway: ## Собрать Gateway
	cd gateway && go build -o bin/gateway cmd/gateway/main.go

# ========== Utilities ==========
init-db: ## Инициализировать базы данных
	@echo "${BLUE}Initializing databases...${NC}"
	docker exec -i news-postgres psql -U postgres < scripts/init-db.sql

clean: ## Очистить сгенерированные файлы
	@echo "${BLUE}Cleaning...${NC}"
	find . -name "*.out" -delete
	find . -type d -name "bin" -exec rm -rf {} +
	find . -name "coverage.html" -delete

health: ## Проверить здоровье всех сервисов
	@echo "${BLUE}Checking service health...${NC}"
	@curl -s http://localhost:8080/health || echo "${RED}Gateway: DOWN${NC}"
	@echo ""

# ========== Quick Start ==========
setup: deps proto ## Первоначальная настройка проекта
	@echo "${GREEN}Project setup complete!${NC}"

start: docker-up ## Быстрый старт всех сервисов
	@echo "${GREEN}All services started!${NC}"
	@echo "Frontend: http://localhost:3000"
	@echo "API Gateway: http://localhost:8080"
	@echo "Consul: http://localhost:8500"

stop: docker-down ## Остановить все сервисы

restart: stop start ## Перезапустить все сервисы

# ========== Development ==========
dev: ## Запустить в режиме разработки
	@echo "${BLUE}Starting development environment...${NC}"
	docker-compose up -d postgres redis minio rabbitmq consul
	@echo "${GREEN}Infrastructure started!${NC}"
	@echo "Run services with: make run-auth, make run-news, etc."

# ========== Monitoring ==========
monitoring-up: ## Запустить Prometheus и Grafana
	docker-compose up -d prometheus grafana
	@echo "${GREEN}Monitoring started!${NC}"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3001 (admin/admin)"

# Default target
.DEFAULT_GOAL := help
