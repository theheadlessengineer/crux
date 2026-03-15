# ADR-001: Selection of Go as Implementation Language for Microservice Skeleton Generator

## Status

Accepted

## Date

2026-03-09

## Context

We are building an internal CLI tool (working name: `crux`) that generates production-ready, enterprise-compliant microservice skeletons through an interactive terminal interface. The tool must support:

- Plugin-based architecture for extensibility
- Interactive Terminal User Interface (TUI) with themes, filtering, and search
- Template rendering for multiple languages and frameworks
- Cross-platform distribution (macOS, Linux, Windows)
- Embedded resources (templates, schemas, default configurations)
- Fast execution and minimal resource footprint
- Single-command installation with zero runtime dependencies
- Concurrent plugin execution
- JSON/YAML processing for configuration and lockfiles
- Integration with external tools (Git, Docker, Terraform, CI/CD systems)

The tool is the foundation of our Internal Developer Platform (IDP) and will be used by all engineering teams. Distribution friction, maintenance burden, and execution performance are critical success factors.

## Decision Drivers

1. **Distribution Model**: Single binary with no runtime dependencies required
2. **Cross-Platform Support**: Must run on macOS, Linux, and Windows without modification
3. **TUI Capabilities**: Rich terminal interface with real-time updates, themes, and interactive controls
4. **Template Engine**: Built-in or well-supported templating for code generation
5. **Plugin Architecture**: Ability to load and execute plugins safely
6. **Performance**: Fast startup time, low memory footprint, responsive UI
7. **Compilation Speed**: Fast build times for development iteration
8. **Ecosystem Maturity**: Production-ready libraries for CLI, TUI, file I/O, and system integration
9. **Team Familiarity**: Learning curve and existing expertise
10. **Enterprise Adoption**: Track record in similar tooling domains
11. **Long-Term Maintenance**: Language stability, backward compatibility, community support

## Options Considered

### Option 1: Go

**Strengths:**
- Compiles to single static binary with no runtime dependencies
- Native cross-compilation via GOOS/GOARCH environment variables
- Mature TUI ecosystem: bubbletea (Elm architecture), lipgloss (styling), bubbles (components)
- Built-in template engine (text/template, html/template)
- Embed package for bundling resources at compile time
- Goroutines for concurrent plugin execution
- Fast compilation (sub-second incremental builds)
- Strong standard library for file I/O, JSON/YAML, subprocess management
- Proven in infrastructure tooling: Kubernetes, Terraform, Docker, Vault, Hugo
- Excellent backward compatibility guarantees (Go 1 compatibility promise)
- Simple deployment: single binary, no package manager required

**Weaknesses:**
- More verbose syntax compared to Python or Ruby
- Error handling via explicit error returns (no exceptions)
- Smaller template library ecosystem compared to Node.js

### Option 2: Rust

**Strengths:**
- Compiles to single static binary
- Fastest execution speed and lowest memory footprint
- Memory safety guarantees at compile time
- Excellent TUI libraries: ratatui, crossterm (proven in llmfit)
- Strong type system prevents entire classes of runtime errors
- Growing ecosystem with high-quality crates

**Weaknesses:**
- Steeper learning curve (ownership, lifetimes, borrow checker)
- Slower compilation times (minutes for large projects)
- Smaller team familiarity
- More complex error handling (Result types, unwrap/expect patterns)
- Less mature template engine ecosystem

### Option 3: Node.js with TypeScript

**Strengths:**
- Massive ecosystem: npm registry with extensive CLI and template libraries
- Native JSON/YAML handling
- Rich template engines: Handlebars, EJS, Nunjucks
- Team likely familiar with JavaScript/TypeScript
- Fast prototyping and hot reload during development
- Excellent TUI libraries: ink (React-based), blessed

**Weaknesses:**
- Requires Node.js runtime installation (not single binary)
- Larger distribution size (node_modules, bundled runtime)
- Slower execution compared to compiled languages
- Dependency management complexity (package-lock.json, security vulnerabilities)
- Version compatibility issues across Node.js releases

