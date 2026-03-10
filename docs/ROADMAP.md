# Crux Development Roadmap - Epics & User Stories

## Overview

This document defines the high-level epics and user stories for building Crux, the microservice skeleton generator CLI. Each epic can be broken down into smaller tasks during sprint planning.

## Epic Summary

**Phase 0 (Foundation):** 1 epic  
**Phase 1 (Pilot):** 10 epics  
**Phase 2 (MVP):** 10 epics  
**Phase 3 (Core):** 10 epics  
**Phase 4 (Future):** 2 epics  

**Total:** 33 active epics

---

## Phase 0: Foundation

### Epic 0.1: Project Setup & Infrastructure

**As a** platform engineer  
**I want** the project structure and development environment set up  
**So that** the team can start development with consistent tooling

**Stories:**
- Initialize Go module and project structure
- Set up CI/CD pipeline (GitHub Actions)
- Configure linting (golangci-lint) and formatting (gofmt)
- Set up pre-commit hooks
- Create Makefile with common commands
- Configure dependency management (go.mod)
- Set up test infrastructure
- Create initial README and contribution guidelines

**Acceptance Criteria:**
- `make build` produces binary
- `make test` runs all tests
- `make lint` passes
- CI pipeline runs on every PR
- Pre-commit hooks enforce standards

---

## Phase 1: Pilot

### Epic 1.1: Tier 1 Standards Generation

**As a** user  
**I want** all 38 Tier 1 standards automatically generated  
**So that** my service is production-ready from day one

**Stories:**
- Generate health endpoints (`/health`, `/ready`, `/live`, `/metrics`, `/version`)
- Generate structured JSON logging with trace correlation
- Generate W3C `traceparent` propagation
- Generate OpenTelemetry SDK wiring
- Generate RFC 7807 error format handler
- Generate graceful shutdown handlers (SIGTERM/SIGINT)
- Generate security headers middleware (CORS, XSS, CSRF)
- Generate input sanitization middleware
- Generate non-root Dockerfile
- Generate read-only root filesystem K8s config
- Generate network policies (ingress + egress default-deny)
- Generate secrets via Vault/AWS Secrets Manager config
- Generate secret rotation watcher
- Generate least-privilege IAM/ServiceAccount
- Generate pre-commit hooks (lint, format, secret scan)
- Generate SBOM generation in CI
- Generate dependency license scan in CI
- Generate DAST slot in CI (disabled by default)
- Generate `resilience.yaml` with defaults
- Generate `slo.yaml` with targets
- Generate `alerts.yaml` with 4 baseline alerts
- Generate Grafana dashboard JSON
- Generate cost allocation tags
- Generate `cost-budget.yaml`
- Generate `data-classification.yaml` stub
- Generate `log-retention.yaml`
- Generate `catalog-entry.yaml`
- Generate API breaking change check (oasdiff)
- Generate decommission checklist
- Generate runbook template
- Generate ADR folder with first entry
- Generate Makefile with standard targets
- Generate conventional commits config
- Generate CHANGELOG.md
- Generate capacity model template
- Generate TODO.md with intentional placeholders
- Generate .editorconfig
- Generate license headers

**Acceptance Criteria:**
- All 38 standards are generated for every service
- Generated service passes all Tier 1 validation checks
- Standards cannot be disabled or removed
- Documentation explains each standard

---

### Epic 1.2: CLI Framework & Commands

**As a** user  
**I want** a well-structured CLI interface  
**So that** I can easily interact with Crux

**Stories:**
- Set up Cobra CLI framework
- Implement `crux new` command skeleton
- Implement `crux version` command
- Implement `crux system` command
- Implement `crux validate` command
- Add global flags (--verbose, --output, --config)
- Implement help text for all commands
- Add command aliases
- Implement exit codes (0=success, 1=error, 2=validation)
- Write CLI integration tests

**Acceptance Criteria:**
- All commands have help text
- `crux --help` shows all commands
- `crux <command> --help` shows command-specific help
- Exit codes are consistent
- Flags work across all commands

---

### Epic 1.3: Interactive Prompt Engine & Decision Graph

