package unifi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestGetSiteStats(t *testing.T) {
	t.Run("decodes site stats", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/sites/test-site-id/statistics/site" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"num_user":  42,
					"num_guest": 3,
				},
			})
		})
		stats, err := client.GetSiteStats(context.Background(), "")
		if err != nil {
			t.Fatalf("GetSiteStats: %v", err)
		}
		if stats.NumUser != 42 {
			t.Errorf("got NumUser %d, want 42", stats.NumUser)
		}
	})
}

func TestGetDeviceStats(t *testing.T) {
	t.Run("decodes device stats", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/sites/test-site-id/statistics/devices/dev-42" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"id": "dev-42", "mac": "ff:ee:dd:00:00:01"},
			})
		})
		dev, err := client.GetDeviceStats(context.Background(), "", "dev-42")
		if err != nil {
			t.Fatalf("GetDeviceStats: %v", err)
		}
		if dev.ID != "dev-42" {
			t.Errorf("got ID %q, want %q", dev.ID, "dev-42")
		}
	})
}

func TestGetClientStats(t *testing.T) {
	t.Run("decodes client stats", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/sites/test-site-id/statistics/clients/cl-99" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"mac": "11:22:33:44:55:66", "hostname": "tv", "ip": "10.0.0.50"},
			})
		})
		cl, err := client.GetClientStats(context.Background(), "", "cl-99")
		if err != nil {
			t.Fatalf("GetClientStats: %v", err)
		}
		if cl.Hostname != "tv" {
			t.Errorf("got hostname %q, want %q", cl.Hostname, "tv")
		}
	})
}

func TestListEvents(t *testing.T) {
	t.Run("decodes event list without limit", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/s/test-site-id/stat/event" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			if r.URL.RawQuery != "" {
				t.Errorf("expected no query string, got %q", r.URL.RawQuery)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(legacyOK([]map[string]any{
				{"key": "EVT_WU_Connected", "msg": "device connected"},
			}))
		})
		events, err := client.ListEvents(context.Background(), "", 0)
		if err != nil {
			t.Fatalf("ListEvents: %v", err)
		}
		if len(events) != 1 {
			t.Fatalf("got %d events, want 1", len(events))
		}
		if events[0].Key != "EVT_WU_Connected" {
			t.Errorf("got key %q, want %q", events[0].Key, "EVT_WU_Connected")
		}
	})

	t.Run("sends limit query parameter", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.RawQuery, "_limit=10") {
				t.Errorf("expected _limit=10 in query, got %q", r.URL.RawQuery)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(legacyOK([]any{}))
		})
		_, err := client.ListEvents(context.Background(), "", 10)
		if err != nil {
			t.Fatalf("ListEvents with limit: %v", err)
		}
	})
}

func TestListAlarms(t *testing.T) {
	t.Run("decodes alarm list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/s/test-site-id/stat/alarm" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(legacyOK([]map[string]any{
				{"key": "EVT_IPS_IpsAlert", "msg": "IPS alert"},
			}))
		})
		alarms, err := client.ListAlarms(context.Background(), "", false)
		if err != nil {
			t.Fatalf("ListAlarms: %v", err)
		}
		if len(alarms) != 1 {
			t.Fatalf("got %d alarms, want 1", len(alarms))
		}
	})

	t.Run("sends archived query parameter", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.RawQuery != "archived=true" {
				t.Errorf("expected archived=true query, got %q", r.URL.RawQuery)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(legacyOK([]any{}))
		})
		_, err := client.ListAlarms(context.Background(), "", true)
		if err != nil {
			t.Fatalf("ListAlarms archived: %v", err)
		}
	})
}
