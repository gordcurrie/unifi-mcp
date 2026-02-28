package unifi

import (
	"context"
	"fmt"
)

// ListWLANs returns all configured WLANs from GET /api/s/{site}/rest/wlanconf.
// Pass an empty site to use the client default.
func (c *Client) ListWLANs(ctx context.Context, site string) ([]WLAN, error) {
	s := c.site(site)
	data, err := c.get(ctx, fmt.Sprintf("/api/s/%s/rest/wlanconf", s))
	if err != nil {
		return nil, fmt.Errorf("ListWLANs %s: %w", s, err)
	}
	wlans, err := decodeLegacy[WLAN](data)
	if err != nil {
		return nil, fmt.Errorf("ListWLANs %s: %w", s, err)
	}
	return wlans, nil
}

// SetWLANEnabled enables or disables the WLAN with the given ID.
// Pass an empty site to use the client default.
func (c *Client) SetWLANEnabled(ctx context.Context, site, wlanID string, enabled bool) error {
	s := c.site(site)
	_, err := c.put(ctx, fmt.Sprintf("/api/s/%s/rest/wlanconf/%s", s, wlanID), wlanEnabledRequest{Enabled: enabled})
	if err != nil {
		return fmt.Errorf("SetWLANEnabled %s enabled=%v: %w", wlanID, enabled, err)
	}
	return nil
}

// ListNetworks returns all configured networks from GET /api/s/{site}/rest/networkconf.
// Pass an empty site to use the client default.
func (c *Client) ListNetworks(ctx context.Context, site string) ([]NetworkConf, error) {
	s := c.site(site)
	data, err := c.get(ctx, fmt.Sprintf("/api/s/%s/rest/networkconf", s))
	if err != nil {
		return nil, fmt.Errorf("ListNetworks %s: %w", s, err)
	}
	networks, err := decodeLegacy[NetworkConf](data)
	if err != nil {
		return nil, fmt.Errorf("ListNetworks %s: %w", s, err)
	}
	return networks, nil
}

// ListFirewallRules returns user-defined firewall rules from GET /api/s/{site}/rest/firewallrule.
// Pass an empty site to use the client default.
func (c *Client) ListFirewallRules(ctx context.Context, site string) ([]FirewallRule, error) {
	s := c.site(site)
	data, err := c.get(ctx, fmt.Sprintf("/api/s/%s/rest/firewallrule", s))
	if err != nil {
		return nil, fmt.Errorf("ListFirewallRules %s: %w", s, err)
	}
	rules, err := decodeLegacy[FirewallRule](data)
	if err != nil {
		return nil, fmt.Errorf("ListFirewallRules %s: %w", s, err)
	}
	return rules, nil
}

// ListPortForwards returns configured port forwarding rules from GET /api/s/{site}/rest/portforward.
// Pass an empty site to use the client default.
func (c *Client) ListPortForwards(ctx context.Context, site string) ([]PortForward, error) {
	s := c.site(site)
	data, err := c.get(ctx, fmt.Sprintf("/api/s/%s/rest/portforward", s))
	if err != nil {
		return nil, fmt.Errorf("ListPortForwards %s: %w", s, err)
	}
	rules, err := decodeLegacy[PortForward](data)
	if err != nil {
		return nil, fmt.Errorf("ListPortForwards %s: %w", s, err)
	}
	return rules, nil
}
