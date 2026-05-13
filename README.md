# go-platform

Монорепозиторий микросервисов на **Go** с общим пакетом `pkg`, оркестрацией через **Docker Compose** и инфраструктурой **PostgreSQL**, **Redis**, **Kafka** и **ClickHouse**.

- **Модуль:** `github.com/palach608/go-platform`
- **Go:** 1.26.2

## Сервисы

| Сервис | Порт | Роль |
|--------|------|------|
| **chat** | 8080 | Комнаты, сообщения, WebSocket; Redis; consumer Kafka |
| **auth** | 8081 | Регистрация, логин, JWT, операции с пользователем |
| **notes** | 8082 | CRUD заметок в PostgreSQL; producer событий в Kafka |
| **analytics** | 8083 | Consumer Kafka → ClickHouse; REST со статистикой |

События: сервис **notes** публикует в топик `note.created`. **chat** и **analytics** читают его в группах `chat-group` и `analytics-group`.

## Стек

- HTTP: [chi](https://github.com/go-chi/chi), Swagger: [http-swagger](https://github.com/swaggo/http-swagger)
- DI: [fx](https://github.com/uber-go/fx)
- PostgreSQL: [pgx](https://github.com/jackc/pgx), миграции — [golang-migrate](https://github.com/golang-migrate/migrate) (notes, chat), [goose](https://github.com/pressly/goose) (auth в Docker)
- Redis, Kafka ([kafka-go](https://github.com/segmentio/kafka-go)), ClickHouse ([clickhouse-go](https://github.com/ClickHouse/clickhouse-go))
- JWT: [golang-jwt](https://github.com/golang-jwt/jwt)

## Запуск через Docker

Из каталога `deployments`:

```powershell
docker compose up --build
```

При первом старте Postgres выполняется `deployments/init-db.sql` (БД `auth_db`, `notes_db`, `chat_db`). Миграции auth в контейнере выполняет goose; notes и chat — golang-migrate при старте приложения.

## Swagger

| Сервис | UI | Спецификация |
|--------|-----|--------------|
| auth | http://localhost:8081/swagger/index.html | http://localhost:8081/api/swagger.yaml |
| chat | http://localhost:8080/swagger/index.html | http://localhost:8080/api-docs/swagger.yaml |
| notes | http://localhost:8082/swagger/index.html | http://localhost:8082/api-docs/swagger.yaml |
| analytics | http://localhost:8083/swagger/index.html | http://localhost:8083/api-docs/swagger.yaml |

`JWT_SECRET` должен совпадать у **auth**, **notes** и **chat** — в `deployments/docker-compose.yaml` задано одно значение для всех.

## Локальный запуск (без Docker)

Нужны экземпляры Postgres (три базы из `deployments/init-db.sql`), Redis, Kafka и ClickHouse с DSN/адресами под вашу среду. Переменные окружения — как в таблице ниже (ориентир: `deployments/docker-compose.yaml`).

Из корня репозитория:

```powershell
go run ./services/auth/cmd/main.go
go run ./services/notes/cmd/server/main.go
go run ./services/chat/cmd/server/main.go
go run ./services/analytics/cmd/main.go
```

Для **notes** и **chat** вызывается `godotenv.Load()` — параметры можно задать в `.env` в рабочей директории процесса.

## Переменные окружения

| Переменная | Сервисы | Пример (как в compose) |
|------------|---------|-------------------------|
| `DB_DSN` | auth, notes, chat | `postgres://postgres:4019@localhost:5432/auth_db?sslmode=disable` (для каждой БД свой URL) |
| `HTTP_PORT` | все | `8081` / `8082` / `8080` / `8083` |
| `JWT_SECRET` | auth, notes, chat | общий секрет для выдачи и проверки токена |
| `KAFKA_BROKERS` | notes, chat, analytics | `localhost:9092` или `kafka:9092` внутри сети compose |
| `REDIS_DSN` | chat | `redis://localhost:6379` |
| `CLICKHOUSE_DSN` | analytics | `clickhouse://localhost:9000/default` |

Дефолты для **auth** задаются в `services/auth/internal/config/config.go` (в т.ч. `DB_DSN` и `HTTP_PORT`, если переменные не заданы).

## HTTP-маршруты (обзор)

**auth** (`/auth`):

- `POST /auth/register`, `POST /auth/login`
- `PUT /auth/user`, `DELETE /auth/user/{id}`

**notes** (`/api/v1/notes`, часть ручек за JWT middleware):

- `GET /api/v1/notes` — без токена
- `POST`, `GET/{id}`, `PUT/{id}`, `DELETE/{id}` — с `Authorization: Bearer …`

**chat** (`/api/v1/...`, за JWT):

- `POST /api/v1/rooms`, `GET /api/v1/rooms`, `GET /api/v1/rooms/{id}`, `DELETE /api/v1/rooms/{id}`
- `GET /api/v1/rooms/{id}/messages`
- `GET /api/v1/ws/{roomID}` — WebSocket

Статика UI чата: `GET /*` отдаёт файлы из `./web` внутри образа/сборки.

**analytics** (без JWT в роутере):

- `GET /api/v1/stats/top-users`
- `GET /api/v1/stats/notes-per-day`

## Структура репозитория

```
pkg/                 JWT middleware, Kafka producer/consumer, события
services/
  auth/              пользователи, JWT
  notes/             заметки, Kafka producer
  chat/              комнаты, сообщения, WebSocket, Redis, Kafka consumer
  analytics/         Kafka consumer, ClickHouse, REST stats
migrations/auth/     SQL-миграции goose для auth
deployments/         Dockerfile.*, docker-compose.yaml, init-db.sql
```

## Сборка

```powershell
go build ./...
```

Образы собираются из `deployments/Dockerfile.*` (контекст сборки — корень репозитория).
