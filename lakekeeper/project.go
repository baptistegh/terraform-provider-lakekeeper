package lakekeeper

import (
	"context"
	"fmt"
)

type Project struct {
	ID   string `json:"project-id"`
	Name string `json:"project-name"`
}

type Projects struct {
	Projects []Project `json:"projects"`
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
