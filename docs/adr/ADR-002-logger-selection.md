# ADR-002 — Logger Selection: slog (standard library)

**Status:** Accepted
**Date:** 2026-03-15
**Deciders:** Platform Engineering

---

## Context

US-1102 requires a structured JSON logger for every generated service. The story
explicitly names two candidates: `slog` (Go standard library, Go 1.21+) and
`zerolog` (third-party).

## Decision

Use `log/slog` with `slog.NewJSONHandler`.

## Rationale

| Criterion | slog | zerolog |
|---|---|---|
| Dependency | None — stdlib | External dependency |
| JSON output | `NewJSONHandler` | Native |
| Performance | Adequate for services | Marginally faster |
| Dynamic level | `LevelVar` | Supported |
| Go version req | 1.21+ (project uses 1.26) | Any |

The project already targets Go 1.26. Using stdlib eliminates an external
dependency, reduces supply-chain risk, and is sufficient for the logging
throughput of a typical microservice. `zerolog` offers no meaningful advantage
that justifies the added dependency.

## Consequences

- No additional `go.mod` entry required for logging.
- `slog.Logger` is the canonical logger type passed through the application.
- No global `slog.SetDefault` calls — logger is injected via constructor.
- If profiling reveals logging as a bottleneck in a specific service, that
  service may adopt `zerolog` locally — this is a per-service decision, not
  a platform decision.
