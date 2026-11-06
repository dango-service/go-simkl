# go-simkl

[![Go Reference](https://pkg.go.dev/badge/github.com/dango-service/go-simkl.svg)](https://pkg.go.dev/github.com/dango-service/go-simkl)
[![CI](https://github.com/dango-service/go-simkl/actions/workflows/ci.yml/badge.svg)](https://github.com/dango-service/go-simkl/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A lightweight, standalone Go client library for the [Simkl](https://simkl.com) API (`https://api.simkl.com`).
Designed with **zero external dependencies** (standard library only), direct streaming deserialization bounded by memory limits, and strict efficiency for background sync daemons like [Dango](https://github.com/dango-service/dango).

---

## Features

- **Zero External Dependencies**: Built entirely upon the Go standard library (`net/http`, `encoding/json`, `context`, `time`, `net/url`).
- **Streaming JSON**: Direct streaming decoding with `json.NewDecoder` bounded by `io.LimitReader` (10 MiB limit) to prevent memory exhaustion DoS.
- **Dual Rate Limit Aware**: Automatic inspection of Simkl quota headers (`X-RateLimit-*`) and native handling of HTTP `429 Too Many Requests`, distinguishing between per-second burst limits (`10 GET/sec`, `1 POST/sec`) and daily quota exhaustion (`user_limit_exceeded`, `Retry-After`).
- **Credential Protection**: Redacted string formatting (`[REDACTED]`) on `Token` to prevent secret leakage in logs.
- **Modular Services**:
  - `Search`: Full-text search and cross-catalog external ID resolution (`/search/id`, `/redirect`) across MyAnimeList, AniList, Kitsu, TMDB, TVDB, and IMDb.
  - `Anime`: Anime metadata (`/anime/{id}`) and episode lists (`/anime/episodes/{id}`).
  - `Sync`: Watchlist tracking (`/sync/all-items`), activity polling checks (`/sync/activities`), add to list, watch history, and ratings.
  - `Scrobble`: Real-time media playback sessions (`/scrobble/start`, `/scrobble/stop`, `/scrobble/pause`).
  - `User`: Profile and account settings (`/users/settings`).
  - `Auth`: Full Auth V2 OAuth2 support (Authorization Code grant, token refresh, token revoke) and Device/PIN flow for headless scrobblers.
- **100% Deterministic Testing**: Complete test suite using `net/http/httptest` with zero real-network requests.

---

## Installation

```bash
go get github.com/dango-service/go-simkl
```

---

## Usage Examples

### 1. Client Initialization

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/dango-service/go-simkl"
)

func main() {
    client := simkl.NewClient("YOUR_SIMKL_CLIENT_ID",
        simkl.WithAppInfo("dango", "1.0.0"),
        simkl.WithToken("USER_ACCESS_TOKEN"),
    )

    // Check user activities
    activities, _, err := client.Sync.GetActivities(context.Background())
    if err != nil {
        log.Fatalf("failed to fetch activities: %v", err)
    }

    fmt.Printf("Last Anime Activity: %s\n", activities.Anime.WatchedAt)
}
```

### 2. Search & Cross-Catalog ID Resolution

```go
// Search anime by title
results, _, err := client.Search.Search(ctx, "anime", "Frieren", simkl.WithSearchLimit(5))

// Resolve external ID (e.g. MyAnimeList ID 52991) to Simkl ID
matches, _, err := client.Search.FindByExternalID(ctx, "mal", "52991")
if len(matches) > 0 {
    fmt.Printf("Simkl ID: %d\n", matches[0].IDs.PrimaryID())
}
```

### 3. Watchlist & Watch History

```go
// Retrieve watching anime list
watchlist, _, err := client.Sync.GetAllItems(ctx, "anime", "watching")

// Mark episode as watched
_, _, err = client.Sync.AddHistory(ctx, &simkl.HistoryInput{
    Anime: []*simkl.SyncItem{
        {
            IDs: simkl.IDs{Simkl: 40164},
            Episodes: []simkl.SyncEpisode{
                {Number: 1},
            },
        },
    },
})
```

### 4. Real-Time Scrobbling

```go
// Start scrobble
_, _, err := client.Scrobble.Start(ctx, &simkl.ScrobbleInput{
    Anime: &simkl.SyncItem{
        IDs: simkl.IDs{Simkl: 40164},
    },
    Episode:  &simkl.SyncEpisode{Number: 1},
    Progress: 10.0,
})

// Stop scrobble
_, _, err = client.Scrobble.Stop(ctx, &simkl.ScrobbleInput{
    Anime: &simkl.SyncItem{
        IDs: simkl.IDs{Simkl: 40164},
    },
    Episode:  &simkl.SyncEpisode{Number: 1},
    Progress: 100.0,
})
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
