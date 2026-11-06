# Architectural Invariants  -  go-simkl

Structural rules and invariants that must be upheld at all times in `go-simkl`.

---

## 1. Independence and Decoupling

- **Completely Agnostic of Dango**: `go-simkl` is an autonomous library.
  Under no circumstances may it import `Dango` packages, reference Plex, SQLite,
  media matching algorithms, or webhook handlers.
- **Zero Global State**: Multiple `Client` instances must be able to run concurrently
  and safely within the same process without interference.

---

## 2. Layer Separation: Transport vs Typed Public API

1. **Internal Transport Layer (`doRequest`)**:
   - Handles HTTP communication (`GET`, `POST`, `DELETE`).
   - Injects required query parameters (`client_id`, `app-name`, `app-version`) on every call.
   - Sets required headers (`User-Agent`, `simkl-api-key`, `Authorization: Bearer <token>`, `Content-Type`).
   - Automatically omits `Authorization` header on public cached catalog lookups as recommended by Simkl.
   - Parses Simkl quota headers (`X-RateLimit-*`, `Retry-After`).
   - Deserializes JSON via streaming `json.NewDecoder` bounded by `io.LimitReader`.
   - Normalizes HTTP status codes and error payloads into Go error types.
2. **Typed Public Services**:
   - `SearchService`: Full-text search and cross-catalog external ID resolution.
   - `AnimeService`: Anime details and episode lists.
   - `SyncService`: Watchlist tracking, activity polling check, history sync, rating sync.
   - `ScrobbleService`: Real-time media playback scrobbling sessions.
   - `UserService`: User account settings and profiles.
   - `AuthService`: OAuth2 V2 code exchange, token refresh, token revoke, and Device/PIN flow.

---

## 3. Rate Limit Handling

- Simkl enforces two distinct limits:
  1. **Per-Second Rate**: 10 GET/s, 1 POST/s (`429 rate_limit`). Transient, clears in ~1s.
  2. **Daily User Quota**: 500-10,000 req/day (`429 user_limit_exceeded`). Resets at midnight New York time.
- The transport layer extracts `X-RateLimit-Limit`, `X-RateLimit-Remaining`, and `Retry-After` into `RateLimitInfo`.
- The library does not block or sleep invisibly; it returns enriched `ErrRateLimited` with duration and quota details, allowing consumers to make informed retry and backoff decisions.

---

## 4. Credential Protection

- Access tokens and refresh tokens must never be written to logs in plaintext.
- The `Token` struct implements `fmt.Stringer` and `fmt.GoStringer` returning `[REDACTED]`.
- Tokens are stripped of leading and trailing whitespace upon option injection.
