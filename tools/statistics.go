package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerStatisticsTools(s *mcp.Server, client unifiClient) {
	type siteInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
	}
	type deviceStatsInput struct {
		SiteID   string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		DeviceID string `json:"device_id"         jsonschema:"device ID"`
	}
	type clientStatsInput struct {
		SiteID   string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		ClientID string `json:"client_id"         jsonschema:"client ID"`
	}
	type eventsInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		Limit  int    `json:"limit,omitempty"   jsonschema:"max events to return; 0 = server default (200)"`
	}
	type alarmsInput struct {
		SiteID       string `json:"site_id,omitempty"       jsonschema:"site ID; omit to use default"`
		ArchivedOnly bool   `json:"archived_only,omitempty" jsonschema:"return only archived alarms"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_site_statistics",
		Description: "Get aggregate traffic and client statistics for the site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		stats, err := client.GetSiteStats(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("get_site_statistics: %w", err)
		}
		return jsonResult(stats)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_device_statistics",
		Description: "Get statistics for a specific device.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input deviceStatsInput) (*mcp.CallToolResult, any, error) {
		dev, err := client.GetDeviceStats(ctx, input.SiteID, input.DeviceID)
		if err != nil {
			return nil, nil, fmt.Errorf("get_device_statistics: %w", err)
		}
		return jsonResult(dev)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_client_statistics",
		Description: "Get statistics for a specific connected client.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input clientStatsInput) (*mcp.CallToolResult, any, error) {
		cl, err := client.GetClientStats(ctx, input.SiteID, input.ClientID)
		if err != nil {
			return nil, nil, fmt.Errorf("get_client_statistics: %w", err)
		}
		return jsonResult(cl)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_events",
		Description: "List recent network events. Set limit to control how many are returned (0 = server default).",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input eventsInput) (*mcp.CallToolResult, any, error) {
		events, err := client.ListEvents(ctx, input.SiteID, input.Limit)
		if err != nil {
			return nil, nil, fmt.Errorf("list_events: %w", err)
		}
		return jsonResult(events)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_alarms",
		Description: "List active (or archived) alarms for the site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input alarmsInput) (*mcp.CallToolResult, any, error) {
		alarms, err := client.ListAlarms(ctx, input.SiteID, input.ArchivedOnly)
		if err != nil {
			return nil, nil, fmt.Errorf("list_alarms: %w", err)
		}
		return jsonResult(alarms)
	})
}
