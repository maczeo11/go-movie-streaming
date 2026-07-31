# Security

This is a portfolio demo, not a bank. Here's what's handled properly, and what
is deliberately kept simple — with the honest gap between them called out.

## In place

- **Passwords** hashed with bcrypt (default cost) before storage. The raw
  password is never logged, returned, or stored.
- **No role escalation** — registration ignores any client-supplied `role` and
  always stores `USER`. Admin is created by editing the DB directly.
- **httpOnly cookies** for tokens — JavaScript never sees them, which kills
  XSS token theft as a vector. `Secure` + `SameSite=None` (browsers permit this
  on localhost).
- **JWT secrets read lazily** from the environment, with a dev-only fallback so
  the server boots without a `.env`. The fallback is clearly marked and would
  be removed/rotated before any real deployment.
- **Role checks at two layers** — the `AdminOnly` middleware gates the admin
  route group, and `AdminReviewUpdate` still re-checks the role from the token
  claims (defense in depth).
- **Rate limiting** — token bucket per client IP. Auth endpoints get a tight
  limit to slow brute-force password attempts; public reads get a generous one.
  Over-limit requests get `429` + `Retry-After`.
- **Path traversal guard** on `/media` — the requested path is resolved and
  verified to stay inside `MEDIA_DIR` before the file is served.
- **No user enumeration** — login returns the same `401` for unknown email and
  wrong password. (Register necessarily reveals whether an email exists; a
  production app would add a verification step.)
- **Boundary validation** — all bodies go through `go-playground/validator`
  before touching the DB, with structured `details` on failure.
- **CORS is explicit** — only configured origins are allowed, credentials
  included, rather than `*`.

## Deliberately simple (read the caveats)

- **In-memory rate limiter** — per-IP buckets live in a map in one process.
  Great for a demo; useless across multiple instances behind a load balancer,
  and shared NATs will share a bucket. Swap the `RateLimiter` for Redis (or an
  edge limiter) when this stops being a demo.
- **No token revocation list** — logout clears the stored token in the user
  document, but a stolen token keeps working until it expires. Real fix:
  maintain a denylist or move to short-lived access + refresh rotation with
  reuse detection.
- **Refresh token isn't rotated on a schedule** — `/refreshtoken` issues a new
  pair each time it's called, but there's no sliding-window reuse detection.
- **SameSite=None without HTTPS** — cookies work in browsers on localhost;
  anything beyond localhost should sit behind TLS (the `Secure` flag then earns
  its keep).
- **No CSRF token** — the API relies on cookie auth; modern browsers only send
  those cookies same-site for cross-site requests thanks to `SameSite=None` +
  `Secure`. A same-site subdomain attacker would be out of scope for a demo but
  worth a CSRF token in a real product.
- **bcrypt cost is the default** (10). Fine for this app; tune upward if you
  expect to store thousands of accounts.

## Production checklist

Before shipping this anywhere real:

1. Set strong, unique `SECRET_KEY` / `SECRET_REFRESH_KEY` (generate with
   `openssl rand -hex 32`), never the dev fallback.
2. Put TLS in front (reverse proxy or managed load balancer).
3. Replace the in-memory rate limiter and add a token denylist.
4. Add proper index creation + migrations for schema changes.
5. Run `govulncheck ./...` and `gosec ./...` in CI.
6. Store `OPENAI_API_KEY` in a secret manager, not `.env`.

## Secrets in the repo

`.env` is git-ignored at the repo root; only `.env.example` (with empty secret
values) is committed. The seeded credentials are demo-only and never valid in
production.
