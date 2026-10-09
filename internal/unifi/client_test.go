package unifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewTestServer(t, handler)
	// srv.Client sets srv.URL; its transport dials the in-memory network.
	transport := srv.Client().Transport
	client, err := NewClient(srv.URL, "test-api-key", "11111111-1111-4111-8111-111111111111", false)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	// Swap only the transport so the production client's timeout still applies.
	client.httpClient.Transport = transport
	return client
}

func TestGetInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-API-Key") != "test-api-key" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if r.URL.Path != "/integration/v1/info" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"applicationVersion": "9.0.92",
			})
		})

		info, err := client.GetInfo(context.Background())
		if err != nil {
			t.Fatalf("GetInfo: %v", err)
		}
		if info.ApplicationVersion != "9.0.92" {
			t.Errorf("got version %q, want %q", info.ApplicationVersion, "9.0.92")
		}
	})

	t.Run("api key header is sent", func(t *testing.T) {
		called := false
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			called = true
			got := r.Header.Get("X-API-Key")
			if got != "test-api-key" {
				t.Errorf("X-API-Key header: got %q, want %q", got, "test-api-key")
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"applicationVersion": ""})
		})
		_, _ = client.GetInfo(context.Background())
		if !called {
			t.Error("handler was never called")
		}
	})

	t.Run("non-2xx status returns error", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "internal error", http.StatusInternalServerError)
		})
		_, err := client.GetInfo(context.Background())
		if err == nil {
			t.Error("expected error for 500 response, got nil")
		}
	})
}

func TestListSites(t *testing.T) {
	t.Run("decodes list response", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "site-1", "name": "default"},
					{"id": "site-2", "name": "guest"},
				},
				"totalCount": 2,
			})
		})
		sites, err := client.ListSites(context.Background(), 0, 0)
		if err != nil {
			t.Fatalf("ListSites: %v", err)
		}
		if len(sites.Data) != 2 {
			t.Fatalf("got %d sites, want 2", len(sites.Data))
		}
		if sites.Data[0].ID != "site-1" {
			t.Errorf("got site[0].ID %q, want %q", sites.Data[0].ID, "site-1")
		}
		if sites.TotalCount != 2 {
			t.Errorf("got TotalCount %d, want 2", sites.TotalCount)
		}
	})
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		apiKey  string
		siteID  string
		wantErr bool
	}{
		{"valid", "https://192.168.1.1/proxy/network", "key", "11111111-1111-4111-8111-111111111111", false},
		{"missing base url", "", "key", "11111111-1111-4111-8111-111111111111", true},
		{"missing api key", "https://192.168.1.1/proxy/network", "", "11111111-1111-4111-8111-111111111111", true},
		{"missing site id", "https://192.168.1.1/proxy/network", "key", "", true},
		{"site id not a uuid", "https://192.168.1.1/proxy/network", "key", "default", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.baseURL, tt.apiKey, tt.siteID, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateIDs(t *testing.T) {
	const valid = "5f4d0e88-1234-4678-abcd-ef0123456789"
	tests := []struct {
		name    string
		ids     []string
		wantErr bool
	}{
		{"no ids", nil, false},
		{"canonical lowercase", []string{valid}, false},
		{"canonical uppercase", []string{strings.ToUpper(valid)}, false},
		{"multiple valid", []string{valid, valid}, false},
		{"empty", []string{""}, true},
		{"not a uuid", []string{"dev-1"}, true},
		{"second id invalid", []string{valid, "dev-1"}, true},
		{"braced form", []string{"{" + valid + "}"}, true},
		{"urn form", []string{"urn:uuid:" + valid}, true},
		{"no dashes", []string{strings.ReplaceAll(valid, "-", "")}, true},
		{"path traversal", []string{"../admin"}, true},
		{"reserved characters", []string{"dev/1"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIDs(tt.ids...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateIDs(%q) error = %v, wantErr %v", tt.ids, err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrInvalidID) {
				t.Errorf("error %v does not wrap ErrInvalidID", err)
			}
		})
	}
}

func TestSiteFallback(t *testing.T) {
	client, err := NewClient("https://192.168.1.1/proxy/network", "key", "00000000-0000-4000-8000-000000000002", false)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if got := client.site(""); got != "00000000-0000-4000-8000-000000000002" {
		t.Errorf(`site("") = %q, want %q`, got, "00000000-0000-4000-8000-000000000002")
	}
	if got := client.site("00000000-0000-4000-8000-000000000003"); got != "00000000-0000-4000-8000-000000000003" {
		t.Errorf(`site("00000000-0000-4000-8000-000000000003") = %q, want %q`, got, "00000000-0000-4000-8000-000000000003")
	}
}

func TestGetWithQuery(t *testing.T) {
	cases := []struct {
		name      string
		path      string
		offset    int
		limit     int
		wantQuery string // expected raw query string (empty = no query params)
		wantErr   bool
	}{
		{"zero/zero omits params", "/integration/v1/sites", 0, 0, "", false},
		{"offset only", "/integration/v1/sites", 10, 0, "offset=10", false},
		{"limit only", "/integration/v1/sites", 0, 25, "limit=25", false},
		{"both", "/integration/v1/sites", 10, 25, "limit=25&offset=10", false},
		{"path with existing query", "/integration/v1/sites?foo=bar", 5, 10, "foo=bar&limit=10&offset=5", false},
		{"negative offset returns error", "/integration/v1/sites", -1, 0, "", true},
		{"negative limit returns error", "/integration/v1/sites", 0, -1, "", true},
		{"limit at max is allowed", "/integration/v1/sites", 0, maxPageLimit, fmt.Sprintf("limit=%d", maxPageLimit), false},
		{"limit over max returns error", "/integration/v1/sites", 0, maxPageLimit + 1, "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotQuery string
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				gotQuery = r.URL.RawQuery
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]any{
					"data": []map[string]any{}, "totalCount": 0,
				})
			})
			_, err := client.getWithQuery(context.Background(), tc.path, tc.offset, tc.limit)
			if tc.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotQuery != tc.wantQuery {
				t.Errorf("RawQuery = %q, want %q", gotQuery, tc.wantQuery)
			}
		})
	}
}

func TestPathTraversalRejection(t *testing.T) {
	cases := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"normal path is accepted", "/integration/v1/sites/abc-123/devices/dev-456", false},
		{"dotdot traversal is rejected", "/integration/v1/sites/../admin", true},
		{"encoded dotdot traversal is rejected", "/integration/v1/sites/%2e%2e/admin", true},
		{"uppercase encoded dotdot traversal is rejected", "/integration/v1/sites/%2E%2E/admin", true},
		{"single dot is rejected", "/integration/v1/sites/./devices", true},
		{"encoded single dot is rejected", "/integration/v1/sites/%2e/devices", true},
		{"dotdot in resource id is rejected", "/integration/v1/sites/site-id/devices/../../../other", true},
		{"uuid path is accepted", "/integration/v1/sites/5f4d0e88-1234-5678-abcd-ef0123456789/devices/aabbccdd-1234-5678-abcd-ef0123456789", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]any{})
			})
			_, err := client.get(context.Background(), tc.path)
			if tc.wantErr && err == nil {
				t.Error("expected error for path traversal, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("unexpected error for valid path: %v", err)
			}
		})
	}
}
