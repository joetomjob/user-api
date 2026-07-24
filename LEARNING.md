# Learning journal — user-api

Notes from building this project step by step. Updated as we go, including mistakes and fixes.

---

## Step 1 — Create the project folder (2026-07-23)

**What I did:** Created a new folder `user-api` under `Projects` (sibling of `golang-interview-prep`) so this interview exercise starts from scratch.

**Why:** The old project felt overwhelming; building one small thing at a time is clearer.

**Commands / files:**
- Folder: `/Users/joetomjob/Joe-personal/Projects/user-api`

**Challenges / mistakes:** None.

**What I learned:** A clean empty repo is a better starting point than continuing a large unfinished codebase.

---

## Step 2 — Initialize the Go module (2026-07-23)

**What I did:** Ran `go mod init` in the new folder. Module path: `github.com/joetomjob/user-api`.

**Why:** Every Go project needs a `go.mod` so the toolchain knows the module path and Go version.

**Commands / files:**
- `go mod init github.com/joetomjob/user-api`
- Created `go.mod` (Go 1.26.5)

**Challenges / mistakes:** None.

**What I learned:** `go.mod` is the project’s identity for imports and dependency management.

---

## Step 3 — Hello world entrypoint in `cmd/server` (2026-07-23)

**What I did:** Created `cmd/server/main.go` so the program prints hello and can be run with `go run ./cmd/server`.

**Why:** `cmd/server` is a common Go layout for an executable. Shared library code will live elsewhere later; this folder is only the program entrypoint.

**Commands / files:**
- `cmd/server/main.go`
- `go run ./cmd/server`

**Challenges / mistakes:**
- First tried `package server` because the folder is named `server`. That does not produce a runnable binary.
- Fixed by using `package main`, defining `func main()`, and importing `fmt` to print.

**What I learned:**
- Folder name and package name are not the same thing.
- Only `package main` + `func main()` make an executable you can `go run`.
- `fmt` is needed for printing to the console.

---

## Step 4 — Basic HTTP server with ServeMux (2026-07-23)

**What I did:** Replaced the hello print with an HTTP server in `cmd/server/main.go`. Created a `ServeMux` router, registered a handler for `/`, built an `http.Server` on `:8080`, and started it with `ListenAndServe`.

**Why:** The service needs to accept HTTP requests; this is the smallest real server.

**Commands / files:**
- `cmd/server/main.go`
- `go run ./cmd/server`

**Challenges / mistakes:**
- Had to figure out how to initialize a router (`http.NewServeMux`).
- Had to figure out how to create and start the server (`http.Server` + `ListenAndServe`).
- Had to figure out how to attach a handler function to a path (`mux.HandleFunc`).
- Had to figure out how to set status OK (`w.WriteHeader(http.StatusOK)`) and write `"OK"` to the response body (`fmt.Fprintln`).

**What I learned:**
- `ServeMux` maps paths to handler functions.
- Handler signature is `func(http.ResponseWriter, *http.Request)`.
- `http.Server` holds address, handler, and timeouts; `ListenAndServe` blocks and serves requests.
- Writing the status and body goes through `ResponseWriter`.

**Note:** Also set `ReadTimeout` / `WriteTimeout` / `IdleTimeout` and logged fatal errors from `ListenAndServe` (beyond the minimum for this step — good practice).

---

## Step 5 — Verify with curl (2026-07-23)

**What I did:** Confirmed the running server responds correctly (e.g. `curl localhost:8080` returns `OK`).

**Why:** Prove the server actually works from outside the process, not only that it compiles.

**Commands / files:**
- `curl localhost:8080` (or equivalent)

**Challenges / mistakes:** Covered under Step 4 (router/server/handler/status/body).

**What I learned:** End-to-end check with curl is the quick smoke test before adding more routes.

---

## Step 6–9 — Health JSON route + 404 behavior (2026-07-23)

