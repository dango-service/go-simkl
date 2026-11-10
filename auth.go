package simkl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// AuthService handles OAuth2 Authorization Code, token refresh, token revocation, and Device/PIN flows.
type AuthService struct {
	client *Client
}

// Token represents an OAuth2 access and refresh token pair.
type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// String implements fmt.Stringer, redacting credentials to prevent log leakage.
func (t Token) String() string {
	return "[REDACTED]"
}

// GoString implements fmt.GoStringer, redacting credentials to prevent log leakage.
func (t Token) GoString() string {
	return "[REDACTED]"
}

// DeviceCodeResponse represents the parameters returned by the Device Authorization endpoint.
type DeviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// AuthorizeURL returns the user-facing consent URL to begin the OAuth2 Authorization Code grant.
// If codeChallenge is provided, standard PKCE parameters (S256) are appended to the query.
func (s *AuthService) AuthorizeURL(redirectURI, state string, codeChallenge ...string) string {
	baseURL := "https://simkl.com/oauth/authorize"
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", s.client.clientID)
	if redirectURI != "" {
		params.Set("redirect_uri", redirectURI)
	}
	if state != "" {
		params.Set("state", state)
	}
	if len(codeChallenge) > 0 && strings.TrimSpace(codeChallenge[0]) != "" {
		params.Set("code_challenge", strings.TrimSpace(codeChallenge[0]))
		params.Set("code_challenge_method", "S256")
	}
	return baseURL + "?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for an OAuth2 token pair.
// If clientSecret is empty (e.g. PKCE applications), it is omitted from the request body.
// If codeVerifier is provided, it is included in the payload.
func (s *AuthService) ExchangeCode(ctx context.Context, code, redirectURI, clientSecret string, codeVerifier ...string) (*Token, *Response, error) {
	reqBody := map[string]string{
		"grant_type":   "authorization_code",
		"code":         strings.TrimSpace(code),
		"client_id":    s.client.clientID,
		"redirect_uri": redirectURI,
	}
	if strings.TrimSpace(clientSecret) != "" {
		reqBody["client_secret"] = strings.TrimSpace(clientSecret)
	}
	if len(codeVerifier) > 0 && strings.TrimSpace(codeVerifier[0]) != "" {
		reqBody["code_verifier"] = strings.TrimSpace(codeVerifier[0])
	}
	return s.requestToken(ctx, reqBody)
}

// RefreshToken exchanges an active refresh token for a fresh OAuth2 token pair.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, clientSecret string) (*Token, *Response, error) {
	reqBody := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": strings.TrimSpace(refreshToken),
		"client_id":     s.client.clientID,
		"client_secret": strings.TrimSpace(clientSecret),
	}
	return s.requestToken(ctx, reqBody)
}

// RequestDeviceCode initiates the Device / PIN flow for devices without a web browser (e.g., CLIs, TVs).
func (s *AuthService) RequestDeviceCode(ctx context.Context) (*DeviceCodeResponse, *Response, error) {
	reqBody := map[string]string{
		"client_id": s.client.clientID,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("simkl: marshal device code request: %w", err)
	}

	var res DeviceCodeResponse
	resp, err := s.client.doRequest(ctx, http.MethodPost, "/oauth2/device", bytes.NewReader(bodyBytes), &res)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: request device code: %w", err)
	}
	return &res, resp, nil
}

// PollDeviceToken polls the token endpoint using the device code until user approval or timeout.
func (s *AuthService) PollDeviceToken(ctx context.Context, deviceCode string) (*Token, *Response, error) {
	reqBody := map[string]string{
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
		"device_code": strings.TrimSpace(deviceCode),
		"client_id":   s.client.clientID,
	}
	return s.requestToken(ctx, reqBody)
}

// RevokeToken revokes an active access token or refresh token.
func (s *AuthService) RevokeToken(ctx context.Context, token string) (*Response, error) {
	reqBody := map[string]string{
		"token":     strings.TrimSpace(token),
		"client_id": s.client.clientID,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("simkl: marshal revoke request: %w", err)
	}

	resp, err := s.client.doRequest(ctx, http.MethodPost, "/oauth2/revoke", bytes.NewReader(bodyBytes), nil)
	if err != nil {
		return resp, fmt.Errorf("simkl: revoke token: %w", err)
	}
	return resp, nil
}

func (s *AuthService) requestToken(ctx context.Context, payload map[string]string) (*Token, *Response, error) {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("simkl: marshal token request: %w", err)
	}

	var token Token
	resp, err := s.client.doRequest(ctx, http.MethodPost, "/oauth2/token", bytes.NewReader(bodyBytes), &token)
	if err != nil {
		return nil, resp, fmt.Errorf("simkl: token request: %w", err)
	}
	return &token, resp, nil
}
