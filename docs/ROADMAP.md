# Roadmap

This is a demo that grew a few real muscles. These are the obvious next steps,
in rough priority order — the first three are the ones that would change the
"portfolio piece" into a "production-shaped" thing.

## 1. Repository interfaces + real controller tests

The biggest gap. Today `database.OpenCollection` hands controllers a live Mongo
collection, so the handlers aren't unit-testable. Introduce a `MovieRepository`
/ `UserRepository` interface, implement it over Mongo, and inject it into the
handlers. Then the controller tests stop being "pure helpers only".

## 2. Redis-backed rate limiting

Swap the in-memory token bucket for a shared store so limits survive multiple
instances. Bonus: a per-user limit in addition to per-IP.

## 3. Token revocation + refresh rotation with reuse detection

Short-lived access tokens, refresh tokens that rotate on every use, and a
denylist for logged-out sessions. This is the difference between "clears the
cookie" and "actually revokes the session".

## 4. Watch history / continue watching

The natural next user feature: record play events from the stream endpoint and
surface a "Continue watching" row. Was intentionally cut from the first pass to
keep scope tight.

## 5. Smarter recommendations

The current recommender is "same genres, best ranked". Options, roughly in
ascending effort: exclude movies already watched; score by genre overlap
weighted by user's pick order; or an embedding-based similarity search (Mongo
Atlas Vector Search).

## 6. Search polish

Move from regex on title to a proper text index (or Atlas Search) for typo
tolerance and relevance ranking, plus multi-field search (title + genre +
curator review).

## 7. Media niceties

HLS packaging for adaptive streaming, per-file metadata in Mongo, and an upload
endpoint so admins can add clips without touching the filesystem.

## 8. Observability

Prometheus metrics (`/metrics`) and OpenTelemetry tracing, with the request-ID
middleware feeding trace IDs.

## 9. CI config

A pipeline (GitHub Actions or otherwise) running `vet`, `test`, `gosec`,
`govulncheck`, and the frontend `build` + `lint` on every PR.
