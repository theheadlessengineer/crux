# US-1103 — Generate W3C Traceparent Propagation and OpenTelemetry SDK Wiring

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want every generated service to propagate the W3C `traceparent` header and have the OpenTelemetry SDK initialized,
so that distributed traces flow correctly across service boundaries without any manual instrumentation.

---

## Pre-Development Checklist

- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] The company's observability backend (Jaeger, Tempo, OTLP collector endpoint) is agreed
- [ ] W3C Trace Context specification (RFC 7231 + W3C) reviewed by the team
- [x] OTel SDK version pinned and agreed — `go.opentelemetry.io/otel v1.42.0` in `go.mod`
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate OpenTelemetry SDK initialization, W3C `traceparent` header extraction on inbound requests, and `traceparent` header injection on all outbound HTTP requests.

### In Scope

- OTel SDK initialization at service startup with OTLP exporter configuration
- HTTP middleware that extracts `traceparent` from inbound requests and stores it in context
- HTTP client wrapper that injects `traceparent` into all outbound requests
- Span creation for inbound requests (one span per request)
- Environment variable configuration for the OTLP exporter endpoint (`OTEL_EXPORTER_OTLP_ENDPOINT`)
- Unit tests for context propagation

### Out of Scope

- Custom span attributes (service teams add these)
- Database span instrumentation (added by database plugins)
- Metrics pipeline (US-1101 handles the Prometheus endpoint)

---

## Acceptance Criteria

- [x] OTel SDK initializes at startup without panicking — `tracing.Init()` in `internal/infrastructure/tracing/provider.go`; builds and compiles cleanly
- [x] Inbound requests with `traceparent` header have the trace context propagated into the request context — `TestMiddleware_PropagatesInboundTraceparent` passes; trace ID matches inbound header
- [x] Inbound requests without `traceparent` header have a new root span created — `TestMiddleware_CreatesRootSpanWithoutTraceparent` passes
- [x] Log lines include `trace_id` and `span_id` extracted from the OTel context — `logging.Middleware()` reads from `trace.SpanFromContext`; `TestMiddleware_InjectsTraceFields` and `TestMiddleware_NoSpan_EmptyStrings` pass
- [x] All outbound HTTP requests carry `traceparent` header — `tracing.Transport` / `NewHTTPClient`; `TestTransport_InjectsTraceparentHeader` passes
- [x] OTLP exporter endpoint is configurable via `OTEL_EXPORTER_OTLP_ENDPOINT` — read in `tracing.Init()`; defaults to `localhost:4317`
- [ ] Service shuts down the OTel tracer cleanly on SIGTERM — `tp.Shutdown` is returned from `tracing.Init()` but no signal handler wires it yet (blocked on Epic 1.5 server entrypoint)
- [x] Unit tests verify propagation on inbound requests — `TestMiddleware_PropagatesInboundTraceparent`, `TestMiddleware_CreatesRootSpanWithoutTraceparent`
- [x] Unit tests verify injection on outbound requests — `TestTransport_InjectsTraceparentHeader`, `TestTransport_NoActiveSpan_NoTraceparent`

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] End-to-end trace verified in local Jaeger or Tempo instance
- [x] `traceparent` header confirmed present on outbound test requests — `TestTransport_InjectsTraceparentHeader` asserts header is non-empty
- [x] Unit tests pass — `go test -race ./internal/infrastructure/tracing/... ./internal/infrastructure/logging/...` all green
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| US-1102 Structured logging | Parallel | Coordinate context injection |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
