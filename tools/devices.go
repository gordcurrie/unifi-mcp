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
	type deviceMACInput struct {
		SiteID string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		MAC    string `json:"mac"               jsonschema:"device MAC address"`
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
		Description: "Restart a UniFi device by MAC address.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input deviceMACInput) (*mcp.CallToolResult, any, error) {
		if err := client.RestartDevice(ctx, input.SiteID, input.MAC); err != nil {
			return nil, nil, fmt.Errorf("restart_device: %w", err)
		}
		return textResult("restart command sent to " + input.MAC)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "locate_device",
		Description: "Enable the locate/LED blink on a UniFi device.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input deviceMACInput) (*mcp.CallToolResult, any, error) {
		if err := client.LocateDevice(ctx, input.SiteID, input.MAC); err != nil {
			return nil, nil, fmt.Errorf("locate_device: %w", err)
		}
		return textResult("locate enabled for " + input.MAC)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "unlocate_device",
		Description: "Disable the locate/LED blink on a UniFi device.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input deviceMACInput) (*mcp.CallToolResult, any, error) {
		if err := client.UnlocateDevice(ctx, input.SiteID, input.MAC); err != nil {
			return nil, nil, fmt.Errorf("unlocate_device: %w", err)
		}
		return textResult("locate disabled for " + input.MAC)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "upgrade_device",
		Description: "Trigger a firmware upgrade on a UniFi device.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input deviceMACInput) (*mcp.CallToolResult, any, error) {
		if err := client.UpgradeDevice(ctx, input.SiteID, input.MAC); err != nil {
			return nil, nil, fmt.Errorf("upgrade_device: %w", err)
		}
		return textResult("upgrade initiated for " + input.MAC)
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "run_speed_test",
		Description: "Initiate a speed test from the gateway/UCG.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		if err := client.RunSpeedTest(ctx, input.SiteID); err != nil {
			return nil, nil, fmt.Errorf("run_speed_test: %w", err)
		}
		return textResult("speed test initiated")
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_speed_test_status",
		Description: "Get the most recent speed test result.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input siteInput) (*mcp.CallToolResult, any, error) {
		status, err := client.GetSpeedTestStatus(ctx, input.SiteID)
		if err != nil {
			return nil, nil, fmt.Errorf("get_speed_test_status: %w", err)
		}
		return jsonResult(status)
	})
}
