package unifi

import (
	"context"
	"fmt"
)

// ListActiveClients returns currently connected clients from GET /v1/sites/{siteID}/clients/active.
// Pass an empty siteID to use the client default.
func (c *Client) ListActiveClients(ctx context.Context, siteID string) ([]ActiveClient, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/v1/sites/%s/clients/active", id))
	if err != nil {
		return nil, fmt.Errorf("ListActiveClients %s: %w", id, err)
	}
	clients, err := decodeV1List[ActiveClient](data)
	if err != nil {
		return nil, fmt.Errorf("ListActiveClients %s: %w", id, err)
	}
	return clients, nil
}

// ListKnownClients returns all known/historical clients from GET /v1/sites/{siteID}/clients/history.
// Pass an empty siteID to use the client default.
func (c *Client) ListKnownClients(ctx context.Context, siteID string) ([]KnownClient, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/v1/sites/%s/clients/history", id))
	if err != nil {
		return nil, fmt.Errorf("ListKnownClients %s: %w", id, err)
	}
	clients, err := decodeV1List[KnownClient](data)
	if err != nil {
		return nil, fmt.Errorf("ListKnownClients %s: %w", id, err)
	}
	return clients, nil
}

// BlockClient blocks the client with the given MAC address.
func (c *Client) BlockClient(ctx context.Context, site, mac string) error {
	s := c.site(site)
	_, err := c.postWithBody(ctx, fmt.Sprintf("/api/s/%s/cmd/stamgr", s), clientCmdRequest{Cmd: "block-sta", MAC: mac})
	if err != nil {
		return fmt.Errorf("BlockClient %s: %w", mac, err)
	}
	return nil
}

// UnblockClient removes the block on the client with the given MAC address.
func (c *Client) UnblockClient(ctx context.Context, site, mac string) error {
	s := c.site(site)
	_, err := c.postWithBody(ctx, fmt.Sprintf("/api/s/%s/cmd/stamgr", s), clientCmdRequest{Cmd: "unblock-sta", MAC: mac})
	if err != nil {
		return fmt.Errorf("UnblockClient %s: %w", mac, err)
	}
	return nil
}

// KickClient disconnects (but does not ban) the client with the given MAC address.
func (c *Client) KickClient(ctx context.Context, site, mac string) error {
	s := c.site(site)
	_, err := c.postWithBody(ctx, fmt.Sprintf("/api/s/%s/cmd/stamgr", s), clientCmdRequest{Cmd: "kick-sta", MAC: mac})
	if err != nil {
		return fmt.Errorf("KickClient %s: %w", mac, err)
	}
	return nil
}

// ForgetClient permanently removes the client record from the controller.
func (c *Client) ForgetClient(ctx context.Context, site, mac string) error {
	s := c.site(site)
	_, err := c.postWithBody(ctx, fmt.Sprintf("/api/s/%s/cmd/stamgr", s), clientCmdRequest{Cmd: "forget-sta", MAC: mac})
	if err != nil {
		return fmt.Errorf("ForgetClient %s: %w", mac, err)
	}
	return nil
}
