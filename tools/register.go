// Package tools registers all MCP tools for the unifi-mcp server.
package tools

import (
	"github.com/gordcurrie/unifi-mcp/internal/unifi"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Config controls optional tool categories.
type Config struct {
	// AllowDestructive enables tools that permanently modify or delete data.
	// Controlled by the UNIFI_ALLOW_DESTRUCTIVE env var.
	AllowDestructive bool
}

// RegisterAll registers every enabled tool group with the MCP server.
func RegisterAll(s *mcp.Server, client *unifi.Client, cfg Config) {
	registerSiteTools(s, client)
	registerDeviceTools(s, client)
	registerClientTools(s, client)
	registerStatisticsTools(s, client)
	registerNetworkTools(s, client)
	if cfg.AllowDestructive {
		registerDestructiveTools(s, client)
	}
}
