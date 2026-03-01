package tools

import (
	"context"
	"fmt"

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
		TTLSeconds  int    `json:"ttl_seconds,omitempty"  jsonschema:"TTL in seconds; 0 uses the server default"`
		Enabled     *bool  `json:"enabled"                jsonschema:"true to activate the policy, false to create disabled"`
	}
	type updateDNSPolicyInput struct {
		SiteID      string `json:"site_id,omitempty"      jsonschema:"site ID; omit to use default"`
		PolicyID    string `json:"policy_id"              jsonschema:"DNS policy ID to update"`
		Type        string `json:"type"                   jsonschema:"policy type, e.g. A_RECORD"`
		Domain      string `json:"domain"                 jsonschema:"domain name to resolve"`
		IPv4Address string `json:"ipv4_address,omitempty" jsonschema:"IPv4 address the domain maps to"`
		TTLSeconds  int    `json:"ttl_seconds,omitempty"  jsonschema:"TTL in seconds; 0 uses the server default"`
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
		Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
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
		Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
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
}