**What I did:**
- Replaced the `/` OK handler with a `health` handler on `/health`.
- Returned JSON `{"status":"ok"}` using `encoding/json` and a `map[string]string`.
- Set `Content-Type` to `application/json`.
- Added a `handleOthers` handler on `/*` that writes 404 / `"Not found"`.
- Kept the named handler + `ServeMux` + `http.Server` setup from earlier.

**Why:** Health checks are a standard first real API route; JSON + content type prepare us for user APIs; confirming 404 behavior shows the mux only serves registered paths.

**Commands / files:**
- `cmd/server/main.go`
- Check: `curl localhost:8080/health` and `curl localhost:8080/nope` (or similar)

**Challenges / mistakes:**
- Setting `Content-Type` to `application/json` on the response header was new.
- Encoding a value to JSON in the response body (`encoding/json`) was new.
- Initially added a custom `handleOthers` on `/*` for 404s — not needed. Didn’t know `ServeMux` already returns 404 for paths that aren’t registered.
- Had `WriteHeader` before setting headers at first; moved `WriteHeader` after setting `Content-Type`.

**What I learned / review notes:**
- Expected for these steps: `/health` → JSON + content type, named handler, unknown path → 404.
- Using `json.NewEncoder` is a bit **more** than the plan’s “hardcoded JSON string is fine” — that’s good, not wrong.
- Removed `handleOthers` / `/*` after learning that unmatched routes get a built-in 404 from the mux.
- Set headers before `WriteHeader`.
- Optional refinement (later): register as `GET /health` so only GET is allowed.

---

## Steps 10–15 — User model, POST /users echo + validation (2026-07-23)

**What I did:**
- Added `internal/user/model.go` with a `User` struct (`Name`, `Email`, `Age` + JSON tags).
- Added `internal/user/http.go` with an exported handler (`TestUser`) for `POST /users`.
- Wired it in `main` via `user.TestUser` (cross-package call).
- Decode request JSON into `User`, validate name/email non-empty and age >= 0, echo user JSON on success.
- Invalid method → 405; bad JSON / validation failures → 400 via `http.Error`.

**Why:** Shape the user domain and practice request/response JSON before touching a database.

**Commands / files:**
- `internal/user/model.go`
- `internal/user/http.go`
- `cmd/server/main.go` — `mux.HandleFunc("/users", user.TestUser)`
- Practice: `curl -X POST localhost:8080/users -H 'Content-Type: application/json' -d '{"name":"Ada","email":"a@b.com","age":30}'`

**Challenges / mistakes:**
- Learning how to call a function from another package (import path + exported name starting with a capital letter).
- Decoding JSON from the request body into a struct was new.
- Using `http.Error` for error responses was new.

**What I learned / review notes:**
- Steps 10–15 goal met for decode, echo, and validation.
- **Slightly less:** `User` has no `ID` field yet (plan included it for later DB use — add when we persist).
- **Slightly less:** no shared JSON write helper yet (Step 11); encoding is inline — fine for now, can extract later.
- **Fix soon:** success response uses `Content-Type: application-json` (typo) — should be `application/json`.
- **Naming:** `TestUser` works but something like `CreateUser` / `HandleCreate` will read clearer.
- **Difference:** plan said JSON error bodies; `http.Error` sends plain text. Acceptable for this phase; we can switch to JSON errors later.
- Method check with `http.Error` + 405 is a nice extra (or use `POST /users` on the mux later).

---

## Follow-up — JSON error helper (2026-07-23)

**What I did:** Added `writeJsonError` and used it for bad JSON / validation instead of `http.Error`, returning `{"error":"..."}` with `Content-Type: application/json`. Also fixed success `Content-Type` to `application/json`.

**Challenges / mistakes:** Learned that `http.Error` is plain text only.

