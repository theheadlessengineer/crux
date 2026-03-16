# US-2102 — TUI Enhancements Beyond US-2101

**Epic:** 2.1 Terminal UI (TUI)
**Phase:** 2 — MVP
**Priority:** Should Have
**Status:** Done

---

## User Story

As a user of Crux,
I want the `crux new` TUI to guide me safely through the full service configuration — including plugin selection, per-plugin questions, a review screen, and safe abort — in a single uninterrupted session,
so that I never lose context, never see a blank screen between steps, and always have the opportunity to review and correct my choices before a service is created.

---

## Context

US-2101 defined the core TUI skeleton: question rendering, keyboard navigation, two themes, help overlay, and theme persistence. This story captures everything implemented beyond that scope during the Epic 2.1 build-out.

---

## What Was Implemented

### 1. Six Named Themes (US-2101 specified two)

Six full RGB hex palettes were implemented, matching the llmfit colour system:

| Theme | Description |
|---|---|
| `default` | Terminal-native 256-colour palette |
| `dracula` | Dark purple/pink — `#282a36` background |
| `solarized` | Dark teal — `#002b36` background |
| `nord` | Arctic blue — `#2e3440` background |
| `monokai` | Dark olive — `#272822` background |
| `gruvbox` | Warm dark — `#282828` background |

Each theme defines 13 semantic colour tokens: `BG`, `FG`, `Muted`, `Border`, `Title`, `Accent`, `Cursor`, `Selected`, `Good`, `Warning`, `Error`, `StatusBG`, `StatusFG`.

Theme preference is persisted to `~/.crux/theme` and restored on next run.

**Files:** `internal/presentation/tui/theme.go`

---

### 2. Context-Sensitive Status Bar

The status bar at the bottom of the screen changes its keybinding hints based on the active question type:

- **Text / Number questions:** `type answer  enter confirm  backspace delete  b back (when empty)  ctrl+c quit`
- **Select / Confirm questions:** `↑/↓ navigate  enter confirm  b back  t theme  ? help  ctrl+c quit`
- **MultiSelect questions:** `↑/↓ navigate  space toggle  enter confirm  b back  t theme  ctrl+c quit`
- **Review screen:** `↑/↓ navigate  enter select  ctrl+c abort`
- **Edit-pick mode:** `↑/↓ select answer  enter edit  esc cancel`

The `t` (theme) and `?` (help) keys are intentionally suppressed on text/number questions because those characters are valid input. The status bar reflects this so users are not confused.

**Files:** `internal/presentation/tui/model.go` — `renderStatusBar()`

---

### 3. Plugin Embedding and Single-Session Plugin Flow

**Problem solved:** Plugins were previously loaded from the filesystem using `runtime.Caller` path resolution. This broke when the binary ran outside the repository directory, causing the plugin selection question to never appear.

**Solution:**

- All 9 pilot plugins are embedded into the binary via `//go:embed all:crux-plugin-*` in `data/plugins/embed.go`
- A new `LoadFromFS(fsys fs.FS, cruxVersion string)` function was added to `internal/infrastructure/plugin/loader.go` that reads plugin manifests directly from the embedded FS
- `loadAvailablePlugins` in `new.go` now calls `LoadFromFS(dataplugins.FS, cruxVersion)` — works from any working directory
- Version compatibility check (`checkCompatibility`) was updated to bypass semver validation for non-semver version strings (git hashes, dirty tags), so development builds always load plugins

**All plugin questions are now included in a single session upfront**, gated by `DependsOn` on the `_plugins` multiselect answer. This eliminates the previous two-pass approach that launched multiple Bubbletea programs and caused blank screens between plugin question sets.

`conditionMet` in `internal/domain/prompt/graph.go` was extended to handle `[]string` answers (multiselect), enabling `DependsOn` conditions of the form "show this question if `_plugins` contains `crux-plugin-postgresql`".

**Files:** `data/plugins/embed.go`, `internal/infrastructure/plugin/loader.go`, `internal/presentation/cli/new.go`, `internal/domain/prompt/graph.go`

---

### 4. MultiSelect Question UX

The multiselect question type (used for plugin selection) received dedicated UX treatment:

- Live selection counter: `N selected — space to toggle, enter to confirm` rendered in the accent colour, updating as the user toggles options
- Selected items shown with `●` (accent colour); unselected with `○` (muted)
- Active cursor row highlighted with `❯` prefix and prompt colour on the label
- Status bar shows multiselect-specific hints when this question type is active

**Files:** `internal/presentation/tui/model.go` — `updateMultiSelect()`, `renderQuestionPanel()`

---

### 5. Per-Question Help Text

Every question now displays contextual help text below the input area, separated by a thin `─` divider. Help text is word-wrapped to the panel width and rendered in the muted colour.

