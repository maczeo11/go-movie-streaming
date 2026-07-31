# ADR-002: In-memory token bucket for rate limiting

**Status:** accepted (with known scaling limit)

## Context

The API needed basic abuse protection: tighter limits on auth endpoints to slow
brute-force login attempts, generous limits elsewhere. The options were an
in-memory limiter, a Redis-backed limiter, or an edge/API-gateway limiter.

## Decision

A token-bucket `RateLimiter` in `middleware/`, keyed by client IP, stored in a
mutex-guarded map with periodic pruning of idle buckets. Configured per route
group via `newRateLimiter` in `server.go`, with env-tunable rates.

## Alternatives considered

- **Redis**: correct for multi-instance, but adds an operational dependency the
  demo doesn't need yet.
- **API gateway / nginx `limit_req`**: simplest operationally, but it moves the
  policy out of the codebase and hides it from anyone reading the Go.

## Consequences

- Zero new dependencies, works out of the box.
- Single-process only: behind a load balancer the limits are per-instance, not
  global. Documented in SECURITY.md; the roadmap schedules the Redis swap when
  this stops being a demo.
- Shared NATs / proxy IPs share a bucket — acceptable here, worth a
  `X-Forwarded-For` keying strategy in production.