### Option 4: Python

**Strengths:**
- Highly readable and maintainable syntax
- Extensive standard library
- Mature CLI frameworks: click, typer, rich, textual
- Powerful template engines: Jinja2 (industry standard)
- Large community and extensive documentation

**Weaknesses:**
- Requires Python runtime installation
- Slower execution speed
- Distribution complexity (PyInstaller, cx_Freeze produce large binaries)
- Dependency management challenges (virtualenv, pip, version conflicts)
- GIL limits true concurrency for CPU-bound tasks

### Option 5: Deno with TypeScript

**Strengths:**
- Single binary compilation via deno compile
- Modern JavaScript/TypeScript runtime
- Secure by default (explicit permissions)
- No node_modules, built-in dependency management
- Built-in formatter, linter, test runner

**Weaknesses:**
- Smaller ecosystem compared to Node.js
- Less mature tooling and libraries
- Limited team familiarity
- Fewer production deployments in enterprise environments

## Decision

We will implement the microservice skeleton generator CLI in **Go**.

## Rationale

### Critical Success Factors

**1. Distribution Model**

Go produces a single static binary with zero runtime dependencies. Users can download one file and execute it immediately. This eliminates:
- "Install Node.js/Python first" friction
- Version compatibility issues
- Dependency installation failures
- Corporate firewall/proxy complications with package managers

For an enterprise IDP tool used by all engineering teams, installation friction is a primary adoption barrier. Go's distribution model is optimal.

**2. Cross-Platform Support**

Go's cross-compilation is trivial:
```
GOOS=linux GOARCH=amd64 go build
GOOS=darwin GOARCH=arm64 go build
GOOS=windows GOARCH=amd64 go build
```

No additional tooling, Docker containers, or platform-specific build machines required. CI/CD can produce binaries for all platforms in a single pipeline run.

**3. TUI Capabilities**

The bubbletea framework provides production-grade TUI capabilities with:
- Elm architecture (predictable state management)
- Component composition via bubbles library
- Advanced styling via lipgloss
- Proven in production tools: glow, soft-serve, vhs, charm

This matches the reference implementations (mactop, llmfit) and provides the interactive experience required by the specification.

**4. Template Engine**

Go's built-in text/template package provides:
- No external dependencies
- Familiar syntax (similar to Jinja2, Handlebars)
- Safe execution (no arbitrary code execution)
- Embedded at compile time via embed package

Templates, schemas, and default configurations can be bundled into the binary, eliminating external file dependencies.

**5. Plugin Architecture**

Go supports multiple plugin approaches:
- Process-based plugins (subprocess execution, language-agnostic)
- Go plugin package (shared libraries, Go-only)
- gRPC-based plugins (HashiCorp model, used by Terraform)

Process-based plugins provide maximum flexibility and safety, allowing plugins in any language while maintaining isolation.

**6. Performance**

Go provides:
- Fast startup time (milliseconds)
- Low memory footprint (10-50 MB typical)
- Efficient concurrency via goroutines (thousands of concurrent operations)
- No garbage collection pauses during interactive UI updates

For a CLI tool used hundreds of times per day by developers, startup time and responsiveness are critical.

**7. Compilation Speed**

Go's compilation speed enables rapid development iteration:
- Initial build: 1-5 seconds
- Incremental rebuild: sub-second
- No separate compilation and linking phases

This is critical for developer productivity during the multi-phase delivery plan.

**8. Ecosystem Maturity**

Go has production-ready libraries for all required functionality:
- CLI: cobra (used by Kubernetes, Hugo), urfave/cli
- TUI: bubbletea, lipgloss, bubbles
- Templates: text/template (stdlib)
- YAML: gopkg.in/yaml.v3
- JSON: encoding/json (stdlib)
- File operations: os, io, embed (stdlib)
- Subprocess: os/exec (stdlib)
- HTTP client: net/http (stdlib)

All libraries are actively maintained with large user bases.