**As a** user  
**I want** an intelligent prompt system with conditional logic  
**So that** I only see relevant questions and get smart recommendations

**Stories:**
- Implement question types (confirm, text, number, select, multiselect)
- Implement conditional question logic (depends_on with AND/OR)
- Implement question grouping and ordering
- Implement default value handling
- Implement validation for each question type
- Implement decision dependency graph (DAG)
- Implement auto-additions (e.g., Kafka + PostgreSQL → recommend Outbox)
- Implement warning system (non-blocking, shows)
- Implement error system (blocking, shows)
- Add search/filter for select questions
- Implement question history/back navigation
- Create prompt engine interface
- Write unit tests for prompt engine
- Write tests for decision graph logic

**Acceptance Criteria:**
- All question types work correctly
- Conditional questions show/hide based on answers
- Decision graph resolves dependencies correctly
- Auto-additions suggest complementary plugins
- Warnings are shown but don't block
- Errors block generation with clear messages
- Validation prevents invalid input
- User can navigate back to previous questions

---

### Epic 1.4: Template Engine

**As a** developer  
**I want** a template rendering system  
**So that** I can generate code from templates

**Stories:**
- Implement Go text/template wrapper
- Create template variable namespace
- Implement template helper functions (upper, lower, replace, etc.)
- Implement template loading from embedded files
- Implement template rendering to filesystem
- Add support for conditional blocks
- Add support for loops
- Implement template validation
- Create template testing utilities
- Write unit tests for template engine

**Acceptance Criteria:**
- Templates render with correct variable substitution
- Helper functions work correctly
- Conditional blocks work
- Loops work
- Invalid templates are caught at startup
- Generated files have correct permissions

---

### Epic 1.5: Core Templates (Go + Gin)

**As a** user  
**I want** templates for Go with Gin framework  
**So that** I can generate a working Go service

**Stories:**
- Create main.go template
- Create go.mod template
- Create Dockerfile template
- Create .gitignore template
- Create README.md template
- Create Makefile template
- Create health endpoint template
- Create logging configuration template
- Create error handler template
- Create configuration loader template
- Embed templates in binary
- Write template rendering tests

**Acceptance Criteria:**
- Generated service compiles with `go build`
- Generated service runs with `go run`
- Generated tests pass with `go test`
- Docker build succeeds
- All Tier 1 standards are present

---

### Epic 1.6: Plugin System - Core Infrastructure

**As a** platform engineer  
**I want** a plugin system architecture  
**So that** Crux can be extended without modifying core

**Stories:**
- Define plugin manifest schema (plugin.yaml)
- Implement plugin manifest parser
- Implement plugin manifest validator
- Create plugin loader interface
- Implement plugin discovery mechanism
- Implement plugin version compatibility checking
- Create plugin sandbox/isolation strategy
- Implement plugin lifecycle hooks (pre-generate, post-generate)
- Write plugin system tests
- Document plugin development guide

**Acceptance Criteria:**
- Plugins can be loaded from filesystem
- Invalid plugins are rejected with clear errors
- Plugin versions are validated
- Hooks execute in correct order
- Plugins are isolated from core

---

### Epic 1.7: Pilot Plugins (9 Essential Plugins)

**As a** user  
**I want** essential integration plugins  
**So that** I can build production services

**Stories:**
- Create `crux-plugin-postgresql` (pool, migrations, read replica, audit logging)
- Create `crux-plugin-redis` (caching, distributed lock, TTL strategy)
- Create `crux-plugin-kafka` (producer, consumer, DLQ, outbox, schema registry)
- Create `crux-plugin-auth-jwt` (JWT validation, RBAC/ABAC stubs)
- Create `crux-plugin-kubernetes` (deployment, service, HPA, PDB, network policy)
- Create `crux-plugin-terraform-aws` (RDS, ElastiCache, MSK, IAM)
- Create `crux-plugin-github-actions` (CI/CD pipeline with all Tier 1 checks)
- Create `crux-plugin-prometheus` (metrics, alerts, dashboard)
- Create `crux-plugin-claude-code` (CLAUDE.md, .claude/ config)
- Write plugin tests for all 9 plugins
- Document all plugins with examples
- Test all plugins with pilot language

