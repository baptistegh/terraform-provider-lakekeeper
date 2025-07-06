package lakekeeper

import (
	"errors"
	"net/http"
)

type (
	RoleServiceInterface interface {
		ListRoles(opts *ListRolesOptions, options ...RequestOptionFunc) ([]*Role, error)
		GetRole(id string, projectID string, options ...RequestOptionFunc) (*Role, *http.Response, error)
		CreateRole(opts *CreateRoleOptions, options ...RequestOptionFunc) (*Role, *http.Response, error)
		UpdateRole(id string, opts *UpdateRoleOptions, options ...RequestOptionFunc) (*Role, *http.Response, error)
		DeleteRole(id, projectID string, options ...RequestOptionFunc) (*http.Response, error)
	}

	// RoleService handles communication with role endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs: https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role
	RoleService struct {
		client *Client
	}
)

var _ RoleServiceInterface = (*RoleService)(nil)

// Project represents a lakekeeper role
type Role struct {
	ID          string  `json:"id"`
	ProjectID   string  `json:"project-id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	CreatedAt string  `json:"created-at"`
	UpdatedAt *string `json:"updated-at,omitempty"`
}

// GetRole retrieves information about a role.
//
// Lakekeeper API docs: https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/get_role
func (s *RoleService) GetRole(id string, projectID string, options ...RequestOptionFunc) (*Role, *http.Response, error) {
	if projectID != "" {
		options = append(options, WithProject(id))
	}

	req, err := s.client.NewRequest(http.MethodGet, "/role/"+id, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var role Role

	resp, apiErr := s.client.Do(req, &role)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &role, resp, nil
}

// ListRolesOptions represents ListRoles() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
type ListRolesOptions struct {
	Name      *string `url:"name,omitempty"`
	PageToken *string `url:"pageToken,omitempty"`
	PageSize  *string `url:"pageSize,omitempty"`
	ProjectID *string `url:"projectId,omitempty"`
}

// listRoleResponse represents a response from list_roles API endpoint.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/list_roles
type listRolesResponse struct {
	NextPageToken *string `json:"next-page-token,omitempty"`
	Roles         []*Role `json:"role"`
}

// ListRoles returns all roles in the project that the current user has access to view.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/list_roles
func (s *RoleService) ListRoles(opts *ListRolesOptions, options ...RequestOptionFunc) ([]*Role, error) {
	if opts != nil && opts.ProjectID != nil {
		options = append(options, WithProject(*opts.ProjectID))
	}

	var roles []*Role

	for {
		var r listRolesResponse
		req, err := s.client.NewRequest(http.MethodGet, "/project-list", opts, options)
		if err != nil {
			return nil, err
		}

		_, apiErr := s.client.Do(req, &r)
		if apiErr != nil {
			return nil, apiErr
		}

		roles = append(roles, r.Roles...)

		if r.NextPageToken == nil {
			break
		}
		opts.PageToken = r.NextPageToken
	}

	return roles, nil
}

// CreateRoleOptions represents CreateRole() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/create_role
type CreateRoleOptions struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ProjectID   *string `json:"project-id"`
}

// CreateRole creates a role with the specified name and description.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/create_role
func (s *RoleService) CreateRole(opts *CreateRoleOptions, options ...RequestOptionFunc) (*Role, *http.Response, error) {
	if opts == nil {
		return nil, nil, errors.New("CreateRole needs options to create a role")
	}

	if opts.ProjectID != nil {
		options = append(options, WithProject(*opts.ProjectID))
	}

	req, err := s.client.NewRequest(http.MethodPost, "/role", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var role Role

	resp, apiErr := s.client.Do(req, &role)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &role, resp, nil
}

// CreateRoleOptions represents UpdateRole() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/update_role
type UpdateRoleOptions struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	ProjectID   *string `json:"-"`
}

// UpdateRole update a role with the specified name and description.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/update_role
func (s *RoleService) UpdateRole(id string, opts *UpdateRoleOptions, options ...RequestOptionFunc) (*Role, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("Role ID must be defined to be updated")
	}

	if opts.ProjectID != nil {
		options = append(options, WithProject(*opts.ProjectID))
	}

	req, err := s.client.NewRequest(http.MethodPost, "/role/"+id, opts, options)
	if err != nil {
		return nil, nil, err
	}

	var role Role

	resp, apiErr := s.client.Do(req, &role)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &role, resp, nil
}

// DeleteRole permanently removes a role and all its associated permissions.
//
// Lakekeeper API docs: https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/delete_role
func (s *RoleService) DeleteRole(id string, projectID string, options ...RequestOptionFunc) (*http.Response, error) {
	if projectID != "" {
		options = append(options, WithProject(id))
	}

	req, err := s.client.NewRequest(http.MethodDelete, "/role/"+id, nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
