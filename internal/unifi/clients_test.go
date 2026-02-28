package unifi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListActiveClients(t *testing.T) {
	t.Run("decodes active client list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/sites/test-site-id/clients/active" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"mac": "aa:bb:cc:00:00:01", "hostname": "laptop", "ip": "192.168.1.100"},
					{"mac": "aa:bb:cc:00:00:02", "hostname": "phone", "ip": "192.168.1.101"},
				},
				"totalCount": 2,
			})
		})
		clients, err := client.ListActiveClients(context.Background(), "")
		if err != nil {
			t.Fatalf("ListActiveClients: %v", err)
		}
		if len(clients) != 2 {
			t.Fatalf("got %d clients, want 2", len(clients))
		}
		if clients[0].MAC != "aa:bb:cc:00:00:01" {
			t.Errorf("got MAC %q, want %q", clients[0].MAC, "aa:bb:cc:00:00:01")
		}
		if clients[1].Hostname != "phone" {
			t.Errorf("got hostname %q, want %q", clients[1].Hostname, "phone")
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.ListActiveClients(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestListKnownClients(t *testing.T) {
	t.Run("decodes known client list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/sites/test-site-id/clients/history" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"mac": "dd:ee:ff:00:00:01", "name": "old-tablet"},
				},
				"totalCount": 1,
			})
		})
		clients, err := client.ListKnownClients(context.Background(), "")
		if err != nil {
			t.Fatalf("ListKnownClients: %v", err)
		}
		if len(clients) != 1 {
			t.Fatalf("got %d clients, want 1", len(clients))
		}
		if clients[0].MAC != "dd:ee:ff:00:00:01" {
			t.Errorf("got MAC %q, want %q", clients[0].MAC, "dd:ee:ff:00:00:01")
		}
	})
}

func testStamgrCommand(t *testing.T, invoker func(c *Client) error, wantCmd, wantMAC string) {
	t.Helper()
	var gotCmd, gotMAC string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/s/test-site-id/cmd/stamgr" || r.Method != http.MethodPost {
			http.Error(w, "unexpected path/method", http.StatusBadRequest)
			return
		}
		var req clientCmdRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		gotCmd = req.Cmd
		gotMAC = req.MAC
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(legacyOK([]any{}))
	})
	if err := invoker(client); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCmd != wantCmd {
		t.Errorf("cmd: got %q, want %q", gotCmd, wantCmd)
	}
	if gotMAC != wantMAC {
		t.Errorf("mac: got %q, want %q", gotMAC, wantMAC)
	}
}

func TestBlockClient(t *testing.T) {
	testStamgrCommand(t, func(c *Client) error {
		return c.BlockClient(context.Background(), "", "aa:bb:cc:dd:ee:01")
	}, "block-sta", "aa:bb:cc:dd:ee:01")
}

func TestUnblockClient(t *testing.T) {
	testStamgrCommand(t, func(c *Client) error {
		return c.UnblockClient(context.Background(), "", "aa:bb:cc:dd:ee:02")
	}, "unblock-sta", "aa:bb:cc:dd:ee:02")
}

func TestKickClient(t *testing.T) {
	testStamgrCommand(t, func(c *Client) error {
		return c.KickClient(context.Background(), "", "aa:bb:cc:dd:ee:03")
	}, "kick-sta", "aa:bb:cc:dd:ee:03")
}

func TestForgetClient(t *testing.T) {
	testStamgrCommand(t, func(c *Client) error {
		return c.ForgetClient(context.Background(), "", "aa:bb:cc:dd:ee:04")
	}, "forget-sta", "aa:bb:cc:dd:ee:04")
}
