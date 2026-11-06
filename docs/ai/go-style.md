# Go Style Guide  -  go-simkl

Coding standards and idiom conventions for `go-simkl`.

---

## 1. Core Principles

- **Idiomatic Standard Library Go**: Follow standard Go practices as exemplified by the standard library (`net/http`, `io`, `context`).
- **Zero Third-Party Dependencies**: No external libraries in `go.mod`.
- **Pure Go**: Strict `CGO_ENABLED=0` compatibility.

---

## 2. Functions & Context

1. **`context.Context` First**:
   Any function performing network communication, I/O, or parsing that could block must take `ctx context.Context` as its first parameter.
   ```go
   func (s *AnimeService) GetByID(ctx context.Context, id int) (*Anime, *Response, error)
   ```
2. **Never Ignore Context Cancellation**:
   Always pass the context down to `http.NewRequestWithContext`.

---

## 3. Error Handling & Wrapping

1. **Explicit Error Wrapping**:
   Always wrap errors with `%w` using the library prefix:
   ```go
   return nil, nil, fmt.Errorf("simkl: %s: %w", op, err)
   ```
2. **Sentinel Errors & Typed Errors**:
   - `ErrRateLimited`: Thrown on HTTP 429 with rate limit and retry metadata.
   - `ErrUnauthorized`: Thrown on HTTP 401.
   - `ErrNotFound`: Thrown on HTTP 404.
   - Helper predicates: `IsRateLimited(err) bool`, `IsNotFound(err) bool`.

---

## 4. Construction & Functional Options

- Client constructors use functional options:
  ```go
  func NewClient(clientID string, opts ...ClientOption) *Client
  ```
- Options modify unexported fields safely without global state:
  - `WithHTTPClient(c *http.Client)`
  - `WithBaseURL(rawURL string)`
  - `WithToken(token string)`
  - `WithAppInfo(name, version string)`
  - `WithUserAgent(ua string)`
  - `WithMaxBodyBytes(limit int64)`

---

## 5. Memory Efficiency & Streaming Deserialization

- Always use `json.NewDecoder(io.LimitReader(body, maxBytes))` to decode JSON responses.
- Do not retain raw HTTP response bytes in memory for successful responses.
