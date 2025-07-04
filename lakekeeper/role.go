package lakekeeper

import (
	"context"
	"encoding/json"
	"fmt"
)

type Role struct {
	ID          string  `json:"id"`
	ProjectID   string  `json:"project-id"`
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

type RoleSearchResponse struct {
	NextPageToken string  `json:"next-page-token"`
	Roles         []*Role `json:"roles"`
}

func (r Role) String() string {
	return fmt.Sprintf("{id:%s,project_id:%s,name:%s}", r.ID, r.ProjectID, r.Name)
}

func (client *Client) GetRoleByID(ctx context.Context, roleID string, projectID string) (*Role, *ApiError) {
	if roleID == "" {
		return nil, ApiErrorFromError("role id can not be empty")
	}
	var role Role

	if err := client.getWithProjectID(ctx, "/management/v1/role/"+roleID, projectID, &role, nil); err != nil {
		return nil, err
	}

	// populate project id if it is not in the response (api deprecated field)
	if projectID == "" {
		role.ProjectID = client.defaultProjectID
	} else {
		role.ProjectID = projectID
	}

	return &role, nil
}

func (client *Client) GetRoleByName(ctx context.Context, name, projectID string) (*Role, *ApiError) {
	if name == "" {
		return nil, ApiErrorFromError("could not find role with empty name")
	}

	var roles RoleSearchResponse
	if respErr := client.getWithProjectID(ctx, "/management/v1/role", projectID, &roles, map[string]string{"name": name}); respErr != nil {
		return nil, respErr
	}

	for _, role := range roles.Roles {
		if role.Name == name {
			// populate project id if it is not in the response (api deprecated field)
			if projectID == "" {
				role.ProjectID = client.defaultProjectID
			} else {
				role.ProjectID = projectID
			}
			return role, nil
		}
	}

	return nil, ApiErrorFromError("did not find role with name %s in project %s", name, projectID)
}

func (client *Client) NewRole(ctx context.Context, request *RoleCreateRequest) (*Role, *ApiError) {
	if request.Name == "" {
		return nil, ApiErrorFromError("role name must be defined")
	}

	evaluatedProjectID := client.defaultProjectID
	if request.ProjectID != "" {
		evaluatedProjectID = request.ProjectID
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, ApiErrorFromError("could not marshal role creation request, %v", err)
	}

	resp, apiErr := client.postWithProjectID(ctx, "/management/v1/role", evaluatedProjectID, body)
	if apiErr != nil {
		return nil, apiErr
	}

	var role Role
	if err := json.Unmarshal(resp, &role); err != nil {
		return nil, ApiErrorFromError("role %s is created but the create response could not be decoded, %v", request.Name, err)
	}

	if role.ProjectID == "" {
		role.ProjectID = evaluatedProjectID
	}

	return &role, nil
}

func (client *Client) DeteleteRoleByID(ctx context.Context, roleID, projectID string) *ApiError {
	if roleID == "" {
		return ApiErrorFromError("could not delete role with empty ID")
	}

	if err := client.deleteWithProjectID(ctx, "/management/v1/role/"+roleID, projectID); err != nil {
		return err
	}

	return nil
}

func (client *Client) UpdateRole(ctx context.Context, request *RoleUpdateRequest) (*Role, *ApiError) {
	if request.ID == "" {
		return nil, ApiErrorFromError("role id must be defined")
	}

	if request.Name == "" {
		return nil, ApiErrorFromError("role name must be defined")
	}

	evaluatedProjectID := client.defaultProjectID
	if request.ProjectID != "" {
		evaluatedProjectID = request.ProjectID
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, ApiErrorFromError("could not marshal role update request, %v", err)
	}

	resp, apiErr := client.postWithProjectID(ctx, "/management/v1/role/"+request.ID, evaluatedProjectID, body)
	if apiErr != nil {
		return nil, apiErr
	}

	var role Role
	if err := json.Unmarshal(resp, &role); err != nil {
		return nil, ApiErrorFromError("role %s is updated but the update response could not be decoded, %v", request.Name, err)
	}

	if role.ProjectID == "" {
		role.ProjectID = evaluatedProjectID
	}

	return &role, nil
}
