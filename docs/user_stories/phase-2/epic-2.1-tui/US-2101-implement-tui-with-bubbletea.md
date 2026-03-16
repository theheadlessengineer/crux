# US-2101 — Implement Terminal UI with Bubbletea

**Epic:** 2.1 Terminal UI (TUI)
**Phase:** 2 — MVP
**Priority:** Should Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want a polished terminal interface for the `crux new` prompt flow,
so that the experience is more intuitive and professional than plain text prompts.

---

## Pre-Development Checklist

- [ ] Phase 1 is complete and stable
- [ ] Bubbletea version is agreed and added to `go.mod`
- [ ] TUI layout and component design is agreed (wireframes or prototype)
- [ ] Terminal compatibility requirements are agreed (minimum terminal dimensions)
- [ ] Story estimated and accepted into the sprint

---

## Scope

Implement the Bubbletea-based TUI that replaces the plain text prompt flow for the `crux new` command. The prompt engine (Epic 1.3) remains unchanged — the TUI is a presentation layer over it.

### Required Components

| Component | Description |
|---|---|
| Question display | Renders the current question with options |
| Progress indicator | Shows current question number out of total |
| Summary display | Shows answers so far on the right panel |
| Keyboard navigation | Arrows, vim keys (j/k), Enter, Backspace |
| Theme system | Light and dark themes, cycled with `t` key |
| Help overlay | Triggered by `?`, shows all keybindings |

### In Scope

- Bubbletea application wiring
- All five question types rendered in the TUI
- Keyboard navigation as specified
- Two themes: light and dark
- Theme preference persisted in `~/.crux/theme`
- Help overlay
- TUI tests using Bubbletea test utilities

---

## Acceptance Criteria

- [ ] TUI renders correctly in the agreed terminal emulators (iTerm2, Terminal.app, common Linux terminals)
- [ ] All five question types are rendered and interactive
- [ ] Keyboard navigation works: arrow keys, j/k, Enter to confirm, Backspace to go back
- [ ] Theme cycles with the `t` key
- [ ] Theme preference is saved and restored on next run
- [ ] Help overlay shows all keybindings when `?` is pressed
- [ ] TUI gracefully handles terminal resize
- [ ] TUI tests pass

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] TUI tested in at least three different terminal emulators
- [ ] All question types verified in TUI
- [ ] Theme persistence verified across sessions
- [ ] TUI tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| Epic 1.3 Prompt Engine | Predecessor | TUI renders on top of the engine |
| Bubbletea added to go.mod | Prerequisite | |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
