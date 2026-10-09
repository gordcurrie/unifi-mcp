package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const initializeBody = `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"0"}}}`

func TestNewHTTPHandlerCrossOrigin(t *testing.T) {
	s := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0"}, nil)
	srv := httptest.NewTestServer(t, newHTTPHandler(s))
	// srv.Client sets srv.URL and routes requests over the in-memory network.
	// It has no timeout by default, so set one to keep a hung handler from
	// stalling the test run.
	httpClient := srv.Client()
	httpClient.Timeout = 10 * time.Second

	tests := []struct {
		name       string
		headers    map[string]string
		wantStatus int
	}{
		{
			name:       "non-browser client allowed",
			wantStatus: http.StatusOK,
		},
		{
			name:       "same origin allowed",
			headers:    map[string]string{"Origin": srv.URL},
			wantStatus: http.StatusOK,
		},
		{
			name:       "untrusted origin rejected",
			headers:    map[string]string{"Origin": "https://evil.example"},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "cross-site fetch rejected",
			headers:    map[string]string{"Sec-Fetch-Site": "cross-site"},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, srv.URL, strings.NewReader(initializeBody))
			if err != nil {
				t.Fatalf("NewRequest: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json, text/event-stream")
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			resp, err := httpClient.Do(req) /* #nosec G704 */ //nolint:gosec // G704: URL is the in-memory httptest server
			if err != nil {
				t.Fatalf("Do: %v", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
