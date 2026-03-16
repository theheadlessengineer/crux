package prompt

import "fmt"

// AutoAddition describes a plugin/question suggested when a combination is selected.
type AutoAddition struct {
	TriggerIDs    []string
	TriggerValues []string
	SuggestID     string
	Message       string
}

// Rule is a warning or blocking error triggered by a combination of answers.
type Rule struct {
	ID        string
	Condition DependsOn
	Message   string
	Blocking  bool
}

// DecisionGraph holds the ordered question list plus combination rules.
type DecisionGraph struct {
	questions     []Question
	autoAdditions []AutoAddition
	rules         []Rule
	index         map[string]int
}

// NewDecisionGraph builds and validates a DecisionGraph from a question list.
func NewDecisionGraph(questions []Question, autoAdditions []AutoAddition, rules []Rule) (*DecisionGraph, error) {
	index := make(map[string]int, len(questions))
	for i := range questions {
		q := &questions[i]
		if _, dup := index[q.ID]; dup {
			return nil, fmt.Errorf("duplicate question ID %q", q.ID)
		}
		index[q.ID] = i
	}
	if err := validateDAG(questions, index); err != nil {
		return nil, err
	}
	return &DecisionGraph{
		questions:     questions,
		autoAdditions: autoAdditions,
		rules:         rules,
		index:         index,
	}, nil
}

// Questions returns the ordered question slice.
func (g *DecisionGraph) Questions() []Question { return g.questions }

// IsVisible reports whether question q should be shown given the current answers.
func (g *DecisionGraph) IsVisible(q *Question, answers map[string]Answer) bool {
	if q.DependsOn == nil {
		return true
	}
	dep := q.DependsOn
	if len(dep.And) > 0 {
		for _, c := range dep.And {
			if !conditionMet(c, answers) {
				return false
			}
		}
		return true
	}
	if len(dep.Or) > 0 {
		for _, c := range dep.Or {
			if conditionMet(c, answers) {
				return true
			}
		}
		return false
	}
	return true
}

// EvalAutoAdditions returns auto-addition messages triggered by the current answers.
func (g *DecisionGraph) EvalAutoAdditions(answers map[string]Answer) []string {
	var msgs []string
	for _, aa := range g.autoAdditions {
		triggered := true
		for i, id := range aa.TriggerIDs {
			a, ok := answers[id]
			if !ok {
				triggered = false
				break
			}
			val, _ := a.Value.(string)
			if val != aa.TriggerValues[i] {
				triggered = false
				break
			}
		}
		if triggered {
			msgs = append(msgs, aa.Message)
		}
	}
	return msgs
}

// EvalRules evaluates all rules against the current answers.
// Returns (warnings, blockingErrors).
func (g *DecisionGraph) EvalRules(answers map[string]Answer) (warnings, errs []string) {
	for i := range g.rules {
		r := &g.rules[i]
		stub := &Question{DependsOn: &r.Condition}
		if g.IsVisible(stub, answers) {
			if r.Blocking {
				errs = append(errs, r.Message)
			} else {
				warnings = append(warnings, r.Message)
			}
		}
	}
	return
}

func conditionMet(c Condition, answers map[string]Answer) bool {
	a, ok := answers[c.QuestionID]
	if !ok {
		return false
	}
	switch v := a.Value.(type) {
	case string:
		return v == c.Value
	case bool:
		return fmt.Sprintf("%v", v) == c.Value
	case float64:
		return fmt.Sprintf("%v", v) == c.Value
	case []string:
		for _, s := range v {
			if s == c.Value {
				return true
			}
		}
		return false
	}
	return false
}

func validateDAG(questions []Question, index map[string]int) error {
	deps := make(map[string][]string, len(questions))
	for i := range questions {
		q := &questions[i]
		if q.DependsOn == nil {
			continue
		}
		for _, c := range q.DependsOn.And {
			if _, ok := index[c.QuestionID]; !ok {
				return fmt.Errorf("question %q depends_on unknown ID %q", q.ID, c.QuestionID)
			}
			deps[q.ID] = append(deps[q.ID], c.QuestionID)
		}
		for _, c := range q.DependsOn.Or {
			if _, ok := index[c.QuestionID]; !ok {
				return fmt.Errorf("question %q depends_on unknown ID %q", q.ID, c.QuestionID)
			}
			deps[q.ID] = append(deps[q.ID], c.QuestionID)
		}
	}

	const (
		unvisited = 0
		visiting  = 1
		visited   = 2
	)
	state := make(map[string]int, len(questions))

	var dfs func(id string) error
	dfs = func(id string) error {
		if state[id] == visiting {
			return fmt.Errorf("cyclic dependency detected involving question %q", id)
		}
		if state[id] == visited {
			return nil
		}
		state[id] = visiting
		for _, dep := range deps[id] {
			if err := dfs(dep); err != nil {
				return err
			}
		}
		state[id] = visited
		return nil
	}

	for i := range questions {
		if err := dfs(questions[i].ID); err != nil {
			return err
		}
	}
	return nil
}
