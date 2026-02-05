.PHONY: run build test tidy lint lint-install vet fmt clean deps swagger\
	docker-up docker-down migrate-up migrate-down

APP_NAME := go-tasks-api
BIN_DIR := bin
CMD_DIR := ./cmd/app
MIGRATIONS_DIR := ./migrations

GO ?= go
GOLANGCI_LINT ?= golangci-lint


run:
	@echo "Запуск приложения..."
	@if [ -f .env.local ]; then \
		set -a; \
		. ./.env.local; \
		set +a; \
		$(GO) run $(CMD_DIR); \
	else \
		echo "Файл .env.local не найден. Создайте его из .env.example"; \
		exit 1; \
	fi
	

build: $(BIN_DIR)/$(APP_NAME)

$(BIN_DIR)/$(APP_NAME):
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

fmt:
	$(GO) fmt ./...
	$(GO) vet ./...

lint:
	$(GOLANGCI_LINT) run ./...

lint-install:
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

vet:
	$(GO) vet ./...

clean:
	rm -rf $(BIN_DIR)

deps:
	@echo "Установка зависимостей..."
	$(GO) mod download
	$(GO) mod tidy

swagger:
	swag init -g cmd/app/main.go

docker-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down

migrate-up:
	@if ! docker ps --format '{{.Names}}' | grep -q '^tasks-postgres$$'; then \
		echo "PostgreSQL container not running. Run 'make docker-up' first."; \
		exit 1; \
	fi
	@for f in $(MIGRATIONS_DIR)/*.up.sql; do \
		echo "apply $$f"; \
		docker exec -i tasks-postgres psql -U postgres -d tasks < $$f; \
	done

migrate-down:
	@if ! docker ps --format '{{.Names}}' | grep -q '^tasks-postgres$$'; then \
		echo "PostgreSQL container not running. Run 'make docker-up' first."; \
		exit 1; \
	fi
	@for f in $$(ls -r $(MIGRATIONS_DIR)/*.down.sql); do \
		echo "rollback $$f"; \
		docker exec -i tasks-postgres psql -U postgres -d tasks < $$f; \
	done
