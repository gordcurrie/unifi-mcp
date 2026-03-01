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
		clients, err := client.ListClients(context.Background(), "", 0, 0)
		if err != nil {
			t.Fatalf("ListClients: %v", err)
		}
		if len(clients.Data) != 2 {
			t.Fatalf("got %d clients, want 2", len(clients.Data))
		}
		if clients.Data[0].MAC != "aa:bb:cc:00:00:01" {
			t.Errorf("got MAC %q, want %q", clients.Data[0].MAC, "aa:bb:cc:00:00:01")
		}
		if clients.Data[1].Type != "WIRELESS" {
			t.Errorf("got Type %q, want WIRELESS", clients.Data[1].Type)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.ListClients(context.Background(), "", 0, 0)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAuthorizeGuestClient(t *testing.T) {
	t.Run("posts action and succeeds on 200", func(t *testing.T) {
		var gotBody GuestAuthRequest
		var gotMethod string
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/clients/c-1/actions" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			gotMethod = r.Method
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			w.WriteHeader(http.StatusOK)
		})
		req := GuestAuthRequest{Action: "AUTHORIZE_GUEST_ACCESS", TimeLimitMinutes: 120}
		if err := client.AuthorizeGuestClient(context.Background(), "", "c-1", req); err != nil {
			t.Fatalf("AuthorizeGuestClient: %v", err)
		}
		if gotMethod != http.MethodPost {
			t.Errorf("got method %q, want POST", gotMethod)
		}
		if gotBody.Action != "AUTHORIZE_GUEST_ACCESS" || gotBody.TimeLimitMinutes != 120 {
			t.Errorf("got body %+v, want {Action:AUTHORIZE_GUEST_ACCESS TimeLimitMinutes:120}", gotBody)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		if err := client.AuthorizeGuestClient(context.Background(), "", "c-1", GuestAuthRequest{Action: "AUTHORIZE_GUEST_ACCESS"}); err == nil {
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
