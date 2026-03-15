package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry_Ready(t *testing.T) {
	tests := []struct {
		name      string
		checks    map[string]CheckFunc
		wantReady bool
	}{
		{
			name:      "no checks — ready",
			checks:    nil,
			wantReady: true,
		},
		{
			name:      "all healthy — ready",
			checks:    map[string]CheckFunc{"db": func() Status { return StatusHealthy }},
			wantReady: true,
		},
		{
			name:      "one unhealthy — not ready",
			checks:    map[string]CheckFunc{"db": func() Status { return StatusUnhealthy }},
			wantReady: false,
		},
		{
			name:      "one degraded — not ready",
			checks:    map[string]CheckFunc{"cache": func() Status { return StatusDegraded }},
			wantReady: false,
		},
		{
			name: "mixed — not ready",
			checks: map[string]CheckFunc{
				"db":    func() Status { return StatusHealthy },
				"cache": func() Status { return StatusUnhealthy },
			},
			wantReady: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := NewRegistry()
			for name, fn := range tt.checks {
				reg.Register(name, fn)
			}
			assert.Equal(t, tt.wantReady, reg.Ready())
		})
	}
}

func TestRegistry_Results(t *testing.T) {
	reg := NewRegistry()
	reg.Register("db", func() Status { return StatusHealthy })
	reg.Register("cache", func() Status { return StatusDegraded })

	results := reg.Results()
	assert.Equal(t, StatusHealthy, results["db"])
	assert.Equal(t, StatusDegraded, results["cache"])
}
