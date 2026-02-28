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

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_clients",
		Description: "List all currently connected clients on the network.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		clients, err := client.ListClients(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("list_clients: %w", err)
		}
		return jsonResult(clients)
	})
}
