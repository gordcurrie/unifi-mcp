package unifi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListWLANs(t *testing.T) {
	t.Run("decodes WLAN list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/s/test-site-id/rest/wlanconf" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(legacyOK([]map[string]any{
				{"_id": "wlan-1", "name": "HomeNet", "enabled": true},
				{"_id": "wlan-2", "name": "IoT", "enabled": false},
			}))
		})
		wlans, err := client.ListWLANs(context.Background(), "")
		if err != nil {
			t.Fatalf("ListWLANs: %v", err)
		}
		if len(wlans) != 2 {
			t.Fatalf("got %d WLANs, want 2", len(wlans))
		}
		if wlans[0].Name != "HomeNet" {
			t.Errorf("got Name %q, want %q", wlans[0].Name, "HomeNet")
		}
		if wlans[1].Enabled {
			t.Errorf("expected wlans[1].Enabled false, got true")
		}
	})
}

func TestSetWLANEnabled(t *testing.T) {
	tests := []struct {
		name        string
		wlanID      string
		enabled     bool
		wantEnabled bool
	}{
		{"enable WLAN", "wlan-1", true, true},
		{"disable WLAN", "wlan-2", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotEnabled *bool
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				wantPath := "/api/s/test-site-id/rest/wlanconf/" + tt.wlanID
				if r.URL.Path != wantPath || r.Method != http.MethodPut {
					http.Error(w, "unexpected request", http.StatusBadRequest)
					return
				}
				var req wlanEnabledRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "bad body", http.StatusBadRequest)
					return
				}
				gotEnabled = &req.Enabled
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(legacyOK([]any{}))
			})
			if err := client.SetWLANEnabled(context.Background(), "", tt.wlanID, tt.enabled); err != nil {
				t.Fatalf("SetWLANEnabled: %v", err)
			}
			if gotEnabled == nil {
				t.Fatal("request body was not captured")
			}
			if *gotEnabled != tt.wantEnabled {
				t.Errorf("got enabled=%v, want %v", *gotEnabled, tt.wantEnabled)
			}
		})
	}
}

func TestListNetworks(t *testing.T) {
	t.Run("decodes network list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/s/test-site-id/rest/networkconf" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(legacyOK([]map[string]any{
				{"_id": "net-1", "name": "LAN", "ip_subnet": "192.168.1.0/24"},
				{"_id": "net-2", "name": "IoT_VLAN", "ip_subnet": "10.10.10.0/24", "vlan_id": 10},
			}))
		})
		nets, err := client.ListNetworks(context.Background(), "")
		if err != nil {
			t.Fatalf("ListNetworks: %v", err)
		}
		if len(nets) != 2 {
			t.Fatalf("got %d networks, want 2", len(nets))
		}
		if nets[0].Name != "LAN" {
			t.Errorf("got Name %q, want %q", nets[0].Name, "LAN")
		}
		if nets[1].VLANID != 10 {
			t.Errorf("got VLANID %d, want 10", nets[1].VLANID)
		}
	})
}

func TestListFirewallRules(t *testing.T) {
	t.Run("decodes firewall rule list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/s/test-site-id/rest/firewallrule" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(legacyOK([]map[string]any{
				{"_id": "rule-1", "name": "block-iot-out", "action": "drop", "enabled": true},
			}))
		})
		rules, err := client.ListFirewallRules(context.Background(), "")
		if err != nil {
			t.Fatalf("ListFirewallRules: %v", err)
		}
		if len(rules) != 1 {
			t.Fatalf("got %d rules, want 1", len(rules))
		}
		if rules[0].Name != "block-iot-out" {
			t.Errorf("got Name %q, want %q", rules[0].Name, "block-iot-out")
		}
	})
}

func TestListPortForwards(t *testing.T) {
	t.Run("decodes port forward list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/s/test-site-id/rest/portforward" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(legacyOK([]map[string]any{
				{"_id": "pf-1", "name": "ssh", "dst_port": "22", "fwd": "192.168.1.50", "fwd_port": "22"},
			}))
		})
		rules, err := client.ListPortForwards(context.Background(), "")
		if err != nil {
			t.Fatalf("ListPortForwards: %v", err)
		}
		if len(rules) != 1 {
			t.Fatalf("got %d rules, want 1", len(rules))
		}
		if rules[0].Name != "ssh" {
			t.Errorf("got Name %q, want %q", rules[0].Name, "ssh")
		}
	})
}
