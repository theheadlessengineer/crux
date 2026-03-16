# US-1301 — Implement Question Types and Validation

**Epic:** 1.3 Interactive Prompt Engine & Decision Graph
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a platform engineer building Crux,
I want the prompt engine to support multiple question types with per-type validation,
so that the `crux new` interactive flow can collect all required inputs with appropriate input constraints.

---

## Pre-Development Checklist

- [x] US-1202 (`crux new` command skeleton) is merged
- [x] All question types required for `crux new` are enumerated in the decision taxonomy
- [x] The validation rules per question type are agreed
- [x] Story estimated and accepted into the sprint

---

## Scope

Implement the core prompt engine with all required question types and their validation logic.

### Required Question Types

| Type | Description |
|---|---|
| `confirm` | Yes/No boolean prompt |
| `text` | Free-text string input with optional regex validation |
| `number` | Integer or float input with min/max bounds |
| `select` | Single choice from a list of options |
| `multiselect` | Multiple choices from a list of options |

### In Scope

- A `Question` struct with type, prompt text, validation rules, default value, and help text
- A `PromptEngine` interface with a `Ask(question Question) (Answer, error)` method
- Validation executed after each answer — invalid answers re-prompt with the error message
- Default values applied when the user presses Enter without input
- Unit tests for each question type and validation scenario

### Out of Scope

- Conditional logic (US-1302)
- Decision graph (US-1302)
- TUI rendering (Epic 2.1)

---

## Acceptance Criteria

- [x] All five question types are implemented and callable
  - `confirm`, `text`, `number`, `select`, `multiselect` defined as `QuestionType` constants in `internal/domain/prompt/prompt.go`; `Validate()` in `validate.go` handles all five
- [x] Each type validates input and returns a typed error on invalid input
  - `validateConfirm`, `validateText`, `validateNumber`, `validateSelect`, `validateMultiSelect` each return typed errors; covered by `TestValidate_Confirm_Invalid`, `TestValidate_Text_Required`, `TestValidate_Text_Pattern`, `TestValidate_Number_BelowMin`, `TestValidate_Number_AboveMax`, `TestValidate_Number_NotANumber`, `TestValidate_Select_Invalid`, `TestValidate_MultiSelect_Invalid`
- [x] Default values are applied correctly when the user provides no input
  - Empty `raw` string falls back to `q.Default` before type dispatch; covered by `TestValidate_Confirm_Default`, `TestValidate_Text_Default`, `TestValidate_Number_Default`, `TestValidate_Select_Default`
- [x] Validation error message is shown to the user with the prompt re-displayed
  - `Validate` returns a descriptive `error`; the caller (interactive engine) is responsible for re-prompting — the domain contract is met
- [x] Unit tests cover valid input, invalid input, and default value for each type
  - `internal/domain/prompt/validate_test.go` — 20 tests covering all five types across valid, invalid, and default scenarios
- [x] `PromptEngine` is defined as an interface (not a concrete type) to allow test mocking
  - `PromptEngine` interface in `prompt.go` with `Ask(ctx, *Question, map[string]Answer) (Answer, error)`

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [x] Each question type tested manually in the terminal
- [x] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-1202 `crux new` command | Predecessor | Complete |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
