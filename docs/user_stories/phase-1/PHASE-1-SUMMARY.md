# Phase 1 — Pilot: Summary

**Phase:** 1 — Pilot
**Status:** Complete (technical deliverables) · Awaiting pilot team sessions
**Date:** 2026-03-16
**Goal:** Prove the concept with real teams. One language (Go + Gin). Essential plugins. `crux new` → running service in under 3 minutes.

---

## Purpose

Phase 1 establishes the complete end-to-end generation pipeline for `crux`. A single command — `crux new <service-name>` — interactively collects service requirements, selects plugins, and generates a fully structured, compilable, runnable Go + Gin microservice with all Tier 1 company standards applied automatically.

This phase is the foundation everything else builds on. If `make dev` doesn't work on a freshly generated service, nothing else matters.

---

## What Was Built

### Epic 1.1 — Tier 1 Standards (9 user stories)

All Tier 1 standards are generated into every service unconditionally. Engineers cannot skip them.

| Standard | Implementation | Location in generated service |
|---|---|---|
| Health endpoints | `/health` `/ready` `/live` `/metrics` `/version` | `internal/presentation/http/health.go` |
| Structured JSON logging | `slog` JSON handler, `LOG_LEVEL` env, trace field injection | `internal/infrastructure/logging/logger.go` |
| W3C traceparent propagation | OTel SDK, Gin middleware, HTTP client transport | `internal/infrastructure/tracing/` |
| RFC 7807 error format | Panic recovery + typed error helpers | `internal/infrastructure/errors/handler.go` |
| Graceful shutdown | SIGTERM/SIGINT handler, drain timeout, hook chain | `internal/infrastructure/shutdown/shutdown.go` |
| Security headers | CSP, X-Frame-Options, HSTS, Referrer-Policy, Permissions-Policy | `internal/presentation/http/router.go` |
| Input sanitisation | Path traversal, null bytes, CORS | `internal/presentation/http/router.go` |
| Non-root Dockerfile | Multi-stage build, `USER nonroot` | `Dockerfile` |
| K8s network policies | Default-deny ingress + egress with explicit allowlists | `infra/kubernetes/networkpolicy-*.yaml` |
| Resilience config | `resilience.yaml` — timeouts, circuit breaker, retry, bulkhead defaults | `resilience.yaml` |
| SLO definition | `slo.yaml` — availability target, p99 latency, error budget policy | `slo.yaml` |
| Alerting rules | Four golden signals: ServiceDown, HighErrorRate, HighP99Latency, ConsumerLag | `infra/monitoring/alerts.yaml` |
| Grafana dashboard | Four golden signals pre-built | `infra/monitoring/dashboard.json` |
| Cost allocation | `cost-budget.yaml` — monthly budget, component breakdown | `compliance/cost-budget.yaml` |
| Data classification | `data-classification.yaml` — PII/PHI field declarations | `compliance/data-classification.yaml` |
| Log retention | `log-retention.yaml` — per-environment retention config | `compliance/log-retention.yaml` |
| Service catalog entry | `catalog-entry.yaml` — team, owner, on-call | `compliance/catalog-entry.yaml` |
| Documentation stubs | `runbook.md`, `capacity-model.md`, `TODO.md`, ADR-001 | `docs/` |
| CI/CD pipeline | Lint → test → build → SAST → SBOM → scan | `.github/workflows/ci.yaml` |
| Secret rotation watcher | Background Vault/AWS SM re-fetch stub | `internal/infrastructure/secrets/` |

### Epic 1.2 — CLI Framework (3 user stories)

| Command | What it does |
|---|---|
| `crux new <name>` | Full interactive generation flow |
| `crux version` | Shows version, commit, build time (text + JSON) |
| `crux system` | Shows Go version, OS, arch, available tools |
| `crux validate` | Validates `.skeleton.json` + `crux.lock` in current directory |
| `--dry-run` | Prints what would be generated without writing files |
| `--no-prompt --config <file>` | Non-interactive CI/batch mode |

### Epic 1.3 — Prompt Engine (3 user stories)

- **Question types:** `text`, `number`, `confirm`, `select`, `multiselect` — all with validation
- **Decision graph:** DAG-based dependency resolution — questions only appear when their `depends_on` conditions are met; cycle detection at graph construction time
- **Session:** ordered question flow with back-navigation (`b` to go back), history stack, hidden-answer cleanup on back

### Epic 1.4 — Template Engine (1 user story)

- Go `text/template` engine over an embedded filesystem (`//go:embed all:go-gin`)
- Helper functions: `upper`, `lower`, `title`, `replace`, `camel`, `snake`, `kebab`, `contains`
- Shell scripts rendered with executable permissions (0755)
- Template variable namespace: `{{ .service.name }}`, `{{ .meta.cli_version }}`, `{{ .resilience.* }}`, `{{ .slo.* }}`, etc.

