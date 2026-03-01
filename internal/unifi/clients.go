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

// GetClient returns a single connected client from GET /integration/v1/sites/{siteID}/clients/{clientID}.
// Pass an empty siteID to use the client default.
func (c *Client) GetClient(ctx context.Context, siteID, clientID string) (NetworkClient, error) {
	id := c.site(siteID)
	data, err := c.get(ctx, fmt.Sprintf("/integration/v1/sites/%s/clients/%s", id, clientID))
	if err != nil {
		return NetworkClient{}, fmt.Errorf("GetClient %s %s: %w", id, clientID, err)
	}
	client, err := decodeV1[NetworkClient](data)
	if err != nil {
		return NetworkClient{}, fmt.Errorf("GetClient %s %s: %w", id, clientID, err)
	}
	return client, nil
}

// AuthorizeGuestClient posts an AUTHORIZE_GUEST_ACCESS action for a connected client via
// POST /integration/v1/sites/{siteID}/clients/{clientID}/actions.
// Pass an empty siteID to use the client default.
func (c *Client) AuthorizeGuestClient(ctx context.Context, siteID, clientID string, req GuestAuthRequest) error {
	id := c.site(siteID)
	_, err := c.postWithBody(ctx, fmt.Sprintf("/integration/v1/sites/%s/clients/%s/actions", id, clientID), req)
	if err != nil {
		return fmt.Errorf("AuthorizeGuestClient %s %s: %w", id, clientID, err)
	}
	return nil
}