**Acceptance Criteria:**
- All 9 plugins load successfully
- Questions appear in correct order
- Templates generate correctly for pilot language
- Generated code compiles and runs
- Plugins can be combined
- Documentation is complete

---

### Epic 1.8: Configuration & Lockfile System

**As a** user  
**I want** to save and reuse configurations  
**So that** I can generate similar services quickly

**Stories:**
- Define `.skeleton.json` schema (all decisions + versions + deviations)
- Define `crux.lock` schema (exact plugin versions)
- Implement configuration file format (YAML)
- Implement configuration loader
- Implement configuration validator
- Add `--config` flag to `crux new`
- Implement `.skeleton.json` generation during `crux new`
- Implement `crux.lock` generation during `crux new`
- Implement configuration merging (CLI flags > config file > defaults)
- Implement lockfile reading for `crux upgrade`
- Implement lockfile validation for `crux audit`
- Add deviation tracking in `.skeleton.json`
- Add configuration examples
- Write configuration tests
- Write lockfile tests

**Acceptance Criteria:**
- `crux new --config myconfig.yaml` works
- `.skeleton.json` captures all decisions and metadata
- `crux.lock` captures exact plugin versions
- Configuration validation catches errors
- CLI flags override config file values
- Lockfiles are valid JSON
- Deviations are tracked

---

### Epic 1.9: Complete Skeleton Generation

**As a** user  
**I want** a complete, production-ready project structure  
**So that** I can start coding business logic immediately

**Stories:**
- Generate complete directory structure (app/, infra/, tests/, scripts/, docs/)
- Generate application code (main, config, logging, errors, health)
- Generate Kubernetes manifests (deployment, service, configmap, secret, HPA, PDB)
- Generate network policies (ingress + egress)
- Generate Terraform modules (database, cache, messaging, IAM)
- Generate Docker Compose for local development
- Generate Dockerfile (multi-stage, non-root, health check)
- Generate CI/CD pipeline with all Tier 1 checks
- Generate monitoring (alerts.yaml, dashboard.json)
- Generate SLO definition (slo.yaml)
- Generate resilience config (resilience.yaml)
- Generate cost budget (cost-budget.yaml)
- Generate data classification (data-classification.yaml)
- Generate log retention (log-retention.yaml)
- Generate service catalog entry (catalog-entry.yaml)
- Generate documentation (README, runbook, ADR, capacity model)
- Generate scripts (seed, check_env, snapshot-db, restore-db)
- Generate Makefile with all standard targets
- Generate .gitignore, .editorconfig, .envrc
- Generate TODO.md with intentional placeholders
- Write skeleton generation tests
- Test generated skeleton compiles
- Test generated skeleton runs
- Test generated tests pass

**Acceptance Criteria:**
- Generated service has complete directory structure
- All Tier 1 standards are present
- Service compiles with no errors
- Service runs with `make dev`
- All generated tests pass
- Docker build succeeds
- Kubernetes manifests are valid
- Terraform plan succeeds
- Documentation is complete

**As a** pilot team member  
**I want** to generate a real service with Crux  
**So that** I can validate it works end-to-end

**Stories:**
- Generate test service with pilot team
- Verify generated service compiles
- Verify generated service runs
- Verify generated tests pass
- Verify Docker build works
- Collect pilot feedback
- Fix critical bugs
- Update documentation based on feedback
- Measure generation time (target: <3 minutes)

**Acceptance Criteria:**
- 2 pilot teams successfully generate services
- Generated services deploy to staging
- No P1 bugs remain
- Generation time < 3 minutes
- Pilot teams provide positive feedback

---

### Epic 1.10: End-to-End Pilot Testing

**As a** pilot team member  
**I want** to generate a real service with Crux  
**So that** I can validate it works end-to-end

