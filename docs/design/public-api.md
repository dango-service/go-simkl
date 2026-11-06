# Go Public API Specification  -  `go-simkl`

This document defines the Go public API signatures, models, options, and contracts
for `go-simkl`. The design mirrors the clean architecture of `go-kitsu` and `go-anilist`.

---

## 1. Package and Core Client

```go
package simkl

import (
    "context"
    "net/http"
    "net/url"
)

// Client manages communication with the Simkl API.
type Client struct {
    Search   *SearchService
    Anime    *AnimeService
    Sync     *SyncService
    Scrobble *ScrobbleService
    User     *UserService
    Auth     *AuthService
}

func NewClient(clientID string, opts ...ClientOption) *Client
```

---

## 2. Functional Options

```go
type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption
func WithBaseURL(rawURL string) ClientOption
func WithToken(token string) ClientOption
func WithAppInfo(appName, appVersion string) ClientOption
func WithUserAgent(userAgent string) ClientOption
func WithMaxBodyBytes(limit int64) ClientOption
```

---

## 3. Core Models

```go
type IDs struct {
    Simkl    int    `json:"simkl,omitempty"`
    SimklID  int    `json:"simkl_id,omitempty"`
    IMDb     string `json:"imdb,omitempty"`
    TMDB     string `json:"tmdb,omitempty"`
    TVDB     string `json:"tvdb,omitempty"`
    MAL      string `json:"mal,omitempty"`
    AniDB    string `json:"anidb,omitempty"`
    AniList  string `json:"anilist,omitempty"`
    Kitsu    string `json:"kitsu,omitempty"`
}

type Anime struct {
    Title        string   `json:"title"`
    Year         int      `json:"year,omitempty"`
    Type         string   `json:"type,omitempty"`
    IDs          IDs      `json:"ids"`
    TotalEpisodes int     `json:"total_episodes,omitempty"`
    Status       string   `json:"status,omitempty"`
    Overview     string   `json:"overview,omitempty"`
}

type Episode struct {
    Title       string `json:"title,omitempty"`
    Episode     int    `json:"episode"`
    Season      int    `json:"season,omitempty"`
    Aired       string `json:"aired,omitempty"`
    IDs         IDs    `json:"ids,omitempty"`
}
```

---

## 4. Service Signatures

### SearchService
```go
func (s *SearchService) Search(ctx context.Context, mediaType, query string, opts ...SearchOption) ([]*SearchResult, *Response, error)
func (s *SearchService) FindByExternalID(ctx context.Context, service, id string) ([]*SearchResult, *Response, error)
func (s *SearchService) Redirect(ctx context.Context, service, id string) (int, *Response, error)
```

### AnimeService
```go
func (s *AnimeService) GetByID(ctx context.Context, id int) (*Anime, *Response, error)
func (s *AnimeService) GetEpisodes(ctx context.Context, id int) ([]*Episode, *Response, error)
```

### SyncService
```go
func (s *SyncService) GetActivities(ctx context.Context) (*Activities, *Response, error)
func (s *SyncService) GetAllItems(ctx context.Context, mediaType, status string, opts ...SyncOption) (*WatchlistCollection, *Response, error)
func (s *SyncService) AddToList(ctx context.Context, input *AddToListInput) (*SyncResult, *Response, error)
func (s *SyncService) AddHistory(ctx context.Context, input *HistoryInput) (*SyncResult, *Response, error)
func (s *SyncService) RemoveHistory(ctx context.Context, input *HistoryInput) (*SyncResult, *Response, error)
func (s *SyncService) AddRatings(ctx context.Context, input *RatingsInput) (*SyncResult, *Response, error)
func (s *SyncService) RemoveRatings(ctx context.Context, input *RatingsInput) (*SyncResult, *Response, error)
```

### ScrobbleService
```go
func (s *ScrobbleService) Start(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error)
func (s *ScrobbleService) Stop(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error)
func (s *ScrobbleService) Pause(ctx context.Context, input *ScrobbleInput) (*ScrobbleResponse, *Response, error)
```

### UserService
```go
func (s *UserService) GetSettings(ctx context.Context) (*UserSettings, *Response, error)
```

### AuthService
```go
func (s *AuthService) AuthorizeURL(redirectURI, state string) string
func (s *AuthService) ExchangeCode(ctx context.Context, code, redirectURI, clientSecret string) (*Token, *Response, error)
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, clientSecret string) (*Token, *Response, error)
func (s *AuthService) RequestDeviceCode(ctx context.Context) (*DeviceCodeResponse, *Response, error)
func (s *AuthService) PollDeviceToken(ctx context.Context, deviceCode string) (*Token, *Response, error)
func (s *AuthService) RevokeToken(ctx context.Context, token string) (*Response, error)
```
