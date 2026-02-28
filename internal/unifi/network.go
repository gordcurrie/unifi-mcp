package unifi

import (
	"context"
	"fmt"
)

// ListWiFiBroadcasts returns all WiFi broadcast (SSID) configurations from
// GET /integration/v1/sites/{siteID}/wifi/broadcasts.
// Pass an empty siteID to use the client default.
func (c *Client) ListWiFiBroadcasts(ctx context.Context, siteID string) ([]WiFiBroadcast, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/wifi/broadcasts", id))
	if err != nil {
		return nil, fmt.Errorf("ListWiFiBroadcasts %s: %w", id, err)
	}
	broadcasts, err := decodeV1List[WiFiBroadcast](data)
	if err != nil {
		return nil, fmt.Errorf("ListWiFiBroadcasts %s: %w", id, err)
	}
	return broadcasts, nil
}

// GetWiFiBroadcast returns a single WiFi broadcast configuration from
// GET /integration/v1/sites/{siteID}/wifi/broadcasts/{broadcastID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetWiFiBroadcast(ctx context.Context, siteID, broadcastID string) (WiFiBroadcast, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/wifi/broadcasts/%s", id, broadcastID))
	if err != nil {
		return WiFiBroadcast{}, fmt.Errorf("GetWiFiBroadcast %s %s: %w", id, broadcastID, err)
	}
	broadcast, err := decodeV1[WiFiBroadcast](data)
	if err != nil {
		return WiFiBroadcast{}, fmt.Errorf("GetWiFiBroadcast %s %s: %w", id, broadcastID, err)
	}
	return broadcast, nil
}

// ListNetworks returns all configured networks from GET /integration/v1/sites/{siteID}/networks.
// Pass an empty siteID to use the client default.
func (c *Client) ListNetworks(ctx context.Context, siteID string) ([]NetworkConf, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/networks", id))
	if err != nil {
		return nil, fmt.Errorf("ListNetworks %s: %w", id, err)
	}
	networks, err := decodeV1List[NetworkConf](data)
	if err != nil {
		return nil, fmt.Errorf("ListNetworks %s: %w", id, err)
	}
	return networks, nil
}

// ListFirewallPolicies returns all firewall policies from GET /integration/v1/sites/{siteID}/firewall/policies.
// Pass an empty siteID to use the client default.
func (c *Client) ListFirewallPolicies(ctx context.Context, siteID string) ([]FirewallPolicy, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/firewall/policies", id))
	if err != nil {
		return nil, fmt.Errorf("ListFirewallPolicies %s: %w", id, err)
	}
	policies, err := decodeV1List[FirewallPolicy](data)
	if err != nil {
		return nil, fmt.Errorf("ListFirewallPolicies %s: %w", id, err)
	}
	return policies, nil
}

// ListFirewallZones returns all firewall zones from GET /integration/v1/sites/{siteID}/firewall/zones.
// Pass an empty siteID to use the client default.
func (c *Client) ListFirewallZones(ctx context.Context, siteID string) ([]FirewallZone, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/firewall/zones", id))
	if err != nil {
		return nil, fmt.Errorf("ListFirewallZones %s: %w", id, err)
	}
	zones, err := decodeV1List[FirewallZone](data)
	if err != nil {
		return nil, fmt.Errorf("ListFirewallZones %s: %w", id, err)
	}
	return zones, nil
}

// ListACLRules returns all ACL rules from GET /integration/v1/sites/{siteID}/acl-rules.
// Pass an empty siteID to use the client default.
func (c *Client) ListACLRules(ctx context.Context, siteID string) ([]ACLRule, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/acl-rules", id))
	if err != nil {
		return nil, fmt.Errorf("ListACLRules %s: %w", id, err)
	}
	rules, err := decodeV1List[ACLRule](data)
	if err != nil {
		return nil, fmt.Errorf("ListACLRules %s: %w", id, err)
	}
	return rules, nil
}