**Stories:**
- Generate test service with pilot team
- Verify generated service compiles
- Verify generated service runs
- Verify generated tests pass
- Verify Docker build works
- Verify Kubernetes deployment works
- Verify all Tier 1 standards are present
- Verify monitoring works (metrics, alerts, dashboard)
- Verify CI/CD pipeline works
- Collect pilot feedback
- Fix critical bugs
- Update documentation based on feedback
- Measure generation time (target: <3 minutes)
- Measure `make dev` success rate (target: 100%)

**Acceptance Criteria:**
- 2 pilot teams successfully generate services
- Generated services deploy to staging
- No P1 bugs remain
- Generation time < 3 minutes
- `make dev` success rate: 100%
- Pilot teams provide positive feedback
- Net Promoter Score: positive

---

## Phase 2: MVP

### Epic 2.1: Terminal UI (TUI)

**As a** user  
**I want** a beautiful terminal interface  
**So that** I have a better experience than plain CLI prompts

**Stories:**
- Set up Bubbletea framework
- Implement TUI layout system
- Create question display component
- Create progress indicator component
- Create summary display component
- Implement keyboard navigation
- Implement theme system (2 themes - light and dark)
- Add theme persistence (~/.crux/theme)
- Implement help overlay
- Write TUI tests

**Acceptance Criteria:**
- TUI renders correctly in all terminals
- Keyboard navigation works (arrows, vim keys)
- Themes can be cycled with `t` key
- Theme preference persists
- Help overlay shows all keybindings

---

### Epic 2.2: Multi-Language Support

**As a** user  
**I want** to generate services in multiple languages  
**So that** I can use Crux for all my projects

**Stories:**
- Add Python + FastAPI templates
- Add Java + Spring Boot templates
- Add Node.js + Express templates
- Implement language-specific template selection
- Add language-specific dependency management
- Add language-specific Dockerfiles
- Add language-specific CI pipelines
- Test all language templates
- Document language-specific features

**Acceptance Criteria:**
- All 4 languages generate working services
- Each language has complete template set
- Generated services pass language-specific tests
- Dockerfiles work for all languages

---

### Epic 2.3: Additional Official Plugins (15 MVP Plugins)

**As a** user  
**I want** more integration options  
**So that** I can build services with different tech stacks

**Stories:**
- Create `crux-plugin-resilience` (circuit breaker, bulkhead, timeout, retry)
- Create `crux-plugin-spiffe` (SPIFFE/SPIRE workload identity, OPA policies)
- Create `crux-plugin-multitenant` (tenant context, per-tenant rate limiting)
- Create `crux-plugin-mysql` (connection pool, migrations, read replica)
- Create `crux-plugin-mongodb` (connection, indexes, change streams)
- Create `crux-plugin-rabbitmq` (producer, consumer, DLQ)
- Create `crux-plugin-gitlab-ci` (CI/CD pipeline)
- Create `crux-plugin-datadog` (metrics, traces, logs)
- Create `crux-plugin-terraform-gcp` (Cloud SQL, Memorystore, Pub/Sub)
- Create `crux-plugin-terraform-azure` (Azure SQL, Redis, Service Bus)
- Create `crux-plugin-grpc` (service definition, client, server)
- Create `crux-plugin-claude-api` (AI client with safety guardrails)
- Create `crux-plugin-openai` (AI client with safety guardrails)
- Create `crux-plugin-github-copilot` (copilot config)
- Create `crux-plugin-cursor` (.cursorrules config)
- Create Istio service mesh templates (VirtualService, DestinationRule, PeerAuthentication, AuthorizationPolicy)
- Test all plugins with all 4 languages
- Document all plugins

**Acceptance Criteria:**
- All 15 plugins work with all 4 languages
- Plugins can be combined
- Generated code compiles and runs
- Service mesh templates disable in-app resilience when mesh is selected
- AI plugins include safety guardrails (PII scrubbing, cost limits)
- Documentation is complete

---

### Epic 2.4: Plugin Registry & Discovery

**As a** user  
**I want** to discover and install plugins  
**So that** I can extend Crux functionality

