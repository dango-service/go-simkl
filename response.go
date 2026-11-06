package simkl

import (
	"net/http"
	"strconv"
	"time"
)

// Response wraps the standard http.Response with parsed rate limit details.
type Response struct {
	*http.Response
	RateLimit RateLimitInfo
}

// RateLimitInfo contains quota and throttling metadata returned by the Simkl API.
type RateLimitInfo struct {
	Limit                 int
	Remaining             int
	RetryAfter            time.Duration
	IsPerSecondRateLimit  bool
	IsDailyQuotaExceeded  bool
}

// parseRateLimit extracts rate limit details from response headers and status codes.
func parseRateLimit(h http.Header, statusCode int, errIdentifier string) RateLimitInfo {
	var info RateLimitInfo

	if limitStr := h.Get("X-RateLimit-Limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			info.Limit = val
		}
	}

	if remainingStr := h.Get("X-RateLimit-Remaining"); remainingStr != "" {
		if val, err := strconv.Atoi(remainingStr); err == nil {
			info.Remaining = val
		}
	}

	if retryAfterStr := h.Get("Retry-After"); retryAfterStr != "" {
		if seconds, err := strconv.Atoi(retryAfterStr); err == nil {
			info.RetryAfter = time.Duration(seconds) * time.Second
		} else if parsedTime, err := http.ParseTime(retryAfterStr); err == nil {
			if duration := time.Until(parsedTime); duration > 0 {
				info.RetryAfter = duration
			}
		}
	}

	if statusCode == http.StatusTooManyRequests {
		if errIdentifier == "rate_limit" {
			info.IsPerSecondRateLimit = true
		} else if errIdentifier == "user_limit_exceeded" || errIdentifier == "app_limit_exceeded" {
			info.IsDailyQuotaExceeded = true
		}
	}

	return info
}
