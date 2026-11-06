package simkl

import (
	"context"
	"fmt"
	"net/http"
)

// UserService handles user settings and profile details.
type UserService struct {
	client *Client
}

// UserSettings represents the authenticated user's profile and account settings.
type UserSettings struct {
	User    UserAccount `json:"user"`
	Account AccountInfo `json:"account,omitempty"`
}

// UserAccount represents the user's public profile fields.
type UserAccount struct {
	Name     string `json:"name"`
	JoinedAt string `json:"joined_at,omitempty"`
	Gender   string `json:"gender,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Bio      string `json:"bio,omitempty"`
	ID       int    `json:"id,omitempty"`
}

// AccountInfo represents account configuration options like timezone and country.
type AccountInfo struct {
	Timezone string `json:"timezone,omitempty"`
	Country  string `json:"country,omitempty"`
}

// GetSettings retrieves the profile and settings of the authenticated user.
func (s *UserService) GetSettings(ctx context.Context) (*UserSettings, *Response, error) {
	var settings UserSettings
	resp, err := s.client.doRequest(ctx, http.MethodGet, "/users/settings", nil, &settings)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: get user settings: %w", err)
	}
	return &settings, resp, nil
}
