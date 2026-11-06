package simkl

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SearchService handles searching catalog metadata and cross-catalog external ID resolution.
type SearchService struct {
	client *Client
}

// IDs contains cross-catalog identifiers for media items across services.
type IDs struct {
	Simkl     int    `json:"simkl,omitempty"`
	SimklID   int    `json:"simkl_id,omitempty"`
	IMDb      string `json:"imdb,omitempty"`
	TMDB      string `json:"tmdb,omitempty"`
	TVDB      string `json:"tvdb,omitempty"`
	MAL       string `json:"mal,omitempty"`
	AniDB     string `json:"anidb,omitempty"`
	AniList   string `json:"anilist,omitempty"`
	Kitsu     string `json:"kitsu,omitempty"`
	LiveChart string `json:"livechart,omitempty"`
}

// PrimaryID returns the canonical Simkl ID, accepting either simkl or simkl_id.
func (ids IDs) PrimaryID() int {
	if ids.Simkl != 0 {
		return ids.Simkl
	}
	return ids.SimklID
}

// SearchResult represents a media item returned by search endpoints.
type SearchResult struct {
	Title    string `json:"title"`
	EnTitle  string `json:"en_title,omitempty"`
	Year     int    `json:"year,omitempty"`
	Type     string `json:"type,omitempty"`
	Poster   string `json:"poster,omitempty"`
	Fanart   string `json:"fanart,omitempty"`
	Overview string `json:"overview,omitempty"`
	IDs      IDs    `json:"ids"`
}

// Search searches items by title query. mediaType must be "anime", "tv", "movie", or "all".
func (s *SearchService) Search(ctx context.Context, mediaType, query string, opts ...SearchOption) ([]*SearchResult, *Response, error) {
	if mediaType == "" {
		mediaType = "all"
	}
	params := url.Values{}
	params.Set("q", query)

	for _, opt := range opts {
		if opt != nil {
			opt(params)
		}
	}

	path := fmt.Sprintf("/search/%s?%s", url.PathEscape(mediaType), params.Encode())

	var results []*SearchResult
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &results)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: search %s %q: %w", mediaType, query, err)
	}
	return results, resp, nil
}

// FindByExternalID queries Simkl to resolve an external service ID (e.g., mal, anilist, kitsu, tmdb, imdb).
func (s *SearchService) FindByExternalID(ctx context.Context, service, id string) ([]*SearchResult, *Response, error) {
	service = strings.ToLower(strings.TrimSpace(service))
	id = strings.TrimSpace(id)

	params := url.Values{}
	params.Set(service, id)

	path := fmt.Sprintf("/search/id?%s", params.Encode())

	var results []*SearchResult
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &results)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: find by %s ID %q: %w", service, id, err)
	}
	return results, resp, nil
}

// Redirect resolves an external ID directly to a Simkl ID by checking HTTP redirect or result ID.
func (s *SearchService) Redirect(ctx context.Context, service, id string) (int, *Response, error) {
	results, resp, err := s.FindByExternalID(ctx, service, id)
	if err != nil {
		return 0, resp, err
	}
	if len(results) == 0 {
		return 0, resp, fmt.Errorf("simkl: %w: no match found for %s ID %q", ErrNotFound, service, id)
	}
	simklID := results[0].IDs.PrimaryID()
	if simklID == 0 {
		return 0, resp, fmt.Errorf("simkl: %w: returned item has no simkl ID", ErrNotFound)
	}
	return simklID, resp, nil
}
