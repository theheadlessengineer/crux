// Package testutil provides shared testing utilities and mock helpers.
package testutil

// MockExecutor is a mock implementation for command execution in tests.
type MockExecutor struct {
	ExecuteFunc func(cmd string, args ...string) (string, error)
}

// Execute calls the mock function if set, otherwise returns empty string.
func (m *MockExecutor) Execute(cmd string, args ...string) (string, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(cmd, args...)
	}
	return "", nil
}
