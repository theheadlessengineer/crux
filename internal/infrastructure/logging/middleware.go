package logging

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// Middleware returns a Gin middleware that extracts trace_id and span_id from
// the active OTel span in the request context and stores them as logging Fields.
// It must be registered after the tracing.Middleware so the span is already present.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		sc := span.SpanContext()

		f := Fields{}
		if sc.HasTraceID() {
			f.TraceID = sc.TraceID().String()
		}
		if sc.HasSpanID() {
			f.SpanID = sc.SpanID().String()
		}

		ctx := ContextWithFields(c.Request.Context(), f)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
