package lakekeeper

import (
	"context"
	"encoding/json"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/storage"
	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/storage/credential"
)

type WarehouseNoCreds struct {
	ID                    string                         `json:"id"`
	ProjectID             string                         `json:"-"` // deprecated field from the API
	Name                  string                         `json:"name"`
	Protected             bool                           `json:"protected"`
	Status                string                         `json:"status"`
	StorageProfileWrapper *storage.StorageProfileWrapper `json:"storage-profile"`
	DeleteProfileWrapper  *DeleteProfileWrapper          `json:"delete-profile"`
}

type Warehouse struct {
	WarehouseNoCreds
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

type warehouseListResponse struct {
	Warehouses []*WarehouseNoCreds `json:"warehouses"`
}

func (w *WarehouseNoCreds) IsActive() bool {
	return w.Status == "active"
}

// GetWarehouseByID retrieves a warehouse by its ID.
// If projectID is empty, the client's default project ID is used.
func (client *Client) GetWarehouseByID(ctx context.Context, projectID string, warehouseID string) (*Warehouse, *ApiError) {
	if warehouseID == "" {
		return nil, ApiErrorFromError("warehouse ID cannot be empty")
	}

	var warehouse Warehouse
	if err := client.getWithProjectID(ctx, "/management/v1/warehouse/"+warehouseID, projectID, &warehouse, nil); err != nil {
		return nil, err
	}

	// populate project id if it is not in the response (api deprecated field)
	if warehouse.ProjectID == "" {
		if projectID != "" {
			warehouse.ProjectID = projectID
		} else {
			warehouse.ProjectID = client.defaultProjectID
		}
	}

	return &warehouse, nil
}

func (client *Client) GetWarehouseByName(ctx context.Context, projectID, name string) (*WarehouseNoCreds, *ApiError) {
	if name == "" {
		return nil, ApiErrorFromError("warehouse name must be defined")
	}

	evaluatedProjectID := client.defaultProjectID
	if projectID != "" {
		evaluatedProjectID = projectID
	}

	var resp warehouseListResponse
	params := map[string]string{
		"projectId": evaluatedProjectID,
	}
	if err := client.getWithProjectID(ctx, "/management/v1/warehouse", projectID, &resp, params); err != nil {
		return nil, err
	}

	for _, warehouse := range resp.Warehouses {
		if warehouse.Name == name {
			warehouse.ProjectID = evaluatedProjectID
			return warehouse, nil
		}
	}
	return nil, ApiErrorFromError("could not find warehouse %s in project %s", name, evaluatedProjectID)
}

// DeleteWarehouseByID deletes a warehouse by its ID.
// If projectID is empty, the client's default project ID is used.
func (client *Client) DeleteWarehouseByID(ctx context.Context, projectID, warehouseID string) *ApiError {
	if warehouseID == "" {
		return ApiErrorFromError("warehouse ID cannot be empty")
	}

	if err := client.deleteWithProjectID(ctx, "/management/v1/warehouse/"+warehouseID, projectID); err != nil {
		return err
	}

	return nil
}

func (client *Client) NewWarehouse(ctx context.Context, r *WarehouseCreateOptions) (*Warehouse, *ApiError) {
	if r.Name == "" {
		return nil, ApiErrorFromError("could not create warehouse with an empty name")
	}

	if r.StorageProfile == nil {
		return nil, ApiErrorFromError("storage profile must be defined")
	}

	if r.StorageCredential == nil {
		return nil, ApiErrorFromError("storage credential must be defined")
	}

	var w = warehouseCreateRequest{
		Name:                     r.Name,
		StorageProfileWrapper:    &storage.StorageProfileWrapper{StorageProfile: r.StorageProfile},
		StorageCredentialWrapper: &credential.StorageCredentialWrapper{StorageCredential: r.StorageCredential},
	}

	if r.DeleteProfile != nil {
		w.DeleteProfileWrapper = &DeleteProfileWrapper{DeleteProfile: r.DeleteProfile}
	}

	evaluatedProjectID := client.defaultProjectID
	if r.ProjectID != "" {
		evaluatedProjectID = r.ProjectID
	}

	body, err := json.Marshal(w)
	if err != nil {
		return nil, ApiErrorFromError("could not marshal warehouse creation request, %v", err)
	}

	resp, apiErr := client.postWithProjectID(ctx, "/management/v1/warehouse", evaluatedProjectID, body)
	if apiErr != nil {
		return nil, apiErr
	}
	var wResp WarehouseCreateResponse
	if err := json.Unmarshal(resp, &wResp); err != nil {
		return nil, ApiErrorFromError("warehouse %s is created but could not find its ID, %v", r.Name, err)
	}

	return client.GetWarehouseByID(ctx, evaluatedProjectID, wResp.ID)
}
