package simkl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSyncService_GetActivities(t *testing.T) {
	fixture, err := os.ReadFile("testdata/sync_activities_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sync/activities" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	activities, _, err := client.Sync.GetActivities(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if activities.All != "2026-11-06T15:30:00Z" {
		t.Errorf("expected all '2026-11-06T15:30:00Z', got %q", activities.All)
	}
	if activities.Anime.WatchedAt != "2026-11-06T15:30:00Z" {
		t.Errorf("expected anime watched_at, got %q", activities.Anime.WatchedAt)
	}
}

func TestSyncService_GetAllItems(t *testing.T) {
	fixture, err := os.ReadFile("testdata/sync_all_items_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sync/all-items/anime/watching" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("extended") != "full" {
			t.Errorf("expected extended=full, got %s", r.URL.Query().Get("extended"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	collection, _, err := client.Sync.GetAllItems(context.Background(), "anime", "watching")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(collection.Anime) != 1 {
		t.Fatalf("expected 1 anime entry, got %d", len(collection.Anime))
	}
	entry := collection.Anime[0]
	if entry.Status != "watching" {
		t.Errorf("expected status 'watching', got %q", entry.Status)
	}
	if entry.WatchedEpisodesCount != 12 {
		t.Errorf("expected 12 watched episodes, got %d", entry.WatchedEpisodesCount)
	}
	if entry.Show.IDs.PrimaryID() != 40164 {
		t.Errorf("expected Simkl ID 40164, got %d", entry.Show.IDs.PrimaryID())
	}
}

func TestSyncService_AddToList(t *testing.T) {
	fixture, err := os.ReadFile("testdata/sync_mutation_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/sync/add-to-list" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	input := &AddToListInput{
		To: "watching",
		Anime: []*SyncItem{
			{
				IDs: IDs{Simkl: 40164},
			},
		},
	}
	res, _, err := client.Sync.AddToList(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added == nil {
		t.Fatal("expected added summary, got nil")
	}
}

func TestSyncService_AddHistory(t *testing.T) {
	fixture, err := os.ReadFile("testdata/sync_mutation_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/sync/history" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	input := &HistoryInput{
		Anime: []*SyncItem{
			{
				IDs: IDs{Simkl: 40164},
				Episodes: []SyncEpisode{
					{Number: 1, WatchedAt: "2026-11-06T15:30:00Z"},
				},
			},
		},
	}
	res, _, err := client.Sync.AddHistory(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Added == nil {
		t.Fatal("expected added map, got nil")
	}
}

func TestSyncService_NilInputs(t *testing.T) {
	client := NewClient("test-id")

	if _, _, err := client.Sync.AddToList(context.Background(), nil); err == nil {
		t.Error("expected error for nil input in AddToList")
	}
	if _, _, err := client.Sync.AddHistory(context.Background(), nil); err == nil {
		t.Error("expected error for nil input in AddHistory")
	}
	if _, _, err := client.Sync.RemoveHistory(context.Background(), nil); err == nil {
		t.Error("expected error for nil input in RemoveHistory")
	}
	if _, _, err := client.Sync.AddRatings(context.Background(), nil); err == nil {
		t.Error("expected error for nil input in AddRatings")
	}
	if _, _, err := client.Sync.RemoveRatings(context.Background(), nil); err == nil {
		t.Error("expected error for nil input in RemoveRatings")
	}
}