**Review notes:**
- Method-not-allowed path still uses `http.Error` — switch that to `writeJsonError` too for consistency.
- **Bug:** `writeJsonError` never calls `w.WriteHeader(status)`, so the status stays **200** when the body is written. Add `w.WriteHeader(status)` after setting the header and before `Encode`.

---

## Steps 16–20 — Postgres via Docker + pgx pool (2026-07-23)

**What I did:**
- Added `docker-compose.yml` with Postgres 17 (`admin`/`admin`, DB `users`, port 5432, volume `./data`).
- Started with `docker compose up -d` (detached so the terminal is free).
- Added `.env` / `.env.example` with `DATABASE_URL`, and `.gitignore` for `.env`.
- Added dependencies (`godotenv`, `pgx/v5`).
- In `main`: load env, read `DATABASE_URL`, parse pool config, create pool, `Ping`, log success (fatal on failure).

**Why:** Need a real database before CRUD persistence.

**Commands / files:**
- `docker-compose.yml`, `.env`, `.env.example`, `.gitignore`
- `go get` for packages; `go.mod` / `go.sum`
- `cmd/server/main.go` — `godotenv.Load`, `os.Getenv`, `pgxpool.ParseConfig` + `NewWithConfig`, `Ping`
- `docker compose up -d` / `docker compose ps`

**Challenges / mistakes:**
- How to get a Go package (`go get`).
- How to load a `.env` file (`godotenv`) and read values (`os.Getenv`).
- How to build DB connection config and open a pool.
- How to `Ping` to verify connectivity.
- Chose `NewWithConfig` instead of `New` for more control (MaxConns, MinConns, MaxConnIdleTime).

**What I learned / review notes:**
- `pgxpool.New(ctx, url)` ≈ parse URL + `NewWithConfig` with defaults. `NewWithConfig` is correct when you want to tune the pool — not overkill.
- `-d` on `docker compose up` runs containers in the background.
- **Suggestion:** add `data/` to `.gitignore` so Postgres files under `./data` are not committed.
- Later we can relax “fatal if `.env` missing” so real env vars (CI/production) work without a file; fine for local learning now.

---

## Steps 21–28 — Schema + Create in DB (2026-07-23)

**What I did:**
- Added `migrations/001_create_users.sql` (`SERIAL` id, unique email, unique name, age).
- Added `Repo` with pool + `Create` (`INSERT ... RETURNING id`).
- Added `Handler` holding the pool (`NewHandler`), method `Create` that validates, calls repo, returns 201.
- Wired handler from `main` after creating the pool.
- For DB errors (including duplicate key), currently returning **400** for all — unsure how to detect unique violations.

**Why:** Persist users instead of only echoing JSON.

**Commands / files:**
- `migrations/001_create_users.sql`
- `internal/user/repository.go`, `internal/user/http.go`, `cmd/server/main.go`

**Challenges / mistakes:**
- How to pass the pool / DB into HTTP handlers and the repository (chose `Handler` struct with pool + `NewRepo` inside Create).
- How to detect unique-key violations — not solved yet; all DB errors → 400.

**What I learned / review notes:**
- **Passing the pool:** Putting `*pgxpool.Pool` on a `Handler` and constructing it in `main` is a valid approach. Cleaner next step: put `*Repo` on the Handler (or a Service) so you don’t call `NewRepo` on every request — same idea, one layer of wiring.
- **Bug:** after `writeJsonError` on Create failure, need `return` — otherwise the handler still writes **201** and an id.
- **Route:** interview API is `POST /users`; `/create` works for learning but should become `/users` (and drop or merge the old echo handler).
- **Response:** plan expected **201 + full user including id**; currently encoding only the id — add `ID` to `User` and return the struct.
- **Unique violation (Postgres `23505`):** with pgx, use `errors.As` into `*pgconn.PgError` and check `Code == "23505"` → **409 Conflict**; other DB errors → **500** (don’t leak raw DB strings to clients long-term).
- Schema: unique on `name` is stricter than required (unique `email` is enough); fine if intentional.

