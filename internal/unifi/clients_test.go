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
			if r.URL.Path != "/integration/v1/sites/11111111-1111-4111-8111-111111111111/clients" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "00000000-0000-4000-8000-000000000008", "macAddress": "aa:bb:cc:00:00:01", "type": "WIRED", "ipAddress": "192.168.1.100"},
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
			if r.URL.Path != "/integration/v1/sites/11111111-1111-4111-8111-111111111111/clients/00000000-0000-4000-8000-000000000008/actions" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			gotMethod = r.Method
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
			w.WriteHeader(http.StatusOK)
		})
		req := GuestAuthRequest{Action: "AUTHORIZE_GUEST_ACCESS", TimeLimitMinutes: 120}
		if err := client.AuthorizeGuestClient(context.Background(), "", "00000000-0000-4000-8000-000000000008", req); err != nil {
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
		if err := client.AuthorizeGuestClient(context.Background(), "", "00000000-0000-4000-8000-000000000008", GuestAuthRequest{Action: "AUTHORIZE_GUEST_ACCESS"}); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetClient(t *testing.T) {
	t.Run("decodes single client", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/11111111-1111-4111-8111-111111111111/clients/00000000-0000-4000-8000-000000000009" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "00000000-0000-4000-8000-000000000009", "macAddress": "aa:bb:cc:00:00:99", "type": "WIRELESS", "ipAddress": "10.0.0.5",
			})
		})
		c, err := client.GetClient(context.Background(), "", "00000000-0000-4000-8000-000000000009")
		if err != nil {
			t.Fatalf("GetClient: %v", err)
		}
		if c.ID != "00000000-0000-4000-8000-000000000009" {
			t.Errorf("got ID %q, want 00000000-0000-4000-8000-000000000009", c.ID)
		}
		if c.IP != "10.0.0.5" {
			t.Errorf("got IP %q, want 10.0.0.5", c.IP)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.GetClient(context.Background(), "", "00000000-0000-4000-8000-000000000009")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
