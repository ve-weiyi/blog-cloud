package ai

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamCollectsDeltasAndUsage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("authorization header = %q", got)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)

		fmt.Fprint(w, ": ping\n\n")
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"你\"}}]}\n\n")
		flusher.Flush()
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"好\"}}]}\n\n")
		fmt.Fprint(w, "data: {\"choices\":[],\"usage\":{\"prompt_tokens\":12,\"completion_tokens\":3}}\n\n")
		fmt.Fprint(w, "data: [DONE]\n\n")
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL, APIKey: "test-key", Model: "test-model"})

	var deltas []string
	usage, err := client.Stream(context.Background(), []Message{{Role: "user", Content: "你好"}}, func(s string) error {
		deltas = append(deltas, s)
		return nil
	})
	if err != nil {
		t.Fatalf("stream failed: %v", err)
	}
	if got := strings.Join(deltas, ""); got != "你好" {
		t.Fatalf("deltas = %q", got)
	}
	if usage.PromptTokens != 12 || usage.CompletionTokens != 3 {
		t.Fatalf("usage = %+v", usage)
	}
}

func TestStreamUpstreamError(t *testing.T) {
	tests := []struct {
		status int
		kind   string
	}{
		{status: http.StatusUnauthorized, kind: KindAuth},
		{status: http.StatusTooManyRequests, kind: KindRateLimited},
		{status: http.StatusInternalServerError, kind: KindUpstream},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				fmt.Fprint(w, `{"error":{"message":"boom"}}`)
			}))
			defer server.Close()

			client := NewClient(Config{BaseURL: server.URL, Model: "m"})
			_, err := client.Stream(context.Background(), []Message{{Role: "user", Content: "hi"}}, func(string) error {
				return nil
			})

			var aiErr *Error
			if !errors.As(err, &aiErr) {
				t.Fatalf("want *Error, got %v", err)
			}
			if aiErr.Kind != tt.kind || aiErr.Status != tt.status {
				t.Fatalf("got %+v, want kind %s status %d", aiErr, tt.kind, tt.status)
			}
		})
	}
}

// onDelta 返回错误时应立即停止读取。
func TestStreamStopsOnDeltaError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"a\"}}]}\n\n")
		flusher.Flush()
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"b\"}}]}\n\n")
	}))
	defer server.Close()

	stop := errors.New("stop")
	client := NewClient(Config{BaseURL: server.URL, Model: "m"})
	_, err := client.Stream(context.Background(), []Message{{Role: "user", Content: "hi"}}, func(string) error {
		return stop
	})
	if !errors.Is(err, stop) {
		t.Fatalf("want onDelta error, got %v", err)
	}
}

// onDelta 因取消而返回 ctx 错误时，应归类为取消而不是上游失败。
func TestStreamCanceledFromDelta(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"a\"}}]}\n\n")
		flusher.Flush()
		<-r.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	client := NewClient(Config{BaseURL: server.URL, Model: "m"})
	_, err := client.Stream(ctx, []Message{{Role: "user", Content: "hi"}}, func(string) error {
		cancel()
		return context.Canceled
	})

	var aiErr *Error
	if !errors.As(err, &aiErr) || aiErr.Kind != KindCanceled {
		t.Fatalf("want canceled error, got %v", err)
	}
}

func TestStreamCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		fmt.Fprint(w, "data: {\"choices\":[{\"delta\":{\"content\":\"a\"}}]}\n\n")
		flusher.Flush()
		<-r.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	client := NewClient(Config{BaseURL: server.URL, Model: "m"})
	_, err := client.Stream(ctx, []Message{{Role: "user", Content: "hi"}}, func(string) error {
		cancel()
		return nil
	})

	var aiErr *Error
	if !errors.As(err, &aiErr) || aiErr.Kind != KindCanceled {
		t.Fatalf("want canceled error, got %v", err)
	}
}
