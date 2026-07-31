# Testing

## Running the suite

```bash
cd Server/MagicMovieStream
go test ./...
go test -race ./...    # needs a working C toolchain for the race detector
```

The frontend has a type-check (`tsc -b`) and `oxlint` wired into `npm run
build` / `npm run lint`.

## What's covered

Unit tests target the pure logic — the pieces that don't need a database or a
running server — so they're fast and deterministic.

| Package       | What's tested                                                                 |
| ------------- | ----------------------------------------------------------------------------- |
| `utils`       | JWT round-trip (generate → validate), role/email/user-id claims, wrong-key rejection, garbage input, expired tokens, refresh-token validation |
| `middleware`  | Token bucket: burst passes then blocks, refill over time, per-client isolation; `AdminOnly` allows admins and 403s users/missing role |
| `controllers` | `parsePagination` defaults + query parsing + invalid input rejection; `buildMovieFilter` for q/genre (regex + case-insensitive); `sortField` whitelist; `isVideoExt` |

## What's deliberately not covered (and why)

- **Handlers that touch Mongo** (`GetMovies`, `LoginUser`, …) aren't unit-tested
  because the current `database.OpenCollection` returns a real collection. A
  proper test would inject a fake repository behind an interface — listed as the
  top priority in the roadmap, since it's the difference between "tests exist"
  and "tests guard the business logic".
- **The OpenAI ranking call** — it needs a live key; it's exercised manually and
  its error path is asserted in the handler's response handling.
- **End-to-end flows** — registration → login → recommendations is covered by
  manual testing with the seed data and the Bruno collection.

## Coverage target

The current pure-logic packages sit comfortably above 80% line coverage. As
repository interfaces land, the goal is to bring the controller layer to the
same bar.

## How to add a test

Tests live next to the code they cover (`middleware/rate_limit_test.go`,
`controllers/movie_filter_test.go`, …). The rate-limiter and filter helpers were
written to be testable: no global state, no I/O, just inputs → outputs.
