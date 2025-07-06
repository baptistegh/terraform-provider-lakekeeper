package lakekeeper

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/storage"
)

type (
	WarehouseServiceInterface interface {
		GetWarehouse(id, projectID string, options ...RequestOptionFunc) (*Warehouse, *http.Response, error)
		ListWarehouses(opts *ListWarehousesOptions, options ...RequestOptionFunc) ([]*Warehouse, *http.Response, error)
		CreateWarehouse(opts *CreateWarehouseOptions, options ...RequestOptionFunc) (*Warehouse, *http.Response, error)
		DeleteWarehouse(id string, opts *DeleteWarehouseOptions, options ...RequestOptionFunc) (*http.Response, error)
	}

	// WarehouseService handles communication with warehouse endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse
	WarehouseService struct {
		client *Client
	}
)

var _ WarehouseServiceInterface = (*WarehouseService)(nil)

// Warehouse represents a lakekeeper warehouse
type Warehouse struct {
	ID                    string                         `json:"id"`
	ProjectID             string                         `json:"project-id"`
	Name                  string                         `json:"name"`
	Protected             bool                           `json:"protected"`
	Status                WarehouseStatus                `json:"status"`
	StorageProfileWrapper *storage.StorageProfileWrapper `json:"storage-profile"`
	DeleteProfileWrapper  *DeleteProfileWrapper          `json:"delete-profile"`
}

type WarehouseStatus string

const (
	WarehouseStatusActive   WarehouseStatus = "active"
	WarehouseStatusInactive WarehouseStatus = "inactive"
)

func (w *Warehouse) IsActive() bool {
	return w.Status == WarehouseStatusActive
}

// GetWarehouse retrieves detailed information about a specific warehouse.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse
func (s *WarehouseService) GetWarehouse(id, projectID string, options ...RequestOptionFunc) (*Warehouse, *http.Response, error) {
	if projectID != "" {
		options = append(options, WithProject(projectID))
	}

	req, err := s.client.NewRequest(http.MethodGet, "/warehouse/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var wh Warehouse

	resp, apiErr := s.client.Do(req, &wh)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &wh, resp, nil
}

// ListWarehousesOptions represents ListWarehouses() options
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_warehouses
type ListWarehousesOptions struct {
	WarehouseStatus *WarehouseStatus `url:"warehouseStatus,omitempty"`
	ProjectID       *string          `url:"projectId,omitempty"`
}

// listWarehouseResponse represents the response on list warehouses API action
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_warehouses
type listWarehousesResponse struct {
	Warehouses []*Warehouse `json:"warehouses"`
}

// Returns all warehouses in the project that the current user has access to.
// By default, deactivated warehouses are not included in the results.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_warehouses
func (s *WarehouseService) ListWarehouses(opts *ListWarehousesOptions, options ...RequestOptionFunc) ([]*Warehouse, *http.Response, error) {
	if opts != nil && opts.ProjectID != nil {
		options = append(options, WithProject(*opts.ProjectID))
	}

	req, err := s.client.NewRequest(http.MethodGet, "/warehouse", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var whs listWarehousesResponse

	resp, apiErr := s.client.Do(req, &whs)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return whs.Warehouses, resp, nil
}

// CreateWarehouseOptions represents CreateWarehouse() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/create_warehouse
type CreateWarehouseOptions struct {
	Name              string                           `json:"warehouse-name"`
	ProjectID         string                           `json:"project-id"`
	StorageProfile    storage.StorageProfileWrapper    `json:"storage-profile"`
	StorageCredential storage.StorageCredentialWrapper `json:"storage-credential"`
	DeleteProfile     DeleteProfile                    `json:"delete-profile,omitempty"`
}

// CreateWarehouseOptions represents the response from the API
// on a create_warehouse action.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/create_warehouse
type createWarehouseResponse struct {
	ID string `json:"warehouse-id"`
}

// CreateWarehouse creates a new warehouse in the specified project with
// the provided configuration.
// The project of a warehouse cannot be changed after creation.
// This operation validates the storage configuration.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/create_warehouse
func (s *WarehouseService) CreateWarehouse(opts *CreateWarehouseOptions, options ...RequestOptionFunc) (*Warehouse, *http.Response, error) {
	if opts == nil {
		return nil, nil, errors.New("CreateWarehouse received empty options")
	}

	if opts.ProjectID != "" {
		options = append(options, WithProject(opts.ProjectID))
	}

	req, err := s.client.NewRequest(http.MethodPost, "/warehouse", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var whResp createWarehouseResponse

	resp, apiErr := s.client.Do(req, &whResp)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	warehouse, _, err := s.GetWarehouse(whResp.ID, opts.ProjectID)
	if err != nil {
		return nil, resp, fmt.Errorf("warehouse is created but error occured on get, %w", err)
	}

	return warehouse, resp, nil
}

// DeleteWarehouseOptions represents DeleteWarehouse() options.
//
// force parameters needs to be true to delete protected warehouses.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/delete_warehouse
type DeleteWarehouseOptions struct {
	Force     bool    `url:"force"`
	ProjectID *string `url:"-"`
}

// DeleteWarehouse permanently removes a warehouse and all its associated resources.
// Use the force parameter to delete protected warehouses.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/delete_warehouse
func (s *WarehouseService) DeleteWarehouse(id string, opts *DeleteWarehouseOptions, options ...RequestOptionFunc) (*http.Response, error) {
	if opts != nil && opts.ProjectID != nil {
		options = append(options, WithProject(*opts.ProjectID))
	}

	req, err := s.client.NewRequest(http.MethodDelete, "/warehouse/"+id, opts, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
