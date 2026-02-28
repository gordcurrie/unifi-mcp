package unifi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client, err := NewClient(srv.URL, "test-api-key", "test-site-id", false)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

func TestGetInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-API-Key") != "test-api-key" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if r.URL.Path != "/v1/info" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"applicationVersion": "9.0.92",
					"controllerType":     "UCG",
				},
			})
		})

		info, err := client.GetInfo(context.Background())
		if err != nil {
			t.Fatalf("GetInfo: %v", err)
		}
		if info.ApplicationVersion != "9.0.92" {
			t.Errorf("got version %q, want %q", info.ApplicationVersion, "9.0.92")
		}
		if info.ControllerType != "UCG" {
			t.Errorf("got type %q, want %q", info.ControllerType, "UCG")
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
			_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
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
		sites, err := client.ListSites(context.Background())
		if err != nil {
			t.Fatalf("ListSites: %v", err)
		}
		if len(sites) != 2 {
			t.Fatalf("got %d sites, want 2", len(sites))
		}
		if sites[0].ID != "site-1" {
			t.Errorf("got site[0].ID %q, want %q", sites[0].ID, "site-1")
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
		{"valid", "https://192.168.1.1/proxy/network", "key", "site", false},
		{"missing base url", "", "key", "site", true},
		{"missing api key", "https://192.168.1.1/proxy/network", "", "site", true},
		{"missing site id", "https://192.168.1.1/proxy/network", "key", "", true},
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

func TestSiteFallback(t *testing.T) {
	client, err := NewClient("https://192.168.1.1/proxy/network", "key", "default-site", false)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if got := client.site(""); got != "default-site" {
		t.Errorf(`site("") = %q, want %q`, got, "default-site")
	}
	if got := client.site("override"); got != "override" {
		t.Errorf(`site("override") = %q, want %q`, got, "override")
	}
}
