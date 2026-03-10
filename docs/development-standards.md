# Development Standards

## Code Style

### Formatting

All code must be formatted with `gofmt` before commit.

```bash
gofmt -w .
```

### Naming Conventions

**Packages:**
- Lowercase, single word
- No underscores or camelCase
- Examples: `hardware`, `template`, `scoring`

**Files:**
- Lowercase with underscores
- Group related functionality: `gpu_detector.go`, `cpu_detector.go`
- Test files: `gpu_detector_test.go`

**Variables:**
- camelCase for local: `modelCount`, `gpuInfo`
- PascalCase for exported: `ModelRepository`, `GPUDetector`
- Short names in small scopes: `i`, `err`, `ctx`

**Constants:**
```go
const (
    defaultTimeout = 30 * time.Second
    maxRetries     = 3
    bufferSize     = 4096
)
```

**Interfaces:**
- Noun or adjective: `Reader`, `Detector`, `Runnable`
- Single-method interfaces end in -er: `Renderer`, `Validator`

### Line Length

- Soft limit: 100 characters
- Hard limit: 120 characters

### Import Organization

```go
import (
    // Standard library
    "context"
    "fmt"
    "os"
    
    // External dependencies
    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"
    
    // Internal packages
    "github.com/company/crux/internal/domain/model"
)
```

### Function Length

- Target: 20-30 lines
- Maximum: 50 lines
- Extract helper functions if longer

### Function Parameters

- Maximum 4 parameters
- Use structs for more parameters

```go
type ServiceConfig struct {
    Name      string
    Version   string
    Language  string
    Framework string
}

func CreateService(config ServiceConfig) error
```

### Comments

**Package comments:**
```go
// Package hardware provides cross-platform hardware detection.
package hardware
```

**Function comments:**
```go
// DetectGPU detects available GPU hardware.
// Returns nil if no GPU is found.
func DetectGPU() (*GPUInfo, error)
```

**Inline comments:**
```go
// TODO(username): Add support for AMD GPUs
// FIXME(username): Race condition in concurrent access
// NOTE: This assumes NVIDIA driver is installed
```

### Error Messages

- Lowercase, no punctuation
- Include context
- Use fmt.Errorf with %w for wrapping

```go
return fmt.Errorf("failed to detect GPU: %w", err)
```

## Git Workflow

### Branch Naming

```
feature/add-postgresql-plugin
bugfix/fix-template-rendering
hotfix/security-vulnerability
refactor/improve-scoring-algorithm
docs/update-architecture-guide
```

### Commit Messages

Follow Conventional Commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting, no code change
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(plugin): add PostgreSQL plugin with read replica support

Implements PostgreSQL plugin with connection pooling,
read replica configuration, and migration tool selection.

Closes #123
```

### Pull Request Process

1. Create feature branch from `main`
2. Implement changes with tests
3. Run full test suite locally
4. Update documentation
5. Create PR with description
6. Address review comments
7. Squash commits before merge

### Code Review Guidelines

**Reviewers must check:**
- Correctness: Does it work as intended?
- Tests: Are there adequate tests?
- Design: Does it follow architecture principles?
- Readability: Is it easy to understand?
- Performance: Are there obvious bottlenecks?
- Security: Are there security concerns?

**Review comments should:**
- Be constructive and specific
- Explain the reasoning
- Suggest alternatives
- Use prefixes: `nit:`, `question:`, `suggestion:`, `blocker:`

## Testing Standards

### Test Organization

```
internal/
├── domain/
│   ├── model/
│   │   ├── model.go
│   │   └── model_test.go
```

### Test Naming

```go
func TestDetectGPU(t *testing.T)
func TestDetectGPU_NoGPUFound(t *testing.T)
func TestDetectGPU_MultipleGPUs(t *testing.T)
```

### Table-Driven Tests

```go
func TestScoreModel(t *testing.T) {
    tests := []struct {
        name     string
        model    *Model
        hardware *Hardware
        want     float64
    }{
        {
            name: "perfect fit",
            model: &Model{Params: 7, Context: 8192},
            hardware: &Hardware{VRAM: 16},
            want: 95.0,
        },
        {
            name: "too large",
            model: &Model{Params: 70, Context: 32768},
            hardware: &Hardware{VRAM: 8},
            want: 0.0,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := ScoreModel(tt.model, tt.hardware)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Coverage

**Targets:**
- Unit tests: 80% coverage minimum
- Critical paths: 100% coverage
- Integration tests: All major workflows

**Measure coverage:**
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Mocking

Use interfaces for dependencies:

```go
type GPUDetector interface {
    Detect() (*GPUInfo, error)
}

type MockGPUDetector struct {
    DetectFunc func() (*GPUInfo, error)
}

func (m *MockGPUDetector) Detect() (*GPUInfo, error) {
    if m.DetectFunc != nil {
        return m.DetectFunc()
    }
    return nil, nil
}
```

### Integration Tests

```go
// +build integration

func TestFullWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Test full workflow
}
```

Run with:
```bash
go test -tags=integration ./...
```

### Benchmark Tests

```go
func BenchmarkScoreModel(b *testing.B) {
    model := &Model{Params: 7, Context: 8192}
    hardware := &Hardware{VRAM: 16}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ScoreModel(model, hardware)
    }
}
```

## Continuous Integration

### Pre-commit Hooks

```bash
#!/bin/bash
# .git/hooks/pre-commit

if [ -n "$(gofmt -l .)" ]; then
    echo "Code is not formatted. Run: gofmt -w ."
    exit 1
fi

golangci-lint run
if [ $? -ne 0 ]; then
    exit 1
fi

go test ./...
if [ $? -ne 0 ]; then
    exit 1
fi
```

### CI Pipeline

```yaml
name: CI
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Format check
        run: test -z "$(gofmt -l .)"
      
      - name: Lint
        uses: golangci/golangci-lint-action@v3
      
      - name: Test
        run: go test -race -coverprofile=coverage.out ./...
      
      - name: Build
        run: go build -v ./...
```

## Performance Standards

### Benchmarks

All performance-critical code must have benchmarks.

**Targets:**
- Startup time: < 100ms
- Model scoring: < 1ms per model
- Template rendering: < 100ms for full skeleton
- TUI frame rate: 60 FPS minimum

### Profiling

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## Security Standards

### Input Validation

```go
func ValidateServiceName(name string) error {
    if len(name) == 0 {
        return fmt.Errorf("service name cannot be empty")
    }
    
    if len(name) > 63 {
        return fmt.Errorf("service name too long")
    }
    
    matched, _ := regexp.MatchString(`^[a-z][a-z0-9-]*$`, name)
    if !matched {
        return fmt.Errorf("invalid service name format")
    }
    
    return nil
}
```

### Subprocess Execution

```go
// Never use shell
cmd := exec.Command("git", "clone", url) // Good
cmd := exec.Command("sh", "-c", "git clone "+url) // Bad

// Set timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
cmd := exec.CommandContext(ctx, "git", "clone", url)
```

## Release Process

### Version Bumping

1. Update version in `version.go`
2. Update CHANGELOG.md
3. Create git tag: `git tag -a v1.2.0 -m "Release v1.2.0"`
4. Push tag: `git push origin v1.2.0`

### Semantic Versioning

- MAJOR: Breaking changes
- MINOR: New features, backward compatible
- PATCH: Bug fixes
