# Testing Strategy  -  `go-simkl`

This document defines testing standards for `go-simkl`, ensuring deterministic,
fast, isolated, and reproducible test suites without external network access.

---

## 1. Core Testing Principles

1. **Zero Real Network Traffic**: No unit test may connect to `https://api.simkl.com`. All network traffic is mocked locally via `httptest.Server`.
2. **Determinism and Speed**: The full suite must complete in less than 2 seconds.
3. **Table-Driven Verification**: Tests verify both what the client sends (URL query parameters, HTTP headers, request bodies) and what it deserializes (domain models, rate limit info, errors).
4. **Error Coverage**: Full testing of HTTP 401, 404, 412, and both variants of 429 (`rate_limit` and `user_limit_exceeded`).

---

## 2. Local Test Server (`httptest.Server`)

Example pattern:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Query().Get("client_id") != "test-client-id" {
        http.Error(w, `{"error":"client_id_failed","code":412}`, http.StatusPreconditionFailed)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`[{"title":"Attack on Titan"}]`))
}))
defer server.Close()

client := NewClient("test-client-id", WithBaseURL(server.URL))
```
