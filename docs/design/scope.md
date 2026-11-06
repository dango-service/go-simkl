# Library Scope and Boundaries  -  `go-simkl`

This document defines the functional perimeter of `go-simkl`.
The library does **not** attempt to cover the entirety of Simkl's API,
focusing strictly on the needs of **Dango** and media tracking tools.

---

## 1. Boundary Principle

The official Simkl API exposes over 57 endpoints, including custom user lists, static calendar dumps, DVD release discovery, and trending feeds.
Building an exhaustive client would:
- Multiply compiled binary size and build times.
- Increase runtime allocations with zero added value for scrobblers.
- Introduce an unmanageable maintenance burden.

Therefore, `go-simkl` follows a **strict minimum viable scope**, focusing on
searching, querying, and updating anime, show, and movie watch progress.

---

## 2. In-Scope Modules

1. **Search & Metadata**:
   - Title search across Anime, TV, and Movies.
   - Cross-catalog external ID resolution (`/search/id`) supporting MAL, AniList, Kitsu, TMDB, TVDB, and IMDb.
   - Detailed anime item metadata and episode lists.
2. **Sync & Tracking**:
   - Activity checking (`/sync/activities`) for efficient cache validation.
   - Watchlist queries (`/sync/all-items/{type}/{status}`).
   - Watchlist mutations (`/sync/add-to-list`).
   - Watch history mutations (`/sync/history`, `/sync/history/remove`).
   - Rating mutations (`/sync/ratings`, `/sync/ratings/remove`).
3. **Scrobble**:
   - Media playback session scrobbling (`/scrobble/start`, `/scrobble/stop`, `/scrobble/pause`).
4. **User & Auth**:
   - User profile settings (`/users/settings`).
   - OAuth2 Authorization Code grant and token refresh (`/oauth2/token`).
   - Device & PIN authorization flow (`/oauth2/device`).
   - Token revocation (`/oauth2/revoke`).

---

## 3. Out-of-Scope Modules

- Static calendar and trending dumps (`/calendar/*`, `/discover/*`).
- Custom user list management (`/lists/*`).
- DVD release feeds.
- Forum, comments, and social messaging features.
