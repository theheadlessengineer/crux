# Crux

Standards built in. Not bolted on.

---

## What is Crux

Crux is an internal CLI tool for generating production-ready microservice skeletons. It is the starting point for every service built in this organisation.

A single command generates a fully structured, runnable service with company standards, security configuration, resilience patterns, observability wiring, infrastructure as code, and CI/CD pipelines already in place. Teams write business logic. Crux handles everything else.

Crux is built on a plugin architecture. Every integration beyond the core — databases, caches, message brokers, cloud providers, AI tools, observability backends — is a self-contained, versioned plugin. The core enforces the standard. Plugins extend it.

```bash
crux new payment-service
```

---

## Vision

A world where every engineer in the organisation ships production-grade services with confidence — where the distance between an idea and a running, secure, observable, compliant service is measured in minutes, not weeks.

## Mission

Crux gives every engineering team a single, extensible starting point — embedding company standards, security, resilience, and observability directly into the foundation of every service, so teams can focus entirely on the problems only they can solve.

---

## Current Status

**Phase 0 (Foundation) - In Progress**

The project is currently in the foundation phase. The core infrastructure is in place:
- ✅ Go module initialized
- ✅ CI/CD pipeline configured
- ✅ Linting and formatting standards
- ✅ Pre-commit hooks
- ✅ Test infrastructure
- ✅ Development workflow

**CLI commands will be implemented in Phase 1.** The binary currently compiles but does not yet execute any commands.

---

## Prerequisites

- **Go 1.26 or later** - [Install Go](https://go.dev/doc/install)
- **golangci-lint** - [Installation instructions](https://golangci-lint.run/usage/install/)
- **Git** - Version control

### Verify Prerequisites

```bash
go version        # Should show 1.26 or later
golangci-lint version
git --version
```

---

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd terminal-ms
```

### 2. Install Pre-Commit Hooks

```bash
make hooks
```

This installs hooks that automatically check:
- Code formatting with `gofmt`
- Linting with `golangci-lint`
- Test passage with `go test`

### 3. Verify Your Setup

```bash
make dev
```

This runs the full development workflow: format, lint, and test.

### 4. Build the Binary

```bash
make build
```

The binary will be created at `./bin/crux`.

---

## Development Workflow

### Common Commands

```bash
make help              # Show all available targets
make build             # Compile binary to ./bin/crux
make test              # Run all tests with race detector
make test-integration  # Run integration tests
make coverage          # Run tests and view coverage report
make lint              # Run golangci-lint
make fmt               # Format all Go files
make vet               # Run go vet
make clean             # Remove build artifacts
make dev               # Run format, lint, and test (recommended)
```

### Before Committing

Always run the full development workflow:

```bash
make dev
```

This ensures your code is formatted, passes linting, and all tests pass.

---

## Project Structure

```
crux/
├── cmd/
│   └── crux/              # Main application entrypoint
├── internal/
│   ├── app/
│   │   ├── commands/      # CLI command implementations
│   │   └── config/        # Application configuration
│   ├── domain/
│   │   ├── model/         # Core domain models
│   │   ├── hardware/      # Hardware detection domain
│   │   ├── plugin/        # Plugin domain logic
│   │   └── scoring/       # Model scoring domain
│   ├── infrastructure/
│   │   ├── detector/      # System detection implementations
│   │   ├── repository/    # Data persistence
│   │   ├── template/      # Template rendering
│   │   ├── executor/      # Command execution
│   │   └── filesystem/    # File operations
│   ├── presentation/
│   │   ├── cli/           # CLI interface
│   │   └── tui/           # Terminal UI
│   └── testutil/          # Shared test utilities
├── test/
│   └── integration/       # Integration tests
├── pkg/                   # Public packages (if any)
├── data/
│   ├── templates/         # Service templates
│   └── schemas/           # Configuration schemas
├── docs/                  # Documentation
├── scripts/               # Build and utility scripts
└── .github/workflows/     # CI/CD pipelines
```

**Architecture:** Hexagonal (Ports & Adapters) - See [docs/architecture-principles.md](docs/architecture-principles.md)

---

## Testing

### Run Tests

```bash
# Run all tests
make test

# Run integration tests only
make test-integration

# Generate coverage report
make coverage
```

### Writing Tests

All tests must follow the table-driven pattern. See [docs/TESTING-GUIDE.md](docs/TESTING-GUIDE.md) for detailed guidelines.

**Coverage targets:**
- Unit tests: 80% minimum
- Critical paths: 100%
- Integration tests: All major workflows

---

## Contributing

We welcome contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development workflow
- Commit message conventions
- Pull request process
- Code review guidelines
- Testing requirements

---

## Documentation

- [Architecture Principles](docs/architecture-principles.md) - Design principles and patterns
- [Development Standards](docs/development-standards.md) - Coding standards and conventions
- [Testing Guide](docs/TESTING-GUIDE.md) - How to write and run tests
- [Testing Strategy](docs/testing-strategy.md) - Overall testing approach
- [Roadmap](docs/ROADMAP.md) - Project roadmap and milestones
- [Plugin Development Guide](docs/plugin-development-guide.md) - How to create plugins (Phase 1+)

---

## Maintainers

**Platform Engineering Team**

For questions or support:
- Open an issue for bugs or feature requests
- Start a discussion for questions or ideas
- Reach out to the platform team for guidance

---

## License

Internal use only. All rights reserved.