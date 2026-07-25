# Go User API — Baby-Step Learning Plan

This is a **checklist to work through from scratch**, not a status report of a finished build. Check off steps as you go; leave notes in [`LEARNING.md`](LEARNING.md).

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

1. Create folder `user-api` and `cd` into it.
2. `go mod init github.com/joetomjob/user-api`.
3. `cmd/server/main.go` prints hello; `go run ./cmd/server`.
4. HTTP server on `:8080` returning `OK`.
5. Confirm with curl.

## Phase B — Routing without a framework

6. `GET /health` → JSON `{"status":"ok"}`.
7. `Content-Type: application/json`.
8. Named health handler.
9. Unknown paths → 404 (mux default).

## Phase C — User shape + JSON (no DB yet)

10. `internal/user/model.go` with `User` + JSON tags.
11. JSON write helper (deferred / inline encode OK).
12. `POST /users` decode + echo.
13. Practice with curl.
14. Validation → 400.
15. Invalid JSON → 400.

## Phase D — PostgreSQL locally

16. `docker-compose.yml` for Postgres.
17. `docker compose up -d` + verify.
18. `.env.example` + `.gitignore` (`.env`, `data/`).
19. Add `pgx` (+ `godotenv`).
20. Open pool from `DATABASE_URL`, `Ping`.

## Phase E — Schema + Create

21. `migrations/001_create_users.sql`.
22. Apply migration.
23. Repository holding pool.
24. `Create` with `INSERT ... RETURNING id`.
25–26. Wire Create through HTTP.
27. 201 + full user with id.
28. Unique violation → 409 (`23505`).

## Phase F — Service + remaining CRUD

29. Service layer.
30. Create via service.
31. GetByID + not-found (`ErrNotFound` / `pgx.ErrNoRows`).
32. `GET /users/{id}`.
33. Update + rows affected.
34. `PUT /users/{id}`.
35. Delete + rows affected.
36. `DELETE /users/{id}`.

## Phase G — Request timing (goroutines / channels / WaitGroup)

37. Timing logger worker + buffered channel.
38. WaitGroup; `close(ch)` + `wg.Wait()` after server stops.
39. Timing middleware + non-blocking send (`select`/`default`).
40. Wrap mux; verify with curl.

## Phase H — Error handling polish

41. Centralize API error JSON (`writeJsonError`).
42. Consistent codes: 400 / 404 / 409 / 500.
43. Context timeouts (5s) on DB calls.

## Phase I — Tests

44. Service unit tests + fake repo.
45. Handler POST success + 400.
46. Handler GET 200 + 404.
47. Handler PUT + DELETE tests.
48. Repository integration tests against Postgres.
49. `go test ./...` green.

## Phase J — Docs and cleanup

50. Write `README.md` (run, env, curls, tests).
51. Final reflection in `LEARNING.md`.
52. Cleanup pass.
53. End-to-end curl smoke test.

---

## Current next step

**Step 1:** create folder `user-api` and `cd` into it, then continue through Phase A.
