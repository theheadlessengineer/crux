# crux-plugin-cursor

Generates Cursor IDE workspace rules for `{{ service.name }}` with company standards and service context.

## Generated Files

| File | Purpose |
|---|---|
| `.cursorrules` | Cursor workspace rules |

## Configuration

| Question | Default | Description |
|---|---|---|
| `cursor_team_name` | `Platform Engineering` | Team name for context |
| `cursor_no_touch_paths` | `infra/terraform,.github/workflows` | Paths Cursor must not modify |
