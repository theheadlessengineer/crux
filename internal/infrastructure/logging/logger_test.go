package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) }

// logLine unmarshals one JSON log line from buf.
func logLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var m map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &m))
	return m
}

func TestLogger_JSONStructure(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, "test-service", "test", "1.0.0")
	l.inner.Info("hello")

	line := logLine(t, &buf)
	assert.Equal(t, "INFO", line["level"])
	assert.Equal(t, "hello", line["msg"])
	assert.Equal(t, "test-service", line["service"])
	assert.Equal(t, "test", line["environment"])
	assert.Equal(t, "1.0.0", line["version"])
	assert.NotEmpty(t, line["timestamp"])
}

func TestLogger_WithContext_TraceFields(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, "svc", "dev", "0.1.0")

	ctx := ContextWithFields(context.Background(), Fields{
		TraceID: "abc123",
		SpanID:  "def456",
	})
	l.WithContext(ctx).Info("traced request")

	line := logLine(t, &buf)
	assert.Equal(t, "abc123", line["trace_id"])
	assert.Equal(t, "def456", line["span_id"])
}

func TestLogger_WithContext_NoTrace_EmptyStrings(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, "svc", "dev", "0.1.0")

	l.WithContext(context.Background()).Info("no trace")

	line := logLine(t, &buf)
	assert.Equal(t, "", line["trace_id"], "trace_id must be empty string, not null")
	assert.Equal(t, "", line["span_id"], "span_id must be empty string, not null")
}

func TestLogger_LevelFromEnv(t *testing.T) {
	tests := []struct {
		envVal    string
		wantLevel slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
		{"info", slog.LevelInfo},
		{"", slog.LevelInfo},
		{"invalid", slog.LevelInfo},
	}
	for _, tt := range tests {
		t.Run(tt.envVal, func(t *testing.T) {
			assert.Equal(t, tt.wantLevel, parseLevel(tt.envVal))
		})
	}
}

func TestMiddleware_InjectsTraceFields(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, "svc", "dev", "0.1.0")

	r := gin.New()
	r.Use(Middleware())
	r.GET("/test", func(c *gin.Context) {
		l.WithContext(c.Request.Context()).Info("in handler")
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("X-Trace-Id", "trace-abc")
	req.Header.Set("X-Span-Id", "span-xyz")
	r.ServeHTTP(httptest.NewRecorder(), req)

	line := logLine(t, &buf)
	assert.Equal(t, "trace-abc", line["trace_id"])
	assert.Equal(t, "span-xyz", line["span_id"])
}

func TestMiddleware_NoHeaders_EmptyStrings(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, "svc", "dev", "0.1.0")

	r := gin.New()
	r.Use(Middleware())
	r.GET("/test", func(c *gin.Context) {
		l.WithContext(c.Request.Context()).Info("no headers")
		c.Status(http.StatusOK)
	})

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	line := logLine(t, &buf)
	assert.Equal(t, "", line["trace_id"])
	assert.Equal(t, "", line["span_id"])
}
