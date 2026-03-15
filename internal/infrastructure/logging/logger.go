package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

// contextKey is the key type for values stored in context by this package.
type contextKey struct{}

// Fields holds the trace correlation fields injected per-request.
type Fields struct {
	TraceID string
	SpanID  string
}

// Logger wraps slog.Logger with the service-level fields pre-attached.
type Logger struct {
	inner *slog.Logger
}

// New returns a Logger that writes structured JSON to w.
// service, environment, and version are attached to every log line.
// Log level is read from the LOG_LEVEL environment variable (default: info).
func New(w io.Writer, service, environment, version string) *Logger {
	level := parseLevel(os.Getenv("LOG_LEVEL"))
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			// Rename slog's default "time" key to "timestamp" per log schema.
			if a.Key == slog.TimeKey {
				a.Key = "timestamp"
			}
			return a
		},
	})
	inner := slog.New(handler).With(
		slog.String("service", service),
		slog.String("environment", environment),
		slog.String("version", version),
	)
	return &Logger{inner: inner}
}

// WithContext returns a child logger with trace_id and span_id from ctx injected.
// If no Fields are present in ctx, both fields are set to empty string.
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	f, _ := ctx.Value(contextKey{}).(Fields)
	return l.inner.With(
		slog.String("trace_id", f.TraceID),
		slog.String("span_id", f.SpanID),
	)
}

// ContextWithFields stores trace correlation fields in ctx.
func ContextWithFields(ctx context.Context, f Fields) context.Context {
	return context.WithValue(ctx, contextKey{}, f)
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
