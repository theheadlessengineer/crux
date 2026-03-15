package secrets

import (
	"bytes"
	"context"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// rotatingClient simulates a secret that changes value after firstCallsDone calls.
type rotatingClient struct {
	mu             sync.Mutex
	calls          int
	firstCallsDone int
	initial        string
	rotated        string
}

func (r *rotatingClient) Get(_ context.Context, _ string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls++
	if r.calls > r.firstCallsDone {
		return r.rotated, nil
	}
	return r.initial, nil
}

func testWatcherLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestWatcher_DetectsRotation_CallsNotify(t *testing.T) {
	var buf bytes.Buffer
	client := &rotatingClient{firstCallsDone: 1, initial: "old", rotated: "new"}

	w := NewWatcher(client, []string{"db/password"}, testWatcherLogger(&buf))
	w.interval = 10 * time.Millisecond

	var notified []string
	var mu sync.Mutex
	w.Register(func(key string) {
		mu.Lock()
		notified = append(notified, key)
		mu.Unlock()
	})

	initial := &Config{Values: map[string]string{"db/password": "old"}}
	w.Start(context.Background(), initial)

	require.Eventually(t, func() bool {
		mu.Lock()
		defer mu.Unlock()
		return len(notified) > 0
	}, 500*time.Millisecond, 10*time.Millisecond)

	w.Stop()

	mu.Lock()
	defer mu.Unlock()
	assert.Contains(t, notified, "db/password")
	assert.Contains(t, buf.String(), "secret rotated")
	// Secret value must never appear in logs.
	assert.NotContains(t, buf.String(), "new")
	assert.NotContains(t, buf.String(), "old")
}

func TestWatcher_NoChange_DoesNotNotify(t *testing.T) {
	var buf bytes.Buffer
	client := &mockClient{data: map[string]string{"db/password": "stable"}}

	w := NewWatcher(client, []string{"db/password"}, testWatcherLogger(&buf))
	w.interval = 10 * time.Millisecond

	notified := false
	w.Register(func(_ string) { notified = true })

	initial := &Config{Values: map[string]string{"db/password": "stable"}}
	w.Start(context.Background(), initial)

	time.Sleep(50 * time.Millisecond)
	w.Stop()

	assert.False(t, notified, "notify must not be called when value is unchanged")
}

func TestWatcher_Stop_ShutsDownCleanly(t *testing.T) {
	var buf bytes.Buffer
	client := &mockClient{data: map[string]string{}}

	w := NewWatcher(client, nil, testWatcherLogger(&buf))
	w.interval = 10 * time.Millisecond

	w.Start(context.Background(), &Config{Values: map[string]string{}})

	done := make(chan struct{})
	go func() {
		w.Stop()
		close(done)
	}()

	select {
	case <-done:
		// clean shutdown
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watcher did not stop within timeout")
	}
}

func TestParseRefreshInterval_Default(t *testing.T) {
	t.Setenv("SECRETS_REFRESH_INTERVAL_SECONDS", "")
	assert.Equal(t, 300*time.Second, parseRefreshInterval())
}

func TestParseRefreshInterval_EnvOverride(t *testing.T) {
	t.Setenv("SECRETS_REFRESH_INTERVAL_SECONDS", "60")
	assert.Equal(t, 60*time.Second, parseRefreshInterval())
}

func TestParseRefreshInterval_InvalidEnv_UsesDefault(t *testing.T) {
	t.Setenv("SECRETS_REFRESH_INTERVAL_SECONDS", "bad")
	assert.Equal(t, 300*time.Second, parseRefreshInterval())
}
