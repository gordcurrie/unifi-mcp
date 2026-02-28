package unifi

import (
	"context"
	"fmt"
)

// GetInfo returns application information from GET /v1/info.
func (c *Client) GetInfo(ctx context.Context) (ApplicationInfo, error) {
	data, err := c.get(ctx, "/v1/info")
	if err != nil {
		return ApplicationInfo{}, fmt.Errorf("GetInfo: %w", err)
	}
	info, err := decodeV1Single[ApplicationInfo](data)
	if err != nil {
		return ApplicationInfo{}, fmt.Errorf("GetInfo: %w", err)
	}
	return info, nil
}

// ListSites returns all sites from GET /v1/sites.
func (c *Client) ListSites(ctx context.Context) ([]Site, error) {
	data, err := c.get(ctx, "/v1/sites")
	if err != nil {
		return nil, fmt.Errorf("ListSites: %w", err)
	}
	sites, err := decodeV1List[Site](data)
	if err != nil {
		return nil, fmt.Errorf("ListSites: %w", err)
	}
	return sites, nil
}

// GetSite returns a single site from GET /v1/sites/{siteID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetSite(ctx context.Context, siteID string) (Site, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/v1/sites/%s", id))
	if err != nil {
		return Site{}, fmt.Errorf("GetSite %s: %w", id, err)
	}
	s, err := decodeV1Single[Site](data)
	if err != nil {
		return Site{}, fmt.Errorf("GetSite %s: %w", id, err)
	}
	return s, nil
}
