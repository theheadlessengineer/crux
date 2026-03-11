# Testing Guide

**Last Updated:** March 11, 2026  
**Status:** Active

---

## Overview

This guide provides practical instructions for writing and running tests in the Crux project. All engineers must follow these patterns to ensure consistency and maintainability.

---

## Quick Start

```bash
# Run all tests
make test

# Run integration tests only
make test-integration

# Generate and view coverage report
make coverage

# Run full development workflow (format, lint, test)
make dev
```

---

## Test Organization

Tests live alongside the code they test:

```
internal/
├── domain/
│   ├── model/
│   │   ├── validation.go
│   │   └── validation_test.go      # Unit tests here
│   └── ...
└── testutil/
    └── mocks.go                     # Shared test utilities

test/
└── integration/
    └── service_generation_test.go   # Integration tests here
```

---

## Writing Unit Tests

### Test Naming Convention

```go
// Main test function
func TestValidateServiceName(t *testing.T)

// Specific scenario tests
func TestValidateServiceName_EmptyName(t *testing.T)
func TestValidateServiceName_NameTooLong(t *testing.T)
```

### Table-Driven Test Pattern

**Always use table-driven tests for multiple scenarios:**

```go
func TestValidateServiceName(t *testing.T) {
    tests := []struct {
        name    string        // Test case description
        input   string        // Input to test
        wantErr bool          // Expect error?
        errMsg  string        // Expected error message substring
    }{
        {
            name:    "valid service name",
            input:   "payment-service",
            wantErr: false,
        },
        {
            name:    "empty name",
            input:   "",
            wantErr: true,
            errMsg:  "service name cannot be empty",
        },
        {
            name:    "name too long",
            input:   "this-is-a-very-long-service-name-that-exceeds-maximum",
            wantErr: true,
            errMsg:  "service name too long",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateServiceName(tt.input)

            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

**Key points:**
- Use descriptive test case names
- Test both success and failure paths
- Use `t.Run()` for subtests
- Use `testify/assert` for assertions

---

## Writing Integration Tests

### Build Tag

**All integration tests must use the build tag:**

```go
//go:build integration

package integration

import "testing"

func TestServiceGeneration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    
    // Test implementation here
}
```

### Running Integration Tests

```bash
# Run only integration tests
go test -tags=integration ./...

# Or use the Makefile target
make test-integration

# Skip integration tests in short mode
go test -short ./...
```

---

## Mocking

### Interface-Based Mocks

**Use interface-based mocks, not mocking frameworks:**

```go
// Define the interface
type Executor interface {
    Execute(cmd string, args ...string) (string, error)
}

// Create a mock in testutil package
type MockExecutor struct {
    ExecuteFunc func(cmd string, args ...string) (string, error)
}

func (m *MockExecutor) Execute(cmd string, args ...string) (string, error) {
    if m.ExecuteFunc != nil {
        return m.ExecuteFunc(cmd, args...)
    }
    return "", nil
}
```

### Using Mocks in Tests

```go
func TestCommandExecution(t *testing.T) {
    mock := &testutil.MockExecutor{
        ExecuteFunc: func(cmd string, args ...string) (string, error) {
            if cmd == "git" && args[0] == "version" {
                return "git version 2.39.0", nil
            }
            return "", fmt.Errorf("unexpected command: %s", cmd)
        },
    }

    result, err := mock.Execute("git", "version")
    assert.NoError(t, err)
    assert.Contains(t, result, "git version")
}
```

---

## Coverage

### Coverage Targets

| Type | Target | Enforcement |
|---|---|---|
| Unit tests | 80% minimum | CI blocks below 80% |
| Critical paths | 100% | Manual review required |
| Integration tests | All major workflows | Manual review required |

### Measuring Coverage

```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage summary
go tool cover -func=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# Or use Makefile
make coverage
```

### Coverage Report Example

```
github.com/theheadlessengineer/crux/internal/domain/model/validation.go:8:   ValidateServiceName  100.0%
total:                                                                        (statements)         66.7%
```

---

## Running Tests

### Local Development

```bash
# Run all tests with race detection
make test

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/domain/model/...

# Run a specific test
go test -run TestValidateServiceName ./internal/domain/model/

