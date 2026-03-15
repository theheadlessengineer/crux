package tracing

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Transport is an http.RoundTripper that injects the W3C traceparent header
// into every outbound request using the globally registered propagator.
type Transport struct {
	wrapped http.RoundTripper
}

// NewTransport returns a Transport wrapping base. If base is nil, http.DefaultTransport is used.
func NewTransport(base http.RoundTripper) *Transport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &Transport{wrapped: base}
}

// RoundTrip injects traceparent (and any other registered propagation headers)
// then delegates to the wrapped transport.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid mutating the caller's headers.
	r := req.Clone(req.Context())
	otel.GetTextMapPropagator().Inject(r.Context(), propagation.HeaderCarrier(r.Header))
	return t.wrapped.RoundTrip(r)
}

// NewHTTPClient returns an *http.Client whose transport injects traceparent on
// every outbound request. Pass nil to use http.DefaultTransport as the base.
func NewHTTPClient(base http.RoundTripper) *http.Client {
	return &http.Client{Transport: NewTransport(base)}
}
