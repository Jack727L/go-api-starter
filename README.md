# Go Fiber backend template

A small, production-minded REST API template using [Fiber](https://gofiber.io/), PostgreSQL with [sqlc](https://sqlc.dev/), Redis-backed async jobs, [Sqitch](https://sqitch.org/) migrations, and [testcontainers](https://testcontainers.com/) for tests.

## Features

- **HTTP**: Fiber with CORS, request logging, and validation ([go-playground/validator](https://github.com/go-playground/validator)).
- **Auth**: Register, login, refresh, logout; JWT-style flows with typed handlers under `/auth`.
- **Users**: Authenticated `/users/me` CRUD-style profile routes.
- **Data**: `pgx` connection pool, SQL defined in `db/queries/`, generated Go in `db/sqlc/`.
- **Jobs**: Redis client and job executor (`core/jobs/`) for background-style work.
- **Ops**: `/healthz`, `/healthcheck`, `/readyz` (readiness includes DB ping).
- **Docs**: Swagger annotations; generated files live under `docs/`.

## Requirements

| Tool | Notes |
|------|--------|
| **Go** | 1.24+ (see `go.mod`). |
| **Docker** | For local Postgres/Redis (`docker compose`) and for tests (testcontainers). |
| **Sqitch** | Apply DB migrations from `db/sqitch` before first run. |
| **Optional** | `sqlc` (regenerate DB code), `swag` (regenerate OpenAPI), `reflex` (`make watch`). |

## Quick start

1. **Clone and configure environment**

   ```bash
   cp .env.example .env
   ```

   Edit `.env` if your ports or credentials differ from the defaults.

2. **Start Postgres and Redis**

   ```bash
   docker compose up -d
   ```

3. **Run migrations**

   From the repository root, deploy with a URI that matches `.env` (user, password, host, port, database):

   ```bash
   cd db/sqitch
   sqitch deploy "db:pg://appuser:apppassword@localhost:5432/appdb"
   cd ../..
   ```

4. **Run the API**

   ```bash
   go run .
   ```

   The server listens on **`:3000`**. Set `BACKEND_ENV` via `.env` (e.g. `local-dev`).

## Makefile

| Target | Purpose |
|--------|---------|
| `make build` | Build binary to `/tmp/go-fiber-template`. |
| `make run` | Build and run that binary. |
| `make watch` | Rebuild on `.go` changes (requires [`reflex`](https://github.com/cespare/reflex)). |
| `make test` | Run `./tools/runTests.sh` (Docker required). |
| `make test-verbose` | Same with `-v`. |
| `make docs` | Regenerate Swagger via `./tools/generateDocs.sh` (requires `swag`). |
| `make sqlc` | Regenerate `db/sqlc/` via `./tools/generateSQLC.sh` (requires `sqlc`). |

## API overview

| Area | Routes |
|------|--------|
| Auth | `POST /auth/register`, `POST /auth/login`, `POST /auth/refresh`, `POST /auth/logout` |
| Users | `GET/PUT/DELETE /users/me` (Bearer auth) |
| Health | `GET /healthz`, `GET /healthcheck`, `GET /readyz` |

OpenAPI/Swagger output is generated into `docs/`; regenerate after changing handler annotations.

## Testing

Tests assume a running Docker daemon. Do **not** run `go test` alone for integration-style suites; use the provided script so containers and env are set correctly:

```bash
./tools/runTests.sh              # default: isolated Postgres per run (testcontainers)
./tools/runTests.sh -v           # verbose
./tools/runTests.sh -t TestName # filter
./tools/runTests.sh --async      # uses Compose + Redis (see `tests/testing.docker-compose.yml`)
```

Async mode can run Sqitch against the test database; ensure `sqitch` is installed when using `--async`.

## Docker image

Build a minimal image (multi-stage, distroless runtime):

```bash
docker build -t go-fiber-template .
```

The compose file can be extended with an `app` service (commented example in `docker-compose.yml`).

## Project layout

```
├── core/           # Shared errors, job execution, Redis
├── db/
│   ├── queries/    # SQL for sqlc
│   ├── sqlc/       # Generated + connection helpers (generated — do not hand-edit)
│   └── sqitch/     # Migration plan and deploy scripts
├── docs/           # Swagger (generated)
├── env/            # Small env submodule for configuration helpers
├── middlewares/
├── routers/
├── services/
├── setup/          # Fiber setup, routes, health checks
├── tests/
├── tools/          # generateDocs, generateSQLC, runTests
├── main.go
├── Dockerfile
└── docker-compose.yml
```

## Customizing the module path

Replace the placeholder module `github.com/yourusername/go-fiber-template` in `go.mod`, imports, and Swagger host metadata with your own module path, then run `go mod tidy`.

## License

Add a `LICENSE` file for your project if you distribute this template.
