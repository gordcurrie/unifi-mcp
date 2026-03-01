package tools

import (
	"context"
	"fmt"

	"github.com/gordcurrie/unifi-mcp/internal/unifi"
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

	destructiveTrue := true

	mcp.AddTool(s, &mcp.Tool{
		Name:        "authorize_guest_client",
		Description: "Authorize a connected client for guest network access. Set confirmed=true to proceed. Optional: time_limit_minutes, data_limit_mb, download_bandwidth_kbps, upload_bandwidth_kbps (0 = unlimited).",
		Annotations: &mcp.ToolAnnotations{DestructiveHint: &destructiveTrue},
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input struct {
		SiteID                string `json:"site_id,omitempty"              jsonschema:"site ID; omit to use default"`
		ClientID              string `json:"client_id"                      jsonschema:"client ID to authorize"`
		TimeLimitMinutes      int    `json:"time_limit_minutes,omitempty"   jsonschema:"access duration in minutes; 0 or omit for unlimited"`
		DataLimitMb           int    `json:"data_limit_mb,omitempty"        jsonschema:"data cap in MB; 0 or omit for unlimited"`
		DownloadBandwidthKbps int    `json:"download_bandwidth_kbps,omitempty" jsonschema:"download rate limit in Kbps; 0 or omit for unlimited"`
		UploadBandwidthKbps   int    `json:"upload_bandwidth_kbps,omitempty"   jsonschema:"upload rate limit in Kbps; 0 or omit for unlimited"`
		Confirmed             bool   `json:"confirmed"                      jsonschema:"must be true to confirm the authorization"`
	},
	) (*mcp.CallToolResult, any, error) {
		if !input.Confirmed {
			return errorResult(fmt.Errorf("authorize_guest_client: set confirmed=true to confirm the authorization"))
		}
		if input.ClientID == "" {
			return errorResult(fmt.Errorf("authorize_guest_client: client_id is required"))
		}
		err := client.AuthorizeGuestClient(ctx, input.SiteID, input.ClientID, unifi.GuestAuthRequest{
			Action:                "AUTHORIZE_GUEST_ACCESS",
			TimeLimitMinutes:      input.TimeLimitMinutes,
			DataLimitMb:           input.DataLimitMb,
			DownloadBandwidthKbps: input.DownloadBandwidthKbps,
			UploadBandwidthKbps:   input.UploadBandwidthKbps,
		})
		if err != nil {
			return errorResult(fmt.Errorf("authorize_guest_client: %w", err))
		}
		return textResult(fmt.Sprintf("client %s authorized for guest access", input.ClientID))
	})
}
