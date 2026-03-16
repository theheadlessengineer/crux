# US-1303 — Implement Question History and Back Navigation

**Epic:** 1.3 Interactive Prompt Engine & Decision Graph
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want to navigate back to a previous question and change my answer,
so that I can correct a mistake without restarting the entire `crux new` flow.

---

## Pre-Development Checklist

- [x] US-1302 (Decision graph) is merged
- [x] The back-navigation key binding is agreed (common: Ctrl+B or pressing `b`)
- [x] The behaviour when navigating back past a conditional question is agreed
- [x] Story estimated and accepted into the sprint

---

## Scope

Implement a question history stack within the prompt engine that allows the user to step backward through previously answered questions and change their responses.

### In Scope

- Answer history stack that records each question and its answer in order
- Back navigation triggered by an agreed key combination
- When navigating back, the decision graph re-evaluates question visibility from the changed answer forward
- Questions that are no longer visible due to an answer change have their answers cleared from state
- Unit tests for navigation, answer clearing, and re-evaluation

### Out of Scope

- TUI-specific back navigation (TUI implementation in Epic 2.1 — this covers the engine layer)

---

## Acceptance Criteria

- [x] User can navigate back to the previous question
  - `Session.Back()` in `session.go` pops the last history entry and removes its answer; `TestSession_BackNavigation` verifies the previous question is re-presented
- [x] Changing an answer causes subsequent conditional questions to re-evaluate
  - `clearHidden()` is called after every `Back()`, re-evaluating `IsVisible()` for all remaining answers; `TestSession_BackClearsConditionalAnswers` verifies that a conditional answer is cleared when its dependency changes
- [x] Answers to questions that are no longer visible are cleared from state
  - `clearHidden()` deletes answers whose questions are no longer visible and rebuilds the history slice to match; `TestSession_BackClearsConditionalAnswers` confirms the `framework` answer is absent after `lang` is changed to `python`
- [x] Navigation at the first question produces a clear message (cannot go back further)
  - `Back()` returns `ErrAtFirstQuestion` when history is empty; `TestSession_BackAtFirst_ReturnsError` asserts `errors.Is(err, ErrAtFirstQuestion)`
- [x] Unit tests cover forward navigation, backward navigation, answer clearing, and re-evaluation
  - `internal/domain/prompt/session_test.go` — 5 tests: `TestSession_ForwardNavigation`, `TestSession_BackNavigation`, `TestSession_BackAtFirst_ReturnsError`, `TestSession_BackClearsConditionalAnswers`, `TestSession_AnswersReturnsSnapshot`

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [x] Back navigation tested manually through the full `crux new` flow
- [x] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-1302 Decision graph | Predecessor | Complete |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
