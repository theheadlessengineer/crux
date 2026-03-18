# crux-plugin-gitlab-ci

Generates a GitLab CI/CD pipeline for `{{ service.name }}`: lint, test, SAST, container scan, build, and deploy stages.

## Generated Files

| File | Purpose |
|---|---|
| `.gitlab-ci.yml` | Full GitLab CI/CD pipeline |

## Configuration

| Question | Default | Description |
|---|---|---|
| `gitlab_coverage_threshold` | `80` | Minimum test coverage % |
| `gitlab_container_registry` | `registry.gitlab.com/company/{{ service.name }}` | Container registry URL |
| `gitlab_deploy_env` | `staging` | Auto-deploy target environment |

## Pipeline Stages

1. `lint` — golangci-lint
2. `test` — go test with race detector and coverage
3. `sast` — GitLab SAST template
4. `scan` — Trivy container image scan
5. `build` — Docker build and push
6. `deploy` — kubectl rollout (manual for prod)
