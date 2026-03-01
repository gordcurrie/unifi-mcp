package unifi

import (
	"context"
	"encoding/json"
	"fmt"
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
					{"id": "r-1", "name": "allow-mgmt", "action": "ALLOW", "enabled": true, "type": "IPV4", "index": 0},
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

	t.Run("enables a broadcast", func(t *testing.T) {
		var putEnabled any
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/wifi/broadcasts/bc-2" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id": "bc-2", "name": "GuestWiFi", "type": "STANDARD", "enabled": false,
				})
			case http.MethodPut:
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					http.Error(w, "bad body", http.StatusBadRequest)
					return
				}
				putEnabled = body["enabled"]
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id": "bc-2", "name": "GuestWiFi", "type": "STANDARD", "enabled": true,
				})
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})
		bc, err := client.SetWiFiBroadcastEnabled(context.Background(), "", "bc-2", true)
		if err != nil {
			t.Fatalf("SetWiFiBroadcastEnabled: %v", err)
		}
		if !bc.Enabled {
			t.Error("got Enabled false, want true")
		}
		if putEnabled != true {
			t.Errorf("PUT body enabled = %v, want true", putEnabled)
		}
	})

	t.Run("returns error when PUT fails", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id": "bc-1", "name": "HomeWiFi", "type": "STANDARD", "enabled": true,
				})
			case http.MethodPut:
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
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

func TestListDNSPolicies(t *testing.T) {
	t.Run("decodes policy list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/dns/policies" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "p-1", "type": "A_RECORD", "domain": "nas.home", "ipv4Address": "192.168.1.100", "enabled": true},
					{"id": "p-2", "type": "A_RECORD", "domain": "pi.home", "ipv4Address": "192.168.1.200", "enabled": false},
				},
				"totalCount": 2,
			})
		})
		policies, err := client.ListDNSPolicies(context.Background(), "")
		if err != nil {
			t.Fatalf("ListDNSPolicies: %v", err)
		}
		if len(policies) != 2 {
			t.Fatalf("got %d policies, want 2", len(policies))
		}
		if policies[0].Domain != "nas.home" {
			t.Errorf("got Domain %q, want nas.home", policies[0].Domain)
		}
		if policies[1].Enabled {
			t.Error("expected policies[1].Enabled false")
		}
	})
}

