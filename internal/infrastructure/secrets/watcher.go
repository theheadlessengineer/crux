package secrets

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"
)

// NotifyFunc is called when a secret value changes during rotation.
// key is the secret identifier; components use this to reload affected resources.
type NotifyFunc func(key string)

// Watcher polls the secrets backend on a configurable interval and calls
// registered NotifyFuncs when a secret value changes.
// It integrates with the shutdown.Runner via its Stop method.
type Watcher struct {
	client   Client
	keys     []string
	interval time.Duration
	notify   []NotifyFunc
	logger   *slog.Logger
	stop     chan struct{}
	done     chan struct{}
}

// NewWatcher returns a Watcher. The polling interval is read from
// SECRETS_REFRESH_INTERVAL_SECONDS (default: 300). logger must not be nil.
func NewWatcher(client Client, keys []string, logger *slog.Logger) *Watcher {
	return &Watcher{
		client:   client,
		keys:     keys,
		interval: parseRefreshInterval(),
		logger:   logger,
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
}

// Register appends a NotifyFunc that is called when any watched secret rotates.
func (w *Watcher) Register(fn NotifyFunc) {
	w.notify = append(w.notify, fn)
}

// Start launches the background polling goroutine with the provided initial values.
// It returns immediately; call Stop to shut down cleanly.
func (w *Watcher) Start(ctx context.Context, initial *Config) {
	current := make(map[string]string, len(initial.Values))
	for k, v := range initial.Values {
		current[k] = v
	}

	go func() {
		defer close(w.done)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-w.stop:
				return
			case <-ticker.C:
				w.poll(ctx, current)
			}
		}
	}()
}

// Stop signals the watcher to stop and waits for the goroutine to exit.
// Satisfies the shutdown.Hook signature when wrapped: func(ctx) error { w.Stop(); return nil }.
func (w *Watcher) Stop() {
	close(w.stop)
	<-w.done
}

func (w *Watcher) poll(ctx context.Context, current map[string]string) {
	for _, key := range w.keys {
		val, err := w.client.Get(ctx, key)
		if err != nil {
			w.logger.Warn("secrets watcher: failed to fetch secret",
				slog.String("key", key),
				slog.Any("error", err),
			)
			continue
		}
		if val != current[key] {
			current[key] = val
			// Log rotation detection — never log the value itself.
			w.logger.Info("secrets watcher: secret rotated",
				slog.String("key", key),
			)
			for _, fn := range w.notify {
				fn(key)
			}
		}
	}
}

func parseRefreshInterval() time.Duration {
	if s := os.Getenv("SECRETS_REFRESH_INTERVAL_SECONDS"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
	}
	return 300 * time.Second
}
