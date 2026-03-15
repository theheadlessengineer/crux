package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func securityRouter(middleware ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	for _, m := range middleware {
		r.Use(m)
	}
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestSecurityHeaders_AllHeadersPresent(t *testing.T) {
	w := httptest.NewRecorder()
	securityRouter(SecurityHeaders()).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	want := map[string]string{
		"Content-Security-Policy":   "default-src 'self'",
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
		"Permissions-Policy":        "geolocation=(), microphone=(), camera=()",
	}
	for header, value := range want {
		assert.Equal(t, value, w.Header().Get(header), "header %s", header)
	}
}

func TestCORSMiddleware_DeniedByDefault(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("Origin", "https://evil.example.com")
	securityRouter(CORSMiddleware()).ServeHTTP(w, req)

	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com")
	t.Setenv("CORS_ALLOWED_METHODS", "")
	t.Setenv("CORS_ALLOWED_HEADERS", "")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("Origin", "https://app.example.com")
	securityRouter(CORSMiddleware()).ServeHTTP(w, req)

	assert.Equal(t, "https://app.example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "Origin", w.Header().Get("Vary"))
}

func TestCORSMiddleware_UnlistedOriginDenied(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("Origin", "https://other.example.com")
	securityRouter(CORSMiddleware()).ServeHTTP(w, req)

	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_PreflightReturns204(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/test", http.NoBody)
	req.Header.Set("Origin", "https://app.example.com")
	securityRouter(CORSMiddleware()).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestCORSMiddleware_NoOriginHeader_PassesThrough(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com")

	w := httptest.NewRecorder()
	securityRouter(CORSMiddleware()).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}
