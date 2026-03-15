package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func sanitizationRouter() *gin.Engine {
	r := gin.New()
	r.Use(InputSanitization())
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.POST("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.PUT("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.PATCH("/test", func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestInputSanitization_CleanRequest_PassesThrough(t *testing.T) {
	w := httptest.NewRecorder()
	sanitizationRouter().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInputSanitization_PathTraversal_DotDotSlash(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test/../etc/passwd", http.NoBody)
	sanitizationRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, contentTypeProblemJSON, w.Header().Get("Content-Type"))
}

func TestInputSanitization_PathTraversal_URLEncoded(t *testing.T) {
	w := httptest.NewRecorder()
	// Construct request with raw path containing encoded traversal.
	req := httptest.NewRequest(http.MethodGet, "/test/..%2Fetc%2Fpasswd", http.NoBody)
	sanitizationRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInputSanitization_NullByteInQueryValue(t *testing.T) {
	// Go's URL parser rejects null bytes in the raw URL, but they can appear
	// in percent-encoded form (%00) which decodes to a null byte in query values.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test?q=foo%00bar", http.NoBody)
	sanitizationRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, contentTypeProblemJSON, w.Header().Get("Content-Type"))
}

func TestInputSanitization_BodyTooLarge(t *testing.T) {
	t.Setenv("MAX_REQUEST_BODY_BYTES", "10")

	r := gin.New()
	r.Use(InputSanitization())
	r.POST("/test", func(c *gin.Context) { c.Status(http.StatusOK) })

	body := strings.NewReader("this body is definitely longer than ten bytes")
	req := httptest.NewRequest(http.MethodPost, "/test", body)
	req.ContentLength = int64(body.Len())
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	assert.Equal(t, contentTypeProblemJSON, w.Header().Get("Content-Type"))
}

func TestInputSanitization_MaxBodyBytes_EnvOverride(t *testing.T) {
	t.Setenv("MAX_REQUEST_BODY_BYTES", "512")
	assert.Equal(t, int64(512), parseMaxBodyBytes())
}

func TestInputSanitization_MaxBodyBytes_Default(t *testing.T) {
	t.Setenv("MAX_REQUEST_BODY_BYTES", "")
	assert.Equal(t, int64(defaultMaxBodyBytes), parseMaxBodyBytes())
}

func TestInputSanitization_MissingContentType_Post(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
	// No Content-Type header set.
	sanitizationRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
	assert.Equal(t, contentTypeProblemJSON, w.Header().Get("Content-Type"))
}

func TestInputSanitization_MissingContentType_Put(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/test", strings.NewReader("{}"))
	sanitizationRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

func TestInputSanitization_MissingContentType_Patch(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/test", strings.NewReader("{}"))
	sanitizationRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

func TestInputSanitization_WithContentType_Post_PassesThrough(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	sanitizationRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestInputSanitization_GetWithoutContentType_PassesThrough(t *testing.T) {
	w := httptest.NewRecorder()
	sanitizationRouter().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))
	assert.Equal(t, http.StatusOK, w.Code)
}
