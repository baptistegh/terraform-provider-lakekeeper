//go:build acceptance || flakey || settings
// +build acceptance flakey settings

package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

var testLakekeeperConfig = lakekeeper.Config{
	BaseURL: os.Getenv("LAKEKEEPER_ENDPOINT"),
	ClientCredentials: &lakekeeper.ClientCredentials{
		AuthURL:      os.Getenv("LAKEKEEPER_AUTH_URL"),
		ClientID:     os.Getenv("LAKEKEEPER_CLIENT_ID"),
		ClientSecret: os.Getenv("LAKEKEEPER_CLIENT_SECRET"),
		Scope:        "lakekeeper",
	},
	InitialBootstrap:      true,
	HandleTokenExpiration: true,
}

var TestLakekeeperClient *lakekeeper.Client

func init() {
	client, err := lakekeeper.NewClient(context.Background(), &testLakekeeperConfig)
	if err != nil {
		panic("failed to create test client: " + err.Error())
	}

	TestLakekeeperClient = client
}

// CreateProject is a test helper for creating a project.
func CreateProject(t *testing.T) *lakekeeper.Project {
	t.Helper()

	project, err := TestLakekeeperClient.NewProject(context.Background(), acctest.RandomWithPrefix("acctest"))
	if err != nil {
		t.Fatalf("could not create test project: %v", err)
	}

	t.Cleanup(func() {
		if err := TestLakekeeperClient.DeleteProject(context.Background(), project.ID); err != nil {
			t.Fatalf("could not cleanup test project: %v", err)
		}
	})

	return project
}

// CreateRole is a test helper for creating a role.
func CreateRole(t *testing.T, projectID string) *lakekeeper.Role {
	t.Helper()

	request := lakekeeper.RoleCreateRequest{
		Name:        acctest.RandString(8),
		Description: acctest.RandString(32),
		ProjectID:   projectID,
	}
	role, err := TestLakekeeperClient.NewRole(context.Background(), &request)
	if err != nil {
		t.Fatalf("could not create test role: %v", err)
	}

	t.Cleanup(func() {
		if err := TestLakekeeperClient.DeteleteRoleByID(context.Background(), role.ID, role.ProjectID); err != nil {
			t.Fatalf("could not cleanup test role: %v", err)
		}
	})

	return role
}

// CreateUser is a test helper for creating a user.
func CreateUser(t *testing.T, id string) *lakekeeper.User {
	t.Helper()

	name := acctest.RandomWithPrefix("acctest")
	user, err := TestLakekeeperClient.NewUser(context.Background(), id, fmt.Sprintf("%s@test.com", name), name, "human", false)
	if err != nil {
		t.Fatalf("could not create test user: %v", err)
	}

	t.Cleanup(func() {
		if err := TestLakekeeperClient.DeleteUser(context.Background(), id); err != nil {
			t.Fatalf("could not cleanup test user: %v", err)
		}
	})

	return user
}
