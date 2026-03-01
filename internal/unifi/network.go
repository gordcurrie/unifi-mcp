package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// ListWiFiBroadcasts returns one page of WiFi broadcast (SSID) configurations from
// GET /integration/v1/sites/{siteID}/wifi/broadcasts.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListWiFiBroadcasts(ctx context.Context, siteID string, offset, limit int) (Page[WiFiBroadcast], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/wifi/broadcasts", id), offset, limit)
	if err != nil {
		return Page[WiFiBroadcast]{}, fmt.Errorf("ListWiFiBroadcasts %s: %w", id, err)
	}
	page, err := decodeV1List[WiFiBroadcast](data)
	if err != nil {
		return Page[WiFiBroadcast]{}, fmt.Errorf("ListWiFiBroadcasts %s: %w", id, err)
	}
	return page, nil
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

// ListNetworks returns one page of configured networks from GET /integration/v1/sites/{siteID}/networks.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListNetworks(ctx context.Context, siteID string, offset, limit int) (Page[NetworkConf], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/networks", id), offset, limit)
	if err != nil {
		return Page[NetworkConf]{}, fmt.Errorf("ListNetworks %s: %w", id, err)
	}
	page, err := decodeV1List[NetworkConf](data)
	if err != nil {
		return Page[NetworkConf]{}, fmt.Errorf("ListNetworks %s: %w", id, err)
	}
	return page, nil
}

// ListFirewallPolicies returns one page of firewall policies from GET /integration/v1/sites/{siteID}/firewall/policies.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListFirewallPolicies(ctx context.Context, siteID string, offset, limit int) (Page[FirewallPolicy], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/firewall/policies", id), offset, limit)
	if err != nil {
		return Page[FirewallPolicy]{}, fmt.Errorf("ListFirewallPolicies %s: %w", id, err)
	}
	page, err := decodeV1List[FirewallPolicy](data)
	if err != nil {
		return Page[FirewallPolicy]{}, fmt.Errorf("ListFirewallPolicies %s: %w", id, err)
	}
	return page, nil
}

// ListFirewallZones returns one page of firewall zones from GET /integration/v1/sites/{siteID}/firewall/zones.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListFirewallZones(ctx context.Context, siteID string, offset, limit int) (Page[FirewallZone], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/firewall/zones", id), offset, limit)
	if err != nil {
		return Page[FirewallZone]{}, fmt.Errorf("ListFirewallZones %s: %w", id, err)
	}
	page, err := decodeV1List[FirewallZone](data)
	if err != nil {
		return Page[FirewallZone]{}, fmt.Errorf("ListFirewallZones %s: %w", id, err)
	}
	return page, nil
}

// GetFirewallPolicy returns a single firewall policy from
// GET /integration/v1/sites/{siteID}/firewall/policies/{policyID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetFirewallPolicy(ctx context.Context, siteID, policyID string) (FirewallPolicy, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/firewall/policies/%s", id, policyID))
	if err != nil {
		return FirewallPolicy{}, fmt.Errorf("GetFirewallPolicy %s %s: %w", id, policyID, err)
	}
	policy, err := decodeV1[FirewallPolicy](data)
	if err != nil {
		return FirewallPolicy{}, fmt.Errorf("GetFirewallPolicy %s %s: %w", id, policyID, err)
	}
	return policy, nil
}

// SetFirewallPolicyEnabled enables or disables a firewall policy via
// GET then PUT /integration/v1/sites/{siteID}/firewall/policies/{policyID}.
// It fetches the current resource as a raw map to preserve all API fields,
// sets the enabled flag, and PUTs the full object back.
// Pass an empty siteID to use the client default.
func (c *Client) SetFirewallPolicyEnabled(ctx context.Context, siteID, policyID string, enabled bool) (FirewallPolicy, error) {
	id := c.site(siteID)
	path := fmt.Sprintf("/integration/v1/sites/%s/firewall/policies/%s", id, policyID)

	raw, err := c.get(ctx, path)
	if err != nil {
		return FirewallPolicy{}, fmt.Errorf("SetFirewallPolicyEnabled %s %s: get: %w", id, policyID, err)
	}
	body, err := decodeV1[map[string]any](raw)
	if err != nil {
		return FirewallPolicy{}, fmt.Errorf("SetFirewallPolicyEnabled %s %s: decode: %w", id, policyID, err)
	}
	delete(body, "id")
	delete(body, "metadata")
	body["enabled"] = enabled

	updated, err := c.put(ctx, path, body)
	if err != nil {
		return FirewallPolicy{}, fmt.Errorf("SetFirewallPolicyEnabled %s %s: put: %w", id, policyID, err)
	}
	policy, err := decodeV1[FirewallPolicy](updated)
	if err != nil {
		return FirewallPolicy{}, fmt.Errorf("SetFirewallPolicyEnabled %s %s: decode response: %w", id, policyID, err)
	}
	return policy, nil
}

