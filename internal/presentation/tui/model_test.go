package tui_test

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/domain/prompt"
	"github.com/theheadlessengineer/crux/internal/presentation/tui"
)

func newSession(t *testing.T, questions []prompt.Question) *prompt.Session {
	t.Helper()
	graph, err := prompt.NewDecisionGraph(questions, nil, nil)
	require.NoError(t, err)
	return prompt.NewSession(graph)
}

// sendKeys sends rune key messages to the model.
func sendKeys(m tui.Model, keys ...string) tui.Model {
	var cur tea.Model = m
	for _, k := range keys {
		updated, _ := cur.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)})
		cur = updated
	}
	return cur.(tui.Model)
}

func sendKey(m tui.Model, keyType tea.KeyType) tui.Model {
	updated, _ := m.Update(tea.KeyMsg{Type: keyType})
	return updated.(tui.Model)
}

func TestModel_TextQuestion_AnswerAndAdvance(t *testing.T) {
	questions := []prompt.Question{
		{ID: "team", Type: prompt.QuestionTypeText, Prompt: "Team name", Default: "platform"},
	}
	m := tui.NewModel(newSession(t, questions), "my-service")
	m = sendKeys(m, "p", "a", "y", "m", "e", "n", "t", "s")
	m = sendKey(m, tea.KeyEnter)

	assert.True(t, m.Done())
	assert.Equal(t, "payments", m.Answers()["team"].Value)
}

func TestModel_TextQuestion_DefaultOnEmptyEnter(t *testing.T) {
	questions := []prompt.Question{
		{ID: "team", Type: prompt.QuestionTypeText, Prompt: "Team name", Default: "platform"},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	m = sendKey(m, tea.KeyEnter)

	assert.True(t, m.Done())
	assert.Equal(t, "platform", m.Answers()["team"].Value)
}

func TestModel_TextQuestion_ValidationError(t *testing.T) {
	questions := []prompt.Question{
		{
			ID: "name", Type: prompt.QuestionTypeText, Prompt: "Name",
			Validation: prompt.ValidationRule{Required: true, Pattern: `^[a-z]+$`},
		},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	m = sendKeys(m, "A", "B", "C")
	m = sendKey(m, tea.KeyEnter)

	assert.False(t, m.Done())
	assert.Contains(t, m.View(), "✘")
}

func TestModel_ConfirmQuestion_YesNo(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"yes lowercase", "y", true},
		{"yes uppercase", "Y", true},
		{"no lowercase", "n", false},
		{"no uppercase", "N", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			questions := []prompt.Question{
				{ID: "db", Type: prompt.QuestionTypeConfirm, Prompt: "Add DB?", Default: "y"},
			}
			m := tui.NewModel(newSession(t, questions), "svc")
			m = sendKeys(m, tc.key)
			assert.True(t, m.Done())
			assert.Equal(t, tc.expected, m.Answers()["db"].Value)
		})
	}
}

