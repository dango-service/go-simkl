package simkl

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL      = "https://api.simkl.com"
	defaultUserAgent    = "go-simkl/0.1.0 (+https://github.com/dango-service/go-simkl)"
	defaultMaxBodyBytes = 10 * 1024 * 1024 // 10 MiB
	defaultAppName      = "dango"
	defaultAppVersion  = "1.0.0"
)

// Client coordinates communication with the Simkl API.
type Client struct {
	client     *http.Client
	baseURL    *url.URL
	clientID   string
	appName    string
	appVersion string
	token      string
	userAgent  string
	maxBytes   int64

	Search   *SearchService
	Anime    *AnimeService
	Sync     *SyncService
	Scrobble *ScrobbleService
	User     *UserService
	Auth     *AuthService
}

// NewClient initializes a new Simkl API client with the given client ID.
func NewClient(clientID string, opts ...ClientOption) *Client {
	parsedBaseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:    parsedBaseURL,
		clientID:   strings.TrimSpace(clientID),
		appName:    defaultAppName,
		appVersion: defaultAppVersion,
		userAgent:  defaultUserAgent,
		maxBytes:   defaultMaxBodyBytes,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}

	c.Search = &SearchService{client: c}
	c.Anime = &AnimeService{client: c}
	c.Sync = &SyncService{client: c}
	c.Scrobble = &ScrobbleService{client: c}
	c.User = &UserService{client: c}
	c.Auth = &AuthService{client: c}

	return c
}

// SearchService handles searching catalog metadata and cross-catalog external ID resolution.
type SearchService struct{ client *Client }

// AnimeService handles anime metadata and episode retrieval.
type AnimeService struct{ client *Client }

// SyncService handles watchlist queries, activity checks, history, and ratings.
type SyncService struct{ client *Client }

// ScrobbleService handles real-time playback scrobbling sessions.
type ScrobbleService struct{ client *Client }

// UserService handles user settings and profile details.
type UserService struct{ client *Client }

// AuthService handles OAuth2 Authorization Code, token refresh, and Device/PIN flows.
type AuthService struct{ client *Client }