// DeleteFirewallPolicy deletes a firewall policy via
// DELETE /integration/v1/sites/{siteID}/firewall/policies/{policyID}.
// Pass an empty siteID to use the client default.
func (c *Client) DeleteFirewallPolicy(ctx context.Context, siteID, policyID string) error {
	id := c.site(siteID)
	if err := c.delete(ctx, fmt.Sprintf("/integration/v1/sites/%s/firewall/policies/%s", id, policyID)); err != nil {
		return fmt.Errorf("DeleteFirewallPolicy %s %s: %w", id, policyID, err)
	}
	return nil
}

// GetFirewallZone returns a single firewall zone from
// GET /integration/v1/sites/{siteID}/firewall/zones/{zoneID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetFirewallZone(ctx context.Context, siteID, zoneID string) (FirewallZone, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/firewall/zones/%s", id, zoneID))
	if err != nil {
		return FirewallZone{}, fmt.Errorf("GetFirewallZone %s %s: %w", id, zoneID, err)
	}
	zone, err := decodeV1[FirewallZone](data)
	if err != nil {
		return FirewallZone{}, fmt.Errorf("GetFirewallZone %s %s: %w", id, zoneID, err)
	}
	return zone, nil
}

// CreateFirewallZone creates a new firewall zone via
// POST /integration/v1/sites/{siteID}/firewall/zones.
// Pass an empty siteID to use the client default.
func (c *Client) CreateFirewallZone(ctx context.Context, siteID string, req FirewallZoneRequest) (FirewallZone, error) {
	id := c.site(siteID)
	data, err := c.postWithBody(ctx, fmt.Sprintf("/integration/v1/sites/%s/firewall/zones", id), req)
	if err != nil {
		return FirewallZone{}, fmt.Errorf("CreateFirewallZone %s: %w", id, err)
	}
	zone, err := decodeV1[FirewallZone](data)
	if err != nil {
		return FirewallZone{}, fmt.Errorf("CreateFirewallZone %s: %w", id, err)
	}
	return zone, nil
}

// UpdateFirewallZone replaces a firewall zone via
// PUT /integration/v1/sites/{siteID}/firewall/zones/{zoneID}.
// Pass an empty siteID to use the client default.
func (c *Client) UpdateFirewallZone(ctx context.Context, siteID, zoneID string, req FirewallZoneRequest) (FirewallZone, error) {
	id := c.site(siteID)
	data, err := c.put(ctx, fmt.Sprintf("/integration/v1/sites/%s/firewall/zones/%s", id, zoneID), req)
	if err != nil {
		return FirewallZone{}, fmt.Errorf("UpdateFirewallZone %s %s: %w", id, zoneID, err)
	}
	zone, err := decodeV1[FirewallZone](data)
	if err != nil {
		return FirewallZone{}, fmt.Errorf("UpdateFirewallZone %s %s: %w", id, zoneID, err)
	}
	return zone, nil
}

// DeleteFirewallZone deletes a firewall zone via
// DELETE /integration/v1/sites/{siteID}/firewall/zones/{zoneID}.
// Pass an empty siteID to use the client default.
func (c *Client) DeleteFirewallZone(ctx context.Context, siteID, zoneID string) error {
	id := c.site(siteID)
	if err := c.delete(ctx, fmt.Sprintf("/integration/v1/sites/%s/firewall/zones/%s", id, zoneID)); err != nil {
		return fmt.Errorf("DeleteFirewallZone %s %s: %w", id, zoneID, err)
	}
	return nil
}