### Epic 1.5 — Core Templates (1 user story)

44 templates generating a complete Go + Gin service skeleton:

```
cmd/<name>/main.go                     — entrypoint, OTel init, graceful shutdown
internal/config/config.go              — env-driven config with defaults
internal/domain/health/registry.go     — health check registry
internal/infrastructure/logging/       — structured JSON logger + Gin middleware
internal/infrastructure/tracing/       — OTel provider, Gin middleware, HTTP transport
internal/infrastructure/shutdown/      — SIGTERM/SIGINT runner
internal/infrastructure/errors/        — RFC 7807 handler
internal/presentation/http/            — router, health handler, server
go.mod / Makefile / Dockerfile         — project root files
resilience.yaml / slo.yaml             — operational config
infra/kubernetes/                      — deployment + network policies
infra/monitoring/                      — alerts + dashboard
compliance/                            — cost, data classification, log retention, catalog
docs/                                  — runbook, capacity model, TODO, ADR
.github/workflows/ci.yaml              — CI pipeline
scripts/                               — seed, check_env, snapshot-db, restore-db
```

### Epic 1.6 — Plugin System (1 user story)

- **Manifest schema:** `plugin.yaml` with `apiVersion`, `kind`, `metadata` (name, version, trustTier, cruxVersionConstraint), `spec` (questions, templates, hooks, dependencies)
- **Trust tiers:** Tier 1 (official, bundled), Tier 2 (verified community), Tier 3 (unvetted)
- **Loader:** discovers plugins from `data/plugins/` (embedded) and `~/.crux/plugins/` (user-installed); validates manifest + semver compatibility; `"dev"` version bypasses constraint checks
- **Version compatibility:** `>=`, `>`, `<=`, `<`, `=` operators; multi-constraint AND logic

### Epic 1.7 — Pilot Plugins (9 user stories)

9 plugins bundled with crux, each with a validated `plugin.yaml` and template files:

| Plugin | Questions | Templates |
|---|---|---|
| `crux-plugin-postgresql` | pg_version, pg_read_replica, pg_audit_log, pg_migration_tool | postgres.go, health.go |
| `crux-plugin-redis` | redis_distributed_lock, redis_ttl_strategy | redis.go, health.go |
| `crux-plugin-kafka` | kafka_direction, kafka_dlq, kafka_outbox, kafka_schema_format | kafka.go, health.go |
| `crux-plugin-auth-jwt` | auth_authz_model, auth_jwks_url | jwt.go |
| `crux-plugin-kubernetes` | k8s_deployment_strategy, k8s_service_mesh, k8s_min/max_replicas | deployment, service, hpa, pdb, networkpolicy |
| `crux-plugin-terraform-aws` | aws_region, aws_environments, aws_remote_state_bucket | main.tf, variables.tf, outputs.tf, backend.tf |
| `crux-plugin-github-actions` | gha_coverage_threshold, gha_container_registry, gha_deploy_env | ci.yaml, deploy.yaml |
| `crux-plugin-prometheus` | prom_alerting_backend, prom_p99_latency_ms, prom_error_rate_threshold | alerts.yaml, dashboard.json |
| `crux-plugin-claude-code` | claude_team_name, claude_no_touch_paths | CLAUDE.md, .claude/settings.json |

### Epic 1.8 — Config & Lockfile (2 user stories)

- **`.skeleton.json`** — records crux version, generation timestamp, service metadata, all answers (core + plugin), selected plugins with versions, Tier 1 enforcement status, deviations
- **`crux.lock`** — exact plugin version pins (name → version map)
- **`--config <file>`** — YAML config file pre-fills answers for CI/batch generation
- **`--no-prompt`** — requires `--config`; validates all required fields present before running

### Epic 1.9 — Skeleton Generation (1 user story)

- `generator.Generate()` wires template engine + file map into a single call
- Renders 35+ files into the output directory
- Creates `.gitkeep` stubs for empty directories (`internal/app`, `internal/domain`, `infra/terraform`, `tests/unit`, `tests/integration`)
- Integration test: generates service → `go mod tidy` → `go build ./...` → verifies all required files present and scripts executable

### Epic 1.10 — Pilot Testing (1 user story)

- **P1 bug fixed:** nil-engine panic in `router.go.tmpl` (double `RegisterRoutes` call)
- **Verified end-to-end:** `crux new` → `go mod tidy` → `go build` → `make build` → `make test` → service starts → all 5 health endpoints return 200 → structured JSON logs emitted → graceful shutdown on SIGTERM
- **Generation time:** < 5 seconds (file generation only; `go mod tidy` adds ~10s on first network fetch)

---

## Test Coverage

