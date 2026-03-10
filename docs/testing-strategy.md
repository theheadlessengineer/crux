# Testing Strategy

## Testing Pyramid

```
        /\
       /  \
      / E2E \
     /--------\
    /Integration\
   /--------------\
  /   Unit Tests   \
 /------------------\
```

**Distribution:**
- Unit Tests: 70%
- Integration Tests: 20%
- End-to-End Tests: 10%

## Unit Testing

### Scope

Test individual functions and methods in isolation.

**What to test:**
- Business logic in domain layer
- Scoring algorithms
- Validation functions
- Data transformations
- Error handling

**What not to test:**
- Third-party libraries
- Standard library functions
- Trivial getters/setters

### Structure

```go
func TestFunctionName(t *testing.T) {
    // Arrange
    input := setupTestData()
    
    // Act
    result, err := FunctionName(input)
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### Table-Driven Tests

Use for testing multiple scenarios:

```go
func TestValidateServiceName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {
            name:    "valid name",
            input:   "payment-service",
            wantErr: false,
        },
        {
            name:    "empty name",
            input:   "",
            wantErr: true,
        },
        {
            name:    "too long",
            input:   strings.Repeat("a", 64),
            wantErr: true,
        },
        {
            name:    "invalid characters",
            input:   "Payment_Service",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateServiceName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Test Helpers

Extract common setup into helpers:

```go
func newTestHardware() *Hardware {
    return &Hardware{
        RAM:     32,
        VRAM:    16,
        Cores:   16,
        Backend: "CUDA",
    }
}

func newTestModel(params float64) *Model {
    return &Model{
        Name:    "test-model",
        Params:  params,
        Context: 8192,
    }
}
```

### Mocking

Define interfaces for dependencies:

```go
type HardwareDetector interface {
    DetectRAM() (uint64, error)
    DetectGPU() (*GPUInfo, error)
}

type mockHardwareDetector struct {
    ram     uint64
    ramErr  error
    gpu     *GPUInfo
    gpuErr  error
}

func (m *mockHardwareDetector) DetectRAM() (uint64, error) {
    return m.ram, m.ramErr
}

func (m *mockHardwareDetector) DetectGPU() (*GPUInfo, error) {
    return m.gpu, m.gpuErr
}

func TestSomething(t *testing.T) {
    mock := &mockHardwareDetector{
        ram: 32 * 1024 * 1024 * 1024,
        gpu: &GPUInfo{Name: "Test GPU", VRAM: 16},
    }
    
    // Use mock in test
}
```

### Test Coverage

**Minimum Requirements:**
- Overall: 80%
- Domain layer: 90%
- Critical paths: 100%

**Measure:**
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**CI Enforcement:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | \
  awk '{if ($1 < 80) exit 1}'
```

## Integration Testing

### Scope

Test interaction between multiple components.

**What to test:**
- Template rendering with real templates
- Plugin loading and execution
- File system operations
- Configuration loading
- Repository implementations

### Structure

```go
// +build integration

package integration

import (
    "os"
    "path/filepath"
    "testing"
)

func TestTemplateRendering(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    
    // Setup
    tmpDir := t.TempDir()
    
    // Test
    renderer := template.NewRenderer(tmpDir)
    err := renderer.Render("service.go.tmpl", data)
    
    // Verify
    if err != nil {
        t.Fatalf("render failed: %v", err)
    }
    
    content, err := os.ReadFile(filepath.Join(tmpDir, "service.go"))
    if err != nil {
        t.Fatalf("failed to read output: %v", err)
    }
    
    if !strings.Contains(string(content), "package main") {
        t.Error("expected package declaration")
    }
}
```

### Test Fixtures

Store test data in `testdata/` directories:

```
internal/
├── template/
│   ├── renderer.go
│   ├── renderer_test.go
│   └── testdata/
│       ├── simple.tmpl
│       ├── complex.tmpl
│       └── expected/
│           ├── simple.go
│           └── complex.go
```

Load fixtures:

```go
func loadFixture(t *testing.T, name string) []byte {
    t.Helper()
    data, err := os.ReadFile(filepath.Join("testdata", name))
    if err != nil {
        t.Fatalf("failed to load fixture %s: %v", name, err)
    }
    return data
}
```

### Temporary Directories

Use `t.TempDir()` for file operations:

```go
func TestFileOperations(t *testing.T) {
    tmpDir := t.TempDir() // Automatically cleaned up
    
    // Write files to tmpDir
    // Test operations
}
```

## End-to-End Testing

### Scope

Test complete workflows from CLI invocation to output.

**What to test:**
- `crux new` command creates valid project
- Generated project compiles
- Generated tests pass
- Docker build succeeds
- Kubernetes manifests are valid

### Structure

```go
// +build e2e

package e2e

import (
    "os/exec"
    "path/filepath"
    "testing"
)

func TestCruxNew(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping e2e test")
    }
    
    tmpDir := t.TempDir()
    
    // Run crux new
    cmd := exec.Command("crux", "new", "test-service",
        "--language", "go",
        "--framework", "gin",
        "--output", tmpDir,
    )
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("crux new failed: %v\n%s", err, output)
    }
    
    // Verify structure
    serviceDir := filepath.Join(tmpDir, "test-service")
    
    mustExist(t, filepath.Join(serviceDir, "go.mod"))
    mustExist(t, filepath.Join(serviceDir, "main.go"))
    mustExist(t, filepath.Join(serviceDir, "Dockerfile"))
    
    // Verify it compiles
    cmd = exec.Command("go", "build", "./...")
    cmd.Dir = serviceDir
    if output, err := cmd.CombinedOutput(); err != nil {
        t.Fatalf("go build failed: %v\n%s", err, output)
    }
    
    // Verify tests pass
    cmd = exec.Command("go", "test", "./...")
    cmd.Dir = serviceDir
    if output, err := cmd.CombinedOutput(); err != nil {
        t.Fatalf("go test failed: %v\n%s", err, output)
    }
}