**Stories:**
- Implement `crux plugin search` command
- Implement `crux plugin install` command
- Implement `crux plugin list` command
- Implement `crux plugin update` command
- Implement `crux plugin info` command
- Create plugin registry JSON format
- Set up plugin registry hosting
- Implement plugin caching (~/.crux/plugins/)
- Add plugin tier badges (Tier 1/2/3)
- Write plugin management tests

**Acceptance Criteria:**
- `crux plugin search kafka` finds plugins
- `crux plugin install crux-plugin-kafka` works
- Plugins are cached locally
- Plugin tiers are displayed
- Plugin updates work

---

### Epic 2.5: Validation & Combination Rules

**As a** user  
**I want** Crux to prevent invalid configurations  
**So that** I don't generate broken services

**Stories:**
- Implement validator engine
- Add cross-plugin validators (e.g., outbox requires DB + Kafka)
- Add warning system (non-blocking)
- Add error system (blocking)
- Implement validator expression language
- Add validator tests for all combinations
- Document validation rules
- Add `crux validate` command

**Acceptance Criteria:**
- Invalid combinations are blocked
- Warnings are shown but don't block
- Error messages are clear and actionable
- `crux validate` checks config files

---

### Epic 2.6: Audit & Compliance

**As a** platform engineer  
**I want** to audit all generated services  
**So that** I can ensure compliance

**Stories:**
- Implement `crux audit` command
- Add `.crux.lock` file parsing
- Implement compliance checker
- Add Tier 1 standards verification
- Add plugin version checking
- Generate audit report (table format)
- Add JSON output for automation
- Write audit tests

**Acceptance Criteria:**
- `crux audit` scans all services in directory
- Report shows compliance status
- Outdated services are flagged
- JSON output works for CI/CD

---

### Epic 2.7: Documentation & Examples

**As a** user  
**I want** comprehensive documentation  
**So that** I can learn Crux quickly

**Stories:**
- Write getting started guide
- Write user guide
- Write plugin development guide
- Create example configurations
- Create video tutorials
- Write troubleshooting guide
- Create FAQ
- Set up documentation site

**Acceptance Criteria:**
- Documentation covers all features
- Examples work correctly
- Troubleshooting guide addresses common issues
- Documentation site is searchable

---

### Epic 2.8: Shared Component Library

**As a** platform engineer  
**I want** reusable component libraries  
**So that** all services have consistent implementations

**Stories:**
- Create `@company/logger` (structured JSON, OTel injection, sampling, PII redaction)
- Create `@company/http-client` (retry, circuit breaker, bulkhead, trace propagation)
- Create `@company/kafka-client` (DLQ, outbox relay, schema registry, at-least-once)
- Create `@company/error-handler` (RFC 7807 middleware for all frameworks)
- Create `@company/health-check` (health endpoints implementation)
- Create `@company/auth-middleware` (JWT/mTLS/API key + RBAC/ABAC)
- Create `@company/db-utils` (connection pool, soft delete, audit columns)
- Create `@company/redis-client` (Redlock, TTL strategy, reconnect)
- Create `@company/config-loader` (env validation, Vault/AWS SM fetch, rotation watcher)
- Create `@company/tracing` (OTel SDK, span helpers, context propagation)
- Create `@company/resilience` (circuit breaker, bulkhead, timeout, retry budget, mesh-aware)
- Create `@company/ai-client` (retry, rate limits, cost tracking, PII scrub, prompt versioning)
- Create `@company/data-classifier` (PII/PHI annotation, log redaction, lineage events)
- Create `@company/tenant-context` (tenant ID extraction, context propagation, cache namespacing)
- Create `@company/workload-id` (SPIFFE SVID management, mTLS cert rotation)
- Publish all libraries to internal registry
- Version all libraries with semver
- Write tests for all libraries
- Document all libraries
- Update templates to use shared libraries

**Acceptance Criteria:**
- All 15 libraries are published
- Libraries work across all 4 languages
- CVE fix in library propagates to all services
- Documentation is complete
- Generated services use shared libraries

---

### Epic 2.9: Developer Inner Loop

**As a** developer  
**I want** excellent local development experience  
**So that** I can iterate quickly

