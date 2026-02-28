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

// SetWiFiBroadcastEnabled enables or disables a WiFi broadcast via
// GET then PUT /integration/v1/sites/{siteID}/wifi/broadcasts/{broadcastID}.
// It fetches the current resource as a raw map to preserve all API fields,
// sets the enabled flag, and PUTs the full object back.
// Pass an empty siteID to use the client default.
func (c *Client) SetWiFiBroadcastEnabled(ctx context.Context, siteID, broadcastID string, enabled bool) (WiFiBroadcast, error) {
	id := c.site(siteID)
	path := fmt.Sprintf("/integration/v1/sites/%s/wifi/broadcasts/%s", id, broadcastID)

	// GET current state as raw map for the round-trip PUT.
	raw, err := c.get(ctx, path)
	if err != nil {
		return WiFiBroadcast{}, fmt.Errorf("SetWiFiBroadcastEnabled %s %s: get: %w", id, broadcastID, err)
	}
	body, err := decodeV1[map[string]any](raw)
	if err != nil {
		return WiFiBroadcast{}, fmt.Errorf("SetWiFiBroadcastEnabled %s %s: decode: %w", id, broadcastID, err)
	}
	// Strip read-only fields the API rejects in PUT bodies.
	delete(body, "id")
	delete(body, "metadata")
	body["enabled"] = enabled

	updated, err := c.put(ctx, path, body)
	if err != nil {
		return WiFiBroadcast{}, fmt.Errorf("SetWiFiBroadcastEnabled %s %s: put: %w", id, broadcastID, err)
	}
	bc, err := decodeV1[WiFiBroadcast](updated)
	if err != nil {
		return WiFiBroadcast{}, fmt.Errorf("SetWiFiBroadcastEnabled %s %s: decode response: %w", id, broadcastID, err)
	}
	return bc, nil
}

// ListTrafficMatchingLists returns all traffic matching lists from
// GET /integration/v1/sites/{siteID}/traffic-matching-lists.
// Pass an empty siteID to use the client default.
func (c *Client) ListTrafficMatchingLists(ctx context.Context, siteID string) ([]TrafficMatchingList, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/traffic-matching-lists", id))
	if err != nil {
		return nil, fmt.Errorf("ListTrafficMatchingLists %s: %w", id, err)
	}
	lists, err := decodeV1List[TrafficMatchingList](data)
	if err != nil {
		return nil, fmt.Errorf("ListTrafficMatchingLists %s: %w", id, err)
	}
	return lists, nil
}

// GetTrafficMatchingList returns a single traffic matching list from
// GET /integration/v1/sites/{siteID}/traffic-matching-lists/{listID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetTrafficMatchingList(ctx context.Context, siteID, listID string) (TrafficMatchingList, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/traffic-matching-lists/%s", id, listID))
	if err != nil {
		return TrafficMatchingList{}, fmt.Errorf("GetTrafficMatchingList %s %s: %w", id, listID, err)
	}
	list, err := decodeV1[TrafficMatchingList](data)
	if err != nil {
		return TrafficMatchingList{}, fmt.Errorf("GetTrafficMatchingList %s %s: %w", id, listID, err)
	}
	return list, nil
}

// ListWANs returns all WAN interface definitions from
// GET /integration/v1/sites/{siteID}/wans.
// Pass an empty siteID to use the client default.
func (c *Client) ListWANs(ctx context.Context, siteID string) ([]WAN, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/wans", id))
	if err != nil {
		return nil, fmt.Errorf("ListWANs %s: %w", id, err)
	}
	wans, err := decodeV1List[WAN](data)
	if err != nil {
		return nil, fmt.Errorf("ListWANs %s: %w", id, err)
	}
	return wans, nil
}

// ListVPNTunnels returns all site-to-site VPN tunnels from
// GET /integration/v1/sites/{siteID}/vpn/site-to-site-tunnels.
// Pass an empty siteID to use the client default.
func (c *Client) ListVPNTunnels(ctx context.Context, siteID string) ([]VPNTunnel, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/vpn/site-to-site-tunnels", id))
	if err != nil {
		return nil, fmt.Errorf("ListVPNTunnels %s: %w", id, err)
	}
	tunnels, err := decodeV1List[VPNTunnel](data)
	if err != nil {
		return nil, fmt.Errorf("ListVPNTunnels %s: %w", id, err)
	}
	return tunnels, nil
}

// ListVPNServers returns all VPN server configurations from
// GET /integration/v1/sites/{siteID}/vpn/servers.
// Pass an empty siteID to use the client default.
func (c *Client) ListVPNServers(ctx context.Context, siteID string) ([]VPNServer, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/vpn/servers", id))
	if err != nil {
		return nil, fmt.Errorf("ListVPNServers %s: %w", id, err)
	}
	servers, err := decodeV1List[VPNServer](data)
	if err != nil {
		return nil, fmt.Errorf("ListVPNServers %s: %w", id, err)
	}
	return servers, nil
}
