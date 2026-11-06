package simkl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SyncService handles watchlist queries, activity checks, watch history, and rating sync.
type SyncService struct {
	client *Client
}

// Activities contains the last activity timestamps across all categories for cache invalidation.
type Activities struct {
	All         string             `json:"all"`
	Settings    string             `json:"settings,omitempty"`
	TVShows     ActivityTimestamps `json:"tv_shows,omitempty"`
	Anime       ActivityTimestamps `json:"anime,omitempty"`
	Movies      ActivityTimestamps `json:"movies,omitempty"`
	CustomLists ActivityTimestamps `json:"custom_lists,omitempty"`
}

// ActivityTimestamps represents granular last activity timestamps for a specific media category.
type ActivityTimestamps struct {
	All                   string `json:"all,omitempty"`
	WatchedAt             string `json:"watched_at,omitempty"`
	RatedAt               string `json:"rated_at,omitempty"`
	Watchlist             string `json:"watchlist,omitempty"`
	RemovedFromWatchlist  string `json:"removed_from_watchlist,omitempty"`
}

// WatchlistCollection represents a collection of watchlist entries returned by GetAllItems.
type WatchlistCollection struct {
	Anime  []*WatchlistEntry `json:"anime,omitempty"`
	Shows  []*WatchlistEntry `json:"shows,omitempty"`
	Movies []*WatchlistEntry `json:"movies,omitempty"`
}

// WatchlistEntry represents an individual item within a user's watchlist.
type WatchlistEntry struct {
	Status               string  `json:"status"`
	LastWatchedAt        string  `json:"last_watched_at,omitempty"`
	UserRating           float64 `json:"user_rating,omitempty"`
	LastWatched          string  `json:"last_watched,omitempty"`
	NextToWatch          string  `json:"next_to_watch,omitempty"`
	WatchedEpisodesCount int     `json:"watched_episodes_count,omitempty"`
	TotalEpisodesCount   int     `json:"total_episodes_count,omitempty"`
	AnimeType            string  `json:"anime_type,omitempty"`
	Show                 *Anime  `json:"show,omitempty"`
	Movie                *Anime  `json:"movie,omitempty"`
}

// SyncItem represents a media item targeted for a sync operation.
type SyncItem struct {
	Title    string        `json:"title,omitempty"`
	Year     int           `json:"year,omitempty"`
	IDs      IDs           `json:"ids"`
	Rating   float64       `json:"rating,omitempty"`
	Episodes []SyncEpisode `json:"episodes,omitempty"`
}

// SyncEpisode represents an episode within a history sync payload.
type SyncEpisode struct {
	Number    int    `json:"number"`
	Season    int    `json:"season,omitempty"`
	WatchedAt string `json:"watched_at,omitempty"`
}

// AddToListInput represents the payload for POST /sync/add-to-list.
type AddToListInput struct {
	To     string      `json:"to"` // "watching", "plantowatch", "hold", "completed", "dropped"
	Anime  []*SyncItem `json:"anime,omitempty"`
	Shows  []*SyncItem `json:"shows,omitempty"`
	Movies []*SyncItem `json:"movies,omitempty"`
}

// HistoryInput represents the payload for POST /sync/history and /sync/history/remove.
type HistoryInput struct {
	Anime  []*SyncItem `json:"anime,omitempty"`
	Shows  []*SyncItem `json:"shows,omitempty"`
	Movies []*SyncItem `json:"movies,omitempty"`
}

// RatingsInput represents the payload for POST /sync/ratings and /sync/ratings/remove.
type RatingsInput struct {
	Anime  []*SyncItem `json:"anime,omitempty"`
	Shows  []*SyncItem `json:"shows,omitempty"`
	Movies []*SyncItem `json:"movies,omitempty"`
}

// SyncResult represents the summary response returned by sync mutation endpoints.
type SyncResult struct {
	Added    map[string]any `json:"added,omitempty"`
	NotFound map[string]any `json:"not_found,omitempty"`
}

