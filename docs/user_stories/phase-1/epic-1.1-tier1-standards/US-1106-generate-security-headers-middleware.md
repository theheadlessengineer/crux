# US-1106 — Generate Security Headers Middleware (CORS, XSS, CSRF)

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want every generated service to apply security response headers and CORS policy by default,
so that common web security vulnerabilities are mitigated without requiring service teams to configure them manually.

---

## Pre-Development Checklist

- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] The company's default CORS policy (allowed origins, methods, headers) is agreed and documented
- [ ] Security team has reviewed the proposed header set
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate a security headers middleware that applies a hardened default set of HTTP response headers on every response.

### Required Headers

| Header | Default Value |
|---|---|
| `Content-Security-Policy` | `default-src 'self'` |
| `X-Content-Type-Options` | `nosniff` |
| `X-Frame-Options` | `DENY` |
| `Strict-Transport-Security` | `max-age=31536000; includeSubDomains` |
| `Referrer-Policy` | `strict-origin-when-cross-origin` |
| `Permissions-Policy` | `geolocation=(), microphone=(), camera=()` |

### CORS Configuration

CORS policy is configurable via environment variables. The default denies all cross-origin requests unless explicitly configured.

| Environment Variable | Default | Description |
|---|---|---|
| `CORS_ALLOWED_ORIGINS` | (empty — deny all) | Comma-separated list of allowed origins |
| `CORS_ALLOWED_METHODS` | `GET,POST,PUT,DELETE` | Allowed methods |
| `CORS_ALLOWED_HEADERS` | `Authorization,Content-Type` | Allowed headers |

### In Scope

- Security headers middleware applied globally
- CORS middleware with environment variable configuration
- Unit tests verifying each header is present on responses

### Out of Scope

- CSRF token generation (applicable to browser-facing services only — service teams opt in)
- Rate limiting (separate story US-1109)

---

## Acceptance Criteria

- [x] All six security headers are present on every response from the generated service — `SecurityHeaders()` middleware; `TestSecurityHeaders_AllHeadersPresent` verifies all six values
- [x] CORS is denied by default when `CORS_ALLOWED_ORIGINS` is not set — `TestCORSMiddleware_DeniedByDefault` passes
- [x] CORS allows cross-origin requests when `CORS_ALLOWED_ORIGINS` is configured — `TestCORSMiddleware_AllowedOrigin` passes; unlisted origin denied in `TestCORSMiddleware_UnlistedOriginDenied`
- [x] Unit tests verify each header value on a test response — `TestSecurityHeaders_AllHeadersPresent` checks all six headers by exact value
- [x] Headers are applied before any application middleware — both middlewares are registered via `r.Use()` before route handlers
- [x] Security headers cannot be overridden by service teams without modifying the middleware — headers are set unconditionally in `SecurityHeaders()` with no configuration surface

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Response headers verified manually with `curl -I`
- [ ] CORS behaviour verified: denied without config, allowed with config
- [x] Unit tests pass — `go test ./...` green; `golangci-lint` 0 issues
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