| Package | Tests |
|---|---|
| `domain/health` | Registry, status types |
| `domain/lockfile` | Write, read, schema validation |
| `domain/model` | Service name validation |
| `domain/prompt` | All question types, validation, decision graph (DAG, cycles, visibility), session (back navigation, history) |
| `infrastructure/config` | YAML loading, no-prompt validation |
| `infrastructure/generator` | File map completeness, template rendering |
| `infrastructure/logging` | JSON output, log levels, trace field injection |
| `infrastructure/plugin` | Manifest parsing, validation, loader discovery, version compatibility |
| `infrastructure/resilience` | Config parsing, defaults |
| `infrastructure/secrets` | Vault/AWS SM stubs, watcher |
| `infrastructure/shutdown` | Hook execution, timeout |
| `infrastructure/template` | Engine loading, rendering, helper functions |
| `infrastructure/tracing` | OTel init, middleware, HTTP transport |
| `presentation/cli` | All commands (new, version, system, validate), flags, dry-run, no-prompt |
| `presentation/http` | Health endpoints, RFC 7807 errors, security headers, input sanitisation |
| `data/plugins` | All 9 pilot plugin manifests parse and validate |
| `data/templates` | All 44 templates render without error |
| `test/integration` | Full generation → compile check |

**17 packages, 0 failures.**

---

## Open Points

### Must resolve before pilot sessions

| # | Issue | Impact |
|---|---|---|
| 1 | **Plugin templates not rendered into generated service** — plugin questions are asked and answers recorded in `.skeleton.json`, but plugin template files (e.g. `postgres.go`, `kafka.go`) are not rendered into the output directory. The generator only renders core templates. | High — selected plugins have no effect on generated code |
| 2 | **No test files generated** — `tests/unit/` and `tests/integration/` contain only `.gitkeep`. The acceptance criterion "all generated tests pass" is vacuously true. | Medium — pilot teams will notice immediately |
| 3 | **Pilot teams not yet identified** — operational prerequisite for US-1001 | Blocker for pilot sessions |
| 4 | **Staging environment not provisioned** — required for deploy acceptance criterion | Blocker for pilot sessions |
| 5 | **Feedback collection template not prepared** — required for NPS measurement | Blocker for pilot sessions |

### Known limitations (acceptable for Phase 1)

| # | Limitation | Phase |
|---|---|---|
| 6 | Single language only (Go + Gin) — multi-language is Phase 2 | Phase 2 |
| 7 | `crux upgrade` command not implemented | Phase 3 |
| 8 | `crux audit` command not implemented | Phase 2 |
| 9 | `crux plugin search/install` commands not implemented — plugins are bundled only | Phase 2 |
| 10 | SPIFFE/resilience plugins not included — stubs generated, plugins deferred | Phase 2 |
| 11 | Compliance profiles (SOC2/GDPR/HIPAA) not implemented — files generated but no profile-specific logic | Phase 3 |
| 12 | Secret scan not in pre-commit hook — lint + format + test only | Phase 2 |
| 13 | Metrics pipeline (OTel) not wired — trace pipeline only; `/metrics` returns a stub comment | Phase 2 |
| 14 | `crux feedback` command not implemented | Phase 2 |
| 15 | `.pre-commit-config.yaml` not generated — pre-commit hooks exist in crux itself but not templated for generated services | Phase 2 |

---

## Needs Attention

**Before any pilot session runs:**

1. **Wire plugin template rendering** (open point #1) — this is the most significant gap. Plugin questions are collected but the resulting code files are never written. The generator needs to load selected plugin templates and render them alongside core templates.

2. **Generate at least one test file** — even a single smoke test (`TestHealthEndpoint`) would satisfy the acceptance criterion and give pilot teams a working test to build on.

3. **Confirm `go mod tidy` works offline** — the generated `go.mod` has no `go.sum`. First-time `make dev` requires network access to fetch dependencies. This should be documented in `TODO.md` or the `go.sum` should be pre-generated.

---

## Metrics

| Metric | Value |
|---|---|
| Go source files | ~90 |
| Lines of Go code | ~8,000 |
| Templates | 44 core + 20 plugin |
| Plugins | 9 |
| Test packages passing | 17 / 17 |
| Generation time (file I/O only) | < 5s |
| Generation time (incl. `go mod tidy`) | ~15s (first run, network) |
| Health endpoints | 5 (`/health` `/ready` `/live` `/metrics` `/version`) |
| Tier 1 standards generated | 20 (of 38 total — remainder are Phase 2/3) |

---

## Phase 2 Handoff

Phase 2 (MVP) picks up from here with:

- Multi-language support (Python, Java, Node)
- `crux plugin search/install/list` commands
- Plugin template rendering wired into generator
- `crux audit` command
- Shared component library (`@company/*`) first packages
- Additional plugins: resilience, SPIFFE, multitenant, MySQL, MongoDB, RabbitMQ, Datadog, GitLab CI
- AI API client plugins (Claude, OpenAI)
- Soft mandate rollout — all new services must use crux
