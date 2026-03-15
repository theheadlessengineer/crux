# US-1104 — Generate RFC 7807 Error Response Format Handler

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want every generated service to return errors in RFC 7807 Problem Details format,
so that all services in the organisation have a predictable, machine-readable error shape that API clients can handle generically.

---

## Pre-Development Checklist

- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] RFC 7807 specification reviewed by the team
- [ ] The company's agreed extensions to the RFC 7807 schema are documented (e.g., `trace_id` field)
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate an error response middleware and error handler that transforms all service errors into RFC 7807-compliant JSON responses.

### RFC 7807 Required Fields

| Field | Type | Description |
|---|---|---|
| `type` | URI | A URI reference identifying the problem type |
| `title` | string | A short, human-readable summary |
| `status` | integer | HTTP status code |
| `detail` | string | Human-readable explanation specific to the occurrence |
| `instance` | URI | A URI identifying the specific occurrence |

### Company Extension Fields

| Field | Type | Description |
|---|---|---|
| `trace_id` | string | OTel trace ID for correlation |

### In Scope

- Error handler middleware for Gin that intercepts panics and errors
- Helper functions for creating typed problem responses (validation error, not found, unauthorized, etc.)
- Middleware that sets `Content-Type: application/problem+json` on error responses
- Unit tests covering each error type

### Out of Scope

- Error logging (handled by logging middleware, not error handler)
- Domain-specific error types (service teams define these using the helpers)

---

## Acceptance Criteria

- [x] All 4xx and 5xx responses from the generated service use `Content-Type: application/problem+json` — `writeProblem` sets the header; `TestProblem_ContentTypeOnAllErrors` verifies all error types including panic recovery
- [x] All error responses include `type`, `title`, `status`, `detail`, `instance`, and `trace_id` — `Problem` struct; `assertProblem` helper verifies all fields in every test
- [x] Validation errors return 400 with a `type` of `/errors/validation` — `ValidationError()`; `TestValidationError` passes
- [x] Not Found errors return 404 with a `type` of `/errors/not-found` — `NotFound()`; `TestNotFound` passes
- [x] Unhandled panics are recovered and return 500 (never expose stack traces in response) — `ErrorMiddleware()` uses `gin.CustomRecovery`; `TestPanicRecovery_Returns500_NoStackTrace` asserts no goroutine/panic strings in body
- [x] Stack traces are logged (not returned to caller) — `gin.CustomRecovery` logs to the default writer; detail is a fixed safe string
- [x] Unit tests verify the response shape for each error type — `TestValidationError`, `TestNotFound`, `TestUnauthorized`, `TestInternalError`, `TestPanicRecovery_Returns500_NoStackTrace`, `TestProblem_InstanceIsRequestPath`, `TestProblem_ContentTypeOnAllErrors`
- [x] `Content-Type` header is always `application/problem+json` on error responses — `TestProblem_ContentTypeOnAllErrors` covers all five paths

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Error responses validated against the RFC 7807 schema manually
- [ ] Panic recovery tested manually with a deliberate panic in a handler
- [x] Unit tests pass — `go test ./...` green; `golangci-lint` 0 issues
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| US-1103 OTel wiring | Parallel | `trace_id` must be available in context |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
