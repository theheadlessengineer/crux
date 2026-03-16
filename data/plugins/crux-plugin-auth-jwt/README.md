# crux-plugin-auth-jwt

**Tier:** 1 (Official) | **Version:** 1.0.0 | **Phase:** Pilot

JWT authentication middleware and RBAC/ABAC authorization stubs for crux-generated Go services.

## Questions

| ID | Type | Prompt | Default |
|---|---|---|---|
| `auth_authz_model` | select | Authorization model | `rbac` |
| `auth_jwks_url` | input | JWKS endpoint URL | `` |

## Generated Files

| File | Description |
|---|---|
| `internal/infrastructure/auth/jwt.go` | JWT middleware + RBAC stub |

## Usage in Generated Service

```go
// In router setup:
protected := router.Group("/v1")
protected.Use(auth.Middleware(auth.Config{
    JWKSEndpoint: os.Getenv("JWKS_ENDPOINT"),
    Issuer:       os.Getenv("JWT_ISSUER"),
}, log))
protected.GET("/resource", auth.RequireRole("reader"), handler)
```

## Environment Variables

| Variable | Description |
|---|---|
| `JWKS_ENDPOINT` | JWKS URL for public key retrieval |
| `JWT_ISSUER` | Expected token issuer |
| `JWT_AUDIENCE` | Expected token audience |

## TODO

Replace the JWKS key resolver stub in `jwt.go` with your identity provider's JWKS endpoint before deploying to production.
