package prompt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/domain/prompt"
)

func buildSession(questions ...prompt.Question) *prompt.Session {
	g, err := prompt.NewDecisionGraph(questions, nil, nil)
	if err != nil {
		panic(err)
	}
	return prompt.NewSession(g)
}

func TestSession_ForwardNavigation(t *testing.T) {
	q1 := prompt.Question{ID: "lang", Type: prompt.QuestionTypeText}
	q2 := prompt.Question{ID: "framework", Type: prompt.QuestionTypeText}
	s := buildSession(q1, q2)

	next := s.NextQuestion()
	require.NotNil(t, next)
	assert.Equal(t, "lang", next.ID)
	s.Record(next, prompt.Answer{QuestionID: "lang", Value: "go"})

	next = s.NextQuestion()
	require.NotNil(t, next)
	assert.Equal(t, "framework", next.ID)
	s.Record(next, prompt.Answer{QuestionID: "framework", Value: "gin"})

	assert.Nil(t, s.NextQuestion())
}

func TestSession_BackNavigation(t *testing.T) {
	q1 := prompt.Question{ID: "lang", Type: prompt.QuestionTypeText}
	q2 := prompt.Question{ID: "framework", Type: prompt.QuestionTypeText}
	s := buildSession(q1, q2)

	s.Record(&q1, prompt.Answer{QuestionID: "lang", Value: "go"})
	s.Record(&q2, prompt.Answer{QuestionID: "framework", Value: "gin"})

	require.NoError(t, s.Back())

	next := s.NextQuestion()
	require.NotNil(t, next)
	assert.Equal(t, "framework", next.ID)
}

func TestSession_BackAtFirst_ReturnsError(t *testing.T) {
	s := buildSession(prompt.Question{ID: "q1", Type: prompt.QuestionTypeText})
	assert.ErrorIs(t, s.Back(), prompt.ErrAtFirstQuestion)
}

func TestSession_BackClearsConditionalAnswers(t *testing.T) {
	q1 := prompt.Question{ID: "lang", Type: prompt.QuestionTypeText}
	q2 := prompt.Question{
		ID:   "framework",
		Type: prompt.QuestionTypeText,
		DependsOn: &prompt.DependsOn{
			And: []prompt.Condition{{QuestionID: "lang", Value: "go"}},
		},
	}
	s := buildSession(q1, q2)

	s.Record(&q1, prompt.Answer{QuestionID: "lang", Value: "go"})
	s.Record(&q2, prompt.Answer{QuestionID: "framework", Value: "gin"})

	require.NoError(t, s.Back()) // removes framework
	require.NoError(t, s.Back()) // removes lang

	s.Record(&q1, prompt.Answer{QuestionID: "lang", Value: "python"})

	_, hasFramework := s.Answers()["framework"]
	assert.False(t, hasFramework)
}

func TestSession_AnswersReturnsSnapshot(t *testing.T) {
	q1 := prompt.Question{ID: "lang", Type: prompt.QuestionTypeText}
	s := buildSession(q1)
	s.Record(&q1, prompt.Answer{QuestionID: "lang", Value: "go"})

	snap := s.Answers()
	assert.Equal(t, "go", snap["lang"].Value)
}
