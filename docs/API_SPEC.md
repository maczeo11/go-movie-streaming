# API Spec

Base URL (local): `http://localhost:8080`

The client normally reaches the API through the Vite proxy (`/api` prefix,
stripped before forwarding), so in the browser the calls are same-origin and
cookies "just work". Direct API calls work too, using the `Set-Cookie` headers
or the `Authorization` header interchangeably.

## Authentication

Sessions use two JWT cookies set by `/login` and rotated by `/refreshtoken`:

| Cookie          | TTL      | Purpose                          |
| --------------- | -------- | -------------------------------- |
| `access_token`  | 24 hours | short-lived bearer for requests  |
| `refresh_token` | 7 days   | long-lived, swapped on refresh   |

Both are `HttpOnly`, `Secure` (browsers still allow this on `http://localhost`)
and `SameSite=None`. The access token can also be sent as a
`Authorization: Bearer <token>` header — the middleware accepts either.

Protected endpoints return `401` when the token is missing, invalid or expired.
Admin endpoints return `403` when the token is valid but the role isn't `ADMIN`.

## Errors

Errors follow one shape:

```json
{ "error": "Validation failed", "details": "..." }
```

`details` is only present when there's something useful to say (usually
validator output).

---

## Public

### `GET /health`

Uptime and database connectivity. Returns `200` when Mongo answers a ping,
`503` otherwise.

```json
{
  "status": "ok",
  "uptime_seconds": 42,
  "database": { "connected": true },
  "go_version": "go1.26.3",
  "env": "debug"
}
```

### `GET /movies`

List the catalog with optional search and pagination.

Query params:

| Param  | Meaning                              | Example            |
| ------ | ------------------------------------ | ------------------ |
| `q`    | case-insensitive title search        | `q=matrix`         |
| `genre`| exact match on `genre.genre_name`    | `genre=Sci-Fi`     |
| `page` | 1-based page number                  | `page=2`           |
| `limit`| items per page, 1–50                 | `limit=24`         |
| `sort` | `title`, `imdb_id`, `rating`         | `sort=rating`      |

Response:

```json
{
  "data": [
    {
      "imdb_id": "tt0133093",
      "title": "The Matrix",
      "poster_path": "https://image.tmdb.org/t/p/w500/...",
      "youtube_id": "vKQi3bBA1y8",
      "genre": [{ "genre_id": 1, "genre_name": "Sci-Fi" }],
      "admin_review": "Bullet time still holds up...",
      "ranking": { "ranking_value": 1, "ranking_name": "Must Watch" }
    }
  ],
  "meta": { "page": 1, "limit": 12, "total": 15, "total_pages": 2 }
}
```

### `GET /movie/:imdb_id`

Single movie. `404` when the id isn't in the catalog.

### `GET /genres`

The genre list used to seed pickers and filters.

```json
[{ "genre_id": 1, "genre_name": "Sci-Fi" }]
```

### `GET /media`

Names of the video files currently available to stream.

```json
{ "data": ["sample.mp4"], "count": 1 }
```

### `GET /media/:file`

Streams a file from `MEDIA_DIR` with HTTP `Range` support. Returns `206 Partial
Content` on ranged requests, `404` for unknown files, `403` for traversal
attempts.

---

## Auth routes (rate-limited, 10/min default)

### `POST /user` — register

```json
{
  "first_name": "Jane",
  "last_name": "Doe",
  "email": "jane@example.com",
  "password": "hunter2secure",
  "favourite_genres": [{ "genre_id": 1, "genre_name": "Sci-Fi" }]
}
```

- `201` on success (returns the Mongo insert result).
- `409` if the email is already registered.
- `400` with `details` on validation failure.
- Any `role` field in the body is **ignored**; everyone signs up as `USER`.

### `POST /login` — sign in

```json
{ "email": "jane@example.com", "password": "hunter2secure" }
```

Sets both auth cookies and returns the profile:

```json
{
  "user_id": "…",
  "first_name": "Jane",
  "last_name": "Doe",
  "email": "jane@example.com",
  "role": "USER",
  "favourite_genres": [{ "genre_id": 1, "genre_name": "Sci-Fi" }]
}
```

`401` for a wrong email or password (same message for both, no user-enumeration
tell).

### `POST /refreshtoken`

Reads the `refresh_token` cookie, validates it, and issues a fresh pair (new
cookies + updated user document). `401` when the refresh token is missing,
invalid or expired.

---

## Protected (valid access token required)

### `POST /logout`

```json
{ "user_id": "…" }
```

Clears the stored tokens in the user document and expires both cookies.

### `GET /profile`

The current user's profile (same shape as the login response).

### `PATCH /profile/genres`

```json
{ "favourite_genres": [{ "genre_id": 2, "genre_name": "Action" }] }
```

Replaces the user's favourite genres. `400` when validation fails.

### `GET /recommendedmovies`

Movies whose genres intersect the user's favourites, best-ranked first, capped
by `RECOMMENDED_MOVIE_LIMIT` (default 5). An empty array when the user has no
favourite genres yet.

```json
[
  { "imdb_id": "tt1375666", "title": "Inception", "ranking": { "ranking_value": 1, "ranking_name": "Must Watch" } }
]
```

---

## Admin (valid access token + `ADMIN` role)

### `POST /addmovie`

Body is a full movie object — see the `data` shape under `GET /movies`. `201`
on success, `400` when validation fails (every field is required except
`admin_review`).

### `PATCH /updatereview/:imdb_id`

```json
{ "admin_review": "Masterpiece of cinema!" }
```

Runs the review through OpenAI to classify it into a ranking, then stores both.
Response carries the classification back:

```json
{ "ranking_name": "Must Watch", "admin_review": "Masterpiece of cinema!" }
```

Requires `OPENAI_API_KEY` and `BASE_PROMPT_TEMPLATE` on the server, otherwise
`500`.

---

## Rate limiting

Token-bucket, per client IP, in memory:

| Route group        | Default rate | Default burst |
| ------------------ | ------------ | ------------- |
| Public reads       | 2 req/s      | 120           |
| Auth               | 0.1 req/s    | 10            |

Exceeding the limit returns `429` with a `Retry-After: 60` header. Tunable via
`API_RATE`, `API_BURST`, `AUTH_RATE`, `AUTH_BURST`. See
[docs/SECURITY.md](SECURITY.md) for the single-instance caveat.
