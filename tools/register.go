// Package tools registers all MCP tools for the unifi-mcp server.
package tools

import (
	"github.com/gordcurrie/unifi-mcp/internal/unifi"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers every enabled tool group with the MCP server.
func RegisterAll(s *mcp.Server, client *unifi.Client) {
	registerSiteTools(s, client)
	registerDeviceTools(s, client)
	registerClientTools(s, client)
	registerNetworkTools(s, client)
}
