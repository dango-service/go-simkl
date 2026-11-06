package simkl

import (
	"net/url"
	"strconv"
)

// SearchOption customizes a search query.
type SearchOption func(url.Values)

// WithSearchLimit sets the maximum number of results to return (up to 50).
func WithSearchLimit(limit int) SearchOption {
	return func(v url.Values) {
		if limit > 0 {
			v.Set("limit", strconv.Itoa(limit))
		}
	}
}

// WithSearchExtended requests extended metadata (full overview, ratings, etc.).
func WithSearchExtended(extended bool) SearchOption {
	return func(v url.Values) {
		if extended {
			v.Set("extended", "full")
		}
	}
}