// ListACLRules returns one page of ACL rules from GET /integration/v1/sites/{siteID}/acl-rules.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListACLRules(ctx context.Context, siteID string, offset, limit int) (Page[ACLRule], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/acl-rules", id), offset, limit)
	if err != nil {
		return Page[ACLRule]{}, fmt.Errorf("ListACLRules %s: %w", id, err)
	}
	page, err := decodeV1List[ACLRule](data)
	if err != nil {
		return Page[ACLRule]{}, fmt.Errorf("ListACLRules %s: %w", id, err)
	}
	return page, nil
}

// GetACLRule returns a single ACL rule from
// GET /integration/v1/sites/{siteID}/acl-rules/{ruleID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetACLRule(ctx context.Context, siteID, ruleID string) (ACLRule, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/acl-rules/%s", id, ruleID))
	if err != nil {
		return ACLRule{}, fmt.Errorf("GetACLRule %s %s: %w", id, ruleID, err)
	}
	rule, err := decodeV1[ACLRule](data)
	if err != nil {
		return ACLRule{}, fmt.Errorf("GetACLRule %s %s: %w", id, ruleID, err)
	}
	return rule, nil
}

// CreateACLRule creates a new ACL rule via
// POST /integration/v1/sites/{siteID}/acl-rules.
// Pass an empty siteID to use the client default.
func (c *Client) CreateACLRule(ctx context.Context, siteID string, req ACLRuleRequest) (ACLRule, error) {
	id := c.site(siteID)
	data, err := c.postWithBody(ctx, fmt.Sprintf("/integration/v1/sites/%s/acl-rules", id), req)
	if err != nil {
		return ACLRule{}, fmt.Errorf("CreateACLRule %s: %w", id, err)
	}
	rule, err := decodeV1[ACLRule](data)
	if err != nil {
		return ACLRule{}, fmt.Errorf("CreateACLRule %s: %w", id, err)
	}
	return rule, nil
}

// UpdateACLRule replaces an ACL rule via
// PUT /integration/v1/sites/{siteID}/acl-rules/{ruleID}.
// Pass an empty siteID to use the client default.
func (c *Client) UpdateACLRule(ctx context.Context, siteID, ruleID string, req ACLRuleRequest) (ACLRule, error) {
	id := c.site(siteID)
	data, err := c.put(ctx, fmt.Sprintf("/integration/v1/sites/%s/acl-rules/%s", id, ruleID), req)
	if err != nil {
		return ACLRule{}, fmt.Errorf("UpdateACLRule %s %s: %w", id, ruleID, err)
	}
	rule, err := decodeV1[ACLRule](data)
	if err != nil {
		return ACLRule{}, fmt.Errorf("UpdateACLRule %s %s: %w", id, ruleID, err)
	}
	return rule, nil
}

// DeleteACLRule deletes an ACL rule via
// DELETE /integration/v1/sites/{siteID}/acl-rules/{ruleID}.
// Pass an empty siteID to use the client default.
func (c *Client) DeleteACLRule(ctx context.Context, siteID, ruleID string) error {
	id := c.site(siteID)
	if err := c.delete(ctx, fmt.Sprintf("/integration/v1/sites/%s/acl-rules/%s", id, ruleID)); err != nil {
		return fmt.Errorf("DeleteACLRule %s %s: %w", id, ruleID, err)
	}
	return nil
}

// SetACLRuleEnabled enables or disables an ACL rule via GET then PUT.
// It fetches the current rule, flips the enabled flag, and PUTs back.
// Pass an empty siteID to use the client default.
func (c *Client) SetACLRuleEnabled(ctx context.Context, siteID, ruleID string, enabled bool) (ACLRule, error) {
	id := c.site(siteID)
	existing, err := c.GetACLRule(ctx, siteID, ruleID)
	if err != nil {
		return ACLRule{}, fmt.Errorf("SetACLRuleEnabled %s %s: get: %w", id, ruleID, err)
	}
	req := ACLRuleRequest{
		Type:    existing.Type,
		Name:    existing.Name,
		Action:  existing.Action,
		Enabled: enabled,
	}
	rule, err := c.UpdateACLRule(ctx, siteID, ruleID, req)
	if err != nil {
		return ACLRule{}, fmt.Errorf("SetACLRuleEnabled %s %s: put: %w", id, ruleID, err)
	}
	return rule, nil
}

