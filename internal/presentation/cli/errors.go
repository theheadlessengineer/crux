package cli

// ValidationError signals exit code 2 (validation failure).
type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string { return e.Msg }

// exitError signals a non-zero exit code with a message.
type exitError struct {
	code int
	msg  string
}

func (e *exitError) Error() string { return e.msg }
