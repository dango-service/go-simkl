package simkl

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTransport_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("client_id") != "test-id" {
			t.Errorf("missing or invalid client_id: %s", r.URL.Query().Get("client_id"))
		}
		if r.URL.Query().Get("app-name") != "dango" {
			t.Errorf("missing or invalid app-name: %s", r.URL.Query().Get("app-name"))
		}
		if r.Header.Get("User-Agent") != defaultUserAgent {
			t.Errorf("unexpected User-Agent: %s", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("simkl-api-key") != "test-id" {
			t.Errorf("missing simkl-api-key header: %s", r.Header.Get("simkl-api-key"))
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Authorization header: %s", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-RateLimit-Limit", "1000")
		w.Header().Set("X-RateLimit-Remaining", "995")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL), WithToken("test-token"))

	var res map[string]string
	resp, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, &res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["status"] != "ok" {
		t.Errorf("expected status 'ok', got %q", res["status"])
	}
	if resp.RateLimit.Limit != 1000 || resp.RateLimit.Remaining != 995 {
		t.Errorf("unexpected rate limit info: %+v", resp.RateLimit)
	}
}

func TestTransport_RateLimitError_PerSecond(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"rate_limit","code":429,"message":"Rate limit exceeded"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	resp, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsRateLimited(err) {
		t.Errorf("expected IsRateLimited to be true, got false for %v", err)
	}
	if !resp.RateLimit.IsPerSecondRateLimit {
		t.Error("expected IsPerSecondRateLimit to be true")
	}
}

func TestTransport_RateLimitError_DailyQuota(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-RateLimit-Limit", "500")
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.Header().Set("Retry-After", "3600")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"user_limit_exceeded","code":429,"message":"Daily limit reached"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	resp, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsRateLimited(err) {
		t.Errorf("expected IsRateLimited to be true for %v", err)
	}
	if !resp.RateLimit.IsDailyQuotaExceeded {
		t.Error("expected IsDailyQuotaExceeded to be true")
	}
	if resp.RateLimit.RetryAfter != 3600*time.Second {
		t.Errorf("expected RetryAfter 3600s, got %v", resp.RateLimit.RetryAfter)
	}
}

func TestTransport_RateLimitError_HTTPDate(t *testing.T) {
	targetTime := time.Now().Add(5 * time.Minute).Truncate(time.Second)
	dateStr := targetTime.Format(http.TimeFormat)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", dateStr)
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"user_limit_exceeded","code":429}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	resp, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if resp.RateLimit.RetryAfter <= 0 {
		t.Errorf("expected positive RetryAfter from HTTP-Date, got %v", resp.RateLimit.RetryAfter)
	}
}

func TestTransport_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"user_token_required","code":401,"message":"Access token missing"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	_, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsUnauthorized(err) {
		t.Errorf("expected IsUnauthorized to be true for %v", err)
	}
}

func TestTransport_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not_found","code":404}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	_, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound to be true for %v", err)
	}
}

func TestTransport_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.doRequest(ctx, http.MethodGet, "/test", nil, nil)
	if err == nil {
		t.Fatal("expected context cancelled error, got nil")
	}
	if !errors.Is(err, context.Canceled) && !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("expected context cancellation in error, got %v", err)
	}
}
