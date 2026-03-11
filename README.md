# Crux

Standards built in. Not bolted on.

---

## What is Crux

Crux is an internal CLI tool for generating production-ready microservice
skeletons. It is the starting point for every service built in this
organisation.

A single command generates a fully structured, runnable service with
company standards, security configuration, resilience patterns,
observability wiring, infrastructure as code, and CI/CD pipelines already
in place. Teams write business logic. Crux handles everything else.

Crux is built on a plugin architecture. Every integration beyond the
core — databases, caches, message brokers, cloud providers, AI tools,
observability backends — is a self-contained, versioned plugin. The core
enforces the standard. Plugins extend it.

```
crux new payment-service
```

---

## Vision

A world where every engineer in the organisation ships production-grade
services with confidence — where the distance between an idea and a
running, secure, observable, compliant service is measured in minutes,
not weeks.

## Mission

Crux gives every engineering team a single, extensible starting point —
embedding company standards, security, resilience, and observability
directly into the foundation of every service, so teams can focus entirely
on the problems only they can solve.

---

## Development Setup

### Prerequisites

- Go 1.26 or later
- `golangci-lint` installed and available in PATH
- Git

### Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd terminal-ms
   ```

2. Install pre-commit hooks:
   ```bash
   make hooks
   ```

3. Verify your setup:
   ```bash
   make dev
   ```

### Pre-Commit Hooks

Pre-commit hooks enforce code quality standards before commits reach the remote repository. The hooks check:

- **Code formatting** with `gofmt`
- **Linting** with `golangci-lint`
- **Test passage** with `go test`

After cloning, run `make hooks` to install the hooks locally. The hooks will automatically run on every commit and block the commit if any check fails.

### Available Make Targets

- `make hooks` — Install pre-commit hooks
- `make fmt` — Format all Go files
- `make lint` — Run golangci-lint
- `make vet` — Run go vet
- `make test` — Run all tests with race detector
- `make coverage` — Run tests and open coverage report
- `make build` — Compile binary to ./bin/crux
- `make clean` — Remove build artifacts
- `make dev` — Run format, lint, and test (development workflow)
- `make help` — Show all available targets