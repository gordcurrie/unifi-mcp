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
	type wlanInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		WLANID string `json:"wlan_id"           jsonschema:"WLAN configuration ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_wlans",
		Description: "List all configured wireless networks (SSIDs).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		wlans, err := client.ListWLANs(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("list_wlans: %w", err)
		}
		return jsonResult(wlans)
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
		Name:        "list_firewall_rules",
		Description: "List user-defined firewall rules.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		rules, err := client.ListFirewallRules(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("list_firewall_rules: %w", err)
		}
		return jsonResult(rules)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_port_forwards",
		Description: "List configured port forwarding rules.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		rules, err := client.ListPortForwards(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("list_port_forwards: %w", err)
		}
		return jsonResult(rules)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "enable_wlan",
		Description: "Enable a wireless network by its WLAN configuration ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input wlanInput) (*mcp.CallToolResult, any, error) {
		if err := client.SetWLANEnabled(ctx, input.SiteID, input.WLANID, true); err != nil {
			return nil, nil, fmt.Errorf("enable_wlan: %w", err)
		}
		return textResult("WLAN " + input.WLANID + " enabled")
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "disable_wlan",
		Description: "Disable a wireless network by its WLAN configuration ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input wlanInput) (*mcp.CallToolResult, any, error) {
		if err := client.SetWLANEnabled(ctx, input.SiteID, input.WLANID, false); err != nil {
			return nil, nil, fmt.Errorf("disable_wlan: %w", err)
		}
		return textResult("WLAN " + input.WLANID + " disabled")
	})
}
