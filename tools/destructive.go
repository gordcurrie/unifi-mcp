package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerDestructiveTools(s *mcp.Server, client unifiClient) {
	destructiveHint := true

	type forgetClientInput struct {
		SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		MAC       string `json:"mac"               jsonschema:"client MAC address to forget"`
		Confirmed bool   `json:"confirmed"         jsonschema:"must be true to execute"`
	}
	type reprovisionInput struct {
		SiteID    string `json:"site_id,omitempty" jsonschema:"site ID; omit to use default"`
		MAC       string `json:"mac"               jsonschema:"device MAC address to reprovision"`
		Confirmed bool   `json:"confirmed"         jsonschema:"must be true to execute"`
	}

	mcp.AddTool(s, &mcp.Tool{
		Name:        "forget_client",
		Description: "Permanently remove a client record from the controller. Set confirmed=true to execute.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input forgetClientInput) (*mcp.CallToolResult, any, error) {
		if !input.Confirmed {
			return nil, nil, errors.New("forget_client: set confirmed=true to execute this destructive action")
		}
		if err := client.ForgetClient(ctx, input.SiteID, input.MAC); err != nil {
			return nil, nil, fmt.Errorf("forget_client: %w", err)
		}
		return textResult("client " + input.MAC + " forgotten")
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "force_reprovision_device",
		Description: "Force-reprovision a device, reapplying its full configuration. Set confirmed=true to execute.",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveHint},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input reprovisionInput) (*mcp.CallToolResult, any, error) {
		if !input.Confirmed {
			return nil, nil, errors.New("force_reprovision_device: set confirmed=true to execute this destructive action")
		}
		if err := client.ForceReprovisionDevice(ctx, input.SiteID, input.MAC); err != nil {
			return nil, nil, fmt.Errorf("force_reprovision_device: %w", err)
		}
		return textResult("reprovision initiated for " + input.MAC)
	})
}