// GetActivities fetches the user's latest activity timestamps. Useful for polling optimization.
func (s *SyncService) GetActivities(ctx context.Context) (*Activities, *Response, error) {
	var activities Activities
	resp, err := s.client.doRequest(ctx, http.MethodGet, "/sync/activities", nil, &activities)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: get activities: %w", err)
	}
	return &activities, resp, nil
}

// GetAllItems retrieves all items of a given mediaType ("anime", "shows", "movies") and status ("all", "watching", "plantowatch", "hold", "completed", "dropped").
func (s *SyncService) GetAllItems(ctx context.Context, mediaType, status string, opts ...SyncOption) (*WatchlistCollection, *Response, error) {
	if mediaType == "" {
		mediaType = "anime"
	}
	if status == "" {
		status = "all"
	}

	params := url.Values{}
	params.Set("extended", "full")
	for _, opt := range opts {
		if opt != nil {
			opt(params)
		}
	}

	path := fmt.Sprintf("/sync/all-items/%s/%s?%s", url.PathEscape(mediaType), url.PathEscape(status), params.Encode())

	var collection WatchlistCollection
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &collection)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: get all items %s/%s: %w", mediaType, status, err)
	}
	return &collection, resp, nil
}

// AddToList moves media items into a watchlist status bucket (e.g. watching, plantowatch, completed).
func (s *SyncService) AddToList(ctx context.Context, input *AddToListInput) (*SyncResult, *Response, error) {
	if input == nil {
		return nil, nil, fmt.Errorf("simkl: add to list: input is required")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, nil, fmt.Errorf("simkl: marshal add to list input: %w", err)
	}

	var result SyncResult
	resp, err := s.client.doRequest(ctx, http.MethodPost, "/sync/add-to-list", bytes.NewReader(body), &result)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: add to list: %w", err)
	}
	return &result, resp, nil
}

// AddHistory marks media items or specific episodes as watched.
func (s *SyncService) AddHistory(ctx context.Context, input *HistoryInput) (*SyncResult, *Response, error) {
	if input == nil {
		return nil, nil, fmt.Errorf("simkl: add history: input is required")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, nil, fmt.Errorf("simkl: marshal history input: %w", err)
	}

	var result SyncResult
	resp, err := s.client.doRequest(ctx, http.MethodPost, "/sync/history", bytes.NewReader(body), &result)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: add history: %w", err)
	}
	return &result, resp, nil
}

// RemoveHistory removes items or episodes from the user's watch history.
func (s *SyncService) RemoveHistory(ctx context.Context, input *HistoryInput) (*SyncResult, *Response, error) {
	if input == nil {
		return nil, nil, fmt.Errorf("simkl: remove history: input is required")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, nil, fmt.Errorf("simkl: marshal remove history input: %w", err)
	}

	var result SyncResult
	resp, err := s.client.doRequest(ctx, http.MethodPost, "/sync/history/remove", bytes.NewReader(body), &result)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: remove history: %w", err)
	}
	return &result, resp, nil
}

// AddRatings adds or updates ratings for media items.
func (s *SyncService) AddRatings(ctx context.Context, input *RatingsInput) (*SyncResult, *Response, error) {
	if input == nil {
		return nil, nil, fmt.Errorf("simkl: add ratings: input is required")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, nil, fmt.Errorf("simkl: marshal ratings input: %w", err)
	}

	var result SyncResult
	resp, err := s.client.doRequest(ctx, http.MethodPost, "/sync/ratings", bytes.NewReader(body), &result)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: add ratings: %w", err)
	}
	return &result, resp, nil
}

// RemoveRatings removes ratings from media items.
func (s *SyncService) RemoveRatings(ctx context.Context, input *RatingsInput) (*SyncResult, *Response, error) {
	if input == nil {
		return nil, nil, fmt.Errorf("simkl: remove ratings: input is required")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, nil, fmt.Errorf("simkl: marshal remove ratings input: %w", err)
	}

	var result SyncResult
	resp, err := s.client.doRequest(ctx, http.MethodPost, "/sync/ratings/remove", bytes.NewReader(body), &result)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: remove ratings: %w", err)
	}
	return &result, resp, nil
}
