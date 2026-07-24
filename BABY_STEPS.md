# Go User API — Baby-Step Learning Plan

Build a Go user-management microservice from scratch using **PostgreSQL** and stdlib **`net/http`**. Work one tiny step at a time. Record challenges in [`LEARNING.md`](LEARNING.md). [`README.md`](README.md) (Phase J) explains how to run and test.

## How we work

- Do not jump ahead — each step should leave something you can run or test.
- Prefer figuring out the code yourself; ask for help when stuck.
- After each step, note what you did and what went wrong in `LEARNING.md`.

### `LEARNING.md` entry format

```markdown
## Step N — short title (YYYY-MM-DD)

**What I did:** …
**Why:** … (one sentence)
**Commands / files:** …
**Challenges / mistakes:** …
**What I learned:** …
```

## Target end state

```text
Client → Mux → Timing middleware → Handler → Service → Repo → Postgres
                      ↓
                 Log worker (goroutine + channel + WaitGroup)
```

- CRUD: `POST /users`, `GET /users/{id}`, `PUT /users/{id}`, `DELETE /users/{id}`
- Fields: `name`, `email`, `age` (+ `id` from DB)
- Request timing via goroutine, channel, `sync.WaitGroup`
- Unit tests; README with run/test instructions

Suggested layout:

```text
user-api/
  cmd/server/main.go
  internal/user/       # model, handler, service, repository
  migrations/
  docker-compose.yml
  .env.example
  BABY_STEPS.md        # this plan
  LEARNING.md          # journal + challenges
  README.md            # how to run / test (Phase J)
```

## Defaults

- Module: `github.com/joetomjob/user-api`
- Driver: `pgx` + pool
- IDs: `SERIAL` (UUID also fine)
- No Gin/Chi/ORM

---

## Phase A — Empty project that runs

1. ~~Create folder `user-api` and `cd` into it.~~ **done**
2. ~~`go mod init github.com/joetomjob/user-api`.~~ **done**
3. ~~`cmd/server/main.go` prints hello; `go run ./cmd/server`.~~ **done**
4. ~~HTTP server on `:8080` returning `OK`.~~ **done**
5. ~~Confirm with curl.~~ **done**

## Phase B — Routing without a framework

6. ~~`GET /health` → JSON `{"status":"ok"}`.~~ **done**
7. ~~`Content-Type: application/json`.~~ **done**
8. ~~Named health handler.~~ **done**
9. ~~Unknown paths → 404 (mux default).~~ **done**

## Phase C — User shape + JSON (no DB yet)

10. ~~`internal/user/model.go` with `User` + JSON tags.~~ **done**
11. ~~JSON write helper (deferred / inline encode OK).~~ **done**
12. ~~`POST /users` decode + echo.~~ **done**
13. ~~Practice with curl.~~ **done**
14. ~~Validation → 400.~~ **done**
15. ~~Invalid JSON → 400.~~ **done**

## Phase D — PostgreSQL locally

16. ~~`docker-compose.yml` for Postgres.~~ **done**
17. ~~`docker compose up -d` + verify.~~ **done**
18. ~~`.env.example` + `.gitignore` (`.env`, `data/`).~~ **done**
19. ~~Add `pgx` (+ `godotenv`).~~ **done**
20. ~~Open pool from `DATABASE_URL`, `Ping`.~~ **done**

## Phase E — Schema + Create

21. ~~`migrations/001_create_users.sql`.~~ **done**
22. ~~Apply migration.~~ **done**
23. ~~Repository holding pool.~~ **done**
24. ~~`Create` with `INSERT ... RETURNING id`.~~ **done**
25–26. ~~Wire Create through HTTP.~~ **done**
27. ~~201 + full user with id.~~ **done**
28. ~~Unique violation → 409 (`23505`).~~ **done**

## Phase F — Service + remaining CRUD

29. ~~Service layer.~~ **done**
30. ~~Create via service.~~ **done**
31. ~~GetByID + not-found (`ErrNotFound` / `pgx.ErrNoRows`).~~ **done**
32. ~~`GET /users/{id}`.~~ **done**
33. ~~Update + rows affected.~~ **done**
34. ~~`PUT /users/{id}`.~~ **done**
35. ~~Delete + rows affected.~~ **done**
36. ~~`DELETE /users/{id}`.~~ **done**

## Phase G — Request timing (goroutines / channels / WaitGroup)

37. ~~Timing logger worker + buffered channel.~~ **done**
38. ~~WaitGroup; `close(ch)` + `wg.Wait()` after server stops.~~ **done**
39. ~~Timing middleware + non-blocking send (`select`/`default`).~~ **done**
40. ~~Wrap mux; verify with curl.~~ **done**

## Phase H — Error handling polish

41. ~~Centralize API error JSON (`writeJsonError`).~~ **done**
42. ~~Consistent codes: 400 / 404 / 409 / 500.~~ **done**
43. ~~Context timeouts (5s) on DB calls.~~ **done**

## Phase I — Tests (next)

44. Unit-test service with a **fake** repository interface (no Docker).
45. Handler tests (`httptest`) for `POST` success and 400.
46. Handler tests for `GET` 200 and 404.
47. Handler tests for `PUT` and `DELETE`.
48. Repository integration tests against Postgres.
49. `go test ./...` green.

## Phase J — Docs and cleanup

50. Write `README.md` (run, env, curls, tests).
51. Final reflection in `LEARNING.md`.
52. Cleanup pass.
53. End-to-end curl smoke test.

---

## Current next step

**Step 44:** define a small repository interface, point `Service` at it, write a fake repo in a test, and unit-test validation (e.g. empty name) without Docker.