// GetACLRuleOrdering returns the current ACL rule ordering via
// GET /integration/v1/sites/{siteID}/acl-rules/ordering.
// Pass an empty siteID to use the client default.
func (c *Client) GetACLRuleOrdering(ctx context.Context, siteID string) (ACLRuleOrdering, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/acl-rules/ordering", id))
	if err != nil {
		return ACLRuleOrdering{}, fmt.Errorf("GetACLRuleOrdering %s: %w", id, err)
	}
	ordering, err := decodeV1[ACLRuleOrdering](data)
	if err != nil {
		return ACLRuleOrdering{}, fmt.Errorf("GetACLRuleOrdering %s: %w", id, err)
	}
	return ordering, nil
}

// ReorderACLRules sets the ACL rule ordering via
// PUT /integration/v1/sites/{siteID}/acl-rules/ordering.
// Pass an empty siteID to use the client default.
func (c *Client) ReorderACLRules(ctx context.Context, siteID string, orderedIDs []string) (ACLRuleOrdering, error) {
	id := c.site(siteID)
	body := ACLRuleOrdering{OrderedACLRuleIDs: orderedIDs}
	data, err := c.put(ctx, fmt.Sprintf("/integration/v1/sites/%s/acl-rules/ordering", id), body)
	if err != nil {
		return ACLRuleOrdering{}, fmt.Errorf("ReorderACLRules %s: %w", id, err)
	}
	ordering, err := decodeV1[ACLRuleOrdering](data)
	if err != nil {
		return ACLRuleOrdering{}, fmt.Errorf("ReorderACLRules %s: %w", id, err)
	}
	return ordering, nil
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

// ListTrafficMatchingLists returns one page of traffic matching lists from
// GET /integration/v1/sites/{siteID}/traffic-matching-lists.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListTrafficMatchingLists(ctx context.Context, siteID string, offset, limit int) (Page[TrafficMatchingList], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/traffic-matching-lists", id), offset, limit)
	if err != nil {
		return Page[TrafficMatchingList]{}, fmt.Errorf("ListTrafficMatchingLists %s: %w", id, err)
	}
	page, err := decodeV1List[TrafficMatchingList](data)
	if err != nil {
		return Page[TrafficMatchingList]{}, fmt.Errorf("ListTrafficMatchingLists %s: %w", id, err)
	}
	return page, nil
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

// ListWANs returns one page of WAN interface definitions from
// GET /integration/v1/sites/{siteID}/wans.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListWANs(ctx context.Context, siteID string, offset, limit int) (Page[WAN], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/wans", id), offset, limit)
	if err != nil {
		return Page[WAN]{}, fmt.Errorf("ListWANs %s: %w", id, err)
	}
	page, err := decodeV1List[WAN](data)
	if err != nil {
		return Page[WAN]{}, fmt.Errorf("ListWANs %s: %w", id, err)
	}
	return page, nil
}

// ListVPNTunnels returns one page of site-to-site VPN tunnels from
// GET /integration/v1/sites/{siteID}/vpn/site-to-site-tunnels.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListVPNTunnels(ctx context.Context, siteID string, offset, limit int) (Page[VPNTunnel], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/vpn/site-to-site-tunnels", id), offset, limit)
	if err != nil {
		return Page[VPNTunnel]{}, fmt.Errorf("ListVPNTunnels %s: %w", id, err)
	}
	page, err := decodeV1List[VPNTunnel](data)
	if err != nil {
		return Page[VPNTunnel]{}, fmt.Errorf("ListVPNTunnels %s: %w", id, err)
	}
	return page, nil
}

// ListVPNServers returns one page of VPN server configurations from
// GET /integration/v1/sites/{siteID}/vpn/servers.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListVPNServers(ctx context.Context, siteID string, offset, limit int) (Page[VPNServer], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/vpn/servers", id), offset, limit)
	if err != nil {
		return Page[VPNServer]{}, fmt.Errorf("ListVPNServers %s: %w", id, err)
	}
	page, err := decodeV1List[VPNServer](data)
	if err != nil {
		return Page[VPNServer]{}, fmt.Errorf("ListVPNServers %s: %w", id, err)
	}
	return page, nil
}

