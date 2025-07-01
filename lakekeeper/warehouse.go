package lakekeeper

import (
	"context"
	"fmt"
)

type Warehouse struct {
	ID                    string                `json:"id"`
	Name                  string                `json:"name"`
	ProjectID             string                `json:"project-id"`
	Protected             bool                  `json:"protected"`
	Status                string                `json:"status"`
	StorageProfileWrapper StorageProfileWrapper `json:"storage-profile"`
	DeleteProfile         DeleteProfile         `json:"delete-profile"`
	CreatedAt             string                `json:"created_at"`
	UpdatedAt             string                `json:"updated_at"`
}

type DeleteProfile struct {
	Type string `json:"type"`
}

func (w *Warehouse) IsActive() bool {
	return w.Status == "active"
}

func (client *Client) GetWarehouseByID(ctx context.Context, warehouseID string, projectID string) (*Warehouse, error) {
	if warehouseID == "" {
		return nil, fmt.Errorf("warehouse ID cannot be empty")
	}

	var warehouse Warehouse

	if err := client.getWithProjectID(ctx, "/management/v1/warehouse/"+warehouseID, projectID, &warehouse, nil); err != nil {
		return nil, err
	}

	return &warehouse, nil
}
