package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerNetworkTools(s *mcp.Server, client unifiClient) {
	type siteInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_wifi_broadcasts",
		Description: "List all WiFi broadcast configurations (SSIDs).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		broadcasts, err := client.ListWiFiBroadcasts(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("list_wifi_broadcasts: %w", err)
		}
		return jsonResult(broadcasts)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_networks",
		Description: "List all configured networks (VLANs, LAN segments).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		nets, err := client.ListNetworks(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("list_networks: %w", err)
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
			return nil, nil, fmt.Errorf("list_firewall_policies: %w", err)
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
			return nil, nil, fmt.Errorf("list_firewall_zones: %w", err)
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
			return nil, nil, fmt.Errorf("list_acl_rules: %w", err)
		}
		return jsonResult(rules)
	})
}
