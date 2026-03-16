package prompt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/domain/prompt"
)

func ans(id, val string) prompt.Answer {
	return prompt.Answer{QuestionID: id, Value: val}
}

func answerMap(pairs ...prompt.Answer) map[string]prompt.Answer {
	m := make(map[string]prompt.Answer, len(pairs))
	for _, a := range pairs {
		m[a.QuestionID] = a
	}
	return m
}

func TestGraph_Visible_NoCondition(t *testing.T) {
	q := prompt.Question{ID: "q1", Type: prompt.QuestionTypeText}
	g, err := prompt.NewDecisionGraph([]prompt.Question{q}, nil, nil)
	require.NoError(t, err)
	assert.True(t, g.IsVisible(&q, nil))
}

func TestGraph_Visible_AND_Met(t *testing.T) {
	q1 := prompt.Question{ID: "lang", Type: prompt.QuestionTypeText}
	q2 := prompt.Question{
		ID:   "framework",
		Type: prompt.QuestionTypeText,
		DependsOn: &prompt.DependsOn{
			And: []prompt.Condition{{QuestionID: "lang", Value: "go"}},
		},
	}
	g, err := prompt.NewDecisionGraph([]prompt.Question{q1, q2}, nil, nil)
	require.NoError(t, err)

	assert.True(t, g.IsVisible(&q2, answerMap(ans("lang", "go"))))
	assert.False(t, g.IsVisible(&q2, answerMap(ans("lang", "python"))))
	assert.False(t, g.IsVisible(&q2, nil))
}

func TestGraph_Visible_OR_Met(t *testing.T) {
	q1 := prompt.Question{ID: "db", Type: prompt.QuestionTypeText}
	q2 := prompt.Question{
		ID:   "migrate",
		Type: prompt.QuestionTypeText,
		DependsOn: &prompt.DependsOn{
			Or: []prompt.Condition{
				{QuestionID: "db", Value: "postgres"},
				{QuestionID: "db", Value: "mysql"},
			},
		},
	}
	g, err := prompt.NewDecisionGraph([]prompt.Question{q1, q2}, nil, nil)
	require.NoError(t, err)

	assert.True(t, g.IsVisible(&q2, answerMap(ans("db", "postgres"))))
	assert.True(t, g.IsVisible(&q2, answerMap(ans("db", "mysql"))))
	assert.False(t, g.IsVisible(&q2, answerMap(ans("db", "mongo"))))
}

func TestGraph_CycleDetected(t *testing.T) {
	q1 := prompt.Question{
		ID:   "a",
		Type: prompt.QuestionTypeText,
		DependsOn: &prompt.DependsOn{
			And: []prompt.Condition{{QuestionID: "b", Value: "x"}},
		},
	}
	q2 := prompt.Question{
		ID:   "b",
		Type: prompt.QuestionTypeText,
		DependsOn: &prompt.DependsOn{
			And: []prompt.Condition{{QuestionID: "a", Value: "x"}},
		},
	}
	_, err := prompt.NewDecisionGraph([]prompt.Question{q1, q2}, nil, nil)
	assert.Error(t, err)
}

func TestGraph_UnknownDependsOn(t *testing.T) {
	q := prompt.Question{
		ID:   "q1",
		Type: prompt.QuestionTypeText,
		DependsOn: &prompt.DependsOn{
			And: []prompt.Condition{{QuestionID: "nonexistent", Value: "x"}},
		},
	}
	_, err := prompt.NewDecisionGraph([]prompt.Question{q}, nil, nil)
	assert.Error(t, err)
}

func TestGraph_DuplicateID(t *testing.T) {
	q := prompt.Question{ID: "q1", Type: prompt.QuestionTypeText}
	_, err := prompt.NewDecisionGraph([]prompt.Question{q, q}, nil, nil)
	assert.Error(t, err)
}

func TestGraph_AutoAdditions_Triggered(t *testing.T) {
	q1 := prompt.Question{ID: "db", Type: prompt.QuestionTypeText}
	q2 := prompt.Question{ID: "broker", Type: prompt.QuestionTypeText}
	aa := prompt.AutoAddition{
		TriggerIDs:    []string{"db", "broker"},
		TriggerValues: []string{"postgres", "kafka"},
		SuggestID:     "outbox",
		Message:       "PostgreSQL + Kafka detected — Outbox pattern recommended",
	}
	g, err := prompt.NewDecisionGraph([]prompt.Question{q1, q2}, []prompt.AutoAddition{aa}, nil)
	require.NoError(t, err)

	msgs := g.EvalAutoAdditions(answerMap(ans("db", "postgres"), ans("broker", "kafka")))
	assert.Len(t, msgs, 1)
	assert.Contains(t, msgs[0], "Outbox")

	msgs = g.EvalAutoAdditions(answerMap(ans("db", "postgres"), ans("broker", "rabbitmq")))
	assert.Empty(t, msgs)
}

func TestGraph_Rules_Warning(t *testing.T) {
	q := prompt.Question{ID: "auth", Type: prompt.QuestionTypeText}
	rule := prompt.Rule{
		ID: "no-auth-warn",
		Condition: prompt.DependsOn{
			And: []prompt.Condition{{QuestionID: "auth", Value: "none"}},
		},
		Message:  "No authentication selected — security warning",
		Blocking: false,
	}
	g, err := prompt.NewDecisionGraph([]prompt.Question{q}, nil, []prompt.Rule{rule})
	require.NoError(t, err)

	warnings, errs := g.EvalRules(answerMap(ans("auth", "none")))
	assert.Len(t, warnings, 1)
	assert.Empty(t, errs)
}

func TestGraph_Rules_BlockingError(t *testing.T) {
	q := prompt.Question{ID: "outbox", Type: prompt.QuestionTypeText}
	rule := prompt.Rule{
		ID: "outbox-no-db",
		Condition: prompt.DependsOn{
			And: []prompt.Condition{{QuestionID: "outbox", Value: "true"}},
		},
		Message:  "Outbox pattern requires a database",
		Blocking: true,
	}
	g, err := prompt.NewDecisionGraph([]prompt.Question{q}, nil, []prompt.Rule{rule})
	require.NoError(t, err)

	warnings, errs := g.EvalRules(answerMap(ans("outbox", "true")))
	assert.Empty(t, warnings)
	assert.Len(t, errs, 1)
}
