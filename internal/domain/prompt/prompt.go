// Package prompt defines the core prompt engine domain types and interfaces.
package prompt

import "context"

// QuestionType identifies the kind of input a question expects.
type QuestionType string

const (
	QuestionTypeConfirm     QuestionType = "confirm"
	QuestionTypeText        QuestionType = "text"
	QuestionTypeNumber      QuestionType = "number"
	QuestionTypeSelect      QuestionType = "select"
	QuestionTypeMultiSelect QuestionType = "multiselect"
)

// Option is a single choice in a select or multiselect question.
type Option struct {
	Label string
	Value string
}

// Condition is a single depends_on predicate: answer[ID] == Value.
type Condition struct {
	QuestionID string
	Value      string
}

// DependsOn expresses when a question is visible.
// AND: all Conditions must match. OR: at least one must match.
type DependsOn struct {
	And []Condition
	Or  []Condition
}

// Question is the complete definition of a single prompt step.
type Question struct {
	ID         string
	Type       QuestionType
	Prompt     string
	Help       string
	Default    string   // raw string default; interpreted per type
	Options    []Option // for select / multiselect
	Validation ValidationRule
	DependsOn  *DependsOn
}

// ValidationRule holds per-question validation constraints.
type ValidationRule struct {
	Required bool
	Pattern  string  // regex — text type only
	Min      float64 // number type only
	Max      float64 // number type only; 0 means no upper bound
}

// Answer holds the resolved value for a question.
type Answer struct {
	QuestionID string
	Value      interface{} // bool | string | float64 | []string
}

// PromptEngine is the interface the crux new command uses to collect answers.
// Implementations may be interactive (terminal) or non-interactive (config file / test stub).
type PromptEngine interface {
	// Ask presents a single question and returns the validated answer.
	Ask(ctx context.Context, q *Question, answers map[string]Answer) (Answer, error)
}
