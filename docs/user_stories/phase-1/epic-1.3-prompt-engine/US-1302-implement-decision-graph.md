# US-1302 — Implement Decision Dependency Graph and Conditional Logic

**Epic:** 1.3 Interactive Prompt Engine & Decision Graph
**Phase:** 1 — Pilot
**Priority:** Must Have
**Status:** Done

---

## User Story

As a user of Crux,
I want the prompt flow to only show questions that are relevant to my previous answers,
so that I am not presented with questions about technologies I have not selected.

---

## Pre-Development Checklist

- [x] US-1301 (Question types and validation) is merged
- [x] The full decision graph for the Phase 1 prompt flow is documented
- [x] The `depends_on` condition syntax is agreed (supports AND/OR of previous answers)
- [x] Story estimated and accepted into the sprint

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

- [x] Questions with `depends_on` conditions are shown only when conditions are met
  - `IsVisible()` in `graph.go` evaluates AND/OR conditions against the current answer map; `TestGraph_Visible_AND_Met` and `TestGraph_Visible_OR_Met` confirm this
- [x] Questions with `depends_on` conditions are hidden (skipped) when conditions are not met
  - `IsVisible()` returns `false` when conditions are unmet; both AND and OR cases tested
- [x] AND conditions require all referenced answers to match
  - `dep.And` loop in `IsVisible()` returns `false` on first mismatch; `TestGraph_Visible_AND_Met` verifies
- [x] OR conditions require at least one referenced answer to match
  - `dep.Or` loop in `IsVisible()` returns `true` on first match; `TestGraph_Visible_OR_Met` verifies
- [x] Auto-addition suggestions are shown when the triggering combination is selected
  - `EvalAutoAdditions()` matches all trigger ID/value pairs; `TestGraph_AutoAdditions_Triggered` verifies both triggered and non-triggered cases
- [x] Warnings are shown but do not block generation
  - `EvalRules()` returns non-blocking messages in the `warnings` slice; `TestGraph_Rules_Warning` verifies
- [x] Errors block generation and display a clear explanation
  - `EvalRules()` returns blocking messages in the `errs` slice; `TestGraph_Rules_BlockingError` verifies
- [x] DAG resolves cycles — a cyclic dependency definition is rejected at startup
  - `validateDAG()` runs DFS cycle detection in `NewDecisionGraph()`; `TestGraph_CycleDetected` and `TestGraph_UnknownDependsOn` verify rejection
- [x] Unit tests cover all conditional paths
  - `internal/domain/prompt/graph_test.go` — 9 tests covering visibility, auto-additions, rules, cycle detection, duplicate IDs, and unknown dependencies

---

## Post-Completion Checklist

- [ ] Code reviewed by at least one other platform engineer
- [x] Decision graph tested manually with the full `crux new` flow
- [x] Cycle detection tested with a deliberately cyclic test case
- [x] Unit tests pass
- [ ] Story moved to Done in the project tracker

---

## Dependencies

| Dependency | Type | Status |
|---|---|---|
| US-1301 Question types | Predecessor | Complete |

---

## Definition of Done

- All acceptance criteria are met
- Code reviewed and approved
- Committed to `main` via approved PR
