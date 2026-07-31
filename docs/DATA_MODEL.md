# Data Model

Single MongoDB database (configurable via `DATABASE_NAME`). Four collections:
`users`, `movies`, `genres`, `rankings`.

## users

| Field              | Type                 | Notes                                    |
| ------------------ | -------------------- | ---------------------------------------- |
| `_id`              | ObjectID             |                                          |
| `user_id`          | string               | hex ObjectID string, the auth identifier |
| `first_name`       | string               |                                          |
| `last_name`        | string               |                                          |
| `email`            | string               | unique per application logic             |
| `password`         | string               | bcrypt hash — never the raw password     |
| `role`             | `"USER"` / `"ADMIN"` | always `"USER"` on registration          |
| `created_at`       | datetime             |                                          |
| `updated_at`       | datetime             |                                          |
| `token`            | string               | latest access token (for revocation)     |
| `refresh_token`    | string               | latest refresh token (for revocation)    |
| `favourite_genres` | `Genre[]`            | embedded; drives recommendations         |

A `Genre` is `{ genre_id: int, genre_name: string }`, embedded as a
subdocument.

## movies

| Field          | Type      | Notes                                 |
| -------------- | --------- | ------------------------------------- |
| `_id`          | ObjectID  |                                       |
| `imdb_id`      | string    | unique lookup key, e.g. `tt0133093`   |
| `title`        | string    |                                       |
| `poster_path`  | string    | validated as a URL                    |
| `youtube_id`   | string    | trailer id, embedded in the UI        |
| `genre`        | `Genre[]` | multi-genre                           |
| `admin_review` | string    | free-text curator review              |
| `ranking`      | `Ranking` | `{ ranking_value: int, ranking_name: string }` |

Search uses a case-insensitive regex on `title`; the recommendation query
matches `genre.genre_name` and sorts by `ranking.ranking_value` ascending.

## genres

The canonical genre list. Seeded once, referenced by both users and movies by
`genre_id` / `genre_name`.

## rankings

The ranking scale the AI review classifier is allowed to answer with.

| Value | Name                  |
| ----- | --------------------- |
| 1     | Must Watch            |
| 2     | Highly Recommended    |
| 3     | Recommended           |
| 4     | Average               |
| 5     | Skip It               |
| 999   | Not Ranked (sentinel) |

`999` is excluded from the prompt sent to OpenAI, so the model can never "cheat"
by answering with the sentinel.

## Seeding

`go run ./cmd/seed` wipes and re-inserts `genres`, `rankings` and `movies` from
the JSON files in `seed/`. It's idempotent, which makes dev environments
predictable: `make seed` before you start and you know exactly what you're
looking at.

## Indexing notes

The catalog is tiny (15 seeded movies), so no custom indexes are defined in
code. For a production-sized catalog you'd add at least:
- `users.email` unique
- `movies.imdb_id` unique
- `movies.title` text index (or lowercase collation for search)
- `movies.genre.genre_name` (for the recommendation query)