# Run tests in short mode (skip slow tests)
go test -short ./...
```

### CI Pipeline

The CI pipeline automatically runs:

```yaml
- name: Test
  run: go test -race -coverprofile=coverage.out ./...
```

**Race detection is mandatory** - all tests must pass with `-race` flag.

---

## Test Helpers

### Shared Test Utilities

Place reusable test helpers in `/internal/testutil/`:

```go
// testutil/mocks.go
package testutil

type MockExecutor struct {
    ExecuteFunc func(cmd string, args ...string) (string, error)
}

func (m *MockExecutor) Execute(cmd string, args ...string) (string, error) {
    if m.ExecuteFunc != nil {
        return m.ExecuteFunc(cmd, args...)
    }
    return "", nil
}
```

### Test Fixtures

For test data files:

```
test/
├── fixtures/
│   ├── valid_config.yaml
│   └── invalid_config.yaml
└── integration/
    └── service_generation_test.go
```

---

## Best Practices

### ✅ Do

- Write tests before or alongside implementation
- Use table-driven tests for multiple scenarios
- Test both success and failure paths
- Use descriptive test case names
- Keep tests simple and focused
- Use `testify/assert` for assertions
- Run tests with race detection locally
- Aim for 80%+ coverage

### ❌ Don't

- Skip writing tests for "simple" code
- Use external mocking frameworks
- Write tests that depend on external services (use mocks)
- Commit code that fails tests
- Ignore race detector warnings
- Write tests that depend on execution order
- Use global state in tests
- Leave TODO comments in test code without a tracking issue

---

## Common Patterns

### Testing Error Cases

```go
tests := []struct {
    name    string
    input   string
    wantErr bool
    errMsg  string
}{
    {
        name:    "invalid input",
        input:   "bad-input",
        wantErr: true,
        errMsg:  "invalid input",
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        err := SomeFunction(tt.input)
        if tt.wantErr {
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.errMsg)
        } else {
            assert.NoError(t, err)
        }
    })
}
```

### Testing with Setup/Teardown

```go
func TestWithSetup(t *testing.T) {
    // Setup
    tempDir := t.TempDir() // Automatically cleaned up
    
    // Test
    result := DoSomething(tempDir)
    
    // Assert
    assert.NotNil(t, result)
    
    // Teardown happens automatically
}
```

### Testing Concurrent Code

```go
func TestConcurrentAccess(t *testing.T) {
    var wg sync.WaitGroup
    counter := NewSafeCounter()
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }
    
    wg.Wait()
    assert.Equal(t, 100, counter.Value())
}
```

---

## Troubleshooting

### Tests Pass Locally But Fail in CI

**Possible causes:**
- Race conditions (run with `-race` locally)
- Timing dependencies
- Environment-specific assumptions
- Missing test isolation

**Solution:**
```bash
# Run exactly what CI runs
go test -race -coverprofile=coverage.out ./...
```

### Race Detector Warnings

**Example warning:**
```
WARNING: DATA RACE
Write at 0x00c000100000 by goroutine 7:
```

**Solution:**
- Add proper synchronization (mutexes, channels)
- Avoid shared mutable state
- Use atomic operations where appropriate

### Low Coverage

**Check what's not covered:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Focus on:**
- Error handling paths
- Edge cases
- Validation logic

---

## Examples

### Complete Unit Test Example

See `/internal/domain/model/validation_test.go` for a complete example demonstrating:
- Table-driven tests
- Error case testing
- Multiple scenarios
- Proper assertions

### Integration Test Example

See `/test/integration/service_generation_test.go` for integration test structure.

---

## Getting Help

- Review existing tests in the codebase for patterns
- Check `/docs/development-standards.md` for coding standards
- Ask in the team Slack channel: `#crux-development`
- Refer to Go testing documentation: https://pkg.go.dev/testing

---

## Checklist for New Tests

Before submitting a PR with new code:

- [ ] Unit tests written for all new functions
- [ ] Table-driven pattern used for multiple scenarios
- [ ] Both success and failure paths tested
- [ ] Tests pass with `-race` flag
- [ ] Coverage meets 80% minimum
- [ ] Integration tests added for new workflows
- [ ] Test names follow convention
- [ ] No external dependencies in unit tests
- [ ] All tests pass: `make test`

---

**Remember:** Tests are production code. They must be maintained, reviewed, and held to the same quality standards as application code.
