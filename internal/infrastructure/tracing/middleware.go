package tracing

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Middleware returns a Gin middleware that:
//  1. Extracts the W3C traceparent header from the inbound request and
//     restores the remote span context into the request context.
//  2. Starts a new server span for the request (or a root span when no
//     traceparent is present).
//
// The span is ended when the handler chain returns.
// Downstream handlers and the logger retrieve trace/span IDs via
// trace.SpanFromContext(c.Request.Context()).
func Middleware(tracer trace.Tracer) gin.HandlerFunc {
	propagator := otel.GetTextMapPropagator()
	return func(c *gin.Context) {
		// Extract remote context from the incoming traceparent header.
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// Start a server span for this request.
		ctx, span := tracer.Start(ctx, c.FullPath())
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
