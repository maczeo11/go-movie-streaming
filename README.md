# MagicStream 🍿

A small movie streaming server written in Go with a React frontend. Browse a
curated catalog, watch trailers (and a demo clip streamed straight from the Go
backend), sign in to get genre-based recommendations, and — if you're an admin —
add movies and have the server AI-rank your curator reviews.

Built to show off a realistic Go API: JWT auth with httpOnly cookies, a Mongo
data layer, rate limiting, structured logging, and graceful shutdown. Not a toy
CRUD scaffold — but also not pretending to be Netflix.

## Features

- **Catalog browsing** — search by title, filter by genre, paginate and sort
- **Media streaming** — `/media` serves video files with HTTP Range support, so
  the `<video>` tag can seek without downloading the whole file
- **Auth** — register / login / refresh via JWT access + refresh tokens in
  httpOnly cookies; passwords hashed with bcrypt
- **Recommendations** — a logged-in user's favourite genres drive a ranked
  "picked for you" row
- **AI curator reviews** — admins write a review and the server calls OpenAI to
  map it onto the ranking scale (Must Watch … Skip It)
- **Profile** — view and edit favourite genres, which feeds the recommender
- **Admin area** — add movies to the catalog, re-rank reviews
- **Production habits** — token-bucket rate limiting (tighter on auth routes),
  request-ID tracing, JSON structured logs, health endpoint, graceful shutdown,
  Docker + Makefile, unit tests

## Tech stack

| Layer    | Choice                                                            |
| -------- | ----------------------------------------------------------------- |
| API      | [Gin](https://github.com/gin-gonic/gin) (Go 1.26)                 |
| Database | MongoDB via `mongo-driver/v2`                                     |
| Auth     | `golang-jwt/jwt/v5`, httpOnly cookies, `bcrypt`                   |
| AI       | `langchaingo` (OpenAI) for review ranking                         |
| Logging  | `log/slog` (stdlib)                                               |
| Frontend | React 19 + Vite + TypeScript, plain hand-written CSS              |

## Repo layout

```
.
├── Server/MagicMovieStream      # the Go API
│   ├── cmd/seed                 # idempotent seed command
│   ├── controllers/             # HTTP handlers
│   ├── database/                # Mongo connection + collection access
│   ├── middleware/              # auth, admin, rate limit, request-id, logger
│   ├── model/                   # User / Movie / Genre / Ranking
│   ├── routes/                  # public / auth / protected / admin groups
│   ├── seed/                    # JSON seed data (genres, rankings, movies)
│   ├── utils/                   # JWT helpers
│   ├── media/                   # sample clip streamed by /media
│   ├── bruno/                   # ready-made API collection (Bruno)
│   └── tests next to the code
├── Client                       # React app (Vite)
└── docs/                        # architecture, API, data model, etc.
```

## Quick start

Prereqs: Go 1.26+, Node 20+, and a running MongoDB (or Docker).

```bash
# 1. copy the env template and fill in MONGODB_URI / SECRET_KEY
cp Server/MagicMovieStream/.env.example Server/MagicMovieStream/.env

# 2. seed the database (creates genres, rankings and 15 movies)
cd Server/MagicMovieStream && go run ./cmd/seed

# 3. start the API (http://localhost:8080)
go run .

# 4. in another terminal, start the client (http://localhost:5173)
cd Client
npm install
npm run dev
```

Everything is wired through a Vite dev proxy, so the browser talks to one
origin (`/api` → `:8080`). Open http://localhost:5173, register an account, and
the home page shows the seeded catalog with a "picked for you" row once you've
saved a couple of favourite genres.

> The demo clip (`media/sample.mp4`) is a few seconds of generated test video so
> the streaming endpoint is usable out of the box. Drop your own files into
> `Server/MagicMovieStream/media/` to stream real content — they're git-ignored.

### Docker

```bash
docker compose up --build   # mongo + api
make seed                   # or seed inside the container
```

### Without Docker

You'll need MongoDB locally, then the steps above. `Makefile` has shortcuts:
`make run`, `make test`, `make vet`, `make seed`.

## API at a glance

| Method | Path                | Auth    | Purpose                                  |
| ------ | ------------------- | ------- | ---------------------------------------- |
| GET    | `/health`           | public  | uptime + DB ping                         |
| GET    | `/movies`           | public  | list, `?q&genre&page&limit&sort`         |
| GET    | `/movie/:imdb_id`   | public  | movie detail                             |
| GET    | `/genres`           | public  | genre list                               |
| GET    | `/media`            | public  | available clips                          |
| GET    | `/media/:file`      | public  | stream (Range requests supported)        |
| POST   | `/user`             | limited | register                                 |
| POST   | `/login`            | limited | sign in, sets cookies                    |
| POST   | `/refreshtoken`     | limited | rotate tokens                            |
| POST   | `/logout`           | user    | sign out, clears cookies                 |
| GET    | `/profile`          | user    | current profile                          |
| PATCH  | `/profile/genres`   | user    | update favourite genres                  |
| GET    | `/recommendedmovies`| user    | recommendations                          |
| POST   | `/addmovie`         | admin   | add a movie                              |
| PATCH  | `/updatereview/:id` | admin   | write + AI-rank a review                 |

Full details, request/response examples, and the OpenAI prompt contract live in
[docs/API_SPEC.md](docs/API_SPEC.md).

## Testing

```bash
make test          # or: cd Server/MagicMovieStream && go test ./...
```

Unit tests cover the token lifecycle (generate, validate, expiry, wrong key),
the rate limiter (burst, refill, client isolation), pagination parsing, the
title/genre filter builder, the sort-field whitelist, and the admin guard. See
[docs/TESTING.md](docs/TESTING.md) for the full story.

## Documentation

- [Architecture](docs/ARCHITECTURE.md) — how the pieces fit, request lifecycle
- [API spec](docs/API_SPEC.md) — every endpoint with examples
- [Data model](docs/DATA_MODEL.md) — collections and fields
- [Security](docs/SECURITY.md) — what's handled and what's deliberately simple
- [Deployment](docs/DEPLOYMENT.md) — Docker, env vars, production notes
- [Roadmap](docs/ROADMAP.md) — where it could go next
- [Decision records](docs/DECISIONS/) — why things are the way they are

## Acknowledgements

Poster images come from TMDB's image CDN and trailers are embedded from YouTube;
the API keys to fetch them aren't needed because the seed data ships with the
paths. The curator rankings in the seed are jokes, not opinions.

## License

MIT — see [LICENSE](LICENSE).
