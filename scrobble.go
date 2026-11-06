package simkl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ScrobbleService handles real-time media playback scrobbling sessions.
type ScrobbleService struct {
	client *Client
}

// ScrobbleInput represents a payload for scrobbling playback events.
type ScrobbleInput struct {
	Anime    *SyncItem    `json:"anime,omitempty"`
	Show     *SyncItem    `json:"show,omitempty"`
	Movie    *SyncItem    `json:"movie,omitempty"`
	Episode  *SyncEpisode `json:"episode,omitempty"`
	Progress float64      `json:"progress"` // percentage (0.0 to 100.0)
}

// ScrobbleResponse represents the confirmation from a scrobble action.
type ScrobbleResponse struct {
	Action   string    `json:"action"` // "start", "stop", "pause"
	Progress float64   `json:"progress"`
	Item     *SyncItem `json:"item,omitempty"`
}

// Start notifies Simkl that playback of a media item has started.
func (s *ScrobbleService) Start(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error) {
	return s.doScrobble(ctx, "/scrobble/start", input)
}

// Stop notifies Simkl that playback has stopped or finished.
func (s *ScrobbleService) Stop(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error) {
	return s.doScrobble(ctx, "/scrobble/stop", input)
}

// Pause notifies Simkl that playback has paused.
func (s *ScrobbleService) Pause(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error) {
	return s.doScrobble(ctx, "/scrobble/pause", input)
}

func (s *ScrobbleService) doScrobble(ctx context.Context, path string, input *ScrobbleInput) (*ScrobbleResponse, *Response, error) {
	if input == nil {
		return nil, nil, fmt.Errorf("simkl: scrobble: input is required")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, nil, fmt.Errorf("simkl: marshal scrobble input: %w", err)
	}

	var res ScrobbleResponse
	resp, err := s.client.doRequest(ctx, http.MethodPost, path, bytes.NewReader(body), &res)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: scrobble %s: %w", path, err)
	}
	return &res, resp, nil
}
