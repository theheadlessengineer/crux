package prompt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/domain/prompt"
)

func textQ(def, pattern string, required bool) prompt.Question {
	return prompt.Question{
		ID: "q", Type: prompt.QuestionTypeText, Default: def,
		Validation: prompt.ValidationRule{Required: required, Pattern: pattern},
	}
}
func confirmQ(def string) prompt.Question {
	return prompt.Question{ID: "q", Type: prompt.QuestionTypeConfirm, Default: def}
}
func numberQ(lo, hi float64) prompt.Question {
	return prompt.Question{
		ID: "q", Type: prompt.QuestionTypeNumber,
		Validation: prompt.ValidationRule{Min: lo, Max: hi},
	}
}
func selectQ(opts ...string) prompt.Question {
	options := make([]prompt.Option, len(opts))
	for i, o := range opts {
		options[i] = prompt.Option{Label: o, Value: o}
	}
	return prompt.Question{ID: "q", Type: prompt.QuestionTypeSelect, Options: options}
}
func multiQ(opts ...string) prompt.Question {
	options := make([]prompt.Option, len(opts))
	for i, o := range opts {
		options[i] = prompt.Option{Label: o, Value: o}
	}
	return prompt.Question{ID: "q", Type: prompt.QuestionTypeMultiSelect, Options: options}
}

// --- confirm ---

func TestValidate_Confirm(t *testing.T) {
	q := confirmQ("n")
	tests := []struct {
		raw  string
		want bool
	}{
		{"y", true}, {"yes", true}, {"Y", true}, {"true", true}, {"1", true},
		{"n", false}, {"no", false}, {"false", false}, {"0", false},
	}
	for _, tt := range tests {
		a, err := prompt.Validate(&q, tt.raw)
		require.NoError(t, err, "raw=%q", tt.raw)
		assert.Equal(t, tt.want, a.Value)
	}
}

func TestValidate_Confirm_Default(t *testing.T) {
	q := confirmQ("y")
	a, err := prompt.Validate(&q, "")
	require.NoError(t, err)
	assert.Equal(t, true, a.Value)
}

func TestValidate_Confirm_Invalid(t *testing.T) {
	q := confirmQ("")
	_, err := prompt.Validate(&q, "maybe")
	assert.Error(t, err)
}

// --- text ---

func TestValidate_Text_Valid(t *testing.T) {
	q := textQ("", "", false)
	a, err := prompt.Validate(&q, "hello")
	require.NoError(t, err)
	assert.Equal(t, "hello", a.Value)
}

func TestValidate_Text_Default(t *testing.T) {
	q := textQ("world", "", false)
	a, err := prompt.Validate(&q, "")
	require.NoError(t, err)
	assert.Equal(t, "world", a.Value)
}

func TestValidate_Text_Required(t *testing.T) {
	q := textQ("", "", true)
	_, err := prompt.Validate(&q, "")
	assert.Error(t, err)
}

func TestValidate_Text_Pattern(t *testing.T) {
	q := textQ("", `^[a-z]+$`, false)
	_, err := prompt.Validate(&q, "UPPER")
	assert.Error(t, err)

	a, err := prompt.Validate(&q, "lower")
	require.NoError(t, err)
	assert.Equal(t, "lower", a.Value)
}

// --- number ---

func TestValidate_Number_Valid(t *testing.T) {
	q := numberQ(0, 100)
	a, err := prompt.Validate(&q, "42")
	require.NoError(t, err)
	assert.Equal(t, float64(42), a.Value)
}

func TestValidate_Number_Default(t *testing.T) {
	q := numberQ(0, 0)
	a, err := prompt.Validate(&q, "")
	require.NoError(t, err)
	assert.Equal(t, float64(0), a.Value)
}

func TestValidate_Number_BelowMin(t *testing.T) {
	q := numberQ(10, 100)
	_, err := prompt.Validate(&q, "5")
	assert.Error(t, err)
}

func TestValidate_Number_AboveMax(t *testing.T) {
	q := numberQ(0, 10)
	_, err := prompt.Validate(&q, "99")
	assert.Error(t, err)
}

func TestValidate_Number_NotANumber(t *testing.T) {
	q := numberQ(0, 0)
	_, err := prompt.Validate(&q, "abc")
	assert.Error(t, err)
}

// --- select ---

func TestValidate_Select_Valid(t *testing.T) {
	q := selectQ("go", "python", "java")
	a, err := prompt.Validate(&q, "go")
	require.NoError(t, err)
	assert.Equal(t, "go", a.Value)
}

func TestValidate_Select_Invalid(t *testing.T) {
	q := selectQ("go", "python")
	_, err := prompt.Validate(&q, "rust")
	assert.Error(t, err)
}

func TestValidate_Select_Default(t *testing.T) {
	q := selectQ("go", "python")
	q.Default = "go"
	a, err := prompt.Validate(&q, "")
	require.NoError(t, err)
	assert.Equal(t, "go", a.Value)
}

// --- multiselect ---

func TestValidate_MultiSelect_Valid(t *testing.T) {
	q := multiQ("a", "b", "c")
	a, err := prompt.Validate(&q, "a,c")
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "c"}, a.Value)
}

func TestValidate_MultiSelect_Empty(t *testing.T) {
	q := multiQ("a", "b")
	a, err := prompt.Validate(&q, "")
	require.NoError(t, err)
	assert.Equal(t, []string{}, a.Value)
}

func TestValidate_MultiSelect_Invalid(t *testing.T) {
	q := multiQ("a", "b")
	_, err := prompt.Validate(&q, "a,z")
	assert.Error(t, err)
}
