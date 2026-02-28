package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerSiteTools(s *mcp.Server, client unifiClient) {
	type siteInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_application_info",
		Description: "Return UniFi controller application version and type.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		info, err := client.GetInfo(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("get_application_info: %w", err))
		}
		return jsonResult(info)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_sites",
		Description: "List all sites on the UniFi controller.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		sites, err := client.ListSites(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("list_sites: %w", err))
		}
		return jsonResult(sites)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_site",
		Description: "Get details for a specific site. Omit site_id to use the default site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		site, err := client.GetSite(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("get_site: %w", err))
		}
		return jsonResult(site)
	})
}