func TestModel_ConfirmQuestion_DefaultOnEnter(t *testing.T) {
	questions := []prompt.Question{
		{ID: "db", Type: prompt.QuestionTypeConfirm, Prompt: "Add DB?", Default: "y"},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	m = sendKey(m, tea.KeyEnter)
	assert.True(t, m.Done())
	assert.Equal(t, true, m.Answers()["db"].Value)
}

func TestModel_SelectQuestion_NavigateAndConfirm(t *testing.T) {
	questions := []prompt.Question{
		{
			ID: "lang", Type: prompt.QuestionTypeSelect, Prompt: "Language",
			Options: []prompt.Option{{Label: "Go", Value: "go"}, {Label: "Python", Value: "python"}},
		},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	m = sendKeys(m, "j") // move to Python
	m = sendKey(m, tea.KeyEnter)

	assert.True(t, m.Done())
	assert.Equal(t, "python", m.Answers()["lang"].Value)
}

func TestModel_SelectQuestion_VimKeys(t *testing.T) {
	questions := []prompt.Question{
		{
			ID: "lang", Type: prompt.QuestionTypeSelect, Prompt: "Language",
			Options: []prompt.Option{
				{Label: "Go", Value: "go"},
				{Label: "Python", Value: "python"},
				{Label: "Java", Value: "java"},
			},
		},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	m = sendKeys(m, "j", "j", "k") // → index 1 (Python)
	m = sendKey(m, tea.KeyEnter)

	assert.True(t, m.Done())
	assert.Equal(t, "python", m.Answers()["lang"].Value)
}

func TestModel_MultiSelectQuestion_ToggleAndConfirm(t *testing.T) {
	questions := []prompt.Question{
		{
			ID: "plugins", Type: prompt.QuestionTypeMultiSelect, Prompt: "Plugins",
			Options: []prompt.Option{
				{Label: "PostgreSQL", Value: "postgresql"},
				{Label: "Redis", Value: "redis"},
				{Label: "Kafka", Value: "kafka"},
			},
		},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	m = sendKey(m, tea.KeySpace) // toggle postgresql
	m = sendKeys(m, "j")         // move to redis
	m = sendKey(m, tea.KeySpace) // toggle redis
	m = sendKey(m, tea.KeyEnter) // confirm

	assert.True(t, m.Done())
	val, ok := m.Answers()["plugins"].Value.([]string)
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"postgresql", "redis"}, val)
}

func TestModel_BackNavigation(t *testing.T) {
	questions := []prompt.Question{
		{ID: "q1", Type: prompt.QuestionTypeText, Prompt: "First", Default: "a"},
		{ID: "q2", Type: prompt.QuestionTypeText, Prompt: "Second", Default: "b"},
	}
	m := tui.NewModel(newSession(t, questions), "svc")

	m = sendKey(m, tea.KeyEnter) // accept default for q1
	assert.Equal(t, "a", m.Answers()["q1"].Value)

	// On q2 — press b with empty input to go back
	m = sendKeys(m, "b")
	assert.False(t, m.Done())
	_, hasQ1 := m.Answers()["q1"]
	assert.False(t, hasQ1)
}

func TestModel_HelpOverlay(t *testing.T) {
	questions := []prompt.Question{
		{ID: "q1", Type: prompt.QuestionTypeSelect, Prompt: "Pick one",
			Options: []prompt.Option{{Label: "A", Value: "a"}}},
	}
	m := tui.NewModel(newSession(t, questions), "svc")

	m = sendKeys(m, "?")
	assert.Contains(t, m.View(), "Keyboard Shortcuts")

	m = sendKeys(m, "x") // any key closes help
	assert.NotContains(t, m.View(), "Keyboard Shortcuts")
}

func TestModel_ThemeCycle(t *testing.T) {
	questions := []prompt.Question{
		{ID: "q1", Type: prompt.QuestionTypeSelect, Prompt: "Pick",
			Options: []prompt.Option{{Label: "A", Value: "a"}}},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	m = sendKeys(m, "t")
	assert.False(t, m.Done())
}

func TestModel_WindowResize(t *testing.T) {
	questions := []prompt.Question{
		{ID: "q1", Type: prompt.QuestionTypeText, Prompt: "Name", Default: "x"},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(tui.Model)
	assert.False(t, m.Done())
}

func TestModel_ViewContainsServiceName(t *testing.T) {
	questions := []prompt.Question{
		{ID: "q1", Type: prompt.QuestionTypeText, Prompt: "Team", Default: "platform"},
	}
	m := tui.NewModel(newSession(t, questions), "payment-service")
	assert.Contains(t, m.View(), "payment-service")
}

func TestModel_SummaryShownAfterFirstAnswer(t *testing.T) {
	questions := []prompt.Question{
		{ID: "team", Type: prompt.QuestionTypeText, Prompt: "Team", Default: "platform"},
		{ID: "slo", Type: prompt.QuestionTypeText, Prompt: "SLO", Default: "99.9"},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	m = sendKey(m, tea.KeyEnter) // answer first question

	view := m.View()
	assert.Contains(t, view, "answers")
	assert.Contains(t, view, "Team")
}

func TestNextTheme_Cycles(t *testing.T) {
	// Cycle through all 6 themes and confirm we return to the start.
	start := tui.LoadThemePreference()
	names := tui.ThemeNames()
	cur := start
	for range names {
		cur = tui.NextTheme(&cur)
	}
	assert.Equal(t, start.Name, cur.Name)
}

func TestModel_ProgressIndicator(t *testing.T) {
	questions := []prompt.Question{
		{ID: "q1", Type: prompt.QuestionTypeText, Prompt: "First", Default: "a"},
		{ID: "q2", Type: prompt.QuestionTypeText, Prompt: "Second", Default: "b"},
	}
	m := tui.NewModel(newSession(t, questions), "svc")

	// Before any answer: (1/2)
	assert.Contains(t, m.View(), "1/2")

	m = sendKey(m, tea.KeyEnter) // answer q1

	// After first answer: (2/2)
	assert.Contains(t, m.View(), "2/2")
}

func TestModel_BackspaceNavigatesBack_Confirm(t *testing.T) {
	questions := []prompt.Question{
		{ID: "q1", Type: prompt.QuestionTypeText, Prompt: "First", Default: "a"},
		{ID: "q2", Type: prompt.QuestionTypeConfirm, Prompt: "Confirm?", Default: "y"},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	m = sendKey(m, tea.KeyEnter) // answer q1 with default

	// Now on q2 (confirm) — send backspace to go back
	m = sendKey(m, tea.KeyBackspace)
	assert.False(t, m.Done())
	_, hasQ1 := m.Answers()["q1"]
	assert.False(t, hasQ1)
}

func TestModel_BackspaceNavigatesBack_Select(t *testing.T) {
	questions := []prompt.Question{
		{ID: "q1", Type: prompt.QuestionTypeText, Prompt: "First", Default: "a"},
		{
			ID: "q2", Type: prompt.QuestionTypeSelect, Prompt: "Pick",
			Options: []prompt.Option{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}},
		},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	m = sendKey(m, tea.KeyEnter) // answer q1

	m = sendKey(m, tea.KeyBackspace) // go back from select
	assert.False(t, m.Done())
	_, hasQ1 := m.Answers()["q1"]
	assert.False(t, hasQ1)
}

func TestModel_AnswersPanel_ScrollDown(t *testing.T) {
	// Build enough questions so the answers panel overflows.
	questions := make([]prompt.Question, 10)
	for i := range questions {
		questions[i] = prompt.Question{
			ID:      fmt.Sprintf("q%d", i),
			Type:    prompt.QuestionTypeText,
			Prompt:  fmt.Sprintf("Question %d", i),
			Default: fmt.Sprintf("val%d", i),
		}
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	// Answer all questions with defaults.
	for range questions {
		m = sendKey(m, tea.KeyEnter)
	}
	// Now in review mode — go back to question mode by editing.
	// Instead, just test scroll on a model mid-session.
	m2 := tui.NewModel(newSession(t, questions), "svc")
	for i := 0; i < 5; i++ {
		m2 = sendKey(m2, tea.KeyEnter)
	}

	// Scroll down via ctrl+down.
	updated, _ := m2.Update(tea.KeyMsg{Type: tea.KeyCtrlDown})
	m2 = updated.(tui.Model)
	// View should still render without panic.
	view := m2.View()
	assert.NotEmpty(t, view)
}

func TestModel_AnswersPanel_ScrollClampsAtZero(t *testing.T) {
	questions := []prompt.Question{
		{ID: "q1", Type: prompt.QuestionTypeText, Prompt: "Q1", Default: "a"},
	}
	m := tui.NewModel(newSession(t, questions), "svc")
	// Scroll up when already at top — should not panic.
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlUp})
	m = updated.(tui.Model)
	assert.NotEmpty(t, m.View())
}