// ListDNSPolicies returns one page of DNS policies from
// GET /integration/v1/sites/{siteID}/dns/policies.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListDNSPolicies(ctx context.Context, siteID string, offset, limit int) (Page[DNSPolicy], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/dns/policies", id), offset, limit)
	if err != nil {
		return Page[DNSPolicy]{}, fmt.Errorf("ListDNSPolicies %s: %w", id, err)
	}
	page, err := decodeV1List[DNSPolicy](data)
	if err != nil {
		return Page[DNSPolicy]{}, fmt.Errorf("ListDNSPolicies %s: %w", id, err)
	}
	return page, nil
}

// GetDNSPolicy returns a single DNS policy from
// GET /integration/v1/sites/{siteID}/dns/policies/{policyID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetDNSPolicy(ctx context.Context, siteID, policyID string) (DNSPolicy, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/dns/policies/%s", id, policyID))
	if err != nil {
		return DNSPolicy{}, fmt.Errorf("GetDNSPolicy %s %s: %w", id, policyID, err)
	}
	policy, err := decodeV1[DNSPolicy](data)
	if err != nil {
		return DNSPolicy{}, fmt.Errorf("GetDNSPolicy %s %s: %w", id, policyID, err)
	}
	return policy, nil
}

// CreateDNSPolicy creates a new DNS policy via
// POST /integration/v1/sites/{siteID}/dns/policies.
// Pass an empty siteID to use the client default.
func (c *Client) CreateDNSPolicy(ctx context.Context, siteID string, req DNSPolicyRequest) (DNSPolicy, error) {
	id := c.site(siteID)
	data, err := c.postWithBody(ctx, fmt.Sprintf("/integration/v1/sites/%s/dns/policies", id), req)
	if err != nil {
		return DNSPolicy{}, fmt.Errorf("CreateDNSPolicy %s: %w", id, err)
	}
	policy, err := decodeV1[DNSPolicy](data)
	if err != nil {
		return DNSPolicy{}, fmt.Errorf("CreateDNSPolicy %s: %w", id, err)
	}
	return policy, nil
}

// UpdateDNSPolicy replaces a DNS policy via
// PUT /integration/v1/sites/{siteID}/dns/policies/{policyID}.
// Pass an empty siteID to use the client default.
func (c *Client) UpdateDNSPolicy(ctx context.Context, siteID, policyID string, req DNSPolicyRequest) (DNSPolicy, error) {
	id := c.site(siteID)
	data, err := c.put(ctx, fmt.Sprintf("/integration/v1/sites/%s/dns/policies/%s", id, policyID), req)
	if err != nil {
		return DNSPolicy{}, fmt.Errorf("UpdateDNSPolicy %s %s: %w", id, policyID, err)
	}
	policy, err := decodeV1[DNSPolicy](data)
	if err != nil {
		return DNSPolicy{}, fmt.Errorf("UpdateDNSPolicy %s %s: %w", id, policyID, err)
	}
	return policy, nil
}

// DeleteDNSPolicy deletes a DNS policy via
// DELETE /integration/v1/sites/{siteID}/dns/policies/{policyID}.
// Pass an empty siteID to use the client default.
func (c *Client) DeleteDNSPolicy(ctx context.Context, siteID, policyID string) error {
	id := c.site(siteID)
	if err := c.delete(ctx, fmt.Sprintf("/integration/v1/sites/%s/dns/policies/%s", id, policyID)); err != nil {
		return fmt.Errorf("DeleteDNSPolicy %s %s: %w", id, policyID, err)
	}
	return nil
}

// ListVouchers returns one page of hotspot vouchers for a site via
// GET /integration/v1/sites/{siteID}/hotspot/vouchers.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListVouchers(ctx context.Context, siteID string, offset, limit int) (Page[Voucher], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/hotspot/vouchers", id), offset, limit)
	if err != nil {
		return Page[Voucher]{}, fmt.Errorf("ListVouchers %s: %w", id, err)
	}
	page, err := decodeV1List[Voucher](data)
	if err != nil {
		return Page[Voucher]{}, fmt.Errorf("ListVouchers %s: %w", id, err)
	}
	return page, nil
}

// GetVoucher returns a single hotspot voucher via
// GET /integration/v1/sites/{siteID}/hotspot/vouchers/{voucherID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetVoucher(ctx context.Context, siteID, voucherID string) (Voucher, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/hotspot/vouchers/%s", id, voucherID))
	if err != nil {
		return Voucher{}, fmt.Errorf("GetVoucher %s %s: %w", id, voucherID, err)
	}
	voucher, err := decodeV1[Voucher](data)
	if err != nil {
		return Voucher{}, fmt.Errorf("GetVoucher %s %s: %w", id, voucherID, err)
	}
	return voucher, nil
}

