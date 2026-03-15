package tracing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func init() { gin.SetMode(gin.TestMode) }

// newTestTracer returns a no-op TracerProvider and its Tracer for unit tests.
// Using sdktrace.NewTracerProvider() (no exporter) produces real span contexts
// with valid trace/span IDs without requiring a running collector.
func newTestTracer() (trace.Tracer, *sdktrace.TracerProvider) {
	tp := sdktrace.NewTracerProvider()
	tracer := tp.Tracer("test")
	return tracer, tp
}

// TestMiddleware_PropagatesInboundTraceparent verifies that a valid traceparent
// header on the inbound request is extracted and the span context is restored.
func TestMiddleware_PropagatesInboundTraceparent(t *testing.T) {
	tracer, tp := newTestTracer()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	var capturedSpanCtx trace.SpanContext

	r := gin.New()
	r.Use(Middleware(tracer))
	r.GET("/test", func(c *gin.Context) {
		capturedSpanCtx = trace.SpanFromContext(c.Request.Context()).SpanContext()
		c.Status(http.StatusOK)
	})

	// A valid W3C traceparent header.
	const traceID = "4bf92f3577b34da6a3ce929d0e0e4736"
	const parentSpanID = "00f067aa0ba902b7"
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("traceparent", "00-"+traceID+"-"+parentSpanID+"-01")
	r.ServeHTTP(httptest.NewRecorder(), req)

	assert.True(t, capturedSpanCtx.IsValid(), "span context must be valid")
	assert.Equal(t, traceID, capturedSpanCtx.TraceID().String(), "trace ID must match inbound traceparent")
}

// TestMiddleware_CreatesRootSpanWithoutTraceparent verifies that a new root
// span is created when no traceparent header is present.
func TestMiddleware_CreatesRootSpanWithoutTraceparent(t *testing.T) {
	tracer, tp := newTestTracer()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	var capturedSpanCtx trace.SpanContext

	r := gin.New()
	r.Use(Middleware(tracer))
	r.GET("/test", func(c *gin.Context) {
		capturedSpanCtx = trace.SpanFromContext(c.Request.Context()).SpanContext()
		c.Status(http.StatusOK)
	})

	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	assert.True(t, capturedSpanCtx.IsValid(), "a root span must be created when no traceparent is present")
}

// TestTransport_InjectsTraceparentHeader verifies that the outbound HTTP client
// transport injects a traceparent header derived from the active span context.
func TestTransport_InjectsTraceparentHeader(t *testing.T) {
	tracer, tp := newTestTracer()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	var capturedHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Get("traceparent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Start a span so there is an active trace context to inject.
	ctx, span := tracer.Start(t.Context(), "test-op")
	defer span.End()

	client := NewHTTPClient(nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, http.NoBody)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	assert.NotEmpty(t, capturedHeader, "traceparent header must be injected on outbound requests")
}

// TestTransport_NoActiveSpan_NoTraceparent verifies that no traceparent header
// is injected when there is no active span in the context.
func TestTransport_NoActiveSpan_NoTraceparent(t *testing.T) {
	otel.SetTextMapPropagator(propagation.TraceContext{})

	var capturedHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Get("traceparent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(nil)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, http.NoBody)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	assert.Empty(t, capturedHeader, "traceparent must not be injected when no span is active")
}
