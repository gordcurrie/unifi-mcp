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
	type clientInput struct {
		SiteID   string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		ClientID string `json:"client_id"         jsonschema:"client ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_clients",
		Description: "List all currently connected clients on the network.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		clients, err := client.ListClients(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_clients: %w", err))
		}
		return jsonResult(clients)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client",
		Description: "Get details for a specific connected client by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input clientInput) (*mcp.CallToolResult, any, error) {
		if input.ClientID == "" {
			return errorResult(fmt.Errorf("get_client: client_id is required"))
		}
		c, err := client.GetClient(ctx, input.SiteID, input.ClientID)
		if err != nil {
			return errorResult(fmt.Errorf("get_client: %w", err))
		}
		return jsonResult(c)
	})
}