// voucherCreateResponse is the envelope returned by POST /hotspot/vouchers.
type voucherCreateResponse struct {
	Vouchers []Voucher `json:"vouchers"`
}

// CreateVouchers generates one or more hotspot vouchers via
// POST /integration/v1/sites/{siteID}/hotspot/vouchers.
// Pass an empty siteID to use the client default.
func (c *Client) CreateVouchers(ctx context.Context, siteID string, req VoucherRequest) ([]Voucher, error) {
	id := c.site(siteID)
	data, err := c.postWithBody(ctx, fmt.Sprintf("/integration/v1/sites/%s/hotspot/vouchers", id), req)
	if err != nil {
		return nil, fmt.Errorf("CreateVouchers %s: %w", id, err)
	}
	var resp voucherCreateResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("CreateVouchers %s: decode response: %w", id, err)
	}
	return resp.Vouchers, nil
}

// DeleteVoucher revokes a single hotspot voucher via
// DELETE /integration/v1/sites/{siteID}/hotspot/vouchers/{voucherID}.
// Pass an empty siteID to use the client default.
func (c *Client) DeleteVoucher(ctx context.Context, siteID, voucherID string) error {
	id := c.site(siteID)
	if err := c.delete(ctx, fmt.Sprintf("/integration/v1/sites/%s/hotspot/vouchers/%s", id, voucherID)); err != nil {
		return fmt.Errorf("DeleteVoucher %s %s: %w", id, voucherID, err)
	}
	return nil
}

// ListDeviceTags returns one page of device tags from
// GET /integration/v1/sites/{siteID}/device-tags.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListDeviceTags(ctx context.Context, siteID string, offset, limit int) (Page[DeviceTag], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/device-tags", id), offset, limit)
	if err != nil {
		return Page[DeviceTag]{}, fmt.Errorf("ListDeviceTags %s: %w", id, err)
	}
	page, err := decodeV1List[DeviceTag](data)
	if err != nil {
		return Page[DeviceTag]{}, fmt.Errorf("ListDeviceTags %s: %w", id, err)
	}
	return page, nil
}

// ListDPICategories returns one page of DPI application categories from
// GET /integration/v1/dpi/categories.
// offset and limit control pagination; 0 means use the API default.
func (c *Client) ListDPICategories(ctx context.Context, offset, limit int) (Page[DPICategory], error) {
	data, err := c.getWithQuery(ctx, "/integration/v1/dpi/categories", offset, limit)
	if err != nil {
		return Page[DPICategory]{}, fmt.Errorf("ListDPICategories: %w", err)
	}
	page, err := decodeV1List[DPICategory](data)
	if err != nil {
		return Page[DPICategory]{}, fmt.Errorf("ListDPICategories: %w", err)
	}
	return page, nil
}

// ListDPIApplications returns one page of DPI applications from
// GET /integration/v1/dpi/applications.
// offset and limit control pagination; 0 means use the API default.
func (c *Client) ListDPIApplications(ctx context.Context, offset, limit int) (Page[DPIApplication], error) {
	data, err := c.getWithQuery(ctx, "/integration/v1/dpi/applications", offset, limit)
	if err != nil {
		return Page[DPIApplication]{}, fmt.Errorf("ListDPIApplications: %w", err)
	}
	page, err := decodeV1List[DPIApplication](data)
	if err != nil {
		return Page[DPIApplication]{}, fmt.Errorf("ListDPIApplications: %w", err)
	}
	return page, nil
}

// ListRADIUSProfiles returns one page of RADIUS profiles from
// GET /integration/v1/sites/{siteID}/radius/profiles.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListRADIUSProfiles(ctx context.Context, siteID string, offset, limit int) (Page[RADIUSProfile], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/radius/profiles", id), offset, limit)
	if err != nil {
		return Page[RADIUSProfile]{}, fmt.Errorf("ListRADIUSProfiles %s: %w", id, err)
	}
	page, err := decodeV1List[RADIUSProfile](data)
	if err != nil {
		return Page[RADIUSProfile]{}, fmt.Errorf("ListRADIUSProfiles %s: %w", id, err)
	}
	return page, nil
}
