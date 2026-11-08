package v1

import inventoryV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"

type Client struct {
	invClient inventoryV1.InventoryServiceClient
}

func NewClient(invClient inventoryV1.InventoryServiceClient) *Client {
	return &Client{
		invClient: invClient,
	}
}
