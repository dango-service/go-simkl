package simkl

import (
	"net/url"
	"time"
)

// SyncOption customizes a sync query.
type SyncOption func(url.Values)

// WithDateFrom filters items modified since the specified timestamp.
func WithDateFrom(t time.Time) SyncOption {
	return func(v url.Values) {
		if !t.IsZero() {
			v.Set("date_from", t.UTC().Format(time.RFC3339))
		}
	}
}

// WithNextWatchInfo requests next-episode watch metadata in watchlists.
func WithNextWatchInfo(include bool) SyncOption {
	return func(v url.Values) {
		if include {
			v.Set("next_watch_info", "yes")
		}
	}
}
