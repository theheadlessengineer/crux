# US-1303 — Implement Question History and Back Navigation

**Epic:** 1.3 Interactive Prompt Engine & Decision Graph
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want to navigate back to a previous question and change my answer,
so that I can correct a mistake without restarting the entire `crux new` flow.

---

## Pre-Development Checklist

- [ ] US-1302 (Decision graph) is merged
- [ ] The back-navigation key binding is agreed (common: Ctrl+B or pressing `b`)
- [ ] The behaviour when navigating back past a conditional question is agreed
- [ ] Story estimated and accepted into the sprint

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

- [ ] User can navigate back to the previous question
- [ ] Changing an answer causes subsequent conditional questions to re-evaluate
- [ ] Answers to questions that are no longer visible are cleared from state
- [ ] Navigation at the first question produces a clear message (cannot go back further)
- [ ] Unit tests cover forward navigation, backward navigation, answer clearing, and re-evaluation

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Back navigation tested manually through the full `crux new` flow
- [ ] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-1302 Decision graph | Predecessor | Must be merged |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