func TestGetDNSPolicy(t *testing.T) {
	t.Run("decodes single policy", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/dns/policies/p-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "p-1", "type": "A_RECORD", "domain": "nas.home",
				"ipv4Address": "192.168.1.100", "ttlSeconds": 300, "enabled": true,
			})
		})
		policy, err := client.GetDNSPolicy(context.Background(), "", "p-1")
		if err != nil {
			t.Fatalf("GetDNSPolicy: %v", err)
		}
		if policy.IPv4Address != "192.168.1.100" {
			t.Errorf("got IPv4Address %q, want 192.168.1.100", policy.IPv4Address)
		}
		if policy.TTLSeconds != 300 {
			t.Errorf("got TTLSeconds %d, want 300", policy.TTLSeconds)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "not found", http.StatusNotFound)
		})
		_, err := client.GetDNSPolicy(context.Background(), "", "missing")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestCreateDNSPolicy(t *testing.T) {
	t.Run("posts and decodes created policy", func(t *testing.T) {
		var gotBody DNSPolicyRequest
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost || r.URL.Path != "/integration/v1/sites/test-site-id/dns/policies" {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				http.Error(w, "decode error", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "p-new", "type": gotBody.Type, "domain": gotBody.Domain,
				"ipv4Address": gotBody.IPv4Address, "enabled": gotBody.Enabled,
			})
		})
		req := DNSPolicyRequest{Type: "A_RECORD", Domain: "test.home", IPv4Address: "10.0.0.5", Enabled: true}
		policy, err := client.CreateDNSPolicy(context.Background(), "", req)
		if err != nil {
			t.Fatalf("CreateDNSPolicy: %v", err)
		}
		if policy.ID != "p-new" {
			t.Errorf("got ID %q, want p-new", policy.ID)
		}
		if gotBody.Domain != "test.home" {
			t.Errorf("POST body domain = %q, want test.home", gotBody.Domain)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.CreateDNSPolicy(context.Background(), "", DNSPolicyRequest{Type: "A_RECORD", Domain: "x.home"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUpdateDNSPolicy(t *testing.T) {
	t.Run("puts and decodes updated policy", func(t *testing.T) {
		var gotBody DNSPolicyRequest
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut || r.URL.Path != "/integration/v1/sites/test-site-id/dns/policies/p-1" {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				http.Error(w, "decode error", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "p-1", "type": gotBody.Type, "domain": gotBody.Domain,
				"ipv4Address": gotBody.IPv4Address, "enabled": gotBody.Enabled,
			})
		})
		req := DNSPolicyRequest{Type: "A_RECORD", Domain: "nas.home", IPv4Address: "192.168.1.99", Enabled: true}
		policy, err := client.UpdateDNSPolicy(context.Background(), "", "p-1", req)
		if err != nil {
			t.Fatalf("UpdateDNSPolicy: %v", err)
		}
		if gotBody.IPv4Address != "192.168.1.99" {
			t.Errorf("PUT body IPv4Address = %q, want 192.168.1.99", gotBody.IPv4Address)
		}
		if policy.ID != "p-1" {
			t.Errorf("got ID %q, want p-1", policy.ID)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.UpdateDNSPolicy(context.Background(), "", "p-1", DNSPolicyRequest{Type: "A_RECORD", Domain: "nas.home"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestDeleteDNSPolicy(t *testing.T) {
	t.Run("sends DELETE and succeeds on 204", func(t *testing.T) {
		var gotMethod string
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/dns/policies/p-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			gotMethod = r.Method
			w.WriteHeader(http.StatusNoContent)
		})
		if err := client.DeleteDNSPolicy(context.Background(), "", "p-1"); err != nil {
			t.Fatalf("DeleteDNSPolicy: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("got method %q, want DELETE", gotMethod)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "forbidden", http.StatusForbidden)
		})
		if err := client.DeleteDNSPolicy(context.Background(), "", "p-1"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetFirewallPolicy(t *testing.T) {
	t.Run("decodes single policy", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/firewall/policies/fp-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "fp-1", "name": "Allow LAN", "enabled": true, "index": 100,
				"action":          map[string]any{"type": "ALLOW", "allowReturnTraffic": true},
				"source":          map[string]any{"zoneId": "zone-a"},
				"destination":     map[string]any{"zoneId": "zone-b"},
				"ipProtocolScope": map[string]any{"ipVersion": "IPV4_AND_IPV6"},
				"loggingEnabled":  false,
				"metadata":        map[string]any{"origin": "USER_DEFINED", "configurable": true},
			})
		})
		policy, err := client.GetFirewallPolicy(context.Background(), "", "fp-1")
		if err != nil {
			t.Fatalf("GetFirewallPolicy: %v", err)
		}
		if policy.Name != "Allow LAN" {
			t.Errorf("got Name %q, want Allow LAN", policy.Name)
		}
		if policy.Action.Type != "ALLOW" {
			t.Errorf("got Action.Type %q, want ALLOW", policy.Action.Type)
		}
		if policy.Source.ZoneID != "zone-a" {
			t.Errorf("got Source.ZoneID %q, want zone-a", policy.Source.ZoneID)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "not found", http.StatusNotFound)
		})
		_, err := client.GetFirewallPolicy(context.Background(), "", "missing")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestSetFirewallPolicyEnabled(t *testing.T) {
	t.Run("gets then puts with enabled flag set", func(t *testing.T) {
		var putBody map[string]any
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			path := "/integration/v1/sites/test-site-id/firewall/policies/fp-1"
			if r.URL.Path != path {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodGet {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id": "fp-1", "enabled": false, "name": "Block All",
					"action":          map[string]any{"type": "BLOCK"},
					"source":          map[string]any{"zoneId": "zone-a"},
					"destination":     map[string]any{"zoneId": "zone-b"},
					"ipProtocolScope": map[string]any{"ipVersion": "IPV4_AND_IPV6"},
					"loggingEnabled":  false,
					"metadata":        map[string]any{"origin": "USER_DEFINED", "configurable": true},
				})
				return
			}
			if r.Method == http.MethodPut {
				if err := json.NewDecoder(r.Body).Decode(&putBody); err != nil {
					http.Error(w, "decode error", http.StatusBadRequest)
					return
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id": "fp-1", "enabled": true, "name": "Block All",
					"action":          putBody["action"],
					"source":          putBody["source"],
					"destination":     putBody["destination"],
					"ipProtocolScope": putBody["ipProtocolScope"],
					"loggingEnabled":  false,
				})
				return
			}
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		})
		policy, err := client.SetFirewallPolicyEnabled(context.Background(), "", "fp-1", true)
		if err != nil {
			t.Fatalf("SetFirewallPolicyEnabled: %v", err)
		}
		if !policy.Enabled {
			t.Error("expected policy.Enabled true")
		}
		if en, ok := putBody["enabled"].(bool); !ok || !en {
			t.Errorf("PUT body enabled = %v, want true", putBody["enabled"])
		}
		if _, hasID := putBody["id"]; hasID {
			t.Error("PUT body should not contain id")
		}
		if _, hasMeta := putBody["metadata"]; hasMeta {
			t.Error("PUT body should not contain metadata")
		}
	})

	t.Run("returns error on GET failure", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.SetFirewallPolicyEnabled(context.Background(), "", "fp-1", true)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns error on PUT failure", func(t *testing.T) {
		callCount := 0
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				if r.Method != http.MethodGet {
					t.Fatalf("first request method = %s, want GET", r.Method)
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id": "fp-1", "enabled": false, "name": "Block All",
					"action":          map[string]any{"type": "BLOCK"},
					"source":          map[string]any{"zoneId": "zone-a"},
					"destination":     map[string]any{"zoneId": "zone-b"},
					"ipProtocolScope": map[string]any{"ipVersion": "IPV4_AND_IPV6"},
					"loggingEnabled":  false,
					"metadata":        map[string]any{"origin": "USER_DEFINED", "configurable": true},
				})
				return
			}
			if r.Method != http.MethodPut {
				t.Fatalf("second request method = %s, want PUT", r.Method)
			}
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.SetFirewallPolicyEnabled(context.Background(), "", "fp-1", true)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestDeleteFirewallPolicy(t *testing.T) {
	t.Run("sends DELETE and succeeds on 204", func(t *testing.T) {
		var gotMethod string
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/firewall/policies/fp-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			gotMethod = r.Method
			w.WriteHeader(http.StatusNoContent)
		})
		if err := client.DeleteFirewallPolicy(context.Background(), "", "fp-1"); err != nil {
			t.Fatalf("DeleteFirewallPolicy: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("got method %q, want DELETE", gotMethod)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "forbidden", http.StatusForbidden)
		})
		if err := client.DeleteFirewallPolicy(context.Background(), "", "fp-1"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetFirewallZone(t *testing.T) {
	t.Run("decodes single zone", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/firewall/zones/z-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "z-1", "name": "Internal",
				"networkIds": []string{"net-a", "net-b"},
				"metadata":   map[string]any{"origin": "SYSTEM_DEFINED", "configurable": true},
			})
		})
		zone, err := client.GetFirewallZone(context.Background(), "", "z-1")
		if err != nil {
			t.Fatalf("GetFirewallZone: %v", err)
		}
		if zone.Name != "Internal" {
			t.Errorf("got Name %q, want Internal", zone.Name)
		}
		if len(zone.NetworkIDs) != 2 {
			t.Errorf("got %d network IDs, want 2", len(zone.NetworkIDs))
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "not found", http.StatusNotFound)
		})
		_, err := client.GetFirewallZone(context.Background(), "", "missing")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestCreateFirewallZone(t *testing.T) {
	t.Run("posts and decodes created zone", func(t *testing.T) {
		var gotBody FirewallZoneRequest
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost || r.URL.Path != "/integration/v1/sites/test-site-id/firewall/zones" {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				http.Error(w, "decode error", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "z-new", "name": gotBody.Name, "networkIds": gotBody.NetworkIDs,
				"metadata": map[string]any{"origin": "USER_DEFINED", "configurable": true},
			})
		})
		req := FirewallZoneRequest{Name: "MyZone", NetworkIDs: []string{"net-x"}}
		zone, err := client.CreateFirewallZone(context.Background(), "", req)
		if err != nil {
			t.Fatalf("CreateFirewallZone: %v", err)
		}
		if zone.ID != "z-new" {
			t.Errorf("got ID %q, want z-new", zone.ID)
		}
		if gotBody.Name != "MyZone" {
			t.Errorf("POST body Name = %q, want MyZone", gotBody.Name)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.CreateFirewallZone(context.Background(), "", FirewallZoneRequest{Name: "X", NetworkIDs: []string{}})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUpdateFirewallZone(t *testing.T) {
	t.Run("puts and decodes updated zone", func(t *testing.T) {
		var gotBody FirewallZoneRequest
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut || r.URL.Path != "/integration/v1/sites/test-site-id/firewall/zones/z-1" {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				http.Error(w, "decode error", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "z-1", "name": gotBody.Name, "networkIds": gotBody.NetworkIDs,
				"metadata": map[string]any{"origin": "USER_DEFINED", "configurable": true},
			})
		})
		req := FirewallZoneRequest{Name: "UpdatedZone", NetworkIDs: []string{"net-a", "net-b"}}
		zone, err := client.UpdateFirewallZone(context.Background(), "", "z-1", req)
		if err != nil {
			t.Fatalf("UpdateFirewallZone: %v", err)
		}
		if zone.Name != "UpdatedZone" {
			t.Errorf("got Name %q, want UpdatedZone", zone.Name)
		}
		if gotBody.Name != "UpdatedZone" {
			t.Errorf("PUT body Name = %q, want UpdatedZone", gotBody.Name)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.UpdateFirewallZone(context.Background(), "", "z-1", FirewallZoneRequest{Name: "X", NetworkIDs: []string{}})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestDeleteFirewallZone(t *testing.T) {
	t.Run("sends DELETE and succeeds on 204", func(t *testing.T) {
		var gotMethod string
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/firewall/zones/z-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			gotMethod = r.Method
			w.WriteHeader(http.StatusNoContent)
		})
		if err := client.DeleteFirewallZone(context.Background(), "", "z-1"); err != nil {
			t.Fatalf("DeleteFirewallZone: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("got method %q, want DELETE", gotMethod)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "forbidden", http.StatusForbidden)
		})
		if err := client.DeleteFirewallZone(context.Background(), "", "z-1"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetACLRule(t *testing.T) {
	t.Run("decodes single rule", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/acl-rules/ar-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "ar-1", "type": "IPV4", "name": "test-rule",
				"enabled": true, "action": "ALLOW", "index": 0,
				"metadata": map[string]any{"origin": "USER_DEFINED"},
			})
		})
		rule, err := client.GetACLRule(context.Background(), "", "ar-1")
		if err != nil {
			t.Fatalf("GetACLRule: %v", err)
		}
		if rule.Name != "test-rule" {
			t.Errorf("got Name %q, want test-rule", rule.Name)
		}
		if rule.Action != "ALLOW" {
			t.Errorf("got Action %q, want ALLOW", rule.Action)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "not found", http.StatusNotFound)
		})
		_, err := client.GetACLRule(context.Background(), "", "missing")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestCreateACLRule(t *testing.T) {
	t.Run("posts body and decodes response", func(t *testing.T) {
		var gotBody ACLRuleRequest
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/acl-rules" || r.Method != http.MethodPost {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				http.Error(w, "decode error", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "ar-new", "type": gotBody.Type, "name": gotBody.Name,
				"action": gotBody.Action, "enabled": gotBody.Enabled, "index": 0,
				"metadata": map[string]any{"origin": "USER_DEFINED"},
			})
		})
		req := ACLRuleRequest{Type: "IPV4", Name: "my-rule", Action: "BLOCK", Enabled: true}
		rule, err := client.CreateACLRule(context.Background(), "", req)
		if err != nil {
			t.Fatalf("CreateACLRule: %v", err)
		}
		if rule.ID != "ar-new" {
			t.Errorf("got ID %q, want ar-new", rule.ID)
		}
		if gotBody.Action != "BLOCK" {
			t.Errorf("sent Action %q, want BLOCK", gotBody.Action)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.CreateACLRule(context.Background(), "", ACLRuleRequest{Type: "IPV4", Name: "x", Action: "ALLOW"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUpdateACLRule(t *testing.T) {
	t.Run("puts body and decodes response", func(t *testing.T) {
		var gotBody ACLRuleRequest
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/acl-rules/ar-1" || r.Method != http.MethodPut {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				http.Error(w, "decode error", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "ar-1", "type": gotBody.Type, "name": gotBody.Name,
				"action": gotBody.Action, "enabled": gotBody.Enabled, "index": 0,
			})
		})
		req := ACLRuleRequest{Type: "IPV4", Name: "renamed", Action: "ALLOW", Enabled: false}
		rule, err := client.UpdateACLRule(context.Background(), "", "ar-1", req)
		if err != nil {
			t.Fatalf("UpdateACLRule: %v", err)
		}
		if rule.Name != "renamed" {
			t.Errorf("got Name %q, want renamed", rule.Name)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.UpdateACLRule(context.Background(), "", "ar-1", ACLRuleRequest{Type: "IPV4", Name: "x", Action: "ALLOW"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestDeleteACLRule(t *testing.T) {
	t.Run("sends DELETE and succeeds on 200", func(t *testing.T) {
		var gotMethod string
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/acl-rules/ar-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			gotMethod = r.Method
			w.WriteHeader(http.StatusOK)
		})
		if err := client.DeleteACLRule(context.Background(), "", "ar-1"); err != nil {
			t.Fatalf("DeleteACLRule: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("got method %q, want DELETE", gotMethod)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		if err := client.DeleteACLRule(context.Background(), "", "ar-1"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestSetACLRuleEnabled(t *testing.T) {
	t.Run("gets then puts with enabled flag set", func(t *testing.T) {
		callCount := 0
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				if r.Method != http.MethodGet {
					http.Error(w, fmt.Sprintf("first call method = %s, want GET", r.Method), http.StatusBadRequest)
					return
				}
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id": "ar-1", "type": "IPV4", "name": "my-rule",
					"action": "ALLOW", "enabled": false, "index": 0,
				})
				return
			}
			if r.Method != http.MethodPut {
				http.Error(w, fmt.Sprintf("second call method = %s, want PUT", r.Method), http.StatusBadRequest)
				return
			}
			var body ACLRuleRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "decode error", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "ar-1", "type": body.Type, "name": body.Name,
				"action": body.Action, "enabled": body.Enabled, "index": 0,
			})
		})
		rule, err := client.SetACLRuleEnabled(context.Background(), "", "ar-1", true)
		if err != nil {
			t.Fatalf("SetACLRuleEnabled: %v", err)
		}
		if !rule.Enabled {
			t.Error("expected rule.Enabled true")
		}
	})

	t.Run("returns error on GET failure", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.SetACLRuleEnabled(context.Background(), "", "ar-1", true)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetACLRuleOrdering(t *testing.T) {
	t.Run("decodes ordering response", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/acl-rules/ordering" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"orderedAclRuleIds": []string{"ar-1", "ar-2"},
			})
		})
		ordering, err := client.GetACLRuleOrdering(context.Background(), "")
		if err != nil {
			t.Fatalf("GetACLRuleOrdering: %v", err)
		}
		if len(ordering.OrderedACLRuleIDs) != 2 || ordering.OrderedACLRuleIDs[0] != "ar-1" {
			t.Errorf("got IDs %v, want [ar-1 ar-2]", ordering.OrderedACLRuleIDs)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.GetACLRuleOrdering(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestReorderACLRules(t *testing.T) {
	t.Run("puts ordered IDs and decodes response", func(t *testing.T) {
		var gotBody ACLRuleOrdering
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/acl-rules/ordering" || r.Method != http.MethodPut {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				http.Error(w, "decode error", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"orderedAclRuleIds": gotBody.OrderedACLRuleIDs,
			})
		})
		ordering, err := client.ReorderACLRules(context.Background(), "", []string{"ar-2", "ar-1"})
		if err != nil {
			t.Fatalf("ReorderACLRules: %v", err)
		}
		if len(ordering.OrderedACLRuleIDs) != 2 || ordering.OrderedACLRuleIDs[0] != "ar-2" {
			t.Errorf("got IDs %v, want [ar-2 ar-1]", ordering.OrderedACLRuleIDs)
		}
		if gotBody.OrderedACLRuleIDs[0] != "ar-2" {
			t.Errorf("sent IDs[0] %q, want ar-2", gotBody.OrderedACLRuleIDs[0])
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.ReorderACLRules(context.Background(), "", []string{"ar-1"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestListVouchers(t *testing.T) {
	t.Run("decodes voucher list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/hotspot/vouchers" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "v-1", "code": "ABC123", "name": "guest", "timeLimitMinutes": 60, "status": "VALID"},
				},
				"totalCount": 1,
			})
		})
		vouchers, err := client.ListVouchers(context.Background(), "")
		if err != nil {
			t.Fatalf("ListVouchers: %v", err)
		}
		if len(vouchers) != 1 || vouchers[0].ID != "v-1" || vouchers[0].Code != "ABC123" {
			t.Errorf("got %+v, want [{ID:v-1 Code:ABC123 ...}]", vouchers)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.ListVouchers(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGetVoucher(t *testing.T) {
	t.Run("decodes single voucher", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/hotspot/vouchers/v-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "v-1", "code": "ABC123", "name": "guest", "timeLimitMinutes": 60, "status": "VALID",
			})
		})
		voucher, err := client.GetVoucher(context.Background(), "", "v-1")
		if err != nil {
			t.Fatalf("GetVoucher: %v", err)
		}
		if voucher.ID != "v-1" || voucher.Code != "ABC123" {
			t.Errorf("got %+v, want {ID:v-1 Code:ABC123}", voucher)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "not found", http.StatusNotFound)
		})
		_, err := client.GetVoucher(context.Background(), "", "v-1")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestCreateVouchers(t *testing.T) {
	t.Run("sends count and returns vouchers", func(t *testing.T) {
		var gotBody VoucherRequest
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/hotspot/vouchers" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				http.Error(w, "decode error", http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"vouchers": []map[string]any{
					{"id": "v-2", "code": "XYZ789", "name": gotBody.Name, "timeLimitMinutes": gotBody.TimeLimitMinutes, "expired": false},
				},
			})
		})
		req := VoucherRequest{Count: 1, Name: "test", TimeLimitMinutes: 120}
		vouchers, err := client.CreateVouchers(context.Background(), "", req)
		if err != nil {
			t.Fatalf("CreateVouchers: %v", err)
		}
		if len(vouchers) != 1 || vouchers[0].ID != "v-2" {
			t.Errorf("got %+v, want [{ID:v-2 ...}]", vouchers)
		}
		if gotBody.Count != 1 || gotBody.Name != "test" || gotBody.TimeLimitMinutes != 120 {
			t.Errorf("sent body %+v, want {Count:1 Name:test TimeLimitMinutes:120}", gotBody)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		})
		_, err := client.CreateVouchers(context.Background(), "", VoucherRequest{Count: 1})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestDeleteVoucher(t *testing.T) {
	t.Run("sends DELETE and succeeds on 204", func(t *testing.T) {
		var gotMethod string
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/hotspot/vouchers/v-1" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			gotMethod = r.Method
			w.WriteHeader(http.StatusNoContent)
		})
		if err := client.DeleteVoucher(context.Background(), "", "v-1"); err != nil {
			t.Fatalf("DeleteVoucher: %v", err)
		}
		if gotMethod != http.MethodDelete {
			t.Errorf("got method %q, want DELETE", gotMethod)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "forbidden", http.StatusForbidden)
		})
		if err := client.DeleteVoucher(context.Background(), "", "v-1"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestListDeviceTags(t *testing.T) {
	t.Run("decodes tag list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/device-tags" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "dt-1", "name": "cameras"},
					{"id": "dt-2", "name": "iot"},
				},
				"totalCount": 2,
			})
		})
		tags, err := client.ListDeviceTags(context.Background(), "")
		if err != nil {
			t.Fatalf("ListDeviceTags: %v", err)
		}
		if len(tags) != 2 || tags[0].ID != "dt-1" {
			t.Errorf("got %+v, want [{ID:dt-1 ...}]", tags)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.ListDeviceTags(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestListDPICategories(t *testing.T) {
	t.Run("decodes category list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/dpi/categories" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": 0, "name": "Instant messengers"},
					{"id": 4, "name": "Media streaming services"},
				},
				"totalCount": 2,
			})
		})
		cats, err := client.ListDPICategories(context.Background())
		if err != nil {
			t.Fatalf("ListDPICategories: %v", err)
		}
		if len(cats) != 2 || cats[0].ID != 0 || cats[1].Name != "Media streaming services" {
			t.Errorf("unexpected categories: %+v", cats)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.ListDPICategories(context.Background())
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestListDPIApplications(t *testing.T) {
	t.Run("decodes application list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/dpi/applications" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": 39, "name": "Slack", "categoryId": 17},
					{"id": 40, "name": "Telegram", "categoryId": 0},
				},
				"totalCount": 2,
			})
		})
		apps, err := client.ListDPIApplications(context.Background())
		if err != nil {
			t.Fatalf("ListDPIApplications: %v", err)
		}
		if len(apps) != 2 || apps[0].ID != 39 || apps[0].CategoryID != 17 {
			t.Errorf("unexpected applications: %+v", apps)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.ListDPIApplications(context.Background())
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestListRADIUSProfiles(t *testing.T) {
	t.Run("decodes profile list", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/integration/v1/sites/test-site-id/radius/profiles" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": []map[string]any{
					{"id": "rp-1", "name": "Default", "metadata": map[string]any{"origin": "SYSTEM_DEFINED", "configurable": true}},
				},
				"totalCount": 1,
			})
		})
		profiles, err := client.ListRADIUSProfiles(context.Background(), "")
		if err != nil {
			t.Fatalf("ListRADIUSProfiles: %v", err)
		}
		if len(profiles) != 1 || profiles[0].ID != "rp-1" || profiles[0].Metadata == nil {
			t.Errorf("unexpected profiles: %+v", profiles)
		}
		if profiles[0].Metadata.Origin != "SYSTEM_DEFINED" {
			t.Errorf("got origin %q, want SYSTEM_DEFINED", profiles[0].Metadata.Origin)
		}
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "error", http.StatusInternalServerError)
		})
		_, err := client.ListRADIUSProfiles(context.Background(), "")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
