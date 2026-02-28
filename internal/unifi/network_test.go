package unifi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListWiFiBroadcasts(t *testing.T) {
	t.Run("decodes broadcast list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/wifi/broadcasts" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "b-1", "name": "HomeNet", "enabled": true, "type": "2g"},
					{"id": "b-2", "name": "IoT", "enabled": false, "type": "5g"},
				},
				"totalCount": 2,
			})
		})
		broadcasts, err := client.ListWiFiBroadcasts(context.Background(), "")
		if err != nil {
			t.Fatalf("ListWiFiBroadcasts: %v", err)
		}
		if len(broadcasts) != 2 {
			t.Fatalf("got %d broadcasts, want 2", len(broadcasts))
		}
		if broadcasts[0].Name != "HomeNet" {
			t.Errorf("got Name %q, want HomeNet", broadcasts[0].Name)
		}
		if broadcasts[1].Enabled {
			t.Error("expected broadcasts[1].Enabled false")
		}
	})
}

func TestListNetworks(t *testing.T) {
	t.Run("decodes network list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/networks" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "n-1", "name": "LAN", "vlanId": 0},
					{"id": "n-2", "name": "IoT_VLAN", "vlanId": 10},
				},
				"totalCount": 2,
			})
		})
		nets, err := client.ListNetworks(context.Background(), "")
		if err != nil {
			t.Fatalf("ListNetworks: %v", err)
		}
		if len(nets) != 2 {
			t.Fatalf("got %d networks, want 2", len(nets))
		}
		if nets[1].VLANID != 10 {
			t.Errorf("got VLANID %d, want 10", nets[1].VLANID)
		}
	})
}

func TestListFirewallPolicies(t *testing.T) {
	t.Run("decodes policy list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/firewall/policies" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "p-1", "name": "block-iot-out", "enabled": true, "index": 0},
				},
				"totalCount": 1,
			})
		})
		policies, err := client.ListFirewallPolicies(context.Background(), "")
		if err != nil {
			t.Fatalf("ListFirewallPolicies: %v", err)
		}
		if len(policies) != 1 {
			t.Fatalf("got %d policies, want 1", len(policies))
		}
		if policies[0].Name != "block-iot-out" {
			t.Errorf("got Name %q, want block-iot-out", policies[0].Name)
		}
	})
}

func TestListFirewallZones(t *testing.T) {
	t.Run("decodes zone list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/firewall/zones" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "z-1", "name": "Internal", "networkIds": []string{"net-1", "net-2"}},
					{"id": "z-2", "name": "External"},
				},
				"totalCount": 2,
			})
		})
		zones, err := client.ListFirewallZones(context.Background(), "")
		if err != nil {
			t.Fatalf("ListFirewallZones: %v", err)
		}
		if len(zones) != 2 {
			t.Fatalf("got %d zones, want 2", len(zones))
		}
		if zones[0].Name != "Internal" {
			t.Errorf("got Name %q, want Internal", zones[0].Name)
		}
		if len(zones[0].NetworkIDs) != 2 {
			t.Errorf("got %d networkIDs, want 2", len(zones[0].NetworkIDs))
		}
	})
}

func TestListACLRules(t *testing.T) {
	t.Run("decodes ACL rule list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/acl-rules" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "r-1", "name": "allow-mgmt", "action": "ALLOW", "enabled": true, "type": "INTERNET", "index": 0},
				},
				"totalCount": 1,
			})
		})
		rules, err := client.ListACLRules(context.Background(), "")
		if err != nil {
			t.Fatalf("ListACLRules: %v", err)
		}
		if len(rules) != 1 {
			t.Fatalf("got %d rules, want 1", len(rules))
		}
		if rules[0].Action != "ALLOW" {
			t.Errorf("got Action %q, want ALLOW", rules[0].Action)
		}
	})
}
