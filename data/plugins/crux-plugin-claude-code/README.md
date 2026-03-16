# crux-plugin-claude-code

**Tier:** 1 (Official) | **Version:** 1.0.0 | **Phase:** Pilot

Claude Code AI assistant configuration for crux-generated services. Generates `CLAUDE.md` with service-specific context and `.claude/settings.json`.

## Questions

| ID | Type | Prompt | Default |
|---|---|---|---|
| `claude_team_name` | input | Team name for CLAUDE.md | `Platform Engineering` |
| `claude_no_touch_paths` | input | Paths Claude must not modify | `infra/terraform,.github/workflows` |

## Generated Files

| File | Description |
|---|---|
| `CLAUDE.md` | Service context — architecture, standards, key files, do-not-touch paths |
| `.claude/settings.json` | Claude Code permissions and environment config |

## What CLAUDE.md Contains

- Service overview and team ownership
- Hexagonal architecture explanation with directory map
- Company standards Claude must follow (RFC 7807, structured logging, traceparent)
- Do-not-touch paths
- Key files reference table
- Development commands
