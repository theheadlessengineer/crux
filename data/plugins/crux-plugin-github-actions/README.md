# crux-plugin-github-actions

**Tier:** 1 (Official) | **Version:** 1.0.0 | **Phase:** Pilot

GitHub Actions CI/CD pipeline for crux-generated Go services.

## Questions

| ID | Type | Prompt | Default |
|---|---|---|---|
| `gha_coverage_threshold` | input | Minimum test coverage (%) | `80` |
| `gha_container_registry` | select | Container registry | `ghcr.io` |
| `gha_deploy_env` | select | Auto-deploy on merge | `staging` |

## Generated Files

| File | Description |
|---|---|
| `.github/workflows/ci.yaml` | Lint → Test → Coverage gate → Build → SAST |
| `.github/workflows/deploy.yaml` | Build image → push → deploy to staging |

## CI Pipeline Steps

1. `golangci-lint` — linting
2. `go test -race` — tests with race detector
3. Coverage gate — fails if below threshold
4. `go build` — compilation check
5. `govulncheck` — vulnerability scanning