---

## Polish after Steps 21–28 (2026-07-23)

**What I did:**
- Handler now holds `*Repo`; `main` does `NewRepo(pool)` then `NewHandler(repo)`.
- Create registered on `/users`; echo handler removed.
- `Create` returns full `User` with `Id`; response 201 + JSON user.
- Detect unique violation via `pgconn.PgError` code `23505`; `return` after errors.

**Still open / small fixes:**
- Unique violation → **409** — done.
- Fixed client messages (no raw `err`) — done.
- Other DB errors still **400**; change that branch to **500** (`http.StatusInternalServerError`).
- Method-not-allowed still uses `http.Error` (plain text) — optional.

---

## Steps 29–36 — Service layer + GET/PUT/DELETE (2026-07-23)

**What I did:**
- Added `Service` between Handler and Repo; `main` wires `pool → repo → service → handler`.
- Moved Create validation into the service.
- Implemented GetById / Update / Delete in repo + service + handlers.
- Registered routes with methods: `POST /users`, `GET /users/{id}`, `PUT /users`, `DELETE /users/{id}`.
- Used `Exec` for Update/Delete (no row returned) instead of `QueryRow`.

**Challenges / mistakes:**
- Figuring out `Exec` vs `QueryRow`.
- When two APIs share a path shape, register with the HTTP method in the mux pattern (Go 1.22+).

**What I learned / review notes:**
- **Layering looks right:** Handler → Service → Repo.
- **Method in pattern:** correct. With `POST /users` etc., the extra `r.Method != ...` checks in handlers are redundant (optional cleanup).
- **Exec vs QueryRow:** right instinct — `QueryRow` when you `Scan` a row (`RETURNING` / `SELECT`); `Exec` when you don’t need returned columns. For Update/Delete, also check **rows affected** (or use `RETURNING` + `QueryRow`) so “id not found” isn’t treated as success.
- **Bug:** service replaces repo errors with `errors.New("Failed to...")`, so the original Postgres error (including `23505`) is lost — Create/Update **409 detection in the handler can never fire**. Prefer `return err` from the service for DB failures (or wrap with `%w`), and map sentinel/validation errors separately.
- **PUT path:** requirements say `PUT /users/{id}`; you used `PUT /users` with id in the body. Works, but prefer id from the path for consistency with GET/DELETE.
- **Not-found:** GetById maps every error to 404; Update/Delete don’t detect 0 rows. Introduce something like `ErrNotFound` from repo when `pgx.ErrNoRows` or `RowsAffected() == 0`, then handler returns 404.
- Validation errors from service currently look like generic create/update failures in the handler — later you can distinguish 400 vs 500 vs 409 cleanly.

---

## Phase F polish round 2 (2026-07-23)

**What I did:**
- Service now returns the real repo `err` (409 / `23505` can work again).
- `PUT /users/{id}` with id from `PathValue`.
- Update/Delete check `RowsAffected() == 0`.
- GetById maps `pgx.ErrNoRows` → 404.
- Dropped redundant method checks; routes include the method (including `GET /health`).

**Still small gaps:**
- Update/Delete “no rows” currently become generic handler failures (**400**). Prefer **404** — e.g. a shared `ErrNotFound` and `errors.Is` in the handler (same idea as `pgx.ErrNoRows` on Get).
- Unexpected DB errors on Create/Update/Get/Delete are still often **400**; prefer **500** for those, keep **400** for validation / bad id.

---

## Phase F polish round 3 (2026-07-23)

**What I did:** Added package-level `ErrNotFound`; Update/Delete return it when `RowsAffected() == 0`; handlers map it to **404**.

**Still open:** unexpected DB failures (and GetById non-`ErrNoRows` errors) still use **400** in several `else` branches — prefer **500**. Validation / bad JSON / bad id stay **400**.

---

## Steps 37–40 — Phase G timing middleware (2026-07-24)

