package unifi

import (
	"context"
	"fmt"
)

// ListDevices returns all adopted devices from GET /v1/sites/{siteID}/devices.
// Pass an empty siteID to use the client default.
func (c *Client) ListDevices(ctx context.Context, siteID string) ([]Device, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/v1/sites/%s/devices", id))
	if err != nil {
		return nil, fmt.Errorf("ListDevices %s: %w", id, err)
	}
	devices, err := decodeV1List[Device](data)
	if err != nil {
		return nil, fmt.Errorf("ListDevices %s: %w", id, err)
	}
	return devices, nil
}

// GetDevice returns a single device from GET /v1/sites/{siteID}/devices/{deviceID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetDevice(ctx context.Context, siteID, deviceID string) (Device, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/v1/sites/%s/devices/%s", id, deviceID))
	if err != nil {
		return Device{}, fmt.Errorf("GetDevice %s: %w", deviceID, err)
	}
	dev, err := decodeV1Single[Device](data)
	if err != nil {
		return Device{}, fmt.Errorf("GetDevice %s: %w", deviceID, err)
	}
	return dev, nil
}

// RestartDevice sends a restart command via POST /api/s/{site}/cmd/devmgr.
func (c *Client) RestartDevice(ctx context.Context, site, mac string) error {
	s := c.site(site)
	_, err := c.postWithBody(ctx, fmt.Sprintf("/api/s/%s/cmd/devmgr", s), deviceCmdRequest{Cmd: "restart", MAC: mac})
	if err != nil {
		return fmt.Errorf("RestartDevice %s: %w", mac, err)
	}
	return nil
}

// LocateDevice enables the locate/blink LED on a device.
func (c *Client) LocateDevice(ctx context.Context, site, mac string) error {
	s := c.site(site)
	_, err := c.postWithBody(ctx, fmt.Sprintf("/api/s/%s/cmd/devmgr", s), deviceCmdRequest{Cmd: "set-locate", MAC: mac})
	if err != nil {
		return fmt.Errorf("LocateDevice %s: %w", mac, err)
	}
	return nil
}

// UnlocateDevice disables the locate/blink LED on a device.
func (c *Client) UnlocateDevice(ctx context.Context, site, mac string) error {
	s := c.site(site)
	_, err := c.postWithBody(ctx, fmt.Sprintf("/api/s/%s/cmd/devmgr", s), deviceCmdRequest{Cmd: "unset-locate", MAC: mac})
	if err != nil {
		return fmt.Errorf("UnlocateDevice %s: %w", mac, err)
	}
	return nil
}

// UpgradeDevice triggers a firmware upgrade on the device.
func (c *Client) UpgradeDevice(ctx context.Context, site, mac string) error {
	s := c.site(site)
	_, err := c.postWithBody(ctx, fmt.Sprintf("/api/s/%s/cmd/devmgr", s), deviceCmdRequest{Cmd: "upgrade", MAC: mac})
	if err != nil {
		return fmt.Errorf("UpgradeDevice %s: %w", mac, err)
	}
	return nil
}

// ForceReprovisionDevice forces reprovisioning of the device config.
func (c *Client) ForceReprovisionDevice(ctx context.Context, site, mac string) error {
	s := c.site(site)
	_, err := c.postWithBody(ctx, fmt.Sprintf("/api/s/%s/cmd/devmgr", s), deviceCmdRequest{Cmd: "force-reprovision", MAC: mac})
	if err != nil {
		return fmt.Errorf("ForceReprovisionDevice %s: %w", mac, err)
	}
	return nil
}

// RunSpeedTest initiates a speed test from the USG/UCG.
func (c *Client) RunSpeedTest(ctx context.Context, site string) error {
	s := c.site(site)
	_, err := c.post(ctx, fmt.Sprintf("/api/s/%s/cmd/devmgr/speedtest", s))
	if err != nil {
		return fmt.Errorf("RunSpeedTest %s: %w", s, err)
	}
	return nil
}

// GetSpeedTestStatus returns the most recent speed test result.
func (c *Client) GetSpeedTestStatus(ctx context.Context, site string) (SpeedTestStatus, error) {
	s := c.site(site)
	data, err := c.get(ctx, fmt.Sprintf("/api/s/%s/stat/speedtest-status", s))
	if err != nil {
		return SpeedTestStatus{}, fmt.Errorf("GetSpeedTestStatus %s: %w", s, err)
	}
	results, err := decodeLegacy[SpeedTestStatus](data)
	if err != nil {
		return SpeedTestStatus{}, fmt.Errorf("GetSpeedTestStatus %s: %w", s, err)
	}
	if len(results) == 0 {
		return SpeedTestStatus{}, fmt.Errorf("GetSpeedTestStatus %s: empty response", s)
	}
	return results[0], nil
}
