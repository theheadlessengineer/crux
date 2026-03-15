# US-1107 — Generate Input Sanitization Middleware

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want every generated service to sanitize inbound request data by default,
so that common injection vectors are mitigated at the framework layer without requiring service teams to implement sanitization per handler.

---

## Pre-Development Checklist

- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] Security team has agreed on the sanitization approach (allowlist vs. denylist)
- [ ] The boundary between framework-level and domain-level validation is agreed
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate a middleware layer that performs baseline input sanitization on all inbound HTTP requests.

### In Scope

- Path parameter sanitization: reject requests with directory traversal sequences (`../`, `..%2F`)
- Request size limit enforcement via `Content-Length` checking (configurable max body size, default 1 MB)
- Content-Type enforcement on POST/PUT/PATCH requests (must declare a content type)
- Null byte rejection in request paths and query parameters
- Unit tests for each sanitization rule

### Out of Scope

- Domain-level validation (service teams are responsible for validating business inputs)
- HTML escaping in response bodies (that is a rendering concern)
- SQL injection prevention (handled by parameterized queries in database plugins)

---

## Technical Implementation Notes

Sanitization is performed using allowlist logic per architecture-principles.md. Any input that does not conform to the allowed pattern is rejected with a 400 response in RFC 7807 format.

The maximum request body size is configurable via `MAX_REQUEST_BODY_BYTES` (default: 1048576 — 1 MB).

---

## Acceptance Criteria

- [x] Requests with `../` in the path return 400 — `TestInputSanitization_PathTraversal_DotDotSlash` passes
- [x] Requests with null bytes in query parameters return 400 — null bytes via `%00` percent-encoding in query values are caught after URL parsing; `TestInputSanitization_NullByteInQueryValue` passes. Note: raw null bytes in the URL path/query are rejected by Go's `net/url` parser before the middleware runs.
- [x] Requests with body size exceeding the configured limit return 413 — `TestInputSanitization_BodyTooLarge` passes
- [x] POST/PUT/PATCH requests without `Content-Type` return 415 — `TestInputSanitization_MissingContentType_Post`, `TestInputSanitization_MissingContentType_Put`, `TestInputSanitization_MissingContentType_Patch` all pass
- [x] `MAX_REQUEST_BODY_BYTES` overrides the default body size limit — `TestInputSanitization_MaxBodyBytes_EnvOverride` and `TestInputSanitization_MaxBodyBytes_Default` pass
- [x] Unit tests cover all sanitization rules — 12 tests covering all rules and edge cases
- [x] Clean requests pass through the middleware without modification — `TestInputSanitization_CleanRequest_PassesThrough` and `TestInputSanitization_GetWithoutContentType_PassesThrough` pass

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Each sanitization rule tested manually with `curl`
- [x] Unit tests pass — `go test ./...` green; `golangci-lint` 0 issues
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| US-1104 RFC 7807 error handler | Predecessor | Sanitization errors use this format |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
