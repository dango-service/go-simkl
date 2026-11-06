package simkl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// doRequest performs an HTTP request against the Simkl API, enforces memory limits,
// and decodes JSON responses or structured errors.
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, v any) (*Response, error) {
	relURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("simkl: invalid path %q: %w", path, err)
	}

	u := c.baseURL.ResolveReference(relURL)
	q := u.Query()

	// Inject required application parameters if not already specified.
	if c.clientID != "" && q.Get("client_id") == "" {
		q.Set("client_id", c.clientID)
	}
	if c.appName != "" && q.Get("app-name") == "" {
		q.Set("app-name", c.appName)
	}
	if c.appVersion != "" && q.Get("app-version") == "" {
		q.Set("app-version", c.appVersion)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("simkl: failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.clientID != "" {
		req.Header.Set("simkl-api-key", c.clientID)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("simkl: request cancelled: %w", ctx.Err())
		default:
			return nil, fmt.Errorf("simkl: request failed: %w", err)
		}
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	limitedReader := io.LimitReader(resp.Body, c.maxBytes)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		rateLimit := parseRateLimit(resp.Header, resp.StatusCode, "")
		response := &Response{
			Response:  resp,
			RateLimit: rateLimit,
		}

		if v != nil && resp.StatusCode != http.StatusNoContent {
			if err := json.NewDecoder(limitedReader).Decode(v); err != nil {
				return response, fmt.Errorf("simkl: failed to decode response: %w", err)
			}
		}
		return response, nil
	}

	// Handle error payloads
	var errResp ErrorResponse
	_ = json.NewDecoder(limitedReader).Decode(&errResp)
	errResp.Response = resp
	errResp.RateLimit = parseRateLimit(resp.Header, resp.StatusCode, errResp.ErrorType)

	response := &Response{
		Response:  resp,
		RateLimit: errResp.RateLimit,
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return response, fmt.Errorf("%w: %s", ErrRateLimited, errResp.Error())
	case http.StatusUnauthorized:
		return response, fmt.Errorf("%w: %s", ErrUnauthorized, errResp.Error())
	case http.StatusNotFound:
		return response, fmt.Errorf("%w: %s", ErrNotFound, errResp.Error())
	default:
		return response, &errResp
	}
}
