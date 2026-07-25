# user-api

Go REST API for user management. Uses stdlib `net/http`, PostgreSQL (`pgx`), and `godotenv`.

Module: `github.com/joetomjob/user-api`

## Prerequisites

- Go 1.26+
- Docker / Docker Compose
- `psql` (or any Postgres client) to apply the migration

## Project layout

```text
cmd/server/main.go      # HTTP server entrypoint
internal/user/          # model, handler, service, repository, tests
migrations/             # SQL schema
docker-compose.yml      # local Postgres
.env.example            # sample DATABASE_URL
BABY_STEPS.md           # build plan
LEARNING.md             # learning journal
```

## Setup

### 1. Start Postgres

```bash
docker compose up -d
```

Postgres 17 listens on `localhost:5432` (`admin` / `admin`, database `users`). Data is stored in `./data`.

### 2. Configure env

```bash
cp .env.example .env
```

`.env.example`:

```text
DATABASE_URL=postgres://admin:admin@localhost/users?sslmode=disable
```

### 3. Apply migration

```bash
psql "$DATABASE_URL" -f migrations/001_create_users.sql
```

Or:

```bash
psql postgres://admin:admin@localhost/users?sslmode=disable -f migrations/001_create_users.sql
```

Creates table `users` (`id SERIAL PRIMARY KEY`, unique `email`, unique `name`, `age SMALLINT`).

### 4. Run the server

From the repo root (so `.env` loads):

```bash
go run ./cmd/server
```

Server listens on `:8080`. Each request is timed by middleware; duration is logged asynchronously (goroutine + channel + `WaitGroup`).

## API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/users` | Create user |
| `GET` | `/users/{id}` | Get user by ID |
| `PUT` | `/users/{id}` | Update user |
| `DELETE` | `/users/{id}` | Delete user |

User JSON fields: `name`, `email`, `age` (response also includes `id`).

### Example curls

```bash
# Health
curl -s localhost:8080/health

# Create
curl -s -X POST localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ada","email":"ada@example.com","age":36}'

# Get (replace 1 with a real id)
curl -s localhost:8080/users/1

# Update
curl -s -X PUT localhost:8080/users/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ada Lovelace","email":"ada@example.com","age":37}'

# Delete
curl -s -X DELETE localhost:8080/users/1
```

## Tests

```bash
go test ./...
```

- Service and handler tests use a fake repository (no Docker required).
- Repository integration tests need Postgres and load `../../.env` from `internal/user/`; they skip if the DB is unavailable.

## Further reading

- [BABY_STEPS.md](BABY_STEPS.md) — step-by-step build plan
- [LEARNING.md](LEARNING.md) — notes and challenges from building the project
