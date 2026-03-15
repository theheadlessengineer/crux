# US-1105 — Generate Graceful Shutdown Handlers (SIGTERM/SIGINT)

**Epic:** 1.1 Tier 1 Standards Generation
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want every generated service to handle SIGTERM and SIGINT gracefully,
so that in-flight requests complete and resources are released cleanly when Kubernetes terminates a pod.

---

## Pre-Development Checklist

- [ ] Epic 1.4 Template Engine is in progress or complete
- [ ] The agreed default drain timeout is documented (recommended: 30 seconds)
- [ ] Team understands the Kubernetes pod termination lifecycle (preStop hook, terminationGracePeriodSeconds)
- [ ] Story estimated and accepted into the sprint

---

## Scope

Generate a graceful shutdown handler in the service entrypoint that listens for OS signals and drains active connections before exiting.

### In Scope

- Signal handling for `SIGTERM` and `SIGINT`
- HTTP server `Shutdown(ctx)` called with a configurable drain timeout
- Database connection pool closure after HTTP server drain
- OTel tracer flush and shutdown
- Configurable drain timeout via `SHUTDOWN_TIMEOUT_SECONDS` environment variable (default: 30)
- Structured log message on shutdown start and completion
- Exit code 0 on clean shutdown

### Out of Scope

- Pre-stop lifecycle hook configuration (generated in Kubernetes manifests story)
- Custom shutdown hooks registered by service teams (the framework provides a registry; service teams register to it)

---

## Technical Implementation Notes

The shutdown sequence must follow this order:
1. Stop accepting new connections
2. Wait for in-flight HTTP requests to complete (drain)
3. Close database connections
4. Flush and close OTel tracer
5. Exit with code 0

```go
// Pseudocode — actual implementation in template
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
defer cancel()

if err := srv.Shutdown(ctx); err != nil {
    log.Error("server forced to shutdown", "error", err)
}
```

---

## Acceptance Criteria

- [x] Generated service handles `SIGTERM` without dropping in-flight requests — `Runner.ListenAndServe()` calls `signal.Notify` for `syscall.SIGTERM`; `TestListenAndServe_SIGTERM_RunsHooksAndReturnsNil` passes
- [x] Generated service handles `SIGINT` (Ctrl+C) without dropping in-flight requests — `syscall.SIGINT` also notified; `TestListenAndServe_SIGINT_RunsHooksAndReturnsNil` passes
- [x] Drain timeout defaults to 30 seconds and is overridable via `SHUTDOWN_TIMEOUT_SECONDS` — `parseDrainTimeout()`; `TestParseDrainTimeout_Default`, `TestParseDrainTimeout_EnvOverride`, `TestParseDrainTimeout_InvalidEnv_UsesDefault` all pass
- [x] A structured log message is emitted at shutdown start — `logger.Info("shutdown signal received, draining", ...)` with timeout field; asserted in `TestListenAndServe_SIGTERM_RunsHooksAndReturnsNil`
- [x] A structured log message is emitted at shutdown completion — `logger.Info("shutdown complete")`; asserted in `TestListenAndServe_SIGTERM_RunsHooksAndReturnsNil`
- [x] Service exits with code 0 on clean shutdown — `ListenAndServe()` returns `nil`; caller checks and exits 0
- [x] Service exits with code 1 if shutdown exceeds the drain timeout — `ListenAndServe()` returns `context.DeadlineExceeded`; `TestListenAndServe_TimeoutExceeded_ReturnsContextError` passes; caller checks non-nil and exits 1
- [x] OTel tracer is flushed before exit — `tracing.Init()` returns a `func(context.Context) error` shutdown hook; caller registers it via `runner.Register(otelShutdown)`
- [x] Unit tests verify signal handling behaviour — 8 tests covering SIGTERM, SIGINT, hook ordering, hook error propagation, timeout, and env parsing

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Graceful shutdown verified manually: send SIGTERM to running service, confirm in-flight request completes
- [ ] Log output verified to show shutdown messages
- [x] Unit tests pass — `go test ./...` all green; `golangci-lint` 0 issues
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Required |
| US-1103 OTel wiring | Parallel | OTel flush in shutdown sequence |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
