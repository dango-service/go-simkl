# Resource Budget and Limits  -  go-simkl

Resource utilization constraints and performance limits for `go-simkl`.

---

## 1. Memory Budget

- **Peak RSS**: <= 5 MiB during regular background operations.
- **Maximum JSON Response Buffer**: Defaults to 10 MiB via `io.LimitReader(body, maxBytes)`. This bounds heap allocations and prevents memory exhaustion attacks.
- **Zero In-Memory Caching in Client**: The client does not cache responses in memory; caching strategies are delegated to the consuming application (`Dango`).

---

## 2. Network & Throughput Invariants

- **Connection Reuse**: The client utilizes `http.Transport` default connection pooling to reuse TCP/TLS sessions.
- **Idle Connections**: Response bodies are always drained (`io.Copy(io.Discard, ...)`) and closed properly to ensure TCP connections return to the pool.
- **Dual Rate Limiting**: The client parses both per-second burst limits (10 GET/sec, 1 POST/sec) and daily user quotas (500-10,000 req/day).