**Stories:**
- Generate Tiltfile for local K8s development
- Generate contract testing setup (Pact provider/consumer)
- Generate `.envrc` for direnv secret injection from Vault
- Generate database snapshot scripts (`make snapshot-db`)
- Generate database restore scripts (`make restore-db`)
- Generate `make reset-db` (drop, recreate, migrate, seed)
- Generate Docker Compose with all dependencies
- Generate hot-reload configuration
- Add UI tools to Docker Compose (pgAdmin, RedisInsight, Kafdrop, Jaeger UI)
- Write developer inner loop tests
- Document local development workflow

**Acceptance Criteria:**
- `make dev` starts all dependencies
- Tiltfile enables hot-reload in local K8s
- Contract tests run in CI
- Secrets inject automatically with direnv
- Database state can be snapshot and restored
- All UI tools are accessible
- Documentation covers full workflow

---

### Epic 2.10: Adoption & Rollout Strategy

**As a** platform engineer  
**I want** a structured adoption plan  
**So that** teams adopt Crux successfully

**Stories:**
- Create onboarding materials
- Create lunch & learn presentation
- Create engineering all-hands presentation
- Create video tutorials
- Implement `crux feedback` command
- Set up feedback collection mechanism
- Create adoption metrics dashboard
- Implement usage tracking (anonymous)
- Create soft mandate communication plan
- Create hard mandate communication plan
- Schedule training sessions
- Create FAQ based on pilot feedback
- Document migration path from existing services

**Acceptance Criteria:**
- Onboarding materials are complete
- Training sessions scheduled
- Feedback mechanism works
- Adoption metrics tracked
- Communication plan executed
- FAQ addresses common questions

---

## Phase 3: Core

### Epic 3.1: Upgrade System

**As a** user  
**I want** to upgrade existing services  
**So that** I can get new features and fixes

**Stories:**
- Implement `crux upgrade` command
- Add upgrade strategy (patch/minor/major)
- Implement file conflict detection
- Implement merge strategy
- Add dry-run mode
- Implement rollback mechanism
- Add upgrade tests
- Document upgrade process

**Acceptance Criteria:**
- `crux upgrade` updates existing services
- Conflicts are detected and reported
- Dry-run shows what would change
- Rollback works if upgrade fails

---

### Epic 3.2: Resilience Patterns

**As a** user  
**I want** resilience patterns built-in  
**So that** my services handle failures gracefully

**Stories:**
- Create resilience.yaml schema
- Implement circuit breaker template
- Implement retry with backoff template
- Implement timeout configuration
- Implement bulkhead pattern template
- Add resilience plugin
- Add mesh mode detection (disable in-app when mesh present)
- Write resilience tests
- Document resilience patterns

**Acceptance Criteria:**
- resilience.yaml is generated
- Circuit breaker works
- Retry logic works
- Timeouts are enforced
- Mesh mode disables in-app resilience

---

### Epic 3.3: Zero-Trust & Security

**As a** user  
**I want** security built-in from day one  
**So that** my services are secure by default

**Stories:**
- Create SPIFFE/SPIRE plugin
- Implement workload identity templates
- Create OPA policy templates
- Add mTLS configuration
- Implement secret rotation watcher
- Add security headers middleware
- Add input sanitization middleware
- Create network policy templates
- Write security tests

**Acceptance Criteria:**
- SPIFFE identity is configured
- OPA policies are generated
- mTLS works between services
- Secrets rotate automatically
- Network policies are applied

---

### Epic 3.4: Observability & SLOs

**As a** user  
**I want** full observability  
**So that** I can monitor and debug my services

**Stories:**
- Create slo.yaml schema
- Implement SLO definition templates
- Create alert rules templates
- Create dashboard templates
- Add distributed tracing configuration
- Add structured logging templates
- Implement metrics collection
- Create observability plugin
- Write observability tests

**Acceptance Criteria:**
- slo.yaml is generated
- Alerts are configured
- Dashboards are created
- Tracing works end-to-end
- Logs are structured

---

### Epic 3.5: Compliance Profiles (Full Implementation)

**As a** user  
**I want** complete compliance profiles  
**So that** my services meet all regulatory requirements

