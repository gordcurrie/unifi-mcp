package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerDeviceTools(s *mcp.Server, client unifiClient) {
	destructiveTrue := true

	type siteInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
	}
	type deviceInput struct {
		SiteID   string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		DeviceID string `json:"device_id"         jsonschema:"device ID"`
	}
	type restartDeviceInput struct {
		SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		DeviceID  string `json:"device_id"          jsonschema:"device ID"`
		Confirmed bool   `json:"confirmed"          jsonschema:"must be true to confirm the restart"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_devices",
		Description: "List all adopted devices (APs, switches, gateways) for a site.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		devices, err := client.ListDevices(ctx, input.SiteID)
		if err != nil {
			return errorResult(fmt.Errorf("list_devices: %w", err))
		}
		return jsonResult(devices)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_device",
		Description: "Get details for a specific device by ID.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input deviceInput) (*mcp.CallToolResult, any, error) {
		if input.DeviceID == "" {
			return errorResult(fmt.Errorf("get_device: device_id is required"))
		}
		dev, err := client.GetDevice(ctx, input.SiteID, input.DeviceID)
		if err != nil {
			return errorResult(fmt.Errorf("get_device: %w", err))
		}
		return jsonResult(dev)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "restart_device",
		Description: "Restart a UniFi device by device ID. Set confirmed=true to proceed.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input restartDeviceInput) (*mcp.CallToolResult, any, error) {
		if !input.Confirmed {
			return errorResult(fmt.Errorf("restart_device: set confirmed=true to confirm the restart"))
		}
		if input.DeviceID == "" {
			return errorResult(fmt.Errorf("restart_device: device_id is required"))
		}
		if err := client.RestartDevice(ctx, input.SiteID, input.DeviceID); err != nil {
			return errorResult(fmt.Errorf("restart_device: %w", err))
		}
		return textResult("restart command sent to " + input.DeviceID)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_device_stats",
		Description: "Get the latest statistics (CPU, memory, uptime) for a specific device.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input deviceInput) (*mcp.CallToolResult, any, error) {
		if input.DeviceID == "" {
			return errorResult(fmt.Errorf("get_device_stats: device_id is required"))
		}
		stats, err := client.GetDeviceStats(ctx, input.SiteID, input.DeviceID)
		if err != nil {
			return errorResult(fmt.Errorf("get_device_stats: %w", err))
		}
		return jsonResult(stats)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pending_devices",
		Description: "List devices visible on the network that have not yet been adopted.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		devices, err := client.ListPendingDevices(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("list_pending_devices: %w", err))
		}
		return jsonResult(devices)
	})
}
