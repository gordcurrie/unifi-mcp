// Package unifi provides a client for the UniFi Network API running on a UCG-Max.
// It supports both the v1 REST API (/v1/sites/{siteID}/...) and the
// legacy local API (/api/s/{site}/...), authenticated via X-API-Key.
package unifi

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is a UniFi Network API client.
type Client struct {
	baseURL    string
	apiKey     string
	siteID     string
	httpClient *http.Client
}

// v1SingleResponse wraps the v1 API single-object envelope: {"data": T}.
type v1SingleResponse[T any] struct {
	Data T `json:"data"`
}

// v1ListResponse wraps the v1 API list envelope: {"data": [T], "totalCount": N}.
type v1ListResponse[T any] struct {
	Data       []T `json:"data"`
	TotalCount int `json:"totalCount"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	Count      int `json:"count"`
}

// legacyMeta is the meta field in legacy API responses.
type legacyMeta struct {
	RC  string `json:"rc"`
	Msg string `json:"msg,omitempty"`
}

// legacyResponse wraps legacy API list envelopes: {"data": [T], "meta": {...}}.
type legacyResponse[T any] struct {
	Data []T        `json:"data"`
	Meta legacyMeta `json:"meta"`
}

// NewClient creates a UniFi API client.
// baseURL should be the full proxy/network base, e.g. "https://192.168.1.1/proxy/network".
// siteID is the default site UUID used when tools omit the site_id parameter.
// Set insecure to true to skip TLS verification for self-signed certificates (UCG-Max default).
func NewClient(baseURL, apiKey, siteID string, insecure bool) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("UNIFI_BASE_URL is required")
	}
	if apiKey == "" {
		return nil, errors.New("UNIFI_API_KEY is required")
	}
	if siteID == "" {
		return nil, errors.New("UNIFI_SITE_ID is required")
	}

	transport := http.DefaultTransport
	if insecure {
		if base, ok := http.DefaultTransport.(*http.Transport); ok {
			t := base.Clone()
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = &tls.Config{} //nolint:gosec // populated below
			}
			//nolint:gosec // G402: InsecureSkipVerify is only set when UNIFI_INSECURE=true, explicit user opt-in
			t.TLSClientConfig.InsecureSkipVerify = true // #nosec G402
			transport = t
		} else {
			// Fallback if DefaultTransport has been replaced by something other than *http.Transport.
			//nolint:gosec // G402: InsecureSkipVerify is only set when UNIFI_INSECURE=true, explicit user opt-in
			transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402
			}
		}
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		siteID:  siteID,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}, nil
}

// site returns the provided siteID if non-empty, otherwise the client default.
func (c *Client) site(siteID string) string {
	if siteID != "" {
		return siteID
	}
	return c.siteID
}

// get performs a GET request to the given path (relative to baseURL).
func (c *Client) get(ctx context.Context, path string) ([]byte, error) {
	return c.do(ctx, http.MethodGet, path, nil)
}

// post performs a POST request with no body to the given path.
func (c *Client) post(ctx context.Context, path string) ([]byte, error) {
	return c.do(ctx, http.MethodPost, path, nil)
}

// postWithBody performs a POST request with a JSON-encoded body.
func (c *Client) postWithBody(ctx context.Context, path string, body any) ([]byte, error) {
	return c.do(ctx, http.MethodPost, path, body)
}

// put performs a PUT request with a JSON-encoded body.
func (c *Client) put(ctx context.Context, path string, body any) ([]byte, error) {
	return c.do(ctx, http.MethodPut, path, body)
}

// do executes an HTTP request and returns the raw response body.
func (c *Client) do(ctx context.Context, method, path string, body any) (_ []byte, retErr error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	parsedURL, err := url.ParseRequestURI(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("build request URL: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req) /* #nosec G704 */ //nolint:gosec // G704: URL is constructed from UNIFI_BASE_URL which the user must explicitly supply
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && retErr == nil {
			retErr = fmt.Errorf("close response body: %w", cerr)
		}
	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

// decodeV1Single decodes a v1 single-object response envelope.
func decodeV1Single[T any](data []byte) (T, error) {
	var resp v1SingleResponse[T]
	if err := json.Unmarshal(data, &resp); err != nil {
		var zero T
		return zero, fmt.Errorf("decode v1 response: %w", err)
	}
	return resp.Data, nil
}

// decodeV1List decodes a v1 list response envelope.
func decodeV1List[T any](data []byte) ([]T, error) {
	var resp v1ListResponse[T]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode v1 list response: %w", err)
	}
	return resp.Data, nil
}

// decodeLegacy decodes a legacy API list response envelope and validates meta.rc.
func decodeLegacy[T any](data []byte) ([]T, error) {
	var resp legacyResponse[T]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decode legacy response: %w", err)
	}
	if resp.Meta.RC != "" && resp.Meta.RC != "ok" {
		return nil, fmt.Errorf("controller error: rc=%s msg=%s", resp.Meta.RC, resp.Meta.Msg)
	}
	return resp.Data, nil
}

// checkLegacyRC decodes just the meta envelope from a legacy command response
// and returns an error when meta.rc is not "ok".
func checkLegacyRC(data []byte) error {
	type envelope struct {
		Meta legacyMeta `json:"meta"`
	}
	var resp envelope
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("decode legacy response: %w", err)
	}
	if resp.Meta.RC != "" && resp.Meta.RC != "ok" {
		return fmt.Errorf("controller error: rc=%s msg=%s", resp.Meta.RC, resp.Meta.Msg)
	}
	return nil
}
