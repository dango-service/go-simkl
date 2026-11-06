package simkl

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("my-client-id")

	if c.clientID != "my-client-id" {
		t.Errorf("expected clientID 'my-client-id', got %q", c.clientID)
	}
	if c.baseURL.String() != defaultBaseURL {
		t.Errorf("expected baseURL %q, got %q", defaultBaseURL, c.baseURL.String())
	}
	if c.appName != defaultAppName {
		t.Errorf("expected appName %q, got %q", defaultAppName, c.appName)
	}
	if c.appVersion != defaultAppVersion {
		t.Errorf("expected appVersion %q, got %q", defaultAppVersion, c.appVersion)
	}
	if c.userAgent != defaultUserAgent {
		t.Errorf("expected userAgent %q, got %q", defaultUserAgent, c.userAgent)
	}
	if c.maxBytes != defaultMaxBodyBytes {
		t.Errorf("expected maxBytes %d, got %d", defaultMaxBodyBytes, c.maxBytes)
	}
	if c.Search == nil || c.Anime == nil || c.Sync == nil || c.Scrobble == nil || c.User == nil || c.Auth == nil {
		t.Error("expected all sub-services to be initialized")
	}
}

func TestNewClient_Options(t *testing.T) {
	customHTTP := &http.Client{Timeout: 5 * time.Second}
	c := NewClient("  trimmed-id  ",
		WithHTTPClient(customHTTP),
		WithBaseURL("https://custom.simkl.com/v2"),
		WithToken("  secret-token  "),
		WithAppInfo("my-app", "2.1.0"),
		WithUserAgent("CustomAgent/1.0"),
		WithMaxBodyBytes(2048),
	)

	if c.clientID != "trimmed-id" {
		t.Errorf("expected trimmed clientID, got %q", c.clientID)
	}
	if c.client != customHTTP {
		t.Errorf("expected custom HTTP client")
	}
	if c.baseURL.String() != "https://custom.simkl.com/v2" {
		t.Errorf("expected custom baseURL, got %q", c.baseURL.String())
	}
	if c.token != "secret-token" {
		t.Errorf("expected trimmed token, got %q", c.token)
	}
	if c.appName != "my-app" || c.appVersion != "2.1.0" {
		t.Errorf("expected custom app info, got %s %s", c.appName, c.appVersion)
	}
	if c.userAgent != "CustomAgent/1.0" {
		t.Errorf("expected custom userAgent, got %q", c.userAgent)
	}
	if c.maxBytes != 2048 {
		t.Errorf("expected custom maxBytes 2048, got %d", c.maxBytes)
	}
}
