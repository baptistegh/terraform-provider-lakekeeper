package lakekeeper

import (
	"fmt"
	"net/http"
)

type (
	ProjectServiceInterface interface {
		ListProjects(options ...RequestOptionFunc) (*ListProjectsResponse, *http.Response, error)
		GetProject(id string, options ...RequestOptionFunc) (*Project, *http.Response, error)
		DeleteProject(id string, options ...RequestOptionFunc) (*http.Response, error)
		GetDefaultProject(options ...RequestOptionFunc) (*Project, *http.Response, error)
		CreateProject(opts *CreateProjectOptions, options ...RequestOptionFunc) (*Project, *http.Response, error)
		RenameProject(id string, opts *RenameProjectOptions, options ...RequestOptionFunc) (*http.Response, error)
	}

	// ProjectService handles communication with project endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project
	ProjectService struct {
		client *Client
	}
)

var _ ProjectServiceInterface = (*ProjectService)(nil)

// Project represents a lakekeeper project
type Project struct {
	ID   string `json:"project-id"`
	Name string `json:"project-name"`
}

// GetProject retrieves information about a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/get_default_project
func (s *ProjectService) GetProject(id string, options ...RequestOptionFunc) (*Project, *http.Response, error) {
	options = append(options, WithProject(id))
	req, err := s.client.NewRequest(http.MethodGet, "/project", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var prj Project

	resp, apiErr := s.client.Do(req, &prj)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &prj, resp, nil
}

// GetDefaultProject retrieves information about the user's default project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/get_default_project
//
// Deprecated: This endpoint is deprecated and will be removed in a future version.
func (s *ProjectService) GetDefaultProject(options ...RequestOptionFunc) (*Project, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/default-project", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var prj Project

	resp, apiErr := s.client.Do(req, &prj)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &prj, resp, nil
}

// ListProjectsResponse represents ListProjects() response.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/list_projects
type ListProjectsResponse struct {
	Projects []*Project `json:"projects"`
}

// RenameProjectOptions represents RenameProject() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/rename_project
type RenameProjectOptions struct {
	NewName string `json:"new-name"`
}

// ListProjects lists all projects that the requesting user has access to.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/list_projects
func (s *ProjectService) ListProjects(options ...RequestOptionFunc) (*ListProjectsResponse, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/project-list", nil, options)
	if err != nil {
		return nil, nil, err
	}

	var prjs ListProjectsResponse

	resp, apiErr := s.client.Do(req, &prjs)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	return &prjs, resp, nil
}

// CreateProjectOptions represents CreateProject() options.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
type CreateProjectOptions struct {
	ID   *string `json:"project-id,omitempty"` // Request a specific project ID - optional. If not provided, a new project ID will be generated (recommended)
	Name string  `json:"project-name"`
}

// createProjectResponse represents the response on project creation.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
type createProjectResponse struct {
	ID string `json:"project-id"`
}

// CreateProject creates a new project with the specified configuration.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/create_project
func (s *ProjectService) CreateProject(opts *CreateProjectOptions, options ...RequestOptionFunc) (*Project, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "/project", opts, options)
	if err != nil {
		return nil, nil, err
	}

	var prjResp createProjectResponse

	resp, apiErr := s.client.Do(req, &prjResp)
	if apiErr != nil {
		return nil, resp, apiErr
	}

	project, _, err := s.GetProject(prjResp.ID, options...)
	if err != nil {
		return nil, resp, fmt.Errorf("project is created but could not be read, %w", err)
	}

	return project, resp, nil
}

// RenameProject renames a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/rename_project
func (s *ProjectService) RenameProject(id string, opts *RenameProjectOptions, options ...RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(id))

	req, err := s.client.NewRequest(http.MethodPost, "/project/rename", opts, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}

// DeleteProject delete a project.
//
// Lakekeeper API docs:
// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/delete_default_project
func (s *ProjectService) DeleteProject(id string, options ...RequestOptionFunc) (*http.Response, error) {
	options = append(options, WithProject(id))

	req, err := s.client.NewRequest(http.MethodDelete, "/project", nil, options)
	if err != nil {
		return nil, err
	}

	resp, apiErr := s.client.Do(req, nil)
	if apiErr != nil {
		return resp, apiErr
	}

	return resp, nil
}