func mustExist(t *testing.T, path string) {
    t.Helper()
    if _, err := os.Stat(path); os.IsNotExist(err) {
        t.Errorf("expected file to exist: %s", path)
    }
}
```

### Golden Files

Compare generated output against known-good files:

```go
func TestGeneratedOutput(t *testing.T) {
    generated := generateOutput()
    
    goldenPath := filepath.Join("testdata", "golden", "output.go")
    
    if *update {
        os.WriteFile(goldenPath, generated, 0644)
        return
    }
    
    golden, err := os.ReadFile(goldenPath)
    if err != nil {
        t.Fatalf("failed to read golden file: %v", err)
    }
    
    if !bytes.Equal(generated, golden) {
        t.Errorf("output differs from golden file")
        t.Logf("Run with -update to update golden files")
    }
}
```

Update golden files:
```bash
go test -update ./...
```

## Benchmark Testing

### Purpose

Measure performance and detect regressions.

### Structure

```go
func BenchmarkScoreModel(b *testing.B) {
    model := newTestModel(7)
    hardware := newTestHardware()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ScoreModel(model, hardware)
    }
}

func BenchmarkScoreAllModels(b *testing.B) {
    models := loadAllModels()
    hardware := newTestHardware()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, model := range models {
            ScoreModel(model, hardware)
        }
    }
}
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkScoreModel ./...

# With memory allocation stats
go test -bench=. -benchmem ./...

# Compare before/after
go test -bench=. ./... > old.txt
# Make changes
go test -bench=. ./... > new.txt
benchcmp old.txt new.txt
```

### Performance Targets

**Benchmarks must meet:**
- Startup time: < 100ms
- Model scoring: < 1ms per model
- Template rendering: < 100ms per template
- Plugin loading: < 50ms per plugin

## Test Organization

### Directory Structure

```
crux/
├── internal/
│   ├── domain/
│   │   ├── model/
│   │   │   ├── model.go
│   │   │   ├── model_test.go          # Unit tests
│   │   │   └── testdata/              # Test fixtures
│   │   │       └── models.json
│   │   └── scoring/
│   │       ├── scorer.go
│   │       └── scorer_test.go
│   └── infrastructure/
│       ├── template/
│       │   ├── renderer.go
│       │   ├── renderer_test.go       # Unit tests
│       │   └── renderer_integration_test.go  # Integration tests
│       └── detector/
│           ├── gpu.go
│           └── gpu_test.go
├── test/
│   ├── integration/                   # Integration test suite
│   │   ├── template_test.go
│   │   └── plugin_test.go
│   └── e2e/                          # End-to-end test suite
│       ├── crux_new_test.go
│       └── testdata/
│           └── golden/
└── scripts/
    └── test.sh                        # Test runner script
```

### Test Tags

Use build tags to organize tests:

```go
// +build unit
// Unit tests (default)

// +build integration
// Integration tests

// +build e2e
// End-to-end tests
```

Run specific test types:
```bash
go test ./...                    # Unit tests only
go test -tags=integration ./...  # Integration tests
go test -tags=e2e ./...          # E2E tests
```

## Test Utilities

### Assertions

Create helper functions for common assertions:

```go
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

func assertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func assertError(t *testing.T, err error) {
    t.Helper()
    if err == nil {
        t.Fatal("expected error, got nil")
    }
}
```

### Test Cleanup

Use `t.Cleanup()` for resource cleanup:

```go
func TestWithCleanup(t *testing.T) {
    file, err := os.CreateTemp("", "test")
    assertNoError(t, err)
    
    t.Cleanup(func() {
        os.Remove(file.Name())
    })
    
    // Test with file
}
```

## Continuous Testing

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running tests..."
go test ./...
if [ $? -ne 0 ]; then
    echo "Tests failed. Commit aborted."
    exit 1
fi
```

### CI Pipeline

```yaml
name: Test

on: [push, pull_request]

jobs:
  unit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Unit Tests
        run: go test -race -coverprofile=coverage.out ./...
      
      - name: Coverage Check
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$coverage < 80" | bc -l) )); then
            echo "Coverage $coverage% is below 80%"
            exit 1
          fi
  
  integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Integration Tests
        run: go test -tags=integration ./...
  
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build
        run: go build -o crux ./cmd/crux
      
      - name: E2E Tests
        run: go test -tags=e2e ./...
```

## Test Documentation

### Test Comments

Document complex test scenarios:

```go
// TestScoreModel_MoEArchitecture verifies that Mixture-of-Experts models
// are scored based on active parameters rather than total parameters.
// This is critical for accurate memory estimation.
func TestScoreModel_MoEArchitecture(t *testing.T) {
    // Setup: Mixtral 8x7B with 2 active experts
    model := &Model{
        Name:              "mixtral-8x7b",
        TotalParams:       46.7,
        ActiveParams:      12.9,
        NumExperts:        8,
        NumExpertsPerTok:  2,
    }
    
    // Test scoring uses active params
    score := ScoreModel(model, hardware)
    
    // Verify memory calculation based on active params
    expectedMemory := 12.9 * bytesPerParam
    assertEqual(t, model.EstimatedMemory, expectedMemory)
}
```

## References

- Go Testing Package: https://pkg.go.dev/testing
- Table-Driven Tests: https://go.dev/wiki/TableDrivenTests
- Advanced Testing: https://go.dev/doc/tutorial/add-a-test
