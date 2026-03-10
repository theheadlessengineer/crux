# Architecture Principles

## Core Principles

### 1. Single Responsibility Principle

Each package, struct, and function must have one clear responsibility.

**Package Responsibilities:**
- `hardware`: System detection only
- `models`: Model database and scoring only
- `tui`: User interface rendering only
- `template`: Template rendering only

### 2. Dependency Inversion

High-level modules must not depend on low-level modules. Both depend on abstractions.

```go
// Define interfaces in high-level packages
type HardwareDetector interface {
    DetectRAM() (uint64, error)
    DetectGPU() (*GPUInfo, error)
}

// Implement in low-level packages
type SystemDetector struct {}
func (s *SystemDetector) DetectRAM() (uint64, error) { ... }
```

### 3. Explicit Over Implicit

All behavior must be explicit. No hidden side effects, no global state.

**Rules:**
- No global variables except constants
- No init() functions that modify state
- All dependencies passed via constructors or function parameters
- All errors returned, never swallowed

### 4. Fail Fast

Detect errors at startup, not during execution.

**Implementation:**
- Validate configuration on load
- Check required binaries exist (git, docker) at startup
- Verify template syntax during initialization
- Fail with clear error messages

### 5. Plugin Isolation

Plugins must not affect core stability.

**Rules:**
- Plugins run in separate processes
- Plugin failures do not crash core
- Plugin communication via stdin/stdout or gRPC
- Timeout enforcement on all plugin operations

## Architectural Patterns

### Hexagonal Architecture (Ports and Adapters)

```
Core Domain (Business Logic)
    ↕ Ports (Interfaces)
Adapters (Implementations)
    ↕
External Systems (File, Network, Process)
```

**Implementation:**
- Core domain has no external dependencies
- Ports defined as interfaces in domain layer
- Adapters implement ports in infrastructure layer
- Dependency injection at application startup

### Command Pattern

Each CLI command is a separate struct implementing a common interface.

```go
type Command interface {
    Execute(ctx context.Context, args []string) error
    Validate() error
}

type NewCommand struct {
    config *Config
    detector HardwareDetector
    renderer TemplateRenderer
}
```

### Repository Pattern

Abstract data access behind interfaces.

```go
type ModelRepository interface {
    List() ([]*Model, error)
    FindByName(name string) (*Model, error)
    Search(query string) ([]*Model, error)
}
```

### Factory Pattern

Create complex objects through factories.

```go
type PluginFactory interface {
    Create(manifest *Manifest) (Plugin, error)
}
```

## Package Organization

```
crux/
├── cmd/
│   └── crux/              # Main entry point
│       └── main.go
├── internal/               # Private application code
│   ├── app/                # Application layer
│   │   ├── commands/       # CLI commands
│   │   └── config/         # Configuration management
│   ├── domain/             # Domain layer (core business logic)
│   │   ├── model/          # Model entities
│   │   ├── hardware/       # Hardware entities
│   │   ├── plugin/         # Plugin entities
│   │   └── scoring/        # Scoring logic
│   ├── infrastructure/     # Infrastructure layer
│   │   ├── detector/       # Hardware detection
│   │   ├── repository/     # Data access
│   │   ├── template/       # Template rendering
│   │   ├── executor/       # Process execution
│   │   └── filesystem/     # File operations
│   └── presentation/       # Presentation layer
│       ├── tui/            # Terminal UI
│       └── cli/            # CLI output
├── pkg/                    # Public libraries (if any)
└── data/                   # Embedded resources
    ├── templates/
    └── schemas/
```

## Error Handling Strategy

### Error Types

```go
// Domain errors (recoverable)
type ValidationError struct {
    Field   string
    Message string
}

// Infrastructure errors (may be transient)
type NetworkError struct {
    URL string
    Err error
}

// Fatal errors (unrecoverable)
type FatalError struct {
    Message string
    Err     error
}
```

### Error Wrapping

```go
func processFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", path, err)
    }
    return nil
}
```

## Concurrency Patterns

### Worker Pool

```go
type WorkerPool struct {
    workers   int
    tasks     chan func()
    wg        sync.WaitGroup
}
```

### Context Propagation

```go
func Execute(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    result, err := fetchData(ctx)
    if err != nil {
        return err
    }
    
    return processData(ctx, result)
}
```

## Security Principles

### Input Validation

- Validate all user input at boundaries
- Use allowlists, not denylists
- Sanitize file paths (prevent directory traversal)
- Validate plugin manifests before loading

### Subprocess Execution

```go
// Never use shell execution
cmd := exec.Command("git", "clone", url) // Good
cmd := exec.Command("sh", "-c", "git clone " + url) // Bad
```

### Secret Handling

- Never log secrets
- Never include secrets in error messages
- Use environment variables or secure vaults
- Clear sensitive data from memory after use

## Documentation Standards

### Package Documentation

```go
// Package hardware provides system hardware detection capabilities.
package hardware
```

### Function Documentation

```go
// DetectGPU detects available GPU hardware and returns GPU information.
// Returns nil if no GPU is detected.
func DetectGPU() (*GPUInfo, error)
```

## Versioning Strategy

### Semantic Versioning

- MAJOR: Breaking changes to CLI interface or plugin API
- MINOR: New features, backward compatible
- PATCH: Bug fixes, no new features

### API Stability

- Internal packages can change freely
- Public interfaces (plugin API) follow semver strictly
- Deprecation warnings for 2 minor versions before removal
