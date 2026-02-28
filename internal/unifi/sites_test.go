package unifi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func sitesListHandler(sites []map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/integration/v1/sites" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":       sites,
			"totalCount": len(sites),
		})
	}
}

func TestGetSite(t *testing.T) {
	t.Run("returns site by explicit ID", func(t *testing.T) {
		client := newTestClient(t, sitesListHandler([]map[string]any{
			{"id": "abc123", "name": "default"},
		}))
		site, err := client.GetSite(context.Background(), "abc123")
		if err != nil {
			t.Fatalf("GetSite: %v", err)
		}
		if site.ID != "abc123" {
			t.Errorf("got ID %q, want %q", site.ID, "abc123")
		}
		if site.Name != "default" {
			t.Errorf("got Name %q, want %q", site.Name, "default")
		}
	})

	t.Run("falls back to client default site", func(t *testing.T) {
		client := newTestClient(t, sitesListHandler([]map[string]any{
			{"id": "test-site-id", "name": "home"},
		}))
		site, err := client.GetSite(context.Background(), "")
		if err != nil {
			t.Fatalf("GetSite: %v", err)
		}
		if site.ID != "test-site-id" {
			t.Errorf("got ID %q, want %q", site.ID, "test-site-id")
		}
	})

	t.Run("returns error when site not in list", func(t *testing.T) {
		client := newTestClient(t, sitesListHandler([]map[string]any{
			{"id": "other-site", "name": "other"},
		}))
		_, err := client.GetSite(context.Background(), "missing")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.GetSite(context.Background(), "abc123")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
