package lakekeeper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/storage"
	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/storage/credential"
)

type Warehouse struct {
	ID                       string                               `json:"id"`
	ProjectID                string                               `json:"-"`
	Name                     string                               `json:"name"`
	Protected                bool                                 `json:"protected"`
	Status                   string                               `json:"status"`
	StorageProfileWrapper    *storage.StorageProfileWrapper       `json:"storage-profile"`
	DeleteProfileWrapper     *DeleteProfileWrapper                `json:"delete-profile"`
	StorageCredentialWrapper *credential.StorageCredentialWrapper `json:"storage-credential"`
}

type WarehouseCreateOptions struct {
	Name              string
	ProjectID         string // if empty, client default project id is used
	Protected         bool
	Status            string
	StorageProfile    storage.StorageProfile
	StorageCredential credential.StorageCredential
	DeleteProfile     DeleteProfile
}

type warehouseCreateRequest struct {
	Name                     string                               `json:"warehouse-name"`
	StorageProfileWrapper    *storage.StorageProfileWrapper       `json:"storage-profile"`
	StorageCredentialWrapper *credential.StorageCredentialWrapper `json:"storage-credential"`
	DeleteProfileWrapper     *DeleteProfileWrapper                `json:"delete-profile,omitempty"`
}

type WarehouseCreateResponse struct {
	ID string `json:"warehouse-id"`
}

func (w *Warehouse) IsActive() bool {
	return w.Status == "active"
}

// GetWarehouseByID retrieves a warehouse by its ID.
// If projectID is empty, the client's default project ID is used.
func (client *Client) GetWarehouseByID(ctx context.Context, projectID string, warehouseID string) (*Warehouse, error) {
	if warehouseID == "" {
		return nil, fmt.Errorf("warehouse ID cannot be empty")
	}

	var warehouse Warehouse
	if err := client.getWithProjectID(ctx, "/management/v1/warehouse/"+warehouseID, projectID, &warehouse, nil); err != nil {
		return nil, err
	}

	// populate project id if it is not in the response (api deprecated field)
	if warehouse.ProjectID == "" {
		if projectID == "" {
			warehouse.ProjectID = client.defaultProjectID
		}
		warehouse.ProjectID = projectID
	}

	return &warehouse, nil
}

// DeleteWarehouseByID deletes a warehouse by its ID.
// If projectID is empty, the client's default project ID is used.
func (client *Client) DeleteWarehouseByID(ctx context.Context, projectID, warehouseID string) error {
	if warehouseID == "" {
		return fmt.Errorf("warehouse ID cannot be empty")
	}

	if err := client.deleteWithProjectID(ctx, "/management/v1/warehouse/"+warehouseID, projectID); err != nil {
		return err
	}

	return nil
}

func (client *Client) NewWarehouse(ctx context.Context, r *WarehouseCreateOptions) (*Warehouse, error) {
	if r.Name == "" {
		return nil, fmt.Errorf("could not create warehouse with an empty name")
	}

	if r.StorageProfile == nil {
		return nil, fmt.Errorf("storage profile must be defined")
	}

	if r.StorageCredential == nil {
		return nil, fmt.Errorf("storage credential must be defined")
	}

	var w = warehouseCreateRequest{
		Name:                     r.Name,
		StorageProfileWrapper:    &storage.StorageProfileWrapper{StorageProfile: r.StorageProfile},
		StorageCredentialWrapper: &credential.StorageCredentialWrapper{StorageCredential: r.StorageCredential},
	}

	//b, _ := json.Marshal(w)
	//panic(string(b))

	if r.DeleteProfile != nil {
		w.DeleteProfileWrapper = &DeleteProfileWrapper{DeleteProfile: r.DeleteProfile}
	}

	evaluatedProjectID := client.defaultProjectID
	if r.ProjectID != "" {
		evaluatedProjectID = r.ProjectID
	}

	body, err := json.Marshal(w)
	if err != nil {
		return nil, fmt.Errorf("could not marshal warehouse creation request, %s", err.Error())
	}

	resp, apiErr := client.postWithProjectID(ctx, "/management/v1/warehouse", evaluatedProjectID, body)
	if apiErr != nil {
		return nil, apiErr
	}
	var wResp WarehouseCreateResponse
	if err := json.Unmarshal(resp, &wResp); err != nil {
		return nil, fmt.Errorf("warehouse %s is created but could not find its ID, %s", r.Name, err.Error())
	}

	return client.GetWarehouseByID(ctx, evaluatedProjectID, wResp.ID)
}