**Stories:**
- Implement SOC2 compliance profile (immutable audit logs, access control, change management)
- Implement GDPR compliance profile (full implementation)
- Implement HIPAA compliance profile (PHI handling, BAA reference)
- Create `data-classification.yaml` schema and validation
- Implement GDPR erasure handler (delete/anonymize across DB, cache, events, S3, backups)
- Implement data subject request SLA tracking
- Generate retention policy enforcement in Terraform (S3 lifecycle, RDS backups, Kafka retention)
- Generate data residency OPA/Sentinel enforcement in Terraform
- Generate immutable audit log configuration (S3 object lock)
- Create compliance documentation templates
- Write compliance tests
- Document all compliance requirements

**Acceptance Criteria:**
- All 3 profiles generate correctly
- GDPR erasure handler works end-to-end
- Data residency is enforced at Terraform plan time
- Audit logs are immutable
- Retention policies are enforced
- Compliance documentation is complete
- `make gdpr-erase USER_ID=xxx` works

---

### Epic 3.6: Performance & Optimization

**As a** developer  
**I want** Crux to be fast  
**So that** developers have a great experience

**Stories:**
- Profile startup time
- Optimize template loading
- Optimize plugin loading
- Add parallel plugin execution
- Implement caching strategy
- Add benchmark tests
- Optimize TUI rendering
- Measure and document performance

**Acceptance Criteria:**
- Startup time < 100ms
- Template rendering < 100ms
- Full generation < 3 minutes
- TUI runs at 60 FPS
- Benchmarks pass

---

### Epic 3.7: Community & Ecosystem

**As a** community member  
**I want** to contribute plugins  
**So that** I can extend Crux for my needs

**Stories:**
- Create plugin contribution guide
- Set up plugin review process
- Create plugin template repository
- Implement plugin testing framework
- Set up community plugin registry
- Create plugin showcase
- Write plugin best practices
- Set up plugin CI/CD

**Acceptance Criteria:**
- Contribution guide is clear
- Plugin template works
- Review process is documented
- Community plugins can be published

---

### Epic 3.8: AI Safety Guardrails

**As a** user  
**I want** AI integrations to be safe by default  
**So that** I don't leak PII or exceed cost budgets

**Stories:**
- Implement PII scrubbing before AI API calls (reads `data-classification.yaml`)
- Implement response validation against declared schema
- Implement prompt injection detection
- Implement daily cost circuit breaker with graceful fallback
- Implement audit logging for all AI calls (metadata only, no payload)
- Implement response caching for repeated prompts
- Implement token usage tracking → Prometheus metrics
- Create AI cost dashboard (token usage, latency, cost/day, error rates)
- Implement prompt versioning system
- Write AI safety tests
- Document AI safety guardrails

**Acceptance Criteria:**
- PII is scrubbed before API calls
- Cost circuit breaker prevents runaway spend
- All AI calls are audited
- Response caching reduces costs
- AI cost dashboard shows usage
- Prompt injection attempts are blocked

---

### Epic 3.9: FinOps & Cost Management

**As a** platform engineer  
**I want** cost visibility and optimization  
**So that** cloud spend is controlled

**Stories:**
- Generate cost allocation tags (team, service, environment, cost-centre) in all resources
- Generate `cost-budget.yaml` with expected monthly spend
- Implement cost budget alerts (80% threshold)
- Generate weekly VPA recommendation workflow
- Implement automated right-sizing PR generation
- Implement AI spend tracking per service
- Create platform-level AI cost dashboard (aggregated across services)
- Implement per-team chargeback reporting
- Write FinOps tests
- Document cost management practices

**Acceptance Criteria:**
- All resources have cost allocation tags
- Cost budgets are declared
- Alerts trigger at 80% of budget
- VPA recommendations generate PRs weekly
- AI spend is tracked per service
- Chargeback reports are accurate

---

### Epic 3.10: Governance Workflows

**As a** platform engineer  
**I want** automated governance workflows  
**So that** platform quality is maintained

