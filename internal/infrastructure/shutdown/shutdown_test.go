package shutdown

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

// sendSignal sends sig to the current process after delay.
func sendSignal(sig os.Signal, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(sig)
	}()
}

func TestListenAndServe_SIGTERM_RunsHooksAndReturnsNil(t *testing.T) {
	var buf bytes.Buffer
	r := New(testLogger(&buf))

	called := false
	r.Register(func(_ context.Context) error {
		called = true
		return nil
	})

	sendSignal(syscall.SIGTERM, 20*time.Millisecond)

	err := r.ListenAndServe()
	require.NoError(t, err)
	assert.True(t, called, "hook must be called on SIGTERM")
	assert.Contains(t, buf.String(), "shutdown signal received")
	assert.Contains(t, buf.String(), "shutdown complete")
}

func TestListenAndServe_SIGINT_RunsHooksAndReturnsNil(t *testing.T) {
	var buf bytes.Buffer
	r := New(testLogger(&buf))

	called := false
	r.Register(func(_ context.Context) error {
		called = true
		return nil
	})

	sendSignal(syscall.SIGINT, 20*time.Millisecond)

	err := r.ListenAndServe()
	require.NoError(t, err)
	assert.True(t, called)
}

func TestListenAndServe_HookError_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	r := New(testLogger(&buf))

	hookErr := errors.New("drain failed")
	r.Register(func(_ context.Context) error { return hookErr })

	sendSignal(syscall.SIGTERM, 20*time.Millisecond)

	err := r.ListenAndServe()
	assert.ErrorIs(t, err, hookErr)
}

func TestListenAndServe_HooksRunInOrder(t *testing.T) {
	var buf bytes.Buffer
	r := New(testLogger(&buf))

	var order []int
	r.Register(func(_ context.Context) error { order = append(order, 1); return nil })
	r.Register(func(_ context.Context) error { order = append(order, 2); return nil })
	r.Register(func(_ context.Context) error { order = append(order, 3); return nil })

	sendSignal(syscall.SIGTERM, 20*time.Millisecond)

	require.NoError(t, r.ListenAndServe())
	assert.Equal(t, []int{1, 2, 3}, order)
}

func TestListenAndServe_TimeoutExceeded_ReturnsContextError(t *testing.T) {
	var buf bytes.Buffer
	r := New(testLogger(&buf))
	r.timeout = 30 * time.Millisecond // override to a short timeout for the test

	r.Register(func(ctx context.Context) error {
		// Block until the context is cancelled.
		<-ctx.Done()
		return ctx.Err()
	})

	sendSignal(syscall.SIGTERM, 10*time.Millisecond)

	err := r.ListenAndServe()
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestParseDrainTimeout_Default(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "")
	assert.Equal(t, 30*time.Second, parseDrainTimeout())
}

func TestParseDrainTimeout_EnvOverride(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "10")
	assert.Equal(t, 10*time.Second, parseDrainTimeout())
}

func TestParseDrainTimeout_InvalidEnv_UsesDefault(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "notanumber")
	assert.Equal(t, 30*time.Second, parseDrainTimeout())
}
