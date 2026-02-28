package tools

import (
	"context"

	"github.com/gordcurrie/unifi-mcp/internal/unifi"
)

// unifiClient is the interface the tools layer requires from the UniFi client.
// *unifi.Client satisfies this interface automatically.
type unifiClient interface {
	// Sites
	GetInfo(ctx context.Context) (unifi.ApplicationInfo, error)
	ListSites(ctx context.Context) ([]unifi.Site, error)
	GetSite(ctx context.Context, siteID string) (unifi.Site, error)

	// Devices
	ListDevices(ctx context.Context, siteID string) ([]unifi.Device, error)
	GetDevice(ctx context.Context, siteID, deviceID string) (unifi.Device, error)
	RestartDevice(ctx context.Context, siteID, deviceID string) error
	GetDeviceStats(ctx context.Context, siteID, deviceID string) (unifi.DeviceStats, error)

	// Clients
	ListClients(ctx context.Context, siteID string) ([]unifi.NetworkClient, error)

	// Network
	ListWiFiBroadcasts(ctx context.Context, siteID string) ([]unifi.WiFiBroadcast, error)
	ListNetworks(ctx context.Context, siteID string) ([]unifi.NetworkConf, error)
	ListFirewallPolicies(ctx context.Context, siteID string) ([]unifi.FirewallPolicy, error)
	ListFirewallZones(ctx context.Context, siteID string) ([]unifi.FirewallZone, error)
	ListACLRules(ctx context.Context, siteID string) ([]unifi.ACLRule, error)
}
