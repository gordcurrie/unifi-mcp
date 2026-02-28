package unifi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func legacyOK(data any) map[string]any {
	return map[string]any{
		"data": data,
		"meta": map[string]any{"rc": "ok"},
	}
}

func TestListDevices(t *testing.T) {
	t.Run("decodes device list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/sites/test-site-id/devices" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "dev-1", "mac": "aa:bb:cc:dd:ee:01", "name": "switch1", "type": "usw"},
					{"id": "dev-2", "mac": "aa:bb:cc:dd:ee:02", "name": "ap1", "type": "uap"},
				},
				"totalCount": 2,
			})
		})
		devices, err := client.ListDevices(context.Background(), "")
		if err != nil {
			t.Fatalf("ListDevices: %v", err)
		}
		if len(devices) != 2 {
			t.Fatalf("got %d devices, want 2", len(devices))
		}
		if devices[0].ID != "dev-1" {
			t.Errorf("got devices[0].ID %q, want %q", devices[0].ID, "dev-1")
		}
		if devices[1].Type != "uap" {
			t.Errorf("got devices[1].Type %q, want %q", devices[1].Type, "uap")
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.ListDevices(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetDevice(t *testing.T) {
	t.Run("decodes single device", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/sites/test-site-id/devices/dev-99" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"id": "dev-99", "mac": "aa:bb:cc:00:00:99", "name": "gw"},
			})
		})
		dev, err := client.GetDevice(context.Background(), "", "dev-99")
		if err != nil {
			t.Fatalf("GetDevice: %v", err)
		}
		if dev.ID != "dev-99" {
			t.Errorf("got ID %q, want %q", dev.ID, "dev-99")
		}
	})
}

func testDevmgrCommand(t *testing.T, invoker func(c *Client) error, wantCmd, wantMAC string) {
	t.Helper()
	var gotCmd, gotMAC string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/s/test-site-id/cmd/devmgr" || r.Method != http.MethodPost {
			http.Error(w, "unexpected path/method", http.StatusBadRequest)
			return
		}
		var req deviceCmdRequest
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

func TestRestartDevice(t *testing.T) {
	testDevmgrCommand(t, func(c *Client) error {
		return c.RestartDevice(context.Background(), "", "aa:bb:cc:dd:ee:01")
	}, "restart", "aa:bb:cc:dd:ee:01")
}

func TestLocateDevice(t *testing.T) {
	testDevmgrCommand(t, func(c *Client) error {
		return c.LocateDevice(context.Background(), "", "aa:bb:cc:dd:ee:02")
	}, "set-locate", "aa:bb:cc:dd:ee:02")
}

func TestUnlocateDevice(t *testing.T) {
	testDevmgrCommand(t, func(c *Client) error {
		return c.UnlocateDevice(context.Background(), "", "aa:bb:cc:dd:ee:03")
	}, "unset-locate", "aa:bb:cc:dd:ee:03")
}

func TestUpgradeDevice(t *testing.T) {
	testDevmgrCommand(t, func(c *Client) error {
		return c.UpgradeDevice(context.Background(), "", "aa:bb:cc:dd:ee:04")
	}, "upgrade", "aa:bb:cc:dd:ee:04")
}

func TestForceReprovisionDevice(t *testing.T) {
	testDevmgrCommand(t, func(c *Client) error {
		return c.ForceReprovisionDevice(context.Background(), "", "aa:bb:cc:dd:ee:05")
	}, "force-reprovision", "aa:bb:cc:dd:ee:05")
}

func TestRunSpeedTest(t *testing.T) {
	t.Run("posts to speedtest endpoint", func(t *testing.T) {
		called := false
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/s/test-site-id/cmd/devmgr/speedtest" && r.Method == http.MethodPost {
				called = true
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(legacyOK([]any{}))
				return
			}
			http.Error(w, "unexpected", http.StatusBadRequest)
		})
		if err := client.RunSpeedTest(context.Background(), ""); err != nil {
			t.Fatalf("RunSpeedTest: %v", err)
		}
		if !called {
			t.Error("speedtest endpoint was not called")
		}
	})
}

func TestGetSpeedTestStatus(t *testing.T) {
	t.Run("decodes speed test result", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/s/test-site-id/stat/speedtest-status" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(legacyOK([]map[string]any{
				{"status": "complete", "xput_download": 950.5, "xput_upload": 480.2},
			}))
		})
		status, err := client.GetSpeedTestStatus(context.Background(), "")
		if err != nil {
			t.Fatalf("GetSpeedTestStatus: %v", err)
		}
		if status.Status != "complete" {
			t.Errorf("got status %q, want %q", status.Status, "complete")
		}
	})

	t.Run("returns error when response is empty", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(legacyOK([]any{}))
		})
		_, err := client.GetSpeedTestStatus(context.Background(), "")
		if err == nil {
			t.Error("expected error for empty response, got nil")
		}
	})
}
