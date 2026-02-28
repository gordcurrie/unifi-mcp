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

func TestGetWiFiBroadcast(t *testing.T) {
	t.Run("decodes single broadcast", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/wifi/broadcasts/bc-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "bc-1", "name": "HomeWiFi", "type": "STANDARD", "enabled": true,
			})
		})
		bc, err := client.GetWiFiBroadcast(context.Background(), "", "bc-1")
		if err != nil {
			t.Fatalf("GetWiFiBroadcast: %v", err)
		}
		if bc.ID != "bc-1" {
			t.Errorf("got ID %q, want bc-1", bc.ID)
		}
		if !bc.Enabled {
			t.Error("got Enabled false, want true")
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.GetWiFiBroadcast(context.Background(), "", "bc-1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestSetWiFiBroadcastEnabled(t *testing.T) {
	t.Run("disables a broadcast", func(t *testing.T) {
		var putEnabled any
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/wifi/broadcasts/bc-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id": "bc-1", "name": "HomeWiFi", "type": "STANDARD", "enabled": true,
				})
			case http.MethodPut:
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					http.Error(w, "bad body", http.StatusBadRequest)
					return
				}
				putEnabled = body["enabled"]
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id": "bc-1", "name": "HomeWiFi", "type": "STANDARD", "enabled": false,
				})
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})
		bc, err := client.SetWiFiBroadcastEnabled(context.Background(), "", "bc-1", false)
		if err != nil {
			t.Fatalf("SetWiFiBroadcastEnabled: %v", err)
		}
		if bc.Enabled {
			t.Error("got Enabled true, want false")
		}
		if putEnabled != false {
			t.Errorf("PUT body enabled = %v, want false", putEnabled)
		}
	})

	t.Run("returns error when GET fails", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.SetWiFiBroadcastEnabled(context.Background(), "", "bc-1", false)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestListTrafficMatchingLists(t *testing.T) {
	t.Run("decodes list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/traffic-matching-lists" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "tml-1", "name": "BlockedIPs", "type": "IP", "entries": []string{"10.0.0.1", "10.0.0.2"}},
					{"id": "tml-2", "name": "TrustedPorts", "type": "PORT"},
				},
				"totalCount": 2,
			})
		})
		lists, err := client.ListTrafficMatchingLists(context.Background(), "")
		if err != nil {
			t.Fatalf("ListTrafficMatchingLists: %v", err)
		}
		if len(lists) != 2 {
			t.Fatalf("got %d lists, want 2", len(lists))
		}
		if lists[0].Name != "BlockedIPs" {
			t.Errorf("got Name %q, want BlockedIPs", lists[0].Name)
		}
		if len(lists[0].Entries) != 2 {
			t.Errorf("got %d entries, want 2", len(lists[0].Entries))
		}
	})
}

func TestGetTrafficMatchingList(t *testing.T) {
	t.Run("decodes single list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/traffic-matching-lists/tml-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "tml-1", "name": "BlockedIPs", "type": "IP",
				"entries": []string{"10.0.0.1"},
			})
		})
		list, err := client.GetTrafficMatchingList(context.Background(), "", "tml-1")
		if err != nil {
			t.Fatalf("GetTrafficMatchingList: %v", err)
		}
		if list.ID != "tml-1" {
			t.Errorf("got ID %q, want tml-1", list.ID)
		}
		if list.Type != "IP" {
			t.Errorf("got Type %q, want IP", list.Type)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.GetTrafficMatchingList(context.Background(), "", "tml-1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestListWANs(t *testing.T) {
	t.Run("decodes WAN list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/wans" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "wan-1", "name": "WAN1", "type": "DHCP", "enabled": true, "ipAddress": "203.0.113.10", "state": "CONNECTED"},
				},
				"totalCount": 1,
			})
		})
		wans, err := client.ListWANs(context.Background(), "")
		if err != nil {
			t.Fatalf("ListWANs: %v", err)
		}
		if len(wans) != 1 {
			t.Fatalf("got %d WANs, want 1", len(wans))
		}
		if wans[0].Name != "WAN1" {
			t.Errorf("got Name %q, want WAN1", wans[0].Name)
		}
		if wans[0].State != "CONNECTED" {
			t.Errorf("got State %q, want CONNECTED", wans[0].State)
		}
	})
}

func TestListVPNTunnels(t *testing.T) {
	t.Run("decodes tunnel list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/vpn/site-to-site-tunnels" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "tun-1", "name": "OfficeVPN", "type": "IPSEC", "enabled": true, "state": "ACTIVE"},
				},
				"totalCount": 1,
			})
		})
		tunnels, err := client.ListVPNTunnels(context.Background(), "")
		if err != nil {
			t.Fatalf("ListVPNTunnels: %v", err)
		}
		if len(tunnels) != 1 {
			t.Fatalf("got %d tunnels, want 1", len(tunnels))
		}
		if tunnels[0].Name != "OfficeVPN" {
			t.Errorf("got Name %q, want OfficeVPN", tunnels[0].Name)
		}
	})
}

func TestListVPNServers(t *testing.T) {
	t.Run("decodes server list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/vpn/servers" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "srv-1", "name": "HomeVPN", "type": "WIREGUARD", "enabled": true},
				},
				"totalCount": 1,
			})
		})
		servers, err := client.ListVPNServers(context.Background(), "")
		if err != nil {
			t.Fatalf("ListVPNServers: %v", err)
		}
		if len(servers) != 1 {
			t.Fatalf("got %d servers, want 1", len(servers))
		}
		if servers[0].Type != "WIREGUARD" {
			t.Errorf("got Type %q, want WIREGUARD", servers[0].Type)
		}
	})
}
