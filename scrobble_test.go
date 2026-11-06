package simkl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScrobbleService_StartStopPause(t *testing.T) {
	calledPath := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"action":"scrobbled","progress":50.0}`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	input := &ScrobbleInput{
		Anime: &SyncItem{
			IDs: IDs{Simkl: 40164},
		},
		Episode:  &SyncEpisode{Number: 1},
		Progress: 50.0,
	}

	// Test Start
	res, _, err := client.Scrobble.Start(context.Background(), input)
	if err != nil || calledPath != "/scrobble/start" {
		t.Fatalf("unexpected Start error: %v, path: %s", err, calledPath)
	}
	if res.Progress != 50.0 {
		t.Errorf("expected progress 50.0, got %f", res.Progress)
	}

	// Test Stop
	_, _, err = client.Scrobble.Stop(context.Background(), input)
	if err != nil || calledPath != "/scrobble/stop" {
		t.Fatalf("unexpected Stop error: %v, path: %s", err, calledPath)
	}

	// Test Pause
	_, _, err = client.Scrobble.Pause(context.Background(), input)
	if err != nil || calledPath != "/scrobble/pause" {
		t.Fatalf("unexpected Pause error: %v, path: %s", err, calledPath)
	}

	// Test nil input
	if _, _, err := client.Scrobble.Start(context.Background(), nil); err == nil {
		t.Error("expected error for nil input")
	}
}
