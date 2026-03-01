package unifi

import (
	"context"
	"errors"
	"fmt"
)

// ErrSiteNotFound is returned by GetSite when no site matches the requested ID.
var ErrSiteNotFound = errors.New("site not found")

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

// ListSites returns one page of sites from GET /integration/v1/sites.
// offset and limit control pagination; 0 means use the API default.
func (c *Client) ListSites(ctx context.Context, offset, limit int) (Page[Site], error) {
	data, err := c.getWithQuery(ctx, "/integration/v1/sites", offset, limit)
	if err != nil {
		return Page[Site]{}, fmt.Errorf("ListSites: %w", err)
	}
	page, err := decodeV1List[Site](data)
	if err != nil {
		return Page[Site]{}, fmt.Errorf("ListSites: %w", err)
	}
	return page, nil
}

// GetSite returns a single site by ID from GET /integration/v1/sites (no single-get endpoint).
// Pass an empty siteID to use the client default.
// GetSite paginates through all pages until the site is found.
func (c *Client) GetSite(ctx context.Context, siteID string) (Site, error) {
	id := c.site(siteID)
	const pageSize = 25
	offset := 0
	for {
		page, err := c.ListSites(ctx, offset, pageSize)
		if err != nil {
			return Site{}, fmt.Errorf("GetSite %s: %w", id, err)
		}
		for _, s := range page.Data {
			if s.ID == id {
				return s, nil
			}
		}
		offset += len(page.Data)
		if len(page.Data) == 0 || offset >= page.TotalCount {
			break
		}
	}
	return Site{}, fmt.Errorf("GetSite %s: %w", id, ErrSiteNotFound)
}
