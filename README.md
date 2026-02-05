# Task API

Task API — REST API сервис для управления задачами (CRUDL),
написанный на чистом Go с использованием стандартной библиотеки.

Проект создан как pet-проект для практики:
- net/http
- работы с PostgreSQL
- Docker / docker-compose
- базовой backend-архитектуры

## Стек технологий
- Go 1.22
- net/http (без фреймворков)
- PostgreSQL
- database/sql
- Docker, Docker Compose

Используется только стандартная библиотека Go (кроме драйвера PostgreSQL).

## Быстрый старт

1) Создать `.env` из шаблона:

```
cp .env.example .env
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

## UUID

ID задач — UUID (генерация на сервере).