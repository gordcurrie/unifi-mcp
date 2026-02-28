package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerClientTools(s *mcp.Server, client unifiClient) {
	type siteInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
	}
	type clientMACInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		MAC    string `json:"mac"               jsonschema:"client MAC address"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_active_clients",
		Description: "List all currently connected clients on the network.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		clients, err := client.ListActiveClients(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("list_active_clients: %w", err)
		}
		return jsonResult(clients)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_known_clients",
		Description: "List all historically known clients (connected or not).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		clients, err := client.ListKnownClients(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("list_known_clients: %w", err)
		}
		return jsonResult(clients)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "block_client",
		Description: "Block a client device from the network by MAC address.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input clientMACInput) (*mcp.CallToolResult, any, error) {
		if err := client.BlockClient(ctx, input.SiteID, input.MAC); err != nil {
			return nil, nil, fmt.Errorf("block_client: %w", err)
		}
		return textResult("blocked " + input.MAC)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "unblock_client",
		Description: "Unblock a previously blocked client device.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input clientMACInput) (*mcp.CallToolResult, any, error) {
		if err := client.UnblockClient(ctx, input.SiteID, input.MAC); err != nil {
			return nil, nil, fmt.Errorf("unblock_client: %w", err)
		}
		return textResult("unblocked " + input.MAC)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "kick_client",
		Description: "Disconnect a wireless client (forces reconnect, does not ban).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input clientMACInput) (*mcp.CallToolResult, any, error) {
		if err := client.KickClient(ctx, input.SiteID, input.MAC); err != nil {
			return nil, nil, fmt.Errorf("kick_client: %w", err)
		}
		return textResult("kicked " + input.MAC)
	})
}
