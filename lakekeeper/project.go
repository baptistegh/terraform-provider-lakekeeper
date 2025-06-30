package lakekeeper

import (
	"context"
	"encoding/json"
	"fmt"
)

type Project struct {
	ID   string `json:"project-id"`
	Name string `json:"project-name"`
}

type Projects struct {
	Projects []Project `json:"projects"`
}

type ProjectCreateRequest struct {
	Name string `json:"project-name"`
}

type ProjectCreateResponse struct {
	ID string `json:"project-id"`
}

func (client *Client) ListProjects(ctx context.Context) (*Projects, error) {
	var projects Projects
	err := client.get(ctx, "/management/v1/project-list", &projects, nil)
	if err != nil {
		return nil, err
	}

	return &projects, nil
}

func (client *Client) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	var project Project
	path := fmt.Sprintf("/management/v1/project/%s", id)
	err := client.get(ctx, path, &project, nil)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (client *Client) GetProjectByName(ctx context.Context, name string) (*Project, error) {
	projects, err := client.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range projects.Projects {
		if p.Name == name {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("could not find project with name: %s", name)
}

func (client *Client) GetDefaultProject(ctx context.Context) (*Project, error) {
	var project Project
	err := client.get(ctx, "/management/v1/project", &project, nil)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (client *Client) NewProject(ctx context.Context, name string) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("project name can't be empty")
	}
	body, err := json.Marshal(ProjectCreateRequest{Name: name})
	if err != nil {
		return nil, fmt.Errorf("could not marshall project creation request, %s", err.Error())
	}

	resp, err := client.post(ctx, "/management/v1/project", body)
	if err != nil {
		return nil, err
	}

	var r ProjectCreateResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, fmt.Errorf("could not unmarshall project creation response, %s", err.Error())
	}

	project := &Project{
		ID:   r.ID,
		Name: name,
	}

	return project, nil
}

func (client *Client) DeleteProject(ctx context.Context, id string) error {
	if err := client.delete(ctx, fmt.Sprintf("/management/v1/project/%s", id)); err != nil {
		return fmt.Errorf("could not delete project with id %s, %s", id, err.Error())
	}
	return nil
}
