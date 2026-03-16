# US-1302 — Implement Decision Dependency Graph and Conditional Logic

**Epic:** 1.3 Interactive Prompt Engine & Decision Graph
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** To Do

---

## User Story

As a user of Crux,
I want the prompt flow to only show questions that are relevant to my previous answers,
so that I am not presented with questions about technologies I have not selected.

---

## Pre-Development Checklist

- [ ] US-1301 (Question types and validation) is merged
- [ ] The full decision graph for the Phase 1 prompt flow is documented
- [ ] The `depends_on` condition syntax is agreed (supports AND/OR of previous answers)
- [ ] Story estimated and accepted into the sprint

---

## Scope

Implement conditional question visibility using a directed acyclic graph (DAG) where each question may declare dependencies on previous answers.

### In Scope

- `depends_on` field on a Question supporting AND/OR conditions referencing previous answer IDs
- A `DecisionGraph` that resolves question order and visibility given the current answer set
- Auto-addition logic: when specific combinations of plugins are selected, complementary plugins are suggested (e.g., Kafka + PostgreSQL suggests the Outbox plugin)
- Warning system: non-blocking warnings shown when a selection has known implications
- Error system: blocking errors that prevent generation with a clear explanation
- Unit tests for DAG resolution, conditional visibility, auto-additions, and both warning and error cases

### Out of Scope

- Question back-navigation (separate story)
- TUI rendering of the decision graph

---

## Acceptance Criteria

- [ ] Questions with `depends_on` conditions are shown only when conditions are met
- [ ] Questions with `depends_on` conditions are hidden (skipped) when conditions are not met
- [ ] AND conditions require all referenced answers to match
- [ ] OR conditions require at least one referenced answer to match
- [ ] Auto-addition suggestions are shown when the triggering combination is selected
- [ ] Warnings are shown but do not block generation
- [ ] Errors block generation and display a clear explanation
- [ ] DAG resolves cycles — a cyclic dependency definition is rejected at startup
- [ ] Unit tests cover all conditional paths

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [ ] Decision graph tested manually with the full `crux new` flow
- [ ] Cycle detection tested with a deliberately cyclic test case
- [ ] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-1301 Question types | Predecessor | Must be merged |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
