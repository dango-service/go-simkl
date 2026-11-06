package simkl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAnimeService_GetByID(t *testing.T) {
	fixture, err := os.ReadFile("testdata/anime_by_id_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/anime/40164" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	anime, _, err := client.Anime.GetByID(context.Background(), 40164)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if anime.Title != "Steins;Gate" {
		t.Errorf("expected title 'Steins;Gate', got %q", anime.Title)
	}
	if anime.TotalEpisodes != 24 {
		t.Errorf("expected 24 episodes, got %d", anime.TotalEpisodes)
	}
	if anime.IDs.PrimaryID() != 40164 {
		t.Errorf("expected primary ID 40164, got %d", anime.IDs.PrimaryID())
	}
}

func TestAnimeService_GetByID_EmptyArrayNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	_, _, err := client.Anime.GetByID(context.Background(), 99999999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound, got %v", err)
	}
}

func TestAnimeService_GetEpisodes(t *testing.T) {
	fixture, err := os.ReadFile("testdata/anime_episodes_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/anime/episodes/40164" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	episodes, _, err := client.Anime.GetEpisodes(context.Background(), 40164)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(episodes) != 2 {
		t.Fatalf("expected 2 episodes, got %d", len(episodes))
	}
	if episodes[0].Episode != 1 || episodes[1].Episode != 2 {
		t.Errorf("unexpected episode numbers: %d, %d", episodes[0].Episode, episodes[1].Episode)
	}
}
