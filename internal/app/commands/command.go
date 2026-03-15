// Package commands defines the Command interface and shared types for all CLI commands.
package commands

import "context"

// Command is the interface all CLI commands must implement.
type Command interface {
	Execute(ctx context.Context, args []string) error
	Validate() error
}
