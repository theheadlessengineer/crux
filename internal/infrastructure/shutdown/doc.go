// Package shutdown provides a signal-aware graceful shutdown runner for
// generated services. It listens for SIGTERM and SIGINT, then executes
// registered hooks in registration order within a configurable drain timeout.
package shutdown
