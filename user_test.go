package simkl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUserService_GetSettings(t *testing.T) {
	fixture, err := os.ReadFile("testdata/user_settings_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/settings" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL), WithToken("valid-token"))

	settings, _, err := client.User.GetSettings(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings.User.Name != "dango_watcher" {
		t.Errorf("expected user name 'dango_watcher', got %q", settings.User.Name)
	}
	if settings.Account.Timezone != "America/New_York" {
		t.Errorf("expected timezone 'America/New_York', got %q", settings.Account.Timezone)
	}
}
