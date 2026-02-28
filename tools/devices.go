package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerDeviceTools(s *mcp.Server, client unifiClient) {
	type siteInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
	}
	type deviceInput struct {
		SiteID   string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		DeviceID string `json:"device_id"         jsonschema:"device ID"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_devices",
		Description: "List all adopted devices (APs, switches, gateways) for a site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		devices, err := client.ListDevices(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("list_devices: %w", err)
		}
		return jsonResult(devices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_device",
		Description: "Get details for a specific device by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input deviceInput) (*mcp.CallToolResult, any, error) {
		dev, err := client.GetDevice(ctx, input.SiteID, input.DeviceID)
		if err != nil {
			return nil, nil, fmt.Errorf("get_device: %w", err)
		}
		return jsonResult(dev)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "restart_device",
		Description: "Restart a UniFi device by device ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input deviceInput) (*mcp.CallToolResult, any, error) {
		if err := client.RestartDevice(ctx, input.SiteID, input.DeviceID); err != nil {
			return nil, nil, fmt.Errorf("restart_device: %w", err)
		}
		return textResult("restart command sent to " + input.DeviceID)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_device_stats",
		Description: "Get the latest statistics (CPU, memory, uptime) for a specific device.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input deviceInput) (*mcp.CallToolResult, any, error) {
		stats, err := client.GetDeviceStats(ctx, input.SiteID, input.DeviceID)
		if err != nil {
			return nil, nil, fmt.Errorf("get_device_stats: %w", err)
		}
		return jsonResult(stats)
	})
}
