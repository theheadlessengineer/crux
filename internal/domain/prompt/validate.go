package prompt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Validate parses raw string input for the given question and returns a typed Answer.
func Validate(q *Question, raw string) (Answer, error) {
	if raw == "" {
		raw = q.Default
	}
	switch q.Type {
	case QuestionTypeConfirm:
		return validateConfirm(q, raw)
	case QuestionTypeText:
		return validateText(q, raw)
	case QuestionTypeNumber:
		return validateNumber(q, raw)
	case QuestionTypeSelect:
		return validateSelect(q, raw)
	case QuestionTypeMultiSelect:
		return validateMultiSelect(q, raw)
	default:
		return Answer{}, fmt.Errorf("unknown question type: %s", q.Type)
	}
}

func validateConfirm(q *Question, raw string) (Answer, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "y", "yes", "true", "1":
		return Answer{QuestionID: q.ID, Value: true}, nil
	case "n", "no", "false", "0", "":
		return Answer{QuestionID: q.ID, Value: false}, nil
	default:
		return Answer{}, fmt.Errorf("expected y/n, got %q", raw)
	}
}

func validateText(q *Question, raw string) (Answer, error) {
	if q.Validation.Required && strings.TrimSpace(raw) == "" {
		return Answer{}, fmt.Errorf("value is required")
	}
	if q.Validation.Pattern != "" {
		matched, err := regexp.MatchString(q.Validation.Pattern, raw)
		if err != nil {
			return Answer{}, fmt.Errorf("invalid validation pattern: %w", err)
		}
		if !matched {
			return Answer{}, fmt.Errorf("input does not match required pattern %s", q.Validation.Pattern)
		}
	}
	return Answer{QuestionID: q.ID, Value: raw}, nil
}

func validateNumber(q *Question, raw string) (Answer, error) {
	if raw == "" {
		if q.Validation.Required {
			return Answer{}, fmt.Errorf("value is required")
		}
		return Answer{QuestionID: q.ID, Value: float64(0)}, nil
	}
	n, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return Answer{}, fmt.Errorf("expected a number, got %q", raw)
	}
	if n < q.Validation.Min {
		return Answer{}, fmt.Errorf("value %.g is below minimum %.g", n, q.Validation.Min)
	}
	if q.Validation.Max != 0 && n > q.Validation.Max {
		return Answer{}, fmt.Errorf("value %.g exceeds maximum %.g", n, q.Validation.Max)
	}
	return Answer{QuestionID: q.ID, Value: n}, nil
}

func validateSelect(q *Question, raw string) (Answer, error) {
	raw = strings.TrimSpace(raw)
	for _, opt := range q.Options {
		if opt.Value == raw || opt.Label == raw {
			return Answer{QuestionID: q.ID, Value: opt.Value}, nil
		}
	}
	return Answer{}, fmt.Errorf("%q is not a valid option", raw)
}

func validateMultiSelect(q *Question, raw string) (Answer, error) {
	if raw == "" {
		return Answer{QuestionID: q.ID, Value: []string{}}, nil
	}
	parts := strings.Split(raw, ",")
	selected := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		found := false
		for _, opt := range q.Options {
			if opt.Value == p || opt.Label == p {
				selected = append(selected, opt.Value)
				found = true
				break
			}
		}
		if !found {
			return Answer{}, fmt.Errorf("%q is not a valid option", p)
		}
	}
	return Answer{QuestionID: q.ID, Value: selected}, nil
}
