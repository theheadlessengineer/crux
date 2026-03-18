# crux-plugin-github-copilot

Generates GitHub Copilot workspace instructions for `{{ service.name }}` with company standards and service context.

## Generated Files

| File | Purpose |
|---|---|
| `.github/copilot-instructions.md` | Workspace-level Copilot instructions |

## Configuration

| Question | Default | Description |
|---|---|---|
| `copilot_team_name` | `Platform Engineering` | Team name for context |
| `copilot_no_touch_paths` | `infra/terraform,.github/workflows` | Paths Copilot must not modify |