**What I did:**
- `LogEvent` + buffered channel + `LogWorker` goroutine.
- `WaitGroup` with `Add` / `defer Done`.
- `middleware` wraps mux: timer → `next.ServeHTTP` → non-blocking send (`select`/`default`).
- `server.Handler = middleware(mux, ch)`.

**Challenges / mistakes:**
- Code after `ListenAndServe` never ran until moved earlier.
- Confused why middleware returns `http.Handler` and uses an inner `func(w, r)`.
- Confused about `select`/`default` (non-blocking send).
- Hit `cannot send to receive-only channel` — middleware needs `chan<-`, worker needs `<-chan`.

**Still missing for a clean Step 38 shutdown:** after `ListenAndServe` returns, call `close(ch)` then `wg.Wait()` (remove the useless `Sleep`). Note: Ctrl+C won’t reach that unless you add signal handling later — still worth having for a graceful stop path.

**Phase G status:** functionally complete for the interview requirement (goroutine + channel + WaitGroup + per-request timing). Shutdown close/Wait is the small leftover.

---

## Step 38 shutdown check (2026-07-24)

**What I did:** After `ListenAndServe`, `close(ch)` then `wg.Wait()`, still with a `time.Sleep` before close.

**Review:** Placement of `close` + `Wait` is **correct**. The `Sleep` is **unnecessary** — `Wait` already blocks until the worker finishes after close. Safe to delete the sleep.

---

## Steps 41 — Centralize API errors (2026-07-24)

**What I did:** Already had `writeJsonError` in `internal/user/http.go`; replaced remaining `http.Error` paths (invalid id) so all handler errors go through the same JSON `{"error":"..."}` helper.

**Step 41:** done.

**Step 42 next:** change unexpected-failure `else` branches from **400** to **500** (Create/Get/Update/Delete). Keep 400 for bad JSON / validation / bad id; 404 not-found; 409 unique.

---

## Step 42 — Consistent status codes (2026-07-24)

**What I did:** Unexpected failure branches now use `StatusInternalServerError` (500). Not-found stays 404; unique stays 409; bad JSON / invalid id stay 400.

**Note:** Service validation errors (empty name/email) still look like generic failures → 500 until you add distinct validation errors later. Acceptable for now.

**Step 43 next:** DB context timeouts (3–5s).

---

## Steps 41–43 — Phase H complete (2026-07-24)

**What I did:**
- All errors via `writeJsonError`.
- Status codes: 400 / 404 / 409 / 500 as appropriate.
- Each handler: `ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()` then pass `ctx` to the service.

**What I learned:** `WithTimeout` creates a child context with a deadline; `cancel` releases resources; DB calls that respect `ctx` stop when the deadline hits.

---

## Step 37 — Timing log worker + channel (2026-07-24)

**What I did:** Added `LogEvent` (method, path, duration), a buffered channel, and `LogWorker` that `range`s the channel and prints. Started the worker with `go LogWorker(ch)` and sent a sample event from `main`.

**Challenges / mistakes:**
- First put channel/worker/sample **after** `ListenAndServe` — that code never ran because `ListenAndServe` blocks. Moved setup **before** the server starts; then sample logs appeared.

**What I learned:** Anything after `ListenAndServe` only runs when the server stops. Start background workers before listening.

---

## Step 38 — WaitGroup + close (2026-07-24)

**What I did:** `wg.Add(1)`, pass `&wg` into `LogWorker`, `defer wg.Done()`, after sample send `close(ch)` then `wg.Wait()`.

**What I learned:** Closing the channel ends the worker’s `range` loop; `Wait` blocks until `Done` runs.

**Note for next steps:** Right now `close(ch)` happens *before* the HTTP server starts. That’s fine for the sample. Once middleware sends real request events, keep the channel open while the server runs, and move `close` + `Wait` to after `ListenAndServe` returns (shutdown).
