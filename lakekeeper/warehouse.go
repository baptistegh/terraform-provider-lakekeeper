package lakekeeper

import (
	"context"
	"encoding/json"
	"fmt"
)

type Warehouse struct {
	ID                    string                `json:"id"`
	Name                  string                `json:"name"`
	Protected             bool                  `json:"protected"`
	Status                string                `json:"status"`
	StorageProfileWrapper storageProfileWrapper `json:"storage-profile"`
	DeleteProfile         deleteProfileWrapper  `json:"delete-profile"`
	CreatedAt             string                `json:"created_at"`
	UpdatedAt             string                `json:"updated_at"`
}

type WarehouseCreateRequest struct {
	Name                  string                 `json:"warehouse-name"`
	StorageProfileWrapper *storageProfileWrapper `json:"storage-profile"`
	DeleteProfileWrapper  *deleteProfileWrapper  `json:"delete-profile,omitempty"`
}

type WarehouseCreateResponse struct {
	ID string `json:"warehouse-id"`
}

func (w *Warehouse) IsActive() bool {
	return w.Status == "active"
}

// GetWarehouseByID retrieves a warehouse by its ID.
// If projectID is empty, the client's default project ID is used.
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

// DeleteWarehouseByID deletes a warehouse by its ID.
// If projectID is empty, the client's default project ID is used.
func (client *Client) DeleteWarehouseByID(ctx context.Context, warehouseID string, projectID string) error {
	if warehouseID == "" {
		return fmt.Errorf("warehouse ID cannot be empty")
	}

	if err := client.deleteWithProjectID(ctx, "/management/v1/warehouse/"+warehouseID, projectID); err != nil {
		return err
	}

	return nil
}

func (client *Client) NewWarehouse(ctx context.Context, projectId, name string, protected bool, status string, storageProfile StorageProfile, deleteProfile DeleteProfile) (*Warehouse, error) {
	if name == "" {
		return nil, fmt.Errorf("could not create warehouse with an empty name")
	}

	// var evaluatedProjectID = client.DefaultProjectID
	// if projectId == "" {
	// 	evaluatedProjectID = projectId
	// }

	var w = WarehouseCreateRequest{
		Name:                  name,
		StorageProfileWrapper: &storageProfileWrapper{StorageProfile: storageProfile},
	}

	if deleteProfile != nil {
		w.DeleteProfileWrapper = &deleteProfileWrapper{DeleteProfile: deleteProfile}
	}

	evaluatedProjectID := client.defaultProjectID
	if projectId != "" {
		evaluatedProjectID = projectId
	}

	body, err := json.Marshal(w)
	if err != nil {
		return nil, fmt.Errorf("could not marshal warehouse creation request, %s", err.Error())
	}

	resp, err := client.postWithProjectID(ctx, "/management/v1/warehouse", evaluatedProjectID, body)
	if err != nil {
		return nil, err
	}
	var wResp WarehouseCreateResponse
	if err := json.Unmarshal(resp, &wResp); err != nil {
		return nil, fmt.Errorf("warehouse %s is created but could not find its ID, %s", name, err.Error())
	}

	return client.GetWarehouseByID(ctx, wResp.ID, evaluatedProjectID)
}
