package unifi

import (
	"context"
	"errors"
	"fmt"
)

// GetInfo returns application information from GET /integration/v1/info.
func (c *Client) GetInfo(ctx context.Context) (ApplicationInfo, error) {
	data, err := c.get(ctx, "/integration/v1/info")
	if err != nil {
		return ApplicationInfo{}, fmt.Errorf("GetInfo: %w", err)
	}
	info, err := decodeV1[ApplicationInfo](data)
	if err != nil {
		return ApplicationInfo{}, fmt.Errorf("GetInfo: %w", err)
	}
	return info, nil
}

// ListSites returns all sites from GET /integration/v1/sites.
func (c *Client) ListSites(ctx context.Context) ([]Site, error) {
	data, err := c.get(ctx, "/integration/v1/sites")
	if err != nil {
		return nil, fmt.Errorf("ListSites: %w", err)
	}
	sites, err := decodeV1List[Site](data)
	if err != nil {
		return nil, fmt.Errorf("ListSites: %w", err)
	}
	return sites, nil
}

// GetSite returns a single site by ID from GET /integration/v1/sites (no single-get endpoint).
// Pass an empty siteID to use the client default.
func (c *Client) GetSite(ctx context.Context, siteID string) (Site, error) {
	id := c.site(siteID)
	sites, err := c.ListSites(ctx)
	if err != nil {
		return Site{}, fmt.Errorf("GetSite %s: %w", id, err)
	}
	for _, s := range sites {
		if s.ID == id {
			return s, nil
		}
	}
	return Site{}, fmt.Errorf("GetSite %s: %w", id, errors.New("site not found"))
}
