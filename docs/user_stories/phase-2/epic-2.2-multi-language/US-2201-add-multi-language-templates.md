# US-2201 — Add Python, Java, and Node.js Language Templates

**Epic:** 2.2 Multi-Language Support
**Phase:** 2 — MVP
**Priority:** Should Have
**Status:** Done

---

## User Story

As a user of Crux,
I want to generate services in Python + FastAPI, Java + Spring Boot, and Node.js + Express,
so that teams using these languages can benefit from Crux in the same way as Go teams.

---

## What Was Built

### Language Selection Question

A `language` select question was added as the **first** core question in the `crux new` flow:

| Value | Label |
|---|---|
| `go` | Go + Gin (default) |
| `python` | Python + FastAPI |
| `java` | Java + Spring Boot |
| `node` | Node.js + Express |

The language answer gates all subsequent language-specific behaviour — plugin question variants, template selection, and framework defaults.

---

### Templates

Three new template directories were added under `data/templates/`:

#### `python-fastapi/`

| File | Purpose |
|---|---|
| `main.py` | FastAPI app with lifespan, OTel init, RFC 7807 error handler, SIGTERM handler |
| `app/config.py` | Env-driven settings |
| `app/health.py` | All 5 health endpoints (`/health` `/ready` `/live` `/metrics` `/version`) |
| `app/logging_config.py` | Structured JSON logger matching the org-wide log schema |
| `app/middleware.py` | `SecurityHeadersMiddleware` + `TraceparentMiddleware` (W3C traceparent) |
| `requirements.txt` | FastAPI, uvicorn, opentelemetry-sdk, opentelemetry-instrumentation-fastapi |
| `Makefile` | `build`, `run`, `test`, `lint`, `fmt`, `dev` targets |
| `Dockerfile` | Multi-stage, non-root (`nonroot` user), HEALTHCHECK |
| `github/workflows/ci.yaml` | Python CI: lint → test → docker build → SBOM |

#### `java-spring/`

| File | Purpose |
|---|---|
| `src/main/java/Application.java` | Spring Boot entrypoint |
| `src/main/java/health/HealthController.java` | All 5 health endpoints |
| `src/main/resources/application.yaml` | Structured JSON log pattern, graceful shutdown, Prometheus actuator |
| `pom.xml` | Spring Boot 3.3, Java 21, OTel instrumentation, Micrometer Prometheus |
| `Makefile` | `build`, `run`, `test`, `lint`, `dev` targets |
| `Dockerfile` | Multi-stage (eclipse-temurin:21), non-root, HEALTHCHECK |
| `github/workflows/ci.yaml` | Java CI: test → build JAR → docker build → SBOM |

#### `node-express/`

| File | Purpose |
|---|---|
| `index.js` | Server entrypoint with SIGTERM/SIGINT graceful shutdown |
| `src/app.js` | Express app with middleware chain + RFC 7807 error handler |
| `src/health.js` | All 5 health endpoints |
| `src/logging.js` | Structured JSON logger (zero external deps) |
| `src/middleware.js` | Security headers + W3C traceparent middleware |
| `package.json` | Express, Node 20, dev scripts |
| `Makefile` | `build`, `run`, `test`, `lint`, `dev` targets |
| `Dockerfile` | Multi-stage (node:20-slim), non-root, HEALTHCHECK |
| `github/workflows/ci.yaml` | Node CI: lint → test → docker build → SBOM |

#### Shared (language-agnostic) files

All languages share the same Tier 1 YAML/Markdown files sourced from `go-gin` templates:

- `resilience.yaml`, `slo.yaml`
- `infra/kubernetes/` — deployment, network policies
- `infra/monitoring/` — alerts, dashboard
- `compliance/` — catalog-entry, cost-budget, data-classification, log-retention
- `docs/` — runbook, capacity-model, TODO, ADR-001
- `.editorconfig`, `.commitlintrc.yaml`, `.gitignore`, `.envrc`, `README.md`, `CHANGELOG.md`

---

### Generator Changes (`internal/infrastructure/generator/generator.go`)

- `fileMap()` dispatches on `cfg.Language` → `go` | `python` | `java` | `node`
- `sharedFileMap()` extracts language-agnostic Tier 1 files reused across all languages
- `mergeMaps()` combines shared + language-specific file maps
- `emptyDirs()` returns language-appropriate stub directories per language
- `data/templates/embed.go` updated to embed all four language directories

### CLI Changes (`internal/presentation/cli/new.go`)

- `language` added as the first `coreQuestion` (select, default `go`)
- `buildGeneratorConfig` reads the `language` answer and sets `Framework` via `defaultFramework()`

---

### Language-Aware Plugin Questions

Plugin questions that have different meaningful options per language are resolved at question-build time using two new fields on `QuestionSpec`:

```go
OptionsByLang  map[string][]string  // yaml: options_by_language
DefaultByLang  map[string]string    // yaml: default_by_language
```

For questions with `options_by_language` declared, `runPrompt` emits **one question variant per language**, each gated by `DependsOn: {language == <lang>} AND {_plugins contains <plugin>}`. Only the variant matching the user's language selection is visible.

#### Example — `crux-plugin-postgresql` migration tool

```yaml
- id: pg_migration_tool
  type: select
  prompt: "Database migration tool?"
  options: ["goose", "migrate"]
  default: "goose"
  options_by_language:
    go:     ["goose", "migrate"]
    python: ["alembic", "flyway"]
    java:   ["flyway", "liquibase"]
    node:   ["db-migrate", "knex", "flyway"]
  default_by_language:
    go:     "goose"
    python: "alembic"
    java:   "flyway"
    node:   "db-migrate"
```

A Go user sees `goose / migrate`. A Java user sees `flyway / liquibase`. A Python user sees `alembic / flyway`. The wrong options are never shown.

---

## Acceptance Criteria

- [x] `crux new` with language selection generates a working service for each language
- [x] Generated service in each language starts and passes its health checks
- [x] All Tier 1 standards are present for each language
- [x] Docker build succeeds for each language (non-root user in all Dockerfiles)
- [x] Language-specific CI pipeline generated for each language
- [x] Plugin questions with language-specific options show only the relevant options for the selected language

---

## Files Changed

| File | Change |
|---|---|
| `data/templates/embed.go` | Embed all four language dirs |
| `data/templates/python-fastapi/` | New — 9 templates |
| `data/templates/java-spring/` | New — 7 templates |
| `data/templates/node-express/` | New — 9 templates |
| `internal/infrastructure/generator/generator.go` | Multi-language dispatch, shared file map, `mergeMaps` |
| `internal/infrastructure/generator/generator_test.go` | 9 new tests covering all 3 new languages |
| `internal/presentation/cli/new.go` | `language` core question, `buildGeneratorConfig` language capture, language-aware `pluginQuestionToPrompt`, per-language question variants in `runPrompt` |
| `internal/domain/plugin/plugin.go` | `OptionsByLang`, `DefaultByLang` fields on `QuestionSpec` |
| `data/plugins/crux-plugin-postgresql/plugin.yaml` | `options_by_language` + `default_by_language` on `pg_migration_tool` |

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.4 Template Engine | Predecessor | Supports multi-language templates |
| Phase 1 complete | Prerequisite | Baseline established |

---

## Definition of Done

- All acceptance criteria are met
- All three languages generate working services
- All tests pass (`go test ./...`)
- Committed to `main` via approved PR
