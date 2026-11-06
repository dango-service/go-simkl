package simkl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSearchService_Search(t *testing.T) {
	fixture, err := os.ReadFile("testdata/search_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/anime" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "Attack on Titan" {
			t.Errorf("unexpected q param: %s", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("unexpected limit param: %s", r.URL.Query().Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	results, _, err := client.Search.Search(context.Background(), "anime", "Attack on Titan", WithSearchLimit(10))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	item := results[0]
	if item.Title != "Attack on Titan" {
		t.Errorf("expected title 'Attack on Titan', got %q", item.Title)
	}
	if item.IDs.PrimaryID() != 44414 {
		t.Errorf("expected primary ID 44414, got %d", item.IDs.PrimaryID())
	}
	if item.IDs.MAL != "16498" {
		t.Errorf("expected MAL ID '16498', got %q", item.IDs.MAL)
	}
}

func TestSearchService_FindByExternalID(t *testing.T) {
	fixture, err := os.ReadFile("testdata/search_id_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/id" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("mal") != "52991" {
			t.Errorf("unexpected mal param: %s", r.URL.Query().Get("mal"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	results, _, err := client.Search.FindByExternalID(context.Background(), "mal", "52991")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].IDs.PrimaryID() != 2026856 {
		t.Errorf("expected Simkl ID 2026856, got %d", results[0].IDs.PrimaryID())
	}
}

func TestSearchService_Redirect(t *testing.T) {
	fixture, err := os.ReadFile("testdata/search_id_success.json")
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

	simklID, _, err := client.Search.Redirect(context.Background(), "mal", "52991")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if simklID != 2026856 {
		t.Errorf("expected Simkl ID 2026856, got %d", simklID)
	}
}

func TestSearchService_Redirect_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client := NewClient("test-id", WithBaseURL(server.URL))

	_, _, err := client.Search.Redirect(context.Background(), "mal", "999999999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound, got %v", err)
	}
}
