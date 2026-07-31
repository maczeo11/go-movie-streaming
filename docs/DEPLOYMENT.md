# Deployment

## Local (no Docker)

Requires Go 1.26+, Node 20+, and MongoDB reachable at `MONGODB_URI`.

```bash
cp Server/MagicMovieStream/.env.example Server/MagicMovieStream/.env
# fill in MONGODB_URI and SECRET_KEY/SECRET_REFRESH_KEY

cd Server/MagicMovieStream
go run ./cmd/seed        # load demo data
go run .                 # API on :8080

cd ../../Client
npm install
npm run dev              # UI on :5173, proxies /api to :8080
```

## Docker Compose (recommended for a quick spin-up)

```bash
docker compose up --build
```

This starts MongoDB (`mongo:7`, with a healthcheck) and the API. The compose
file wires `MONGODB_URI: mongodb://mongo:27017/` and mounts
`Server/MagicMovieStream/media` into the container so you can drop clips in
without rebuilding.

Seed inside the container:

```bash
docker compose exec api sh -c 'cp .env.example .env 2>/dev/null; ./magicstream' 
```

…which isn't actually how seeding works — the correct way:

```bash
# seed while the API container is up, using its Mongo connection
docker compose run --rm api sh -c "DATABASE_NAME=magicstream MONGODB_URI=mongodb://mongo:27017/ go run ./cmd/seed"
```

> The multi-stage Dockerfile does a `CGO_ENABLED=0` build, so the image is a
> static binary on `alpine` — no libc surprises at runtime.

## Environment variables

| Variable                 | Default                    | Purpose                                     |
| ------------------------ | -------------------------- | ------------------------------------------- |
| `PORT`                   | `8080`                     | API listen port                             |
| `DATABASE_NAME`          | —                          | Mongo database name (required)              |
| `MONGODB_URI`            | —                          | Mongo connection string (required)          |
| `SECRET_KEY`             | dev fallback               | access-token signing secret                 |
| `SECRET_REFRESH_KEY`     | dev fallback               | refresh-token signing secret                |
| `ALLOWED_ORIGINS`        | `http://localhost:5173`    | comma-separated CORS origins                |
| `MEDIA_DIR`              | `./media`                  | where `/media` streams from                 |
| `OPENAI_API_KEY`         | —                          | required for AI review ranking              |
| `BASE_PROMPT_TEMPLATE`   | —                          | review-classification prompt                |
| `RECOMMENDED_MOVIE_LIMIT`| `5`                        | cap on recommendation count                 |
| `API_RATE` / `API_BURST` | `2` / `120`                | public-read rate limit (req/s, burst)       |
| `AUTH_RATE` / `AUTH_BURST`| `0.1` / `10`              | auth-route rate limit                       |
| `LOG_LEVEL`              | `info`                     | reserved for debug output                   |

## Production notes

- Serve the built SPA (`Client/dist`) from a static host, with `/api/*` routed
  to this API. Or serve it from the same host as the API behind a reverse
  proxy — the client just needs a `VITE_API_BASE` that points at the API.
- Put TLS in front. The auth cookies are `Secure`; over plain HTTP outside
  localhost browsers will drop them and login will mysteriously "not stick".
- Secrets: real values in a secret manager, not `.env`.
- See [SECURITY.md](SECURITY.md) for the deployment-specific hardening list
  (rate limiter, token denylist, indexes).

## CI

The `Makefile` exposes `test`, `vet`, and `fmt`; a CI pipeline would run
`go vet ./...`, `go test ./...`, `gosec ./...` and `govulncheck ./...` on every
PR. No CI config ships in this repo to keep it tool-agnostic.
