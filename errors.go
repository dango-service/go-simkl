package simkl

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrRateLimited is returned when Simkl responds with HTTP 429 Too Many Requests.
	ErrRateLimited = errors.New("simkl: rate limit or daily quota exceeded")

	// ErrUnauthorized is returned when Simkl responds with HTTP 401 Unauthorized.
	ErrUnauthorized = errors.New("simkl: unauthorized or invalid credentials")

	// ErrNotFound is returned when Simkl responds with HTTP 404 Not Found.
	ErrNotFound = errors.New("simkl: resource not found")
)

// ErrorResponse represents an error returned by the Simkl API.
type ErrorResponse struct {
	Response         *http.Response `json:"-"`
	ErrorType        string         `json:"error"`
	Code             int            `json:"code,omitempty"`
	Message          string         `json:"message,omitempty"`
	ErrorDescription string         `json:"error_description,omitempty"`
	RateLimit        RateLimitInfo  `json:"-"`
}

func (r *ErrorResponse) Error() string {
	var statusCode int
	if r.Response != nil {
		statusCode = r.Response.StatusCode
	} else if r.Code != 0 {
		statusCode = r.Code
	}

	msg := r.Message
	if msg == "" && r.ErrorDescription != "" {
		msg = r.ErrorDescription
	}

	if msg != "" {
		return fmt.Sprintf("simkl: HTTP %d: %s (%s)", statusCode, r.ErrorType, msg)
	}
	if r.ErrorType != "" {
		return fmt.Sprintf("simkl: HTTP %d: %s", statusCode, r.ErrorType)
	}
	return fmt.Sprintf("simkl: HTTP %d", statusCode)
}

// IsRateLimited reports whether err was caused by an HTTP 429 rate limit or quota exhaustion.
func IsRateLimited(err error) bool {
	if errors.Is(err, ErrRateLimited) {
		return true
	}
	var errResp *ErrorResponse
	return errors.As(err, &errResp) && errResp.Response != nil && errResp.Response.StatusCode == http.StatusTooManyRequests
}

// IsNotFound reports whether err was caused by an HTTP 404 Not Found response.
func IsNotFound(err error) bool {
	if errors.Is(err, ErrNotFound) {
		return true
	}
	var errResp *ErrorResponse
	return errors.As(err, &errResp) && errResp.Response != nil && errResp.Response.StatusCode == http.StatusNotFound
}

// IsUnauthorized reports whether err was caused by an HTTP 401 Unauthorized response.
func IsUnauthorized(err error) bool {
	if errors.Is(err, ErrUnauthorized) {
		return true
	}
	var errResp *ErrorResponse
	return errors.As(err, &errResp) && errResp.Response != nil && errResp.Response.StatusCode == http.StatusUnauthorized
}
