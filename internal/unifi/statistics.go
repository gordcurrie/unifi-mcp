package unifi

import (
	"context"
	"fmt"
)

// GetSiteStats returns aggregate statistics for the site.
// Pass an empty siteID to use the client default.
func (c *Client) GetSiteStats(ctx context.Context, siteID string) (SiteStats, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/v1/sites/%s/statistics/site", id))
	if err != nil {
		return SiteStats{}, fmt.Errorf("GetSiteStats %s: %w", id, err)
	}
	stats, err := decodeV1Single[SiteStats](data)
	if err != nil {
		return SiteStats{}, fmt.Errorf("GetSiteStats %s: %w", id, err)
	}
	return stats, nil
}

// GetDeviceStats returns statistics for a single device.
// Pass an empty siteID to use the client default.
func (c *Client) GetDeviceStats(ctx context.Context, siteID, deviceID string) (Device, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/v1/sites/%s/statistics/devices/%s", id, deviceID))
	if err != nil {
		return Device{}, fmt.Errorf("GetDeviceStats %s %s: %w", id, deviceID, err)
	}
	dev, err := decodeV1Single[Device](data)
	if err != nil {
		return Device{}, fmt.Errorf("GetDeviceStats %s %s: %w", id, deviceID, err)
	}
	return dev, nil
}

// GetClientStats returns statistics for a single client.
// Pass an empty siteID to use the client default.
func (c *Client) GetClientStats(ctx context.Context, siteID, clientID string) (ActiveClient, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/v1/sites/%s/statistics/clients/%s", id, clientID))
	if err != nil {
		return ActiveClient{}, fmt.Errorf("GetClientStats %s %s: %w", id, clientID, err)
	}
	cl, err := decodeV1Single[ActiveClient](data)
	if err != nil {
		return ActiveClient{}, fmt.Errorf("GetClientStats %s %s: %w", id, clientID, err)
	}
	return cl, nil
}

// ListEvents returns events from GET /api/s/{site}/stat/event.
// Pass limit=0 to use the server default (typically 200).
func (c *Client) ListEvents(ctx context.Context, site string, limit int) ([]Event, error) {
	s := c.site(site)
	path := fmt.Sprintf("/api/s/%s/stat/event", s)
	if limit > 0 {
		path = fmt.Sprintf("%s?_limit=%d", path, limit)
	}
	data, err := c.get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("ListEvents %s: %w", s, err)
	}
	events, err := decodeLegacy[Event](data)
	if err != nil {
		return nil, fmt.Errorf("ListEvents %s: %w", s, err)
	}
	return events, nil
}

// ListAlarms returns alarms from GET /api/s/{site}/stat/alarm.
// Pass archivedOnly=true to return only archived alarms.
func (c *Client) ListAlarms(ctx context.Context, site string, archivedOnly bool) ([]Alarm, error) {
	s := c.site(site)
	path := fmt.Sprintf("/api/s/%s/stat/alarm", s)
	if archivedOnly {
		path += "?archived=true"
	}
	data, err := c.get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("ListAlarms %s: %w", s, err)
	}
	alarms, err := decodeLegacy[Alarm](data)
	if err != nil {
		return nil, fmt.Errorf("ListAlarms %s: %w", s, err)
	}
	return alarms, nil
}
