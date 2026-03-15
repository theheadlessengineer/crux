# US-1102 — Generate Structured JSON Logging with Trace Correlation

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** In Progress

---

## User Story

As a user of Crux,
I want every generated service to produce structured JSON logs with distributed trace identifiers injected automatically,
so that logs are machine-parseable, searchable, and correlated with traces without any manual wiring by the service team.

---

## Pre-Development Checklist

- [x] Epic 0.1 Foundation is complete
- [ ] Epic 1.4 Template Engine is in progress or complete
- [x] The company's log schema (field names, required fields) is agreed and documented
- [x] OpenTelemetry trace propagation story (US-1103) is being developed in parallel — coordinate field names
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate a pre-configured logger in every service skeleton that outputs structured JSON to stdout, automatically injects `trace_id` and `span_id` from the request context, and enforces a consistent field schema.

### Required Log Fields

| Field | Type | Source |
|---|---|---|
| `timestamp` | ISO 8601 string | Log time |
| `level` | string | Log level (info, warn, error, debug) |
| `message` | string | Log message |
| `service` | string | Service name from config |
| `trace_id` | string | Injected from OTel context |
| `span_id` | string | Injected from OTel context |
| `environment` | string | Runtime environment (prod, staging, dev) |
| `version` | string | Service version |

### In Scope

- Logger initialization in the generated service entrypoint
- Middleware that extracts `trace_id` and `span_id` from OTel context and injects into log context
- All generated log calls using the structured logger (no `fmt.Println` or `log.Printf`)
- Log level configurable via environment variable (`LOG_LEVEL`)
- Unit tests confirming JSON structure of log output

### Out of Scope

- Log shipping configuration (that is infrastructure concern)
- PII redaction in logs (Epic 2.8 shared library)
- Log sampling configuration (Epic 3.4)
- Log retention policy (generated as a YAML stub in US-1126, not implemented here)

---

## Technical Implementation Notes

The logger must write to stdout only. Services must not write logs to files. Log collection is handled by the infrastructure layer (Kubernetes DaemonSet, Fluentd, etc.).

The logger must be passed explicitly through the application via dependency injection. No global logger variable is permitted (see architecture-principles.md — explicit over implicit).

For Go + Gin, the recommended logger is `slog` (standard library from Go 1.21) with a JSON handler, or `zerolog` if the team prefers. The choice must be documented in an ADR before implementation begins.

Log level must default to `info` in production and be overridable at runtime via the `LOG_LEVEL` environment variable without requiring a service restart (if the logger supports dynamic levels).

---

## Acceptance Criteria

- [x] Generated service produces JSON-formatted logs on stdout — `slog.NewJSONHandler` writing to `io.Writer` in `internal/infrastructure/logging/logger.go`
- [x] Every log line includes: `timestamp`, `level`, `message`, `service`, `environment`, `version` — `New()` attaches `service`, `environment`, `version`; `ReplaceAttr` renames `time` → `timestamp`; verified by `TestLogger_JSONFields`
- [x] Requests that carry a trace context produce logs with `trace_id` and `span_id` populated — `WithContext()` reads `Fields` from context; `TestLogger_WithContext_TraceFields` passes
- [x] Requests without a trace context produce logs with `trace_id` and `span_id` as empty strings, not null — `Fields{}` zero value used when no context key present; `TestLogger_WithContext_NoFields_EmptyStrings` passes
- [x] Log level is configurable via `LOG_LEVEL` environment variable — `parseLevel(os.Getenv("LOG_LEVEL"))` in `New()`; `TestLogger_LogLevel` passes
- [x] No unstructured log calls (`fmt.Println`, `log.Printf`) exist in the generated code — confirmed by code review of all files under `internal/`
- [x] Logger is injected via constructor — no global logger variable — `New()` returns `*Logger`; no package-level logger variable exists
- [x] Unit tests verify the JSON structure of log output — `TestLogger_JSONFields` unmarshals output and asserts all required keys
- [x] Tests verify that `trace_id` and `span_id` are injected when context is present — `TestLogger_WithContext_TraceFields` and `TestMiddleware_InjectsTraceFields` pass

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Generated service started locally — log output verified to be valid JSON — blocked on Epic 1.5
- [ ] Log output piped to `jq` to confirm parseability — blocked on Epic 1.5
- [ ] Trace correlation verified by checking logs against a real trace ID — blocked on Epic 1.5
- [x] Unit tests pass with `go test ./...`
- [x] Logger choice documented in an ADR — docs/adr/ADR-002-logger-selection.md
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 0.1 Foundation | Predecessor | Must be complete |
| Epic 1.4 Template Engine | Predecessor | Required to render templates |
| US-1103 OTel wiring | Parallel | Coordinate context propagation API |

---

## Definition of Done

- All acceptance criteria are met
- Generated code compiles and log output verified
- Unit tests pass
- Logger choice recorded in ADR
- Code reviewed and approved
- Committed to `main` via approved PR
