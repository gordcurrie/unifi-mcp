// Package tools registers all MCP tools for the unifi-mcp server.
package tools

import (
	"github.com/gordcurrie/unifi-mcp/internal/unifi"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers every enabled tool group with the MCP server.
// Set allowDestructive to true (via UNIFI_ALLOW_DESTRUCTIVE=true) to also
// register tools that permanently delete resources.
func RegisterAll(s *mcp.Server, client *unifi.Client, allowDestructive bool) {
	registerSiteTools(s, client)
	registerDeviceTools(s, client)
	registerClientTools(s, client)
	registerNetworkTools(s, client, allowDestructive)
}
