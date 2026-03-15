package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// problemResponse decodes the response body into a Problem struct.
func problemResponse(t *testing.T, w *httptest.ResponseRecorder) Problem {
	t.Helper()
	var p Problem
	require.NoError(t, json.NewDecoder(w.Body).Decode(&p))
	return p
}

// assertProblem checks the common invariants every RFC 7807 response must satisfy.
func assertProblem(t *testing.T, w *httptest.ResponseRecorder, wantStatus int, wantType string) Problem {
	t.Helper()
	assert.Equal(t, wantStatus, w.Code)
	assert.Equal(t, contentTypeProblemJSON, w.Header().Get("Content-Type"))

	p := problemResponse(t, w)
	assert.Equal(t, wantType, p.Type)
	assert.Equal(t, wantStatus, p.Status)
	assert.NotEmpty(t, p.Title)
	assert.NotEmpty(t, p.Detail)
	assert.NotEmpty(t, p.Instance)
	return p
}

func newRouter(handler gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(ErrorMiddleware())
	r.GET("/test", handler)
	return r
}

func TestValidationError(t *testing.T) {
	r := newRouter(func(c *gin.Context) {
		ValidationError(c, "name is required")
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	assertProblem(t, w, http.StatusBadRequest, "/errors/validation")
}

func TestNotFound(t *testing.T) {
	r := newRouter(func(c *gin.Context) {
		NotFound(c, "resource not found")
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	assertProblem(t, w, http.StatusNotFound, "/errors/not-found")
}

func TestUnauthorized(t *testing.T) {
	r := newRouter(func(c *gin.Context) {
		Unauthorized(c, "missing token")
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	assertProblem(t, w, http.StatusUnauthorized, "/errors/unauthorized")
}

func TestInternalError(t *testing.T) {
	r := newRouter(func(c *gin.Context) {
		InternalError(c, "something went wrong")
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	assertProblem(t, w, http.StatusInternalServerError, "/errors/internal")
}

func TestPanicRecovery_Returns500_NoStackTrace(t *testing.T) {
	r := newRouter(func(_ *gin.Context) {
		panic("deliberate panic")
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	p := assertProblem(t, w, http.StatusInternalServerError, "/errors/internal")
	// Stack trace must never appear in the response body.
	assert.NotContains(t, p.Detail, "goroutine")
	assert.NotContains(t, p.Detail, "panic")
}

func TestProblem_InstanceIsRequestPath(t *testing.T) {
	r := gin.New()
	r.Use(ErrorMiddleware())
	r.GET("/payments/:id", func(c *gin.Context) {
		NotFound(c, "payment not found")
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/payments/42", http.NoBody))

	p := problemResponse(t, w)
	assert.Equal(t, "/payments/42", p.Instance)
}

func TestProblem_ContentTypeOnAllErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler gin.HandlerFunc
	}{
		{"validation", func(c *gin.Context) { ValidationError(c, "bad") }},
		{"not-found", func(c *gin.Context) { NotFound(c, "missing") }},
		{"unauthorized", func(c *gin.Context) { Unauthorized(c, "denied") }},
		{"internal", func(c *gin.Context) { InternalError(c, "oops") }},
		{"panic", func(_ *gin.Context) { panic("boom") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			newRouter(tt.handler).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))
			assert.Equal(t, contentTypeProblemJSON, w.Header().Get("Content-Type"))
		})
	}
}
