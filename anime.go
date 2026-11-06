package simkl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AnimeService handles anime metadata and episode retrieval.
type AnimeService struct {
	client *Client
}

// Anime represents detailed metadata for an anime series or movie.
type Anime struct {
	Title         string   `json:"title"`
	EnTitle       string   `json:"en_title,omitempty"`
	Year          int      `json:"year,omitempty"`
	Type          string   `json:"type,omitempty"`
	AnimeType     string   `json:"anime_type,omitempty"`
	Poster        string   `json:"poster,omitempty"`
	Fanart        string   `json:"fanart,omitempty"`
	Overview      string   `json:"overview,omitempty"`
	TotalEpisodes int      `json:"total_episodes,omitempty"`
	Status        string   `json:"status,omitempty"`
	Rank          int      `json:"rank,omitempty"`
	Genres        []string `json:"genres,omitempty"`
	IDs           IDs      `json:"ids"`
}

// Episode represents an individual anime episode.
type Episode struct {
	Title       string `json:"title,omitempty"`
	Episode     int    `json:"episode"`
	Season      int    `json:"season,omitempty"`
	Aired       string `json:"aired,omitempty"`
	Date        string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`
	IDs         IDs    `json:"ids,omitempty"`
}

// GetByID fetches detailed anime metadata by Simkl ID.
func (s *AnimeService) GetByID(ctx context.Context, id int) (*Anime, *Response, error) {
	if id <= 0 {
		return nil, nil, fmt.Errorf("simkl: invalid anime ID %d", id)
	}

	path := fmt.Sprintf("/anime/%d?extended=full", id)

	var raw json.RawMessage
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &raw)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: get anime %d: %w", id, err)
	}

	// Simkl returns an empty array "[]" when the ID is unknown
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("[]")) {
		return nil, resp, fmt.Errorf("simkl: anime %d: %w", id, ErrNotFound)
	}

	var anime Anime
	if err := json.Unmarshal(trimmed, &anime); err != nil {
		return nil, resp, fmt.Errorf("simkl: decode anime %d: %w", id, err)
	}

	return &anime, resp, nil
}

// GetEpisodes fetches the episode list for an anime by Simkl ID.
func (s *AnimeService) GetEpisodes(ctx context.Context, id int) ([]*Episode, *Response, error) {
	if id <= 0 {
		return nil, nil, fmt.Errorf("simkl: invalid anime ID %d", id)
	}

	path := fmt.Sprintf("/anime/episodes/%d", id)

	var episodes []*Episode
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, &episodes)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: get episodes for anime %d: %w", id, err)
	}
	return episodes, resp, nil
}
