package http

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders returns a Gin middleware that sets the hardened default
// security response headers on every response.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("Content-Security-Policy", "default-src 'self'")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Next()
	}
}

// CORSMiddleware returns a Gin middleware that applies CORS policy driven by
// environment variables. Cross-origin requests are denied by default when
// CORS_ALLOWED_ORIGINS is empty.
//
//	CORS_ALLOWED_ORIGINS  comma-separated origins  (default: deny all)
//	CORS_ALLOWED_METHODS  comma-separated methods  (default: GET,POST,PUT,DELETE)
//	CORS_ALLOWED_HEADERS  comma-separated headers  (default: Authorization,Content-Type)
func CORSMiddleware() gin.HandlerFunc {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	methods := envOr("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE")
	headers := envOr("CORS_ALLOWED_HEADERS", "Authorization,Content-Type")

	allowedOrigins := splitTrimmed(origins)

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && isAllowed(origin, allowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", methods)
			c.Header("Access-Control-Allow-Headers", headers)
			c.Header("Vary", "Origin")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitTrimmed(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func isAllowed(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == origin {
			return true
		}
	}
	return false
}
