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
	ListSites(ctx context.Context, offset, limit int) (unifi.Page[unifi.Site], error)
	GetSite(ctx context.Context, siteID string) (unifi.Site, error)

	// Devices
	ListDevices(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.Device], error)
	GetDevice(ctx context.Context, siteID, deviceID string) (unifi.Device, error)
	RestartDevice(ctx context.Context, siteID, deviceID string) error
	GetDeviceStats(ctx context.Context, siteID, deviceID string) (unifi.DeviceStats, error)
	ListPendingDevices(ctx context.Context, offset, limit int) (unifi.Page[unifi.PendingDevice], error)
	PowerCyclePort(ctx context.Context, siteID, deviceID string, portIdx int) error

	// Clients
	ListClients(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.NetworkClient], error)
	GetClient(ctx context.Context, siteID, clientID string) (unifi.NetworkClient, error)
	AuthorizeGuestClient(ctx context.Context, siteID, clientID string, req unifi.GuestAuthRequest) error

	// Network
	ListWiFiBroadcasts(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.WiFiBroadcast], error)
	GetWiFiBroadcast(ctx context.Context, siteID, broadcastID string) (unifi.WiFiBroadcast, error)
	SetWiFiBroadcastEnabled(ctx context.Context, siteID, broadcastID string, enabled bool) (unifi.WiFiBroadcast, error)
	ListNetworks(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.NetworkConf], error)
	ListFirewallPolicies(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.FirewallPolicy], error)
	GetFirewallPolicy(ctx context.Context, siteID, policyID string) (unifi.FirewallPolicy, error)
	SetFirewallPolicyEnabled(ctx context.Context, siteID, policyID string, enabled bool) (unifi.FirewallPolicy, error)
	DeleteFirewallPolicy(ctx context.Context, siteID, policyID string) error
	ListFirewallZones(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.FirewallZone], error)
	GetFirewallZone(ctx context.Context, siteID, zoneID string) (unifi.FirewallZone, error)
	CreateFirewallZone(ctx context.Context, siteID string, req unifi.FirewallZoneRequest) (unifi.FirewallZone, error)
	UpdateFirewallZone(ctx context.Context, siteID, zoneID string, req unifi.FirewallZoneRequest) (unifi.FirewallZone, error)
	DeleteFirewallZone(ctx context.Context, siteID, zoneID string) error
	ListACLRules(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.ACLRule], error)
	GetACLRule(ctx context.Context, siteID, ruleID string) (unifi.ACLRule, error)
	CreateACLRule(ctx context.Context, siteID string, req unifi.ACLRuleRequest) (unifi.ACLRule, error)
	UpdateACLRule(ctx context.Context, siteID, ruleID string, req unifi.ACLRuleRequest) (unifi.ACLRule, error)
	DeleteACLRule(ctx context.Context, siteID, ruleID string) error
	SetACLRuleEnabled(ctx context.Context, siteID, ruleID string, enabled bool) (unifi.ACLRule, error)
	GetACLRuleOrdering(ctx context.Context, siteID string) (unifi.ACLRuleOrdering, error)
	ReorderACLRules(ctx context.Context, siteID string, orderedIDs []string) (unifi.ACLRuleOrdering, error)
	ListTrafficMatchingLists(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.TrafficMatchingList], error)
	GetTrafficMatchingList(ctx context.Context, siteID, listID string) (unifi.TrafficMatchingList, error)
	ListWANs(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.WAN], error)
	ListVPNTunnels(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.VPNTunnel], error)
	ListVPNServers(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.VPNServer], error)

	// DNS policies
	ListDNSPolicies(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.DNSPolicy], error)
	GetDNSPolicy(ctx context.Context, siteID, policyID string) (unifi.DNSPolicy, error)
	CreateDNSPolicy(ctx context.Context, siteID string, req unifi.DNSPolicyRequest) (unifi.DNSPolicy, error)
	UpdateDNSPolicy(ctx context.Context, siteID, policyID string, req unifi.DNSPolicyRequest) (unifi.DNSPolicy, error)
	DeleteDNSPolicy(ctx context.Context, siteID, policyID string) error

	// Vouchers
	ListVouchers(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.Voucher], error)
	GetVoucher(ctx context.Context, siteID, voucherID string) (unifi.Voucher, error)
	CreateVouchers(ctx context.Context, siteID string, req unifi.VoucherRequest) ([]unifi.Voucher, error)
	DeleteVoucher(ctx context.Context, siteID, voucherID string) error

	// Reference data
	ListDeviceTags(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.DeviceTag], error)
	ListDPICategories(ctx context.Context, offset, limit int) (unifi.Page[unifi.DPICategory], error)
	ListDPIApplications(ctx context.Context, offset, limit int) (unifi.Page[unifi.DPIApplication], error)
	ListRADIUSProfiles(ctx context.Context, siteID string, offset, limit int) (unifi.Page[unifi.RADIUSProfile], error)
}
