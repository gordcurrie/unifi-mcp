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
	RestartDevice(ctx context.Context, site, mac string) error
	LocateDevice(ctx context.Context, site, mac string) error
	UnlocateDevice(ctx context.Context, site, mac string) error
	UpgradeDevice(ctx context.Context, site, mac string) error
	ForceReprovisionDevice(ctx context.Context, site, mac string) error
	RunSpeedTest(ctx context.Context, site string) error
	GetSpeedTestStatus(ctx context.Context, site string) (unifi.SpeedTestStatus, error)

	// Clients
	ListActiveClients(ctx context.Context, siteID string) ([]unifi.ActiveClient, error)
	ListKnownClients(ctx context.Context, siteID string) ([]unifi.KnownClient, error)
	BlockClient(ctx context.Context, site, mac string) error
	UnblockClient(ctx context.Context, site, mac string) error
	KickClient(ctx context.Context, site, mac string) error
	ForgetClient(ctx context.Context, site, mac string) error

	// Statistics
	GetSiteStats(ctx context.Context, siteID string) (unifi.SiteStats, error)
	GetDeviceStats(ctx context.Context, siteID, deviceID string) (unifi.Device, error)
	GetClientStats(ctx context.Context, siteID, clientID string) (unifi.ActiveClient, error)
	ListEvents(ctx context.Context, site string, limit int) ([]unifi.Event, error)
	ListAlarms(ctx context.Context, site string, archivedOnly bool) ([]unifi.Alarm, error)

	// Network
	ListWLANs(ctx context.Context, site string) ([]unifi.WLAN, error)
	SetWLANEnabled(ctx context.Context, site, wlanID string, enabled bool) error
	ListNetworks(ctx context.Context, site string) ([]unifi.NetworkConf, error)
	ListFirewallRules(ctx context.Context, site string) ([]unifi.FirewallRule, error)
	ListPortForwards(ctx context.Context, site string) ([]unifi.PortForward, error)
}
