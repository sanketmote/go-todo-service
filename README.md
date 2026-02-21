# Go Todo Service

[![CI](https://github.com/sanketmote/go-todo-service/actions/workflows/ci.yml/badge.svg)](https://github.com/sanketmote/go-todo-service/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A REST API backend for To-Do CRUD operations, built in Go with MySQL.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Project Structure](#project-structure)
- [License](#license)

## Features

- **CRUD operations** — Create, read, update, and delete todos
- **Filtering & pagination** — List todos with `include_completed`, `sort`, `limit`, `offset`
- **Recurring todos** — Support for `none`, `daily`, `weekly`, `monthly`, `yearly` repeat types
- **Health checks** — `/ping` (liveness) and `/healthz` (detailed, token-protected)
- **Docker ready** — Run app + MySQL with `docker-compose up`

## Tech Stack

| Component | Technology |
|-----------|------------|
| Router | [gorilla/mux](https://github.com/gorilla/mux) |
| Transport | [go-kit](https://github.com/go-kit/kit) |
| Logger & Health | [gokit-wrapper](https://github.com/sanketmote/gokit-wrapper) |
| Database | MySQL 8, [sqlx](https://github.com/jmoiron/sqlx) |

## Prerequisites

- **Go 1.24+**
- **Docker & Docker Compose** (for running MySQL or full stack)

## Quick Start

### Option 1: Docker Compose (recommended)

```bash
git clone https://github.com/sanketmote/go-todo-service.git
cd go-todo-service
docker-compose up
```

- **API:** http://localhost:8080  
- **MySQL:** localhost:3306

### Option 2: Local development

1. Start MySQL (or use existing instance).

2. Set environment variables and run:

```bash
export MYSQL_DSN="root:root@tcp(localhost:3306)/todos?parseTime=true"
export XDRV_HEALTHZ_TOKEN="dev"
go run ./cmd/server
```

3. Or build and run the binary:

```bash
go build -o server ./cmd/server
./server
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP listen address | `:8080` |
| `MYSQL_DSN` / `DB_DSN` | MySQL connection string | *(required)* |
| `XDRV_HEALTHZ_TOKEN` / `HEALTHZ_TOKEN` | Token for `GET /healthz` | — |

## API Reference

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/ping` | Liveness check |
| GET | `/healthz` | Health details (requires `Xdrv-Healthz-Token` header) |
| POST | `/api/v1/todos` | Create todo |
| GET | `/api/v1/todos` | List todos |
| GET | `/api/v1/todos/{id}` | Get todo by ID |
| PUT | `/api/v1/todos/{id}` | Update todo (partial) |
| DELETE | `/api/v1/todos/{id}` | Delete todo |

### List query parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `include_completed` | boolean | `false` | Include completed todos |
| `sort` | string | `due_date_asc` | `due_date_asc`, `due_date_desc`, `created_at_asc`, `created_at_desc` |
| `limit` | integer | `50` | Max items (max 100) |
| `offset` | integer | `0` | Pagination offset |

### cURL examples

```bash
# Health
curl http://localhost:8080/ping
curl -H "Xdrv-Healthz-Token: dev" http://localhost:8080/healthz

# Create todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy groceries","description":"Milk, eggs, bread","due_date":"2025-02-25T18:00:00Z","repeat_type":"none"}'

# List todos
curl http://localhost:8080/api/v1/todos
curl "http://localhost:8080/api/v1/todos?include_completed=true&sort=due_date_desc&limit=10&offset=0"

# Get todo by ID
curl http://localhost:8080/api/v1/todos/1

# Update todo (partial)
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated title","completed":true,"due_date":"2025-02-26T18:00:00Z"}'

# Delete todo
curl -X DELETE http://localhost:8080/api/v1/todos/1
```

### API clients

- **OpenAPI spec:** `docs/openapi.yaml` — OpenAPI 3.0 spec (validated in CI)
- **Postman:** Import `docs/go-todo-service.postman_collection.json` (variables: `baseUrl`, `healthzToken`, `todoId`)

## Project Structure

```
.
├── cmd/
│   └── server/          # Application entrypoint
├── internal/
│   ├── config/          # Configuration loading
│   ├── env/             # Environment helpers
│   ├── handler/         # HTTP middleware (recovery, logging, request ID)
│   ├── repository/      # Todo repository (MySQL)
│   ├── spec/
│   │   ├── datalayer/   # DB models and DDL
│   │   └── todo/        # API spec (paths, request/response, constants)
│   ├── svcerror/        # Service errors and HTTP mapping
│   ├── todo/            # Business logic, endpoints, transport
│   └── upgrade/         # Database migrations
├── docs/
│   ├── openapi.yaml     # OpenAPI 3.0 spec
│   └── go-todo-service.postman_collection.json
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

## Tests

```bash
go test ./...
```

## License

Licensed under the [Apache License 2.0](LICENSE).
