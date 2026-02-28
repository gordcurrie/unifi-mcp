package unifi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListClients(t *testing.T) {
	t.Run("decodes client list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/clients" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "c-1", "macAddress": "aa:bb:cc:00:00:01", "type": "WIRED", "ipAddress": "192.168.1.100"},
					{"id": "c-2", "macAddress": "aa:bb:cc:00:00:02", "type": "WIRELESS"},
				},
				"totalCount": 2,
			})
		})
		clients, err := client.ListClients(context.Background(), "")
		if err != nil {
			t.Fatalf("ListClients: %v", err)
		}
		if len(clients) != 2 {
			t.Fatalf("got %d clients, want 2", len(clients))
		}
		if clients[0].MAC != "aa:bb:cc:00:00:01" {
			t.Errorf("got MAC %q, want %q", clients[0].MAC, "aa:bb:cc:00:00:01")
		}
		if clients[1].Type != "WIRELESS" {
			t.Errorf("got Type %q, want WIRELESS", clients[1].Type)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.ListClients(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetClient(t *testing.T) {
	t.Run("decodes single client", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/clients/c-99" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "c-99", "macAddress": "aa:bb:cc:00:00:99", "type": "WIRELESS", "ipAddress": "10.0.0.5",
			})
		})
		c, err := client.GetClient(context.Background(), "", "c-99")
		if err != nil {
			t.Fatalf("GetClient: %v", err)
		}
		if c.ID != "c-99" {
			t.Errorf("got ID %q, want c-99", c.ID)
		}
		if c.IP != "10.0.0.5" {
			t.Errorf("got IP %q, want 10.0.0.5", c.IP)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.GetClient(context.Background(), "", "c-99")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
