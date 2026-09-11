package tracelogx

import (
	"testing"
)

func newTestTracer(rules []traceRule) *TraceEnforcer {
	return &TraceEnforcer{
		rules:        rules,
		policyLoaded: true,
	}
}

func TestIsTraceLog(t *testing.T) {
	tr := newTestTracer([]traceRule{
		mustRule("/api/v1/users", "POST"),
		mustRule("/api/v1/users/:id", "GET"),
		mustRule("/api/v1/posts", "*"),
	})

	cases := []struct {
		name    string
		url     string
		method  string
		want    bool
		wantErr bool
	}{
		{"matched", "/api/v1/users", "POST", true, false},
		{"not matched method", "/api/v1/users", "GET", false, false},
		{"param path", "/api/v1/users/123", "GET", true, false},
		{"wildcard method GET", "/api/v1/posts", "GET", true, false},
		{"wildcard method DELETE", "/api/v1/posts", "DELETE", true, false},
		{"query string stripped", "/api/v1/users?page=1", "POST", true, false},
		{"no rule matched", "/api/v1/other", "GET", false, false},
		{"empty url", "", "GET", false, true},
		{"empty method", "/api/v1/users", "", false, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ok, err := tr.IsTraceLog(c.url, c.method)
			if (err != nil) != c.wantErr {
				t.Errorf("IsTraceLog(%q, %q) error=%v, wantErr=%v", c.url, c.method, err, c.wantErr)
			}
			if ok != c.want {
				t.Errorf("IsTraceLog(%q, %q) = %v, want %v", c.url, c.method, ok, c.want)
			}
		})
	}
}

func mustRule(path, method string) traceRule {
	r, ok := newTraceRule(path, method)
	if !ok {
		panic("invalid rule: " + method + " " + path)
	}
	return r
}