**9. Enterprise Adoption**

Go is the de facto standard for infrastructure tooling:
- Kubernetes (container orchestration)
- Terraform (infrastructure as code)
- Docker (containerization)
- Vault (secrets management)
- Hugo (static site generation)
- Prometheus (monitoring)

Platform engineering teams are familiar with Go. The language choice aligns with existing tooling in the ecosystem.

**10. Long-Term Maintenance**

Go provides exceptional stability:
- Go 1 compatibility promise: code written for Go 1.0 (2012) still compiles
- Minimal breaking changes between versions
- Long-term support for releases
- Clear deprecation policies
- Strong backward compatibility culture

For a foundational IDP tool with multi-year lifespan, language stability is critical.

### Trade-offs Accepted

**Verbosity**: Go's explicit error handling and lack of generics (pre-1.18) results in more verbose code compared to Python or Ruby. We accept this trade-off for compile-time safety and explicit control flow.

**Template Ecosystem**: Go's template library ecosystem is smaller than Node.js. However, the built-in text/template package is sufficient for our needs, and custom template functions can be added as required.

**Learning Curve**: For teams unfamiliar with Go, there is a learning curve around goroutines, channels, and error handling patterns. However, Go's simplicity (25 keywords, minimal language features) makes it one of the easier compiled languages to learn.

## Consequences

### Positive

- Zero installation friction for end users (single binary download)
- Fast execution and responsive UI
- Predictable cross-platform behavior
- Simple CI/CD pipeline (no runtime dependencies to manage)
- Strong type safety catches errors at compile time
- Efficient resource usage (low memory, fast startup)
- Alignment with existing infrastructure tooling ecosystem
- Long-term language stability and backward compatibility
- Concurrent plugin execution via goroutines
- Embedded resources eliminate external file dependencies

### Negative

- Team members unfamiliar with Go will require training
- More verbose code compared to dynamic languages
- Smaller template library ecosystem (mitigated by built-in text/template)
- Plugin development in Go requires compilation step (mitigated by supporting process-based plugins in any language)

### Neutral

- Go's opinionated formatting (gofmt) enforces consistent code style
- Standard library focus means fewer third-party dependencies to manage
- Explicit error handling increases code volume but improves reliability

## Implementation Notes

### Required Libraries

- **CLI Framework**: cobra (github.com/spf13/cobra)
- **TUI Framework**: bubbletea (github.com/charmbracelet/bubbletea)
- **TUI Styling**: lipgloss (github.com/charmbracelet/lipgloss)
- **TUI Components**: bubbles (github.com/charmbracelet/bubbles)
- **YAML Processing**: gopkg.in/yaml.v3
- **Template Engine**: text/template (stdlib)
- **Embedded Resources**: embed (stdlib)

### Development Environment

- Go 1.26 or later (for improved type inference and standard library enhancements)
- gofmt for code formatting
- golangci-lint for static analysis
- go test for unit testing
- go build for compilation

### Distribution Strategy

- GitHub Releases with pre-built binaries for macOS (amd64, arm64), Linux (amd64, arm64), Windows (amd64)
- Homebrew tap for macOS/Linux installation
- Installation script (curl | sh pattern) for automated download
- Internal artifact repository for enterprise distribution

## References

- Go Language Specification: https://go.dev/ref/spec
- Bubbletea TUI Framework: https://github.com/charmbracelet/bubbletea
- Cobra CLI Framework: https://github.com/spf13/cobra
- mactop (reference implementation): https://github.com/metaspartan/mactop
- llmfit (reference implementation): https://github.com/AlexsJones/llmfit
- Microservice Skeleton CLI Specification: microservice-skeleton-cli-scope-v3.md

## Approval

- Platform Engineering Lead: [Pending]
- Architecture Guild: [Pending]
- Security Team: [Pending]

## Revision History

| Date | Version | Author | Changes |
|------|---------|--------|---------|
| 2026-03-09 | 1.0 | Platform Team | Initial decision |
