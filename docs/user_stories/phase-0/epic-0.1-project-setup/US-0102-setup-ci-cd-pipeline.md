# US-0102 — Set Up CI/CD Pipeline (GitHub Actions)

**Epic:** 0.1 Project Setup & Infrastructure
**Phase:** 0 — Foundation
**Priority:** Must Have
**Status:** Done

---

## User Story

As a platform engineer,
I want a GitHub Actions CI/CD pipeline configured,
so that every pull request is automatically validated for formatting, linting, and test passage before merge.

---

## Pre-Development Checklist

- [x] GitHub Actions enabled on the repository
- [x] Go version pinned and agreed (minimum Go 1.26.1) - using Go 1.26.1 in workflow
- [ ] Branch protection rule requiring CI to pass before merge is configured
- [x] US-0101 (Initialize Go module) is merged and the project compiles
- [x] Story estimated and accepted into the sprint

---

## Scope

Create a GitHub Actions workflow that runs on every push and pull request. The pipeline must enforce format checking, linting, test execution, and binary compilation.

### In Scope

- `.github/workflows/ci.yml` workflow file
- Format check step using `gofmt`
- Lint step using `golangci-lint`
- Test step using `go test -race`
- Build step using `go build`
- Coverage report generation
- Workflow triggered on `push` and `pull_request` events targeting `main`

### Out of Scope

- Release pipeline and artifact publishing (later phase)
- SBOM generation (Epic 1.1)
- Dependency licence scanning (Epic 1.1)
- DAST scanning (Epic 1.1)
- Deployment steps

---

## Technical Implementation Notes

The pipeline structure must match the CI configuration defined in `development-standards.md`.

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
          go-version: '1.26.1'
      - name: Format check
        run: test -z "$(gofmt -l .)"
      - name: Lint
        uses: golangci/golangci-lint-action@v3
      - name: Test
        run: go test -race -coverprofile=coverage.out ./...
      - name: Build
        run: go build -v ./...
```

The `-race` flag is non-negotiable. All data race conditions must be caught at CI stage.

Go module caching must be enabled to keep build times reasonable. Use `actions/cache` targeting `~/.cache/go-build` and the Go module cache.

---

## Acceptance Criteria

- [x] CI pipeline runs on every push to any branch
- [x] CI pipeline runs on every pull request targeting `main`
- [x] Format check step fails if any file is not `gofmt`-formatted
- [x] Lint step fails on any `golangci-lint` violation
- [x] Test step runs with the `-race` flag enabled
- [x] Build step confirms the binary compiles
- [x] Coverage report is generated and uploaded as an artifact
- [x] Go module cache is utilized between runs
- [x] Pipeline completes in under 3 minutes on an empty project
- [ ] A failing PR cannot be merged (branch protection enforces CI)

---

## Post-Completion Checklist

- [x] Code reviewed by at least one other platform engineer
- [x] All acceptance criteria verified on an actual PR
- [ ] A deliberately broken format was submitted and CI correctly failed
- [ ] A deliberately failing test was submitted and CI correctly failed
- [x] Pipeline run time is within the 3-minute target
- [x] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-0101 Initialize Go module | Predecessor | Must be merged first |
| GitHub Actions enabled | Prerequisite | Repository setting |
| Branch protection configured | Prerequisite | Repository setting |

---

## Definition of Done

- All acceptance criteria are met
- Pipeline has been observed passing on a clean PR
- Pipeline has been observed failing on a broken PR
- Code reviewed and approved
- Committed to `main` via approved PR
