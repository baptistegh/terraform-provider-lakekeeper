package lakekeeper

import (
	"context"
	"encoding/json"
	"fmt"
)

type Role struct {
	ID          string  `json:"id"`
	ProjectID   string  `json:"-"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	CreatedAt string  `json:"created-at"`
	UpdatedAt *string `json:"updated-at,omitempty"`
}

type RoleCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	ProjectID string `json:"-"`
}

type RoleUpdateRequest struct {
	ID        string `json:"-"`
	ProjectID string `json:"-"`

	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (client *Client) GetRoleByID(ctx context.Context, roleID string, projectID string) (*Role, error) {
	if roleID == "" {
		return nil, fmt.Errorf("role id can not be empty")
	}
	var role Role

	if err := client.getWithProjectID(ctx, "/management/v1/role/"+roleID, projectID, &role, nil); err != nil {
		return nil, err
	}

	// populate project id if it is not in the response (api deprecated field)
	if role.ProjectID == "" {
		if projectID == "" {
			role.ProjectID = client.defaultProjectID
		} else {
			role.ProjectID = projectID
		}
	}

	return &role, nil
}

func (client *Client) NewRole(ctx context.Context, request *RoleCreateRequest) (*Role, error) {
	if request.Name == "" {
		return nil, fmt.Errorf("role name must be defined")
	}

	evaluatedProjectID := client.defaultProjectID
	if request.ProjectID != "" {
		evaluatedProjectID = request.ProjectID
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("could not marshal role creation request, %v", err)
	}

	resp, err := client.postWithProjectID(ctx, "/management/v1/role", evaluatedProjectID, body)
	if err != nil {
		return nil, err
	}

	var role Role
	if err := json.Unmarshal(resp, &role); err != nil {
		return nil, fmt.Errorf("role %s is created but the create response could not be decoded, %v", request.Name, err)
	}

	if role.ProjectID == "" {
		role.ProjectID = evaluatedProjectID
	}

	return &role, nil
}

func (client *Client) DeteleteRoleByID(ctx context.Context, roleID, projectID string) error {
	if roleID == "" {
		return fmt.Errorf("could not delete role with empty ID")
	}

	if err := client.deleteWithProjectID(ctx, "/management/v1/role/"+roleID, projectID); err != nil {
		return err
	}

	return nil
}

func (client *Client) UpdateRole(ctx context.Context, request *RoleUpdateRequest) (*Role, error) {
	if request.ID == "" {
		return nil, fmt.Errorf("role id must be defined")
	}

	if request.Name == "" {
		return nil, fmt.Errorf("role name must be defined")
	}

	evaluatedProjectID := client.defaultProjectID
	if request.ProjectID != "" {
		evaluatedProjectID = request.ProjectID
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("could not marshal role update request, %v", err)
	}

	resp, err := client.postWithProjectID(ctx, "/management/v1/role/"+request.ID, evaluatedProjectID, body)
	if err != nil {
		return nil, err
	}

	var role Role
	if err := json.Unmarshal(resp, &role); err != nil {
		return nil, fmt.Errorf("role %s is updated but the update response could not be decoded, %v", request.Name, err)
	}

	if role.ProjectID == "" {
		role.ProjectID = evaluatedProjectID
	}

	return &role, nil
}
