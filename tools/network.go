package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/gordcurrie/unifi-mcp/internal/unifi"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerNetworkTools(s *mcp.Server, client unifiClient, allowDestructive bool) {
	type siteInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
	}
	type broadcastInput struct {
		SiteID      string `json:"site_id,omitempty"   jsonschema:"site ID; omit to use default"`
		BroadcastID string `json:"broadcast_id"         jsonschema:"WiFi broadcast ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_wifi_broadcasts",
		Description: "List all WiFi broadcast configurations (SSIDs).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		broadcasts, err := client.ListWiFiBroadcasts(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_wifi_broadcasts: %w", err))
		}
		return jsonResult(broadcasts)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_wifi_broadcast",
		Description: "Get details for a specific WiFi broadcast (SSID) by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input broadcastInput) (*mcp.CallToolResult, any, error) {
		if input.BroadcastID == "" {
			return errorResult(fmt.Errorf("get_wifi_broadcast: broadcast_id is required"))
		}
		bc, err := client.GetWiFiBroadcast(ctx, input.SiteID, input.BroadcastID)
		if err != nil {
			return errorResult(fmt.Errorf("get_wifi_broadcast: %w", err))
		}
		return jsonResult(bc)
	})

	destructiveTrue := true

	type setBroadcastInput struct {
		SiteID      string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		BroadcastID string `json:"broadcast_id"       jsonschema:"WiFi broadcast ID"`
		Enabled     *bool  `json:"enabled"            jsonschema:"true to enable the broadcast, false to disable"`
		Confirmed   bool   `json:"confirmed"          jsonschema:"must be true to confirm the change"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "set_wifi_broadcast_enabled",
		Description: "Enable or disable a WiFi broadcast (SSID). Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input setBroadcastInput) (*mcp.CallToolResult, any, error) {
		if !input.Confirmed {
			return errorResult(fmt.Errorf("set_wifi_broadcast_enabled: set confirmed=true to confirm the change"))
		}
		if input.BroadcastID == "" {
			return errorResult(fmt.Errorf("set_wifi_broadcast_enabled: broadcast_id is required"))
		}
		if input.Enabled == nil {
			return errorResult(fmt.Errorf("set_wifi_broadcast_enabled: enabled is required"))
		}
		bc, err := client.SetWiFiBroadcastEnabled(ctx, input.SiteID, input.BroadcastID, *input.Enabled)
		if err != nil {
			return errorResult(fmt.Errorf("set_wifi_broadcast_enabled: %w", err))
		}
		return jsonResult(bc)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_networks",
		Description: "List all configured networks (VLANs, LAN segments).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		nets, err := client.ListNetworks(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_networks: %w", err))
		}
		return jsonResult(nets)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_firewall_policies",
		Description: "List all firewall policies for a site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		policies, err := client.ListFirewallPolicies(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_firewall_policies: %w", err))
		}
		return jsonResult(policies)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_firewall_zones",
		Description: "List all firewall zones for a site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		zones, err := client.ListFirewallZones(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_firewall_zones: %w", err))
		}
		return jsonResult(zones)
	})

	type trafficMatchingListInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		ListID string `json:"list_id"           jsonschema:"traffic matching list ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_traffic_matching_lists",
		Description: "List all traffic matching lists (IP/port sets used by firewall policies).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		lists, err := client.ListTrafficMatchingLists(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_traffic_matching_lists: %w", err))
		}
		return jsonResult(lists)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_traffic_matching_list",
		Description: "Get details for a specific traffic matching list by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input trafficMatchingListInput) (*mcp.CallToolResult, any, error) {
		if input.ListID == "" {
			return errorResult(fmt.Errorf("get_traffic_matching_list: list_id is required"))
		}
		list, err := client.GetTrafficMatchingList(ctx, input.SiteID, input.ListID)
		if err != nil {
			return errorResult(fmt.Errorf("get_traffic_matching_list: %w", err))
		}
		return jsonResult(list)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_wans",
		Description: "List all WAN interface definitions for a site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		wans, err := client.ListWANs(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_wans: %w", err))
		}
		return jsonResult(wans)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_vpn_tunnels",
		Description: "List all site-to-site VPN tunnels for a site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		tunnels, err := client.ListVPNTunnels(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_vpn_tunnels: %w", err))
		}
		return jsonResult(tunnels)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_vpn_servers",
		Description: "List all VPN server configurations for a site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		servers, err := client.ListVPNServers(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_vpn_servers: %w", err))
		}
		return jsonResult(servers)
	})

	type dnsPolicyInput struct {
		SiteID   string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		PolicyID string `json:"policy_id"         jsonschema:"DNS policy ID"`
	}
	type createDNSPolicyInput struct {
		SiteID      string `json:"site_id,omitempty"      jsonschema:"site ID; omit to use default"`
		Type        string `json:"type"                   jsonschema:"policy type, e.g. A_RECORD"`
		Domain      string `json:"domain"                 jsonschema:"domain name to resolve"`
		IPv4Address string `json:"ipv4_address,omitempty" jsonschema:"IPv4 address the domain maps to"`
		TTLSeconds  int    `json:"ttl_seconds"            jsonschema:"TTL in seconds; required by the API (send 0 to use the server default)"`
		Enabled     *bool  `json:"enabled"                jsonschema:"true to activate the policy, false to create disabled"`
	}
	type updateDNSPolicyInput struct {
		SiteID      string `json:"site_id,omitempty"      jsonschema:"site ID; omit to use default"`
		PolicyID    string `json:"policy_id"              jsonschema:"DNS policy ID to update"`
		Type        string `json:"type"                   jsonschema:"policy type, e.g. A_RECORD"`
		Domain      string `json:"domain"                 jsonschema:"domain name to resolve"`
		IPv4Address string `json:"ipv4_address,omitempty" jsonschema:"IPv4 address the domain maps to"`
		TTLSeconds  int    `json:"ttl_seconds"            jsonschema:"TTL in seconds; required by the API (send 0 to use the server default)"`
		Enabled     *bool  `json:"enabled"                jsonschema:"true to activate the policy, false to disable"`
	}
	type deleteDNSPolicyInput struct {
		SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		PolicyID  string `json:"policy_id"         jsonschema:"DNS policy ID to delete"`
		Confirmed bool   `json:"confirmed"         jsonschema:"must be true to confirm the deletion"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_dns_policies",
		Description: "List all local DNS policies (A-record overrides) for a site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		policies, err := client.ListDNSPolicies(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_dns_policies: %w", err))
		}
		return jsonResult(policies)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_dns_policy",
		Description: "Get details for a specific DNS policy by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input dnsPolicyInput) (*mcp.CallToolResult, any, error) {
		if input.PolicyID == "" {
			return errorResult(fmt.Errorf("get_dns_policy: policy_id is required"))
		}
		policy, err := client.GetDNSPolicy(ctx, input.SiteID, input.PolicyID)
		if err != nil {
			return errorResult(fmt.Errorf("get_dns_policy: %w", err))
		}
		return jsonResult(policy)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_dns_policy",
		Description: "Create a new local DNS A-record policy mapping a domain to an IP address.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input createDNSPolicyInput) (*mcp.CallToolResult, any, error) {
		if input.Type == "" {
			return errorResult(fmt.Errorf("create_dns_policy: type is required"))
		}
		if input.Domain == "" {
			return errorResult(fmt.Errorf("create_dns_policy: domain is required"))
		}
		if input.Enabled == nil {
			return errorResult(fmt.Errorf("create_dns_policy: enabled is required"))
		}
		req := unifi.DNSPolicyRequest{
			Type:        input.Type,
			Domain:      input.Domain,
			IPv4Address: input.IPv4Address,
			TTLSeconds:  input.TTLSeconds,
			Enabled:     *input.Enabled,
		}
		policy, err := client.CreateDNSPolicy(ctx, input.SiteID, req)
		if err != nil {
			return errorResult(fmt.Errorf("create_dns_policy: %w", err))
		}
		return jsonResult(policy)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_dns_policy",
		Description: "Update an existing local DNS policy by ID.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input updateDNSPolicyInput) (*mcp.CallToolResult, any, error) {
		if input.PolicyID == "" {
			return errorResult(fmt.Errorf("update_dns_policy: policy_id is required"))
		}
		if input.Type == "" {
			return errorResult(fmt.Errorf("update_dns_policy: type is required"))
		}
		if input.Domain == "" {
			return errorResult(fmt.Errorf("update_dns_policy: domain is required"))
		}
		if input.Enabled == nil {
			return errorResult(fmt.Errorf("update_dns_policy: enabled is required"))
		}
		req := unifi.DNSPolicyRequest{
			Type:        input.Type,
			Domain:      input.Domain,
			IPv4Address: input.IPv4Address,
			TTLSeconds:  input.TTLSeconds,
			Enabled:     *input.Enabled,
		}
		policy, err := client.UpdateDNSPolicy(ctx, input.SiteID, input.PolicyID, req)
		if err != nil {
			return errorResult(fmt.Errorf("update_dns_policy: %w", err))
		}
		return jsonResult(policy)
	})

	if allowDestructive {
		mcp.AddTool(s, &mcp.Tool{
			Name:        "delete_dns_policy",
			Description: "Permanently delete a DNS policy by ID. Requires UNIFI_ALLOW_DESTRUCTIVE=true. Set confirmed=true to proceed.",
			Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
		}, func(ctx context.Context, _ *mcp.CallToolRequest, input deleteDNSPolicyInput) (*mcp.CallToolResult, any, error) {
			if !input.Confirmed {
				return errorResult(fmt.Errorf("delete_dns_policy: set confirmed=true to confirm the deletion"))
			}
			if input.PolicyID == "" {
				return errorResult(fmt.Errorf("delete_dns_policy: policy_id is required"))
			}
			if err := client.DeleteDNSPolicy(ctx, input.SiteID, input.PolicyID); err != nil {
				return errorResult(fmt.Errorf("delete_dns_policy: %w", err))
			}
			return textResult(fmt.Sprintf("DNS policy %s deleted", input.PolicyID))
		})
	}

	// ── Firewall policies ────────────────────────────────────────────────────

	type firewallPolicyInput struct {
		SiteID   string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		PolicyID string `json:"policy_id"          jsonschema:"firewall policy ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_firewall_policy",
		Description: "Get details for a specific firewall policy by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input firewallPolicyInput) (*mcp.CallToolResult, any, error) {
		if input.PolicyID == "" {
			return errorResult(fmt.Errorf("get_firewall_policy: policy_id is required"))
		}
		policy, err := client.GetFirewallPolicy(ctx, input.SiteID, input.PolicyID)
		if err != nil {
			return errorResult(fmt.Errorf("get_firewall_policy: %w", err))
		}
		return jsonResult(policy)
	})

	type setFirewallPolicyEnabledInput struct {
		SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		PolicyID  string `json:"policy_id"          jsonschema:"firewall policy ID"`
		Enabled   *bool  `json:"enabled"            jsonschema:"true to enable the policy, false to disable"`
		Confirmed bool   `json:"confirmed"          jsonschema:"must be true to confirm the change"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "set_firewall_policy_enabled",
		Description: "Enable or disable a firewall policy. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input setFirewallPolicyEnabledInput) (*mcp.CallToolResult, any, error) {
		if !input.Confirmed {
			return errorResult(fmt.Errorf("set_firewall_policy_enabled: set confirmed=true to confirm the change"))
		}
		if input.PolicyID == "" {
			return errorResult(fmt.Errorf("set_firewall_policy_enabled: policy_id is required"))
		}
		if input.Enabled == nil {
			return errorResult(fmt.Errorf("set_firewall_policy_enabled: enabled is required"))
		}
		policy, err := client.SetFirewallPolicyEnabled(ctx, input.SiteID, input.PolicyID, *input.Enabled)
		if err != nil {
			return errorResult(fmt.Errorf("set_firewall_policy_enabled: %w", err))
		}
		return jsonResult(policy)
	})

	if allowDestructive {
		mcp.AddTool(s, &mcp.Tool{
			Name:        "delete_firewall_policy",
			Description: "Permanently delete a firewall policy by ID. Requires UNIFI_ALLOW_DESTRUCTIVE=true. Set confirmed=true to proceed.",
			Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
		}, func(ctx context.Context, _ *mcp.CallToolRequest, input struct {
			SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
			PolicyID  string `json:"policy_id"          jsonschema:"firewall policy ID"`
			Confirmed bool   `json:"confirmed"          jsonschema:"must be true to confirm the deletion"`
		},
		) (*mcp.CallToolResult, any, error) {
			if !input.Confirmed {
				return errorResult(fmt.Errorf("delete_firewall_policy: set confirmed=true to confirm the deletion"))
			}
			if input.PolicyID == "" {
				return errorResult(fmt.Errorf("delete_firewall_policy: policy_id is required"))
			}
			if err := client.DeleteFirewallPolicy(ctx, input.SiteID, input.PolicyID); err != nil {
				return errorResult(fmt.Errorf("delete_firewall_policy: %w", err))
			}
			return textResult(fmt.Sprintf("Firewall policy %s deleted", input.PolicyID))
		})
	}

	// ── Firewall zones ───────────────────────────────────────────────────────

	type firewallZoneInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		ZoneID string `json:"zone_id"            jsonschema:"firewall zone ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_firewall_zone",
		Description: "Get details for a specific firewall zone by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input firewallZoneInput) (*mcp.CallToolResult, any, error) {
		if input.ZoneID == "" {
			return errorResult(fmt.Errorf("get_firewall_zone: zone_id is required"))
		}
		zone, err := client.GetFirewallZone(ctx, input.SiteID, input.ZoneID)
		if err != nil {
			return errorResult(fmt.Errorf("get_firewall_zone: %w", err))
		}
		return jsonResult(zone)
	})

	type firewallZoneMutateInput struct {
		SiteID     string  `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		Name       string  `json:"name"               jsonschema:"zone name"`
		NetworkIDs *string `json:"network_ids,omitempty" jsonschema:"comma-separated list of network IDs to assign to this zone; omit for no networks"`
	}

	splitIDs := func(s *string) []string {
		if s == nil || *s == "" {
			return []string{}
		}
		parts := strings.Split(*s, ",")
		result := make([]string, 0, len(parts))
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_firewall_zone",
		Description: "Create a new firewall zone.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input firewallZoneMutateInput) (*mcp.CallToolResult, any, error) {
		if input.Name == "" {
			return errorResult(fmt.Errorf("create_firewall_zone: name is required"))
		}
		zone, err := client.CreateFirewallZone(ctx, input.SiteID, unifi.FirewallZoneRequest{
			Name:       input.Name,
			NetworkIDs: splitIDs(input.NetworkIDs),
		})
		if err != nil {
			return errorResult(fmt.Errorf("create_firewall_zone: %w", err))
		}
		return jsonResult(zone)
	})

	type updateFirewallZoneInput struct {
		SiteID     string  `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		ZoneID     string  `json:"zone_id"            jsonschema:"firewall zone ID"`
		Name       string  `json:"name"               jsonschema:"zone name"`
		NetworkIDs *string `json:"network_ids,omitempty" jsonschema:"comma-separated list of network IDs to assign to this zone; omit to preserve existing assignments; set to empty string to clear all networks"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_firewall_zone",
		Description: "Update an existing firewall zone by ID. network_ids replaces the full list; omit to preserve existing assignments.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input updateFirewallZoneInput) (*mcp.CallToolResult, any, error) {
		if input.ZoneID == "" {
			return errorResult(fmt.Errorf("update_firewall_zone: zone_id is required"))
		}
		if input.Name == "" {
			return errorResult(fmt.Errorf("update_firewall_zone: name is required"))
		}
		var networkIDs []string
		if input.NetworkIDs == nil {
			// Preserve existing assignments when network_ids is omitted.
			existing, err := client.GetFirewallZone(ctx, input.SiteID, input.ZoneID)
			if err != nil {
				return errorResult(fmt.Errorf("update_firewall_zone: fetch existing zone: %w", err))
			}
			networkIDs = existing.NetworkIDs
		} else {
			networkIDs = splitIDs(input.NetworkIDs)
		}
		zone, err := client.UpdateFirewallZone(ctx, input.SiteID, input.ZoneID, unifi.FirewallZoneRequest{
			Name:       input.Name,
			NetworkIDs: networkIDs,
		})
		if err != nil {
			return errorResult(fmt.Errorf("update_firewall_zone: %w", err))
		}
		return jsonResult(zone)
	})

	if allowDestructive {
		mcp.AddTool(s, &mcp.Tool{
			Name:        "delete_firewall_zone",
			Description: "Permanently delete a firewall zone by ID. Requires UNIFI_ALLOW_DESTRUCTIVE=true. Set confirmed=true to proceed.",
			Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
		}, func(ctx context.Context, _ *mcp.CallToolRequest, input struct {
			SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
			ZoneID    string `json:"zone_id"            jsonschema:"firewall zone ID"`
			Confirmed bool   `json:"confirmed"          jsonschema:"must be true to confirm the deletion"`
		},
		) (*mcp.CallToolResult, any, error) {
			if !input.Confirmed {
				return errorResult(fmt.Errorf("delete_firewall_zone: set confirmed=true to confirm the deletion"))
			}
			if input.ZoneID == "" {
				return errorResult(fmt.Errorf("delete_firewall_zone: zone_id is required"))
			}
			if err := client.DeleteFirewallZone(ctx, input.SiteID, input.ZoneID); err != nil {
				return errorResult(fmt.Errorf("delete_firewall_zone: %w", err))
			}
			return textResult(fmt.Sprintf("Firewall zone %s deleted", input.ZoneID))
		})
	}

	// ── ACL Rules ────────────────────────────────────────────────────────────

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_acl_rules",
		Description: "List all ACL rules for a site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		rules, err := client.ListACLRules(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_acl_rules: %w", err))
		}
		return jsonResult(rules)
	})

	type aclRuleInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		RuleID string `json:"rule_id"            jsonschema:"ACL rule ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_acl_rule",
		Description: "Get details for a specific ACL rule by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input aclRuleInput) (*mcp.CallToolResult, any, error) {
		if input.RuleID == "" {
			return errorResult(fmt.Errorf("get_acl_rule: rule_id is required"))
		}
		rule, err := client.GetACLRule(ctx, input.SiteID, input.RuleID)
		if err != nil {
			return errorResult(fmt.Errorf("get_acl_rule: %w", err))
		}
		return jsonResult(rule)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_acl_rule_ordering",
		Description: "Get the current ACL rule evaluation order.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		ordering, err := client.GetACLRuleOrdering(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("get_acl_rule_ordering: %w", err))
		}
		return jsonResult(ordering)
	})

	// ACL write tools are intentionally gated on allowDestructive (UNIFI_ALLOW_DESTRUCTIVE=true).
	//
	// Unlike firewall zones (organisational containers), any ACL mutation directly
	// controls which traffic is allowed or blocked. A misplaced BLOCK rule — or a
	// reorder that promotes one — can cause a complete network outage. Requiring the
	// explicit opt-in flag ensures that an AI session without it cannot issue any
	// ACL write at all, not just delete. The confirmed:true field per-call is a
	// secondary guard; the flag is the primary one.
	if allowDestructive {
		type aclRuleMutateInput struct {
			SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
			Type      string `json:"type"               jsonschema:"rule type: IPV4 or MAC"`
			Name      string `json:"name"               jsonschema:"rule name"`
			Action    string `json:"action"             jsonschema:"rule action: ALLOW or BLOCK"`
			Enabled   *bool  `json:"enabled"            jsonschema:"true to enable the rule, false to disable"`
			Confirmed bool   `json:"confirmed"          jsonschema:"must be true to confirm the change"`
		}

		mcp.AddTool(s, &mcp.Tool{
			Name:        "create_acl_rule",
			Description: "Create a new ACL rule. type must be IPV4 or MAC; action must be ALLOW or BLOCK. Requires UNIFI_ALLOW_DESTRUCTIVE=true. Set confirmed=true to proceed.",
			Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
		}, func(ctx context.Context, _ *mcp.CallToolRequest, input aclRuleMutateInput) (*mcp.CallToolResult, any, error) {
			if !input.Confirmed {
				return errorResult(fmt.Errorf("create_acl_rule: set confirmed=true to confirm the change"))
			}
			if input.Type == "" {
				return errorResult(fmt.Errorf("create_acl_rule: type is required (IPV4 or MAC)"))
			}
			if input.Name == "" {
				return errorResult(fmt.Errorf("create_acl_rule: name is required"))
			}
			if input.Action == "" {
				return errorResult(fmt.Errorf("create_acl_rule: action is required (ALLOW or BLOCK)"))
			}
			if input.Enabled == nil {
				return errorResult(fmt.Errorf("create_acl_rule: enabled is required"))
			}
			rule, err := client.CreateACLRule(ctx, input.SiteID, unifi.ACLRuleRequest{
				Type:    input.Type,
				Name:    input.Name,
				Action:  input.Action,
				Enabled: *input.Enabled,
			})
			if err != nil {
				return errorResult(fmt.Errorf("create_acl_rule: %w", err))
			}
			return jsonResult(rule)
		})

		type updateACLRuleInput struct {
			SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
			RuleID    string `json:"rule_id"            jsonschema:"ACL rule ID"`
			Type      string `json:"type"               jsonschema:"rule type: IPV4 or MAC"`
			Name      string `json:"name"               jsonschema:"rule name"`
			Action    string `json:"action"             jsonschema:"rule action: ALLOW or BLOCK"`
			Enabled   *bool  `json:"enabled"            jsonschema:"true to enable the rule, false to disable"`
			Confirmed bool   `json:"confirmed"          jsonschema:"must be true to confirm the change"`
		}

		mcp.AddTool(s, &mcp.Tool{
			Name:        "update_acl_rule",
			Description: "Update an existing ACL rule by ID. type must be IPV4 or MAC; action must be ALLOW or BLOCK. Requires UNIFI_ALLOW_DESTRUCTIVE=true. Set confirmed=true to proceed.",
			Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
		}, func(ctx context.Context, _ *mcp.CallToolRequest, input updateACLRuleInput) (*mcp.CallToolResult, any, error) {
			if !input.Confirmed {
				return errorResult(fmt.Errorf("update_acl_rule: set confirmed=true to confirm the change"))
			}
			if input.RuleID == "" {
				return errorResult(fmt.Errorf("update_acl_rule: rule_id is required"))
			}
			if input.Type == "" {
				return errorResult(fmt.Errorf("update_acl_rule: type is required (IPV4 or MAC)"))
			}
			if input.Name == "" {
				return errorResult(fmt.Errorf("update_acl_rule: name is required"))
			}
			if input.Action == "" {
				return errorResult(fmt.Errorf("update_acl_rule: action is required (ALLOW or BLOCK)"))
			}
			if input.Enabled == nil {
				return errorResult(fmt.Errorf("update_acl_rule: enabled is required"))
			}
			rule, err := client.UpdateACLRule(ctx, input.SiteID, input.RuleID, unifi.ACLRuleRequest{
				Type:    input.Type,
				Name:    input.Name,
				Action:  input.Action,
				Enabled: *input.Enabled,
			})
			if err != nil {
				return errorResult(fmt.Errorf("update_acl_rule: %w", err))
			}
			return jsonResult(rule)
		})

		mcp.AddTool(s, &mcp.Tool{
			Name:        "set_acl_rule_enabled",
			Description: "Enable or disable an ACL rule. Requires UNIFI_ALLOW_DESTRUCTIVE=true. Set confirmed=true to proceed.",
			Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
		}, func(ctx context.Context, _ *mcp.CallToolRequest, input struct {
			SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
			RuleID    string `json:"rule_id"            jsonschema:"ACL rule ID"`
			Enabled   *bool  `json:"enabled"            jsonschema:"true to enable the rule, false to disable"`
			Confirmed bool   `json:"confirmed"          jsonschema:"must be true to confirm the change"`
		},
		) (*mcp.CallToolResult, any, error) {
			if !input.Confirmed {
				return errorResult(fmt.Errorf("set_acl_rule_enabled: set confirmed=true to confirm the change"))
			}
			if input.RuleID == "" {
				return errorResult(fmt.Errorf("set_acl_rule_enabled: rule_id is required"))
			}
			if input.Enabled == nil {
				return errorResult(fmt.Errorf("set_acl_rule_enabled: enabled is required"))
			}
			rule, err := client.SetACLRuleEnabled(ctx, input.SiteID, input.RuleID, *input.Enabled)
			if err != nil {
				return errorResult(fmt.Errorf("set_acl_rule_enabled: %w", err))
			}
			return jsonResult(rule)
		})

		mcp.AddTool(s, &mcp.Tool{
			Name:        "reorder_acl_rules",
			Description: "Set the ACL rule evaluation order. Provide rule_ids as a comma-separated list of rule IDs in the desired order. Requires UNIFI_ALLOW_DESTRUCTIVE=true. Set confirmed=true to proceed.",
			Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
		}, func(ctx context.Context, _ *mcp.CallToolRequest, input struct {
			SiteID    string  `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
			RuleIDs   *string `json:"rule_ids"           jsonschema:"comma-separated list of ACL rule IDs in the desired evaluation order"`
			Confirmed bool    `json:"confirmed"          jsonschema:"must be true to confirm the change"`
		},
		) (*mcp.CallToolResult, any, error) {
			if !input.Confirmed {
				return errorResult(fmt.Errorf("reorder_acl_rules: set confirmed=true to confirm the change"))
			}
			if input.RuleIDs == nil || *input.RuleIDs == "" {
				return errorResult(fmt.Errorf("reorder_acl_rules: rule_ids is required"))
			}
			ordering, err := client.ReorderACLRules(ctx, input.SiteID, splitIDs(input.RuleIDs))
			if err != nil {
				return errorResult(fmt.Errorf("reorder_acl_rules: %w", err))
			}
			return jsonResult(ordering)
		})

		mcp.AddTool(s, &mcp.Tool{
			Name:        "delete_acl_rule",
			Description: "Permanently delete an ACL rule by ID. Requires UNIFI_ALLOW_DESTRUCTIVE=true. Set confirmed=true to proceed.",
			Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
		}, func(ctx context.Context, _ *mcp.CallToolRequest, input struct {
			SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
			RuleID    string `json:"rule_id"            jsonschema:"ACL rule ID"`
			Confirmed bool   `json:"confirmed"          jsonschema:"must be true to confirm the deletion"`
		},
		) (*mcp.CallToolResult, any, error) {
			if !input.Confirmed {
				return errorResult(fmt.Errorf("delete_acl_rule: set confirmed=true to confirm the deletion"))
			}
			if input.RuleID == "" {
				return errorResult(fmt.Errorf("delete_acl_rule: rule_id is required"))
			}
			if err := client.DeleteACLRule(ctx, input.SiteID, input.RuleID); err != nil {
				return errorResult(fmt.Errorf("delete_acl_rule: %w", err))
			}
			return textResult(fmt.Sprintf("ACL rule %s deleted", input.RuleID))
		})
	}

	// ── Hotspot Vouchers ─────────────────────────────────────────────────────

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_vouchers",
		Description: "List all hotspot vouchers for a site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		vouchers, err := client.ListVouchers(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_vouchers: %w", err))
		}
		return jsonResult(vouchers)
	})

	type voucherInput struct {
		SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		VoucherID string `json:"voucher_id"         jsonschema:"voucher ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_voucher",
		Description: "Get details for a specific hotspot voucher by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input voucherInput) (*mcp.CallToolResult, any, error) {
		if input.VoucherID == "" {
			return errorResult(fmt.Errorf("get_voucher: voucher_id is required"))
		}
		voucher, err := client.GetVoucher(ctx, input.SiteID, input.VoucherID)
		if err != nil {
			return errorResult(fmt.Errorf("get_voucher: %w", err))
		}
		return jsonResult(voucher)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_vouchers",
		Description: "Generate one or more hotspot vouchers. count is required (minimum 1). time_limit_minutes and data_limit_mb are optional (0 = unlimited). Set confirmed=true to proceed.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input struct {
		SiteID           string `json:"site_id,omitempty"       jsonschema:"site ID; omit to use default"`
		Count            int    `json:"count"                   jsonschema:"number of vouchers to generate (minimum 1)"`
		Name             string `json:"name,omitempty"          jsonschema:"optional label for the vouchers"`
		TimeLimitMinutes int    `json:"time_limit_minutes,omitempty" jsonschema:"access duration in minutes; 0 or omit for unlimited"`
		DataLimitMb      int    `json:"data_limit_mb,omitempty" jsonschema:"data cap in MB; 0 or omit for unlimited"`
		Confirmed        bool   `json:"confirmed"               jsonschema:"must be true to confirm the creation"`
	},
	) (*mcp.CallToolResult, any, error) {
		if !input.Confirmed {
			return errorResult(fmt.Errorf("create_vouchers: set confirmed=true to confirm the creation"))
		}
		if input.Count < 1 {
			return errorResult(fmt.Errorf("create_vouchers: count must be at least 1"))
		}
		vouchers, err := client.CreateVouchers(ctx, input.SiteID, unifi.VoucherRequest{
			Count:            input.Count,
			Name:             input.Name,
			TimeLimitMinutes: input.TimeLimitMinutes,
			DataLimitMb:      input.DataLimitMb,
		})
		if err != nil {
			return errorResult(fmt.Errorf("create_vouchers: %w", err))
		}
		return jsonResult(vouchers)
	})

	if allowDestructive {
		mcp.AddTool(s, &mcp.Tool{
			Name:        "delete_voucher",
			Description: "Permanently revoke a hotspot voucher by ID. Requires UNIFI_ALLOW_DESTRUCTIVE=true. Set confirmed=true to proceed.",
			Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
		}, func(ctx context.Context, _ *mcp.CallToolRequest, input struct {
			SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
			VoucherID string `json:"voucher_id"         jsonschema:"voucher ID"`
			Confirmed bool   `json:"confirmed"          jsonschema:"must be true to confirm the deletion"`
		},
		) (*mcp.CallToolResult, any, error) {
			if !input.Confirmed {
				return errorResult(fmt.Errorf("delete_voucher: set confirmed=true to confirm the deletion"))
			}
			if input.VoucherID == "" {
				return errorResult(fmt.Errorf("delete_voucher: voucher_id is required"))
			}
			if err := client.DeleteVoucher(ctx, input.SiteID, input.VoucherID); err != nil {
				return errorResult(fmt.Errorf("delete_voucher: %w", err))
			}
			return textResult(fmt.Sprintf("Voucher %s deleted", input.VoucherID))
		})
	}
}
