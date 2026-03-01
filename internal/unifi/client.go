// Package unifi provides a client for the UniFi Network Integration API running on a UCG-Max.
// It targets the v1 REST API at /proxy/network/integration/v1/..., authenticated via X-API-Key.
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

// maxResponseBytes caps how much data we read from any single API response.
// 10 MiB is well above any realistic UniFi response payload.
const maxResponseBytes = 10 << 20 // 10 MiB

// Client is a UniFi Network API client.
type Client struct {
	baseURL    string
	apiKey     string
	siteID     string
	httpClient *http.Client
}

// NewClient creates a UniFi Integration API client.
// baseURL should be the full proxy/network base, e.g. "https://192.168.1.1/proxy/network".
// siteID is the site UUID (from Settings → Sites) used when tools omit the site_id parameter.
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
				t.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			}
			//nolint:gosec // G402: InsecureSkipVerify is only set when UNIFI_INSECURE=true, explicit user opt-in
			t.TLSClientConfig.InsecureSkipVerify = true // #nosec G402
			transport = t
		} else {
			// Fallback if DefaultTransport has been replaced by something other than *http.Transport.
			//nolint:gosec // G402: InsecureSkipVerify is only set when UNIFI_INSECURE=true, explicit user opt-in
			transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // #nosec G402
					MinVersion:         tls.VersionTLS12,
				},
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

// getWithQuery performs a GET request appending optional offset/limit query parameters.
// A value of 0 means "omit the parameter and let the API use its default".
// Negative values are invalid and return an error.
func (c *Client) getWithQuery(ctx context.Context, path string, offset, limit int) ([]byte, error) {
	if offset < 0 || limit < 0 {
		return nil, fmt.Errorf("getWithQuery: offset and limit must be >= 0 (got offset=%d, limit=%d)", offset, limit)
	}
	if offset == 0 && limit == 0 {
		return c.get(ctx, path)
	}
	q := url.Values{}
	if offset > 0 {
		q.Set("offset", fmt.Sprintf("%d", offset))
	}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	return c.do(ctx, http.MethodGet, path+"?"+q.Encode(), nil)
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

// delete performs a DELETE request to the given path.
func (c *Client) delete(ctx context.Context, path string) error {
	_, err := c.do(ctx, http.MethodDelete, path, nil)
	return err
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

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

// decodeV1 unmarshals a raw integration v1 response (no envelope) directly into T.
func decodeV1[T any](data []byte) (T, error) {
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("decode v1 response: %w", err)
	}
	return result, nil
}

// decodeV1List decodes an integration v1 paginated list envelope into a Page[T].
func decodeV1List[T any](data []byte) (Page[T], error) {
	var page Page[T]
	if err := json.Unmarshal(data, &page); err != nil {
		return page, fmt.Errorf("decode v1 list response: %w", err)
	}
	return page, nil
}
