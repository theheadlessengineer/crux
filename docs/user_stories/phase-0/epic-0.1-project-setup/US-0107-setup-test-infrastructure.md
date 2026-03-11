# US-0107 — Set Up Test Infrastructure

**Epic:** 0.1 Project Setup & Infrastructure
**Phase:** 0 — Foundation
**Priority:** Must Have
**Status:** Done

---

## User Story

As a platform engineer,
I want a test infrastructure in place with agreed patterns and helpers,
so that engineers write consistent, reliable tests from the first feature story.

---

## Pre-Development Checklist

- [x] US-0106 (Configure dependency management) is merged — `testify` must be available
- [x] Team has reviewed and agreed on the test patterns in development-standards.md
- [x] Coverage targets agreed: 80% unit, 100% critical paths, all major integration workflows
- [x] Story estimated and accepted into the sprint

---

## Scope

Establish the test helper packages, confirm the testing patterns, and verify the test pipeline runs correctly with an example test.

### In Scope

- A `testhelpers` or `testutil` internal package with shared mock utilities
- An example unit test in the domain layer demonstrating the agreed table-driven pattern
- An example integration test stub with the correct build tag (`//go:build integration`)
- Confirmation that `go test -race ./...` runs successfully
- Confirmation that `go test -tags=integration ./...` runs the integration tests
- Coverage reporting via `go test -coverprofile=coverage.out ./...`

### Out of Scope

- Test fixtures for specific features (added by feature stories)
- End-to-end test infrastructure (later phase)
- Test database infrastructure (Epic 1.7)

---

## Technical Implementation Notes

All tests follow the table-driven pattern defined in development-standards.md. Test naming must follow the `TestFunctionName_Scenario` convention:

```go
func TestValidateServiceName(t *testing.T)
func TestValidateServiceName_EmptyName(t *testing.T)
func TestValidateServiceName_NameTooLong(t *testing.T)
```

Integration tests must use the build tag:
```go
//go:build integration
```

Run integration tests with:
```bash
go test -tags=integration ./...
```

The mock pattern must use interface-based mocks, not mocking frameworks, as defined in development-standards.md:

```go
type MockDetector struct {
    DetectFunc func() (*Info, error)
}

func (m *MockDetector) Detect() (*Info, error) {
    if m.DetectFunc != nil {
        return m.DetectFunc()
    }
    return nil, nil
}
```

---

## Acceptance Criteria

- [x] `go test -race ./...` runs successfully on the stub project with exit code 0
- [x] At least one example unit test exists and passes, demonstrating table-driven pattern
- [x] At least one integration test stub exists with the correct build tag
- [x] `go test -tags=integration ./...` runs the integration test
- [x] `go test -coverprofile=coverage.out ./...` generates a coverage report
- [x] A shared `testutil` package exists with at least one mock helper
- [x] Makefile `make test` and `make coverage` targets are wired to the correct commands
- [x] Test naming convention is documented in development-standards.md

---

## Post-Completion Checklist

- [x] Code reviewed by at least one other platform engineer
- [x] Test patterns demonstrated to the team (brief walkthrough)
- [x] CI pipeline confirmed to run tests with race detection
- [x] All acceptance criteria verified
- [x] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-0106 Configure dependency management | Predecessor | ✅ Done |
| US-0105 Create Makefile | Predecessor | ✅ Done |

---

## Definition of Done

- All acceptance criteria are met
- Test patterns confirmed and socialized with the team
- Code reviewed and approved
- Committed to `main` via approved PR
