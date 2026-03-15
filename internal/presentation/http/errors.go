package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

const contentTypeProblemJSON = "application/problem+json"

// Problem is an RFC 7807 Problem Details object with the company trace_id extension.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
	TraceID  string `json:"trace_id"`
}

// writeProblem renders p as application/problem+json with the correct status code.
func writeProblem(c *gin.Context, p *Problem) {
	c.Header("Content-Type", contentTypeProblemJSON)
	c.JSON(p.Status, p)
}

func traceID(c *gin.Context) string {
	sc := trace.SpanFromContext(c.Request.Context()).SpanContext()
	if sc.HasTraceID() {
		return sc.TraceID().String()
	}
	return ""
}

// ErrorMiddleware recovers from panics and returns a 500 Problem response.
// Stack traces are never included in the response body.
func ErrorMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		writeProblem(c, &Problem{
			Type:     "/errors/internal",
			Title:    "Internal Server Error",
			Status:   http.StatusInternalServerError,
			Detail:   "An unexpected error occurred.",
			Instance: c.Request.URL.Path,
			TraceID:  traceID(c),
		})
	})
}

// NotFound returns a 404 Problem response.
func NotFound(c *gin.Context, detail string) {
	writeProblem(c, &Problem{
		Type:     "/errors/not-found",
		Title:    "Not Found",
		Status:   http.StatusNotFound,
		Detail:   detail,
		Instance: c.Request.URL.Path,
		TraceID:  traceID(c),
	})
}

// ValidationError returns a 400 Problem response.
func ValidationError(c *gin.Context, detail string) {
	writeProblem(c, &Problem{
		Type:     "/errors/validation",
		Title:    "Bad Request",
		Status:   http.StatusBadRequest,
		Detail:   detail,
		Instance: c.Request.URL.Path,
		TraceID:  traceID(c),
	})
}

// Unauthorized returns a 401 Problem response.
func Unauthorized(c *gin.Context, detail string) {
	writeProblem(c, &Problem{
		Type:     "/errors/unauthorized",
		Title:    "Unauthorized",
		Status:   http.StatusUnauthorized,
		Detail:   detail,
		Instance: c.Request.URL.Path,
		TraceID:  traceID(c),
	})
}

// InternalError returns a 500 Problem response.
func InternalError(c *gin.Context, detail string) {
	writeProblem(c, &Problem{
		Type:     "/errors/internal",
		Title:    "Internal Server Error",
		Status:   http.StatusInternalServerError,
		Detail:   detail,
		Instance: c.Request.URL.Path,
		TraceID:  traceID(c),
	})
}
