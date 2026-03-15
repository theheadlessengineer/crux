package http

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const defaultMaxBodyBytes = 1 << 20 // 1 MB

// InputSanitization returns a Gin middleware that enforces baseline input
// sanitization on every inbound request:
//   - Rejects path traversal sequences (../ and ..%2F) in the URL path
//   - Rejects null bytes in the URL path and query string
//   - Enforces a maximum request body size (MAX_REQUEST_BODY_BYTES, default 1 MB)
//   - Requires Content-Type on POST, PUT, and PATCH requests
//
// All rejections use RFC 7807 format via the package-level error helpers.
func InputSanitization() gin.HandlerFunc {
	maxBytes := parseMaxBodyBytes()

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Path traversal check.
		if strings.Contains(path, "../") || strings.Contains(path, "..%2F") ||
			strings.Contains(path, "..%2f") {
			ValidationError(c, "path traversal sequences are not permitted")
			c.Abort()
			return
		}

		// Null byte check in parsed query parameter values.
		for _, vals := range c.Request.URL.Query() {
			for _, v := range vals {
				if strings.ContainsRune(v, 0) {
					ValidationError(c, "null bytes are not permitted in the request")
					c.Abort()
					return
				}
			}
		}

		// Body size limit.
		if c.Request.ContentLength > maxBytes {
			c.Header("Content-Type", contentTypeProblemJSON)
			c.JSON(http.StatusRequestEntityTooLarge, &Problem{
				Type:     "/errors/payload-too-large",
				Title:    "Payload Too Large",
				Status:   http.StatusRequestEntityTooLarge,
				Detail:   "Request body exceeds the maximum allowed size.",
				Instance: c.Request.URL.Path,
				TraceID:  traceID(c),
			})
			c.Abort()
			return
		}

		// Content-Type required on mutating methods.
		method := c.Request.Method
		if (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) &&
			c.GetHeader("Content-Type") == "" {
			c.Header("Content-Type", contentTypeProblemJSON)
			c.JSON(http.StatusUnsupportedMediaType, &Problem{
				Type:     "/errors/unsupported-media-type",
				Title:    "Unsupported Media Type",
				Status:   http.StatusUnsupportedMediaType,
				Detail:   "Content-Type header is required for this request.",
				Instance: c.Request.URL.Path,
				TraceID:  traceID(c),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func parseMaxBodyBytes() int64 {
	if s := os.Getenv("MAX_REQUEST_BODY_BYTES"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil && n > 0 {
			return n
		}
	}
	return defaultMaxBodyBytes
}
