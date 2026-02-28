package unifi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListDevices(t *testing.T) {
	t.Run("decodes device list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/devices" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "dev-1", "macAddress": "aa:bb:cc:dd:ee:01", "name": "switch1", "state": "ONLINE"},
					{"id": "dev-2", "macAddress": "aa:bb:cc:dd:ee:02", "name": "ap1", "state": "OFFLINE"},
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
			t.Errorf("got devices[0].ID %q, want dev-1", devices[0].ID)
		}
		if devices[1].State != "OFFLINE" {
			t.Errorf("got devices[1].State %q, want OFFLINE", devices[1].State)
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
			if r.URL.Path != "/integration/v1/sites/test-site-id/devices/dev-99" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "dev-99", "macAddress": "aa:bb:cc:00:00:99", "state": "ONLINE",
			})
		})
		dev, err := client.GetDevice(context.Background(), "", "dev-99")
		if err != nil {
			t.Fatalf("GetDevice: %v", err)
		}
		if dev.ID != "dev-99" {
			t.Errorf("got ID %q, want dev-99", dev.ID)
		}
	})
}

func TestRestartDevice(t *testing.T) {
	t.Run("posts restart action", func(t *testing.T) {
		var gotAction string
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/devices/dev-1/actions" || r.Method != http.MethodPost {
				http.Error(w, "unexpected", http.StatusBadRequest)
				return
			}
			var req deviceActionRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "bad body", http.StatusBadRequest)
				return
			}
			gotAction = req.Action
			w.WriteHeader(http.StatusNoContent)
		})
		if err := client.RestartDevice(context.Background(), "", "dev-1"); err != nil {
			t.Fatalf("RestartDevice: %v", err)
		}
		if gotAction != "RESTART" {
			t.Errorf("got action %q, want RESTART", gotAction)
		}
	})
}

func TestGetDeviceStats(t *testing.T) {
	t.Run("decodes stats", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/devices/dev-1/statistics/latest" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uptimeSec": 12345, "cpuUtilizationPct": 5.2, "memoryUtilizationPct": 42.0,
			})
		})
		stats, err := client.GetDeviceStats(context.Background(), "", "dev-1")
		if err != nil {
			t.Fatalf("GetDeviceStats: %v", err)
		}
		if stats.UptimeSec != 12345 {
			t.Errorf("got UptimeSec %d, want 12345", stats.UptimeSec)
		}
	})
}