**Stories:**
- Implement plugin review process automation
- Implement standards deviation workflow (issue → approval → tracking)
- Implement service decommission checklist enforcement
- Implement CVE response SLA tracking
- Implement plugin deprecation process (90-day notice, migration guide)
- Create governance dashboard
- Implement automated compliance reporting
- Write governance tests
- Document governance processes

**Acceptance Criteria:**
- Plugin reviews follow defined process
- Deviations are tracked and reviewed quarterly
- CVE SLA is enforced
- Deprecations follow 90-day notice
- Governance dashboard shows compliance status

---

## Phase 4: Platform Maturity (Ongoing)

### Epic 4.1: IDP Integration

**As a** platform engineer  
**I want** Crux integrated with our IDP  
**So that** it's part of a cohesive platform

**Stories:**
- Create Backstage integration
- Implement service catalog sync
- Add portal UI (web interface)
- Create API for programmatic access
- Implement SSO integration
- Add usage analytics
- Create admin dashboard

**Acceptance Criteria:**
- Backstage shows Crux-generated services
- Portal UI works
- API is documented
- Analytics track usage

---

### Epic 4.2: Advanced Features

**As a** power user  
**I want** advanced features  
**So that** I can customize Crux for complex scenarios

**Stories:**
- Add custom template support
- Implement template inheritance
- Add plugin composition
- Create plugin SDK
- Add scripting support
- Implement webhooks
- Add event system

**Acceptance Criteria:**
- Custom templates work
- Plugin SDK is documented
- Webhooks trigger correctly

---

## Epic Prioritization

### Must Have (Phase 1 - Pilot)
- Epic 0.1: Foundation
- Epic 1.1: Tier 1 Standards Generation
- Epic 1.2: CLI Framework
- Epic 1.3: Prompt Engine & Decision Graph
- Epic 1.4: Template Engine
- Epic 1.5: Core Templates (Go + framework)
- Epic 1.6: Plugin System
- Epic 1.7: Pilot Plugins (9 essential)
- Epic 1.8: Configuration & Lockfile System
- Epic 1.9: Complete Skeleton Generation
- Epic 1.10: Pilot Testing

### Should Have (Phase 2 - MVP)
- Epic 2.1: TUI
- Epic 2.2: Multi-Language (4 languages)
- Epic 2.3: Additional Plugins (15 MVP plugins)
- Epic 2.4: Plugin Registry
- Epic 2.5: Validation
- Epic 2.6: Audit
- Epic 2.7: Documentation
- Epic 2.8: Shared Component Library
- Epic 2.9: Developer Inner Loop
- Epic 2.10: Adoption & Rollout

### Could Have (Phase 3 - Core)
- Epic 3.1: Upgrade System
- Epic 3.2: Resilience
- Epic 3.3: Zero-Trust
- Epic 3.4: Observability
- Epic 3.5: Compliance (Full Implementation)
- Epic 3.6: Performance
- Epic 3.7: Community
- Epic 3.8: AI Safety Guardrails
- Epic 3.9: FinOps & Cost Management
- Epic 3.10: Governance Workflows

### Won't Have (Phase 4 - Future)
- Epic 4.1: IDP Integration
- Epic 4.2: Advanced Features

---

## Success Metrics

### Phase 1 (Pilot)
- 2 pilot teams successfully generate services
- Generation time < 3 minutes
- 0 P1 bugs
- Positive pilot feedback

### Phase 2 (MVP)
- 10+ services generated across 3+ teams
- All 4 languages tested in production
- 80%+ of new services use Crux
- 0 critical CVEs unpatched > 48h

### Phase 3 (Core)
- `crux upgrade` used on 5+ services
- 0 services with critical compliance gaps
- 80%+ adoption of shared CI workflows
- 5+ community plugins published

### Phase 4 (Platform)
- Full IDP integration
- Self-service infrastructure provisioning
- 100% of new services use Crux
- Platform maturity Level 3+

---

## Notes

- Each epic should be estimated in story points during sprint planning
- Stories can be further broken down into tasks
- Dependencies between epics should be tracked
- Regular retrospectives should inform roadmap adjustments