Core questions have help text defined inline in `coreQuestions()`. Plugin questions read `help:` from the plugin manifest's `QuestionSpec`. The `_plugins` multiselect question has its own help text explaining what plugins add and how to install more later.

**`QuestionSpec`** in `internal/domain/plugin/plugin.go` was extended with a `help` YAML field so plugin authors can document their questions.

**Files:** `internal/domain/plugin/plugin.go`, `internal/presentation/cli/new.go`, `internal/presentation/tui/model.go` — `renderQuestionPanel()`, `wrapText()`

---

### 6. Post-Completion Review Screen

When all questions are answered, the TUI transitions to a full-screen review screen instead of immediately quitting. The review screen shows:

- Every answered question with its value (multiselect values joined with `, `)
- Three action options navigated with `↑/↓` and confirmed with `enter`:
  1. **Confirm & create service** — proceeds to generation
  2. **Change plugin selection** — navigates back to the `_plugins` question; returns to review when re-answered
  3. **Edit an answer** — enters edit-pick mode

**Edit-pick mode:** The cursor moves through the answered question list. Pressing `enter` on any question navigates back to it for re-answering. After re-answering, the session continues forward and returns to the review screen automatically. `esc` or `b` cancels and returns to the action menu.

Navigation back to a specific question is handled by `navigateTo(id string)`, which walks the session backwards using `session.Back()` until the target question is the next unanswered one, keeping the local answers map in sync.

**Files:** `internal/presentation/tui/model.go` — `updateReview()`, `renderReview()`, `navigateTo()`

---

### 7. Abort Confirmation Overlay

Pressing `ctrl+c` at any point during the TUI (including during the review screen) shows a centred confirmation overlay instead of immediately quitting:

```
╭──────────────────────────────────────╮
│  Abort service creation?             │
│                                      │
│  [ Y ] Yes, abort   [ any key ] Continue │
╰──────────────────────────────────────╯
```

- `y` / `Y` → sets `aborted = true`, quits the TUI, `Run()` returns `ErrAborted`
- Any other key → dismisses the overlay, session resumes exactly where it was

`new.go` catches `ErrAborted` with `errors.Is`, prints `Aborted.`, and exits with code 0. No files are written.

**Files:** `internal/presentation/tui/model.go` — `renderQuitConfirm()`, `ErrAborted`; `internal/presentation/cli/new.go`

---

### 8. `session.Progress()` and `session.Graph()` on the Domain

Two methods were added to `prompt.Session` to support TUI rendering without coupling the presentation layer to session internals:

- `Progress() (answered, total int)` — counts currently visible answered and total questions; used by the header bar progress indicator `(N/total)`
- `Graph() *DecisionGraph` — exposes the question graph read-only; used by the answers panel and review screen to walk questions in order

**Files:** `internal/domain/prompt/session.go`

---

## Acceptance Criteria (all met)

- [x] Plugin selection question always appears regardless of working directory
- [x] All plugin questions appear in the same TUI session as core questions — no blank screen between sets
- [x] MultiSelect shows live selection count and clear toggle/confirm instructions
- [x] Every question shows contextual help text below the input
- [x] After all questions are answered, a review screen is shown before any files are written
- [x] User can change plugin selection from the review screen
- [x] User can edit any individual answer from the review screen
- [x] After editing, the session returns to the review screen automatically
- [x] `ctrl+c` shows a confirmation overlay before aborting
- [x] Aborting exits cleanly with code 0 and writes no files
- [x] Status bar shows context-appropriate keybinding hints per question type
- [x] Six themes available and cycle correctly with `t`
- [x] Theme preference persisted and restored across sessions

---

## Files Changed

| File | Change |
|---|---|
| `internal/presentation/tui/model.go` | New file — full TUI model |
| `internal/presentation/tui/theme.go` | New file — 6 themes, persistence |
| `internal/domain/prompt/session.go` | Added `Progress()`, `Graph()` |
| `internal/domain/prompt/graph.go` | `conditionMet` handles `[]string` |
| `internal/domain/plugin/plugin.go` | Added `Help` field to `QuestionSpec` |
| `internal/infrastructure/plugin/loader.go` | Added `LoadFromFS`, fixed version bypass |
| `internal/presentation/cli/new.go` | Single-session flow, embedded plugins, abort handling |
| `data/plugins/embed.go` | New file — embeds all 9 pilot plugins |

---

## Dependencies

| Dependency | Type |
|---|---|
| US-2101 TUI skeleton | Predecessor |
| Epic 1.3 Prompt Engine | Extended (`Progress`, `Graph`, `conditionMet`) |
| Epic 1.5 Plugin System | Extended (`LoadFromFS`, `Help` field, embedding) |
