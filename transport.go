package gh

import (
	"net/http"
	"net/http/httputil"
)

// Printf is a (fmt.Printf|log.Printf|testing.T.Logf)-like function.
type Printf func(format string, args ...any)

// NoopPrintf is a Printf function that does nothing.
func NoopPrintf(string, ...any) {}

// transport is an http.RoundTripper with debug logging.
type transport struct {
	t      http.RoundTripper
	debugf Printf
}

// NewTransport returns a new http.RoundTripper that wraps the source with debug logging.
// Both should be set.
func NewTransport(source http.RoundTripper, debugf Printf) http.RoundTripper {
	if source == nil {
		panic("source is nil")
	}
	if debugf == nil {
		panic("debugf is nil")
	}

	return &transport{
		t:      source,
		debugf: debugf,
	}
}

// RoundTrip implements the http.RoundTripper interface.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	b, dumpErr := httputil.DumpRequestOut(req, true)
	if dumpErr == nil {
		t.debugf("Request:\n%s", b)
	} else {
		t.debugf("DumpRequestOut failed: %s", dumpErr)
	}

	resp, err := t.t.RoundTrip(req)
	if resp == nil {
		return nil, err
	}

	b, dumpErr = httputil.DumpResponse(resp, true)
	if dumpErr == nil {
		t.debugf("Request:\n%s", b)
	} else {
		t.debugf("DumpResponse failed: %s", dumpErr)
	}

	return resp, err
}

// check interfaces
var (
	_ http.RoundTripper = (*transport)(nil)
)
