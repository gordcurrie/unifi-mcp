package unifi

import (
	"context"
	"fmt"
)

// ListClients returns all currently connected clients from GET /integration/v1/sites/{siteID}/clients.
// Pass an empty siteID to use the client default.
func (c *Client) ListClients(ctx context.Context, siteID string) ([]NetworkClient, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/clients", id))
	if err != nil {
		return nil, fmt.Errorf("ListClients %s: %w", id, err)
	}
	clients, err := decodeV1List[NetworkClient](data)
	if err != nil {
		return nil, fmt.Errorf("ListClients %s: %w", id, err)
	}
	return clients, nil
}
