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
	ListPendingDevices(ctx context.Context) ([]unifi.PendingDevice, error)
	PowerCyclePort(ctx context.Context, siteID, deviceID string, portIdx int) error

	// Clients
	ListClients(ctx context.Context, siteID string) ([]unifi.NetworkClient, error)
	GetClient(ctx context.Context, siteID, clientID string) (unifi.NetworkClient, error)

	// Network
	ListWiFiBroadcasts(ctx context.Context, siteID string) ([]unifi.WiFiBroadcast, error)
	GetWiFiBroadcast(ctx context.Context, siteID, broadcastID string) (unifi.WiFiBroadcast, error)
	SetWiFiBroadcastEnabled(ctx context.Context, siteID, broadcastID string, enabled bool) (unifi.WiFiBroadcast, error)
	ListNetworks(ctx context.Context, siteID string) ([]unifi.NetworkConf, error)
	ListFirewallPolicies(ctx context.Context, siteID string) ([]unifi.FirewallPolicy, error)
	ListFirewallZones(ctx context.Context, siteID string) ([]unifi.FirewallZone, error)
	ListACLRules(ctx context.Context, siteID string) ([]unifi.ACLRule, error)
	ListTrafficMatchingLists(ctx context.Context, siteID string) ([]unifi.TrafficMatchingList, error)
	GetTrafficMatchingList(ctx context.Context, siteID, listID string) (unifi.TrafficMatchingList, error)
	ListWANs(ctx context.Context, siteID string) ([]unifi.WAN, error)
	ListVPNTunnels(ctx context.Context, siteID string) ([]unifi.VPNTunnel, error)
	ListVPNServers(ctx context.Context, siteID string) ([]unifi.VPNServer, error)

	// DNS policies
	ListDNSPolicies(ctx context.Context, siteID string) ([]unifi.DNSPolicy, error)
	GetDNSPolicy(ctx context.Context, siteID, policyID string) (unifi.DNSPolicy, error)
	CreateDNSPolicy(ctx context.Context, siteID string, req unifi.DNSPolicyRequest) (unifi.DNSPolicy, error)
	UpdateDNSPolicy(ctx context.Context, siteID, policyID string, req unifi.DNSPolicyRequest) (unifi.DNSPolicy, error)
	DeleteDNSPolicy(ctx context.Context, siteID, policyID string) error
}
