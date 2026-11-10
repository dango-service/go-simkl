package simkl

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAuthService_AuthorizeURL(t *testing.T) {
	client := NewClient("my-client-id")
	url := client.Auth.AuthorizeURL("https://example.com/callback", "random-state")

	if !strings.Contains(url, "client_id=my-client-id") {
		t.Errorf("missing client_id in URL: %s", url)
	}
	if !strings.Contains(url, "response_type=code") {
		t.Errorf("missing response_type in URL: %s", url)
	}
	if !strings.Contains(url, "redirect_uri=https%3A%2F%2Fexample.com%2Fcallback") {
		t.Errorf("missing encoded redirect_uri in URL: %s", url)
	}
	if !strings.Contains(url, "state=random-state") {
		t.Errorf("missing state in URL: %s", url)
	}

	pkceURL := client.Auth.AuthorizeURL("https://example.com/callback", "random-state", "test-challenge-s256")
	if !strings.Contains(pkceURL, "code_challenge=test-challenge-s256") {
		t.Errorf("missing code_challenge in PKCE URL: %s", pkceURL)
	}
	if !strings.Contains(pkceURL, "code_challenge_method=S256") {
		t.Errorf("missing code_challenge_method in PKCE URL: %s", pkceURL)
	}
}

func TestAuthService_ExchangeCode_PKCE(t *testing.T) {
	fixture, err := os.ReadFile("testdata/token_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	token, _, err := client.Auth.ExchangeCode(context.Background(), "auth-code", "https://example.com/callback", "", "my-code-verifier")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token.AccessToken != "simkl_at_test_access_token_1234567890" {
		t.Errorf("unexpected access token: %s", token.AccessToken)
	}
}

func TestAuthService_ExchangeCode(t *testing.T) {
	fixture, err := os.ReadFile("testdata/token_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/oauth2/token" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	token, _, err := client.Auth.ExchangeCode(context.Background(), "auth-code", "https://example.com/callback", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token.AccessToken != "simkl_at_test_access_token_1234567890" {
		t.Errorf("unexpected access token: %s", token.AccessToken)
	}
	if token.RefreshToken != "simkl_rt_test_refresh_token_0987654321" {
		t.Errorf("unexpected refresh token: %s", token.RefreshToken)
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	fixture, err := os.ReadFile("testdata/token_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	token, _, err := client.Auth.RefreshToken(context.Background(), "simkl_rt_test", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
}

func TestAuthService_DeviceCode(t *testing.T) {
	fixture, err := os.ReadFile("testdata/device_code_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth2/device" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	deviceResp, _, err := client.Auth.RequestDeviceCode(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deviceResp.UserCode != "BDWP-HQPK" {
		t.Errorf("expected user code 'BDWP-HQPK', got %q", deviceResp.UserCode)
	}
	if deviceResp.VerificationURIComplete != "https://simkl.com/pin?code=BDWP-HQPK" {
		t.Errorf("unexpected verification URI: %s", deviceResp.VerificationURIComplete)
	}
}

func TestAuthService_PollDeviceToken(t *testing.T) {
	fixture, err := os.ReadFile("testdata/token_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	token, _, err := client.Auth.PollDeviceToken(context.Background(), "secret_device_code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
}

func TestAuthService_RevokeToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth2/revoke" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"revoked"}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	_, err := client.Auth.RevokeToken(context.Background(), "token-to-revoke")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestToken_CredentialRedaction(t *testing.T) {
	tok := Token{
		AccessToken:  "secret-access-token",
		RefreshToken: "secret-refresh-token",
	}

	strOutput := fmt.Sprintf("%s", tok)
	if strOutput != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", strOutput)
	}

	goStrOutput := fmt.Sprintf("%#v", tok)
	if goStrOutput != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", goStrOutput)
	}
}
