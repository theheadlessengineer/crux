package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theheadlessengineer/crux/internal/domain/health"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupRouter(h *HealthHandler) *gin.Engine {
	r := gin.New()
	h.RegisterRoutes(&r.RouterGroup)
	return r
}

func newHandler(checks map[string]health.CheckFunc) *HealthHandler {
	reg := health.NewRegistry()
	for name, fn := range checks {
		reg.Register(name, fn)
	}
	return NewHealthHandler(reg, BuildInfo{
		Service:   "test-service",
		Version:   "1.0.0",
		Commit:    "abc1234",
		BuildTime: "2026-03-15T00:00:00Z",
	})
}

func TestHealth(t *testing.T) {
	tests := []struct {
		name           string
		checks         map[string]health.CheckFunc
		wantStatus     int
		wantBodyStatus health.Status
	}{
		{
			name:           "all healthy",
			checks:         map[string]health.CheckFunc{"db": func() health.Status { return health.StatusHealthy }},
			wantStatus:     http.StatusOK,
			wantBodyStatus: health.StatusHealthy,
		},
		{
			name:           "one degraded",
			checks:         map[string]health.CheckFunc{"db": func() health.Status { return health.StatusDegraded }},
			wantStatus:     http.StatusOK,
			wantBodyStatus: health.StatusDegraded,
		},
		{
			name:           "no checks registered",
			checks:         nil,
			wantStatus:     http.StatusOK,
			wantBodyStatus: health.StatusHealthy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupRouter(newHandler(tt.checks))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", http.NoBody))

			assert.Equal(t, tt.wantStatus, w.Code)
			var body map[string]any
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
			assert.Equal(t, string(tt.wantBodyStatus), body["status"])
		})
	}
}

func TestReady(t *testing.T) {
	tests := []struct {
		name       string
		checks     map[string]health.CheckFunc
		wantStatus int
	}{
		{
			name:       "all ready",
			checks:     map[string]health.CheckFunc{"db": func() health.Status { return health.StatusHealthy }},
			wantStatus: http.StatusOK,
		},
		{
			name:       "dependency not ready",
			checks:     map[string]health.CheckFunc{"db": func() health.Status { return health.StatusUnhealthy }},
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name:       "no dependencies registered — always ready",
			checks:     nil,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupRouter(newHandler(tt.checks))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/ready", http.NoBody))
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestLive(t *testing.T) {
	r := setupRouter(newHandler(nil))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/live", http.NoBody))

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "healthy", body["status"])
}

func TestMetrics(t *testing.T) {
	r := setupRouter(newHandler(nil))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVersion(t *testing.T) {
	r := setupRouter(newHandler(nil))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/version", http.NoBody))

	assert.Equal(t, http.StatusOK, w.Code)
	var info BuildInfo
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &info))
	assert.Equal(t, "test-service", info.Service)
	assert.Equal(t, "1.0.0", info.Version)
	assert.Equal(t, "abc1234", info.Commit)
	assert.Equal(t, "2026-03-15T00:00:00Z", info.BuildTime)
}

func TestEndpointsRequireNoAuth(t *testing.T) {
	r := setupRouter(newHandler(nil))
	paths := []string{"/health", "/ready", "/live", "/metrics", "/version"}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, http.NoBody)
			// Deliberately no Authorization header
			r.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusUnauthorized, w.Code)
			assert.NotEqual(t, http.StatusForbidden, w.Code)
		})
	}
}
