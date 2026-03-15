package logging

import (
	"github.com/gin-gonic/gin"
)

// Middleware returns a Gin middleware that extracts trace_id and span_id from
// request headers and stores them in the request context via ContextWithFields.
// Header names follow the OTel/W3C traceparent convention; US-1103 will wire
// the actual OTel SDK — this middleware provides the context plumbing now.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		f := Fields{
			TraceID: c.GetHeader("X-Trace-Id"),
			SpanID:  c.GetHeader("X-Span-Id"),
		}
		ctx := ContextWithFields(c.Request.Context(), f)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
