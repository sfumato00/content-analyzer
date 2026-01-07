# Repository Guidelines

## Project Structure & Module Organization
- `backend/cmd/api/main.go` is the API entry point.
- `backend/internal/` contains core packages (auth, handlers, server, config, database, cache, middleware).
- `backend/migrations/` holds SQL migrations (`000001_name.up.sql` / `.down.sql`).
- `scripts/` and `backend/scripts/` hold helper scripts (setup verification, migrations).
- Root files like `docker-compose.yml`, `Makefile`, `README.md`, and `SETUP.md` document and run the stack.

## Build, Test, and Development Commands
- `make dev-setup`: install Go deps and start postgres/redis.
- `make run`: run the API locally.
- `make test`: run all Go tests.
- `make test-coverage`: generate `backend/coverage.html`.
- `make fmt` / `make lint`: format Go code and run `go vet`.
- `make docker-up` / `make docker-down`: manage all services.
- `make migrate-up` / `make migrate-down`: run or rollback migrations.

## Coding Style & Naming Conventions
- Use Go formatting: run `gofmt` via `make fmt` (tabs are expected).
- Package names are lowercase and short; keep exported identifiers in Go PascalCase.
- Tests live alongside code as `*_test.go` (example: `backend/internal/auth/jwt_test.go`).
- Migration names should be descriptive (example: `add_users_table`).

## Testing Guidelines
- Primary runner: `go test ./...` from `backend/` or `make test`.
- Keep tests close to the unit under test and use clear test case names.
- Coverage is optional but encouraged; use `make test-coverage` when changing core logic.

## Commit & Pull Request Guidelines
- Recent commits use short, imperative messages; a conventional format like `feat(auth): add login` is preferred.
- PRs should include a brief summary, linked issue (if any), and the test command run.
- If you change APIs or migrations, note endpoints or migration names in the PR body.

## Security & Configuration
- Never commit `.env`; use `.env.example` as the template.
- Required config: `GEMINI_API_KEY`, `DATABASE_URL`, `REDIS_URL`, `JWT_SECRET`.
