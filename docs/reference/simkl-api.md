# Simkl API Reference  -  `go-simkl`

Technical reference for interacting with the Simkl API (`https://api.simkl.com`).

---

## 1. Endpoints Overview

| Category | HTTP Method | Path | Description |
| :--- | :--- | :--- | :--- |
| **Search** | `GET` | `/search/{type}` | Search anime, TV shows, and movies by title |
| **Search** | `GET` | `/search/id` | Find item by external ID (MAL, AniList, Kitsu, etc.) |
| **Redirect** | `GET` | `/redirect` | Direct external ID resolution to Simkl record |
| **Anime** | `GET` | `/anime/{id}` | Detailed anime item metadata |
| **Anime** | `GET` | `/anime/episodes/{id}` | Episode list for an anime |
| **Sync** | `GET` | `/sync/activities` | Last activity timestamps (crucial for polling) |
| **Sync** | `GET` | `/sync/all-items/{type}/{status}` | Get all items by type and status |
| **Sync** | `POST` | `/sync/add-to-list` | Move items to watchlist status |
| **Sync** | `POST` | `/sync/history` | Add items or episodes to watch history |
| **Sync** | `POST` | `/sync/history/remove` | Remove items from watch history |
| **Sync** | `POST` | `/sync/ratings` | Rate items |
| **Sync** | `POST` | `/sync/ratings/remove` | Remove ratings |
| **Scrobble** | `POST` | `/scrobble/start` | Start scrobbling playback session |
| **Scrobble** | `POST` | `/scrobble/stop` | Stop scrobbling playback session |
| **Scrobble** | `POST` | `/scrobble/pause` | Pause scrobbling playback session |
| **Users** | `GET` | `/users/settings` | Authenticated user settings and profile |
| **OAuth2** | `GET` | `/oauth2/authorize` | Authorization code consent URL |
| **OAuth2** | `POST` | `/oauth2/token` | Exchange code or refresh token |
| **OAuth2** | `POST` | `/oauth2/device` | Request device code & user PIN |
| **OAuth2** | `POST` | `/oauth2/revoke` | Revoke active access or refresh token |

---

## 2. Required Request Parameters & Headers

### URL Parameters (Required on every call)
- `client_id`: Developer's application Client ID from developer settings.
- `app-name`: Short identifier of the client application (e.g., `dango`).
- `app-version`: Semantic version of the application (e.g., `1.0.0`).

### HTTP Headers
- `User-Agent`: Descriptive user agent string (`go-simkl/0.1.0`).
- `simkl-api-key`: Injected with the `client_id`.
- `Content-Type`: `application/json` for requests with a body.
- `Authorization`: `Bearer <access_token>` on authenticated endpoints (omitted on public catalog lookups).

---

## 3. Rate Limits & Quotas

- **Per-Second Burst**: 10 GET/s, 1 POST/s. Returns HTTP 429 with `{"error": "rate_limit"}`.
- **Daily Quota**: 500 to 10,000 req/day per user. Returns HTTP 429 with `{"error": "user_limit_exceeded"}` and `Retry-After: <seconds>`.
- **Headers**:
  - `X-RateLimit-Limit`: Daily quota limit.
  - `X-RateLimit-Remaining`: Daily quota remaining.
  - `Retry-After`: Seconds until quota reset.
