# ADR-001: Lazy JWT secret loading

**Status:** accepted

## Context

The original code read `SECRET_KEY` / `SECRET_REFRESH_KEY` at package-init time:

```go
var SECRET_KEY string = os.Getenv("SECRET_KEY")
```

Package-level initializers run when the package is imported — which happens
*before* `main()` calls `godotenv.Load(".env")`. When running purely from a
`.env` file the keys were empty strings, so every token was signed with an empty
key. Tokens validated fine against each other, which masked the bug until you
restarted the server (invalidating every session) or deployed to a host that
didn't have the env vars exported.

## Decision

Replace the init-time reads with lazy getters:

```go
func getSecretKey() string {
    if secret := os.Getenv("SECRET_KEY"); secret != "" {
        return secret
    }
    return "local_dev_secret_key_do_not_use_in_prod"
}
```

called at signing/validation time. A clearly-labeled dev fallback means the
server still boots without any config — convenient for a demo, and obviously
wrong for production.

## Consequences

- Tokens are signed with the actual configured secret.
- Signing and validation can never disagree on the key within a process.
- The dev fallback could lull someone into deploying without secrets; the name
  makes the intent unmistakable and SECURITY.md calls it out.
