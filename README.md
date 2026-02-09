# Task API

Task API — REST API сервис для управления задачами (CRUDL),
написанный на чистом Go с использованием стандартной библиотеки.

Проект создан как pet-проект для практики:
- net/http
- работы с PostgreSQL
- Docker / docker-compose
- базовой backend-архитектуры

## Стек технологий
- Go 1.23 (toolchain 1.24)
- net/http
- PostgreSQL
- database/sql
- Docker, Docker Compose
- Swagger
- GoMock
- UUID

## Быстрый старт

1) Создать `.env` из шаблона:

```
cp .env.example deployments/.env
```

2) Поднять сервисы:

```
make docker-up
```

3) Накатить миграции:

```
make migrate-up
```

## Swagger

Генерация документации:

```
make swagger
```

Swagger UI:

- http://localhost:8080/swagger/index.html

## Переменные окружения

### Сервер
- `SERVER_HOST` (по умолчанию `0.0.0.0`)
- `SERVER_PORT` (по умолчанию `8080`)
- `SERVER_READ_TIMEOUT_SECONDS` (по умолчанию `5`)
- `SERVER_WRITE_TIMEOUT_SECONDS` (по умолчанию `10`)
- `SERVER_IDLE_TIMEOUT_SECONDS` (по умолчанию `60`)

### PostgreSQL
- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`
- `POSTGRES_SSLMODE`

## Линтер

```
make lint
```

### Покрытие

По последнему запуску `go test ./... -cover`:
- `internal/api`: 70.7%
- `internal/service`: 84.9%

## UUID

ID задач — UUID (генерация на сервере).
