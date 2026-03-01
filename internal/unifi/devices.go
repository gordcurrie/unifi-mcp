package unifi

import (
	"context"
	"fmt"
)

// ListDevices returns all adopted devices from GET /integration/v1/sites/{siteID}/devices.
// Pass an empty siteID to use the client default. offset and limit control pagination; 0 means use the API default.
func (c *Client) ListDevices(ctx context.Context, siteID string, offset, limit int) (Page[Device], error) {
	id := c.site(siteID)
	data, err := c.getWithQuery(ctx, fmt.Sprintf("/integration/v1/sites/%s/devices", id), offset, limit)
	if err != nil {
		return Page[Device]{}, fmt.Errorf("ListDevices %s: %w", id, err)
	}
	page, err := decodeV1List[Device](data)
	if err != nil {
		return Page[Device]{}, fmt.Errorf("ListDevices %s: %w", id, err)
	}
	return page, nil
}

// GetDevice returns a single device from GET /integration/v1/sites/{siteID}/devices/{deviceID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetDevice(ctx context.Context, siteID, deviceID string) (Device, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/devices/%s", id, deviceID))
	if err != nil {
		return Device{}, fmt.Errorf("GetDevice %s %s: %w", id, deviceID, err)
	}
	dev, err := decodeV1[Device](data)
	if err != nil {
		return Device{}, fmt.Errorf("GetDevice %s %s: %w", id, deviceID, err)
	}
	return dev, nil
}

// RestartDevice sends a RESTART action via POST /integration/v1/sites/{siteID}/devices/{deviceID}/actions.
// Pass an empty siteID to use the client default.
func (c *Client) RestartDevice(ctx context.Context, siteID, deviceID string) error {
	id := c.site(siteID)
	_, err := c.postWithBody(ctx,
		fmt.Sprintf("/integration/v1/sites/%s/devices/%s/actions", id, deviceID),
		deviceActionRequest{Action: "RESTART"},
	)
	if err != nil {
		return fmt.Errorf("RestartDevice %s %s: %w", id, deviceID, err)
	}
	return nil
}

// PowerCyclePort power-cycles a single PoE port on a switch via
// POST /integration/v1/sites/{siteID}/devices/{deviceID}/interfaces/ports/{portIdx}/actions.
// Pass an empty siteID to use the client default.
func (c *Client) PowerCyclePort(ctx context.Context, siteID, deviceID string, portIdx int) error {
	id := c.site(siteID)
	_, err := c.postWithBody(ctx,
		fmt.Sprintf("/integration/v1/sites/%s/devices/%s/interfaces/ports/%d/actions", id, deviceID, portIdx),
		deviceActionRequest{Action: "POWER_CYCLE"},
	)
	if err != nil {
		return fmt.Errorf("PowerCyclePort %s %s port %d: %w", id, deviceID, portIdx, err)
	}
	return nil
}

// ListPendingDevices returns devices visible on the network but not yet adopted from
// GET /integration/v1/pending-devices. This endpoint is not site-scoped.
// offset and limit control pagination; 0 means use the API default.
func (c *Client) ListPendingDevices(ctx context.Context, offset, limit int) (Page[PendingDevice], error) {
	data, err := c.getWithQuery(ctx, "/integration/v1/pending-devices", offset, limit)
	if err != nil {
		return Page[PendingDevice]{}, fmt.Errorf("ListPendingDevices: %w", err)
	}
	page, err := decodeV1List[PendingDevice](data)
	if err != nil {
		return Page[PendingDevice]{}, fmt.Errorf("ListPendingDevices: %w", err)
	}
	return page, nil
}

// GetDeviceStats returns the latest statistics for a device from
// GET /integration/v1/sites/{siteID}/devices/{deviceID}/statistics/latest.
// Pass an empty siteID to use the client default.
func (c *Client) GetDeviceStats(ctx context.Context, siteID, deviceID string) (DeviceStats, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/devices/%s/statistics/latest", id, deviceID))
	if err != nil {
		return DeviceStats{}, fmt.Errorf("GetDeviceStats %s %s: %w", id, deviceID, err)
	}
	stats, err := decodeV1[DeviceStats](data)
	if err != nil {
		return DeviceStats{}, fmt.Errorf("GetDeviceStats %s %s: %w", id, deviceID, err)
	}
	return stats, nil
}
