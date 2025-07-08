//go:build acceptance
// +build acceptance

package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/api"
	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper"
	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/storage"
	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/types/storage/profile"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

var testLakekeeperConfig = api.Config{
	BaseURL: os.Getenv("LAKEKEEPER_ENDPOINT"),
	OIDCClientConfig: api.OIDCClientConfig{
		AuthURL:      os.Getenv("LAKEKEEPER_AUTH_URL"),
		ClientID:     os.Getenv("LAKEKEEPER_CLIENT_ID"),
		ClientSecret: os.Getenv("LAKEKEEPER_CLIENT_SECRET"),
		Scopes:       []string{"lakekeeper"},
	},
	InitialBootstrap: true,
	EarlyAuthFail:    true,
}

var TestLakekeeperClient *lakekeeper.Client

func init() {
	client, err := testLakekeeperConfig.NewLakekeeperClient(context.Background())
	if err != nil {
		panic("failed to create test client: " + err.Error())
	}

	TestLakekeeperClient = client
}

// CreateProject is a test helper for creating a project.
func CreateProject(t *testing.T) *lakekeeper.Project {
	t.Helper()

	opts := lakekeeper.CreateProjectOptions{
		Name: acctest.RandomWithPrefix("acctest"),
	}

	resp, _, err := TestLakekeeperClient.Project.CreateProject(&opts)
	if err != nil {
		t.Fatalf("could not create test project: %v", err)
	}

	t.Cleanup(func() {
		if _, err := TestLakekeeperClient.Project.DeleteProject(resp.ID); err != nil {
			t.Fatalf("could not cleanup test project: %v", err)
		}
	})

	return resp
}

// CreateWarehouse is a test helper for creating a warehouse.
func CreateWarehouse(t *testing.T, projectID, keyPrefix string) *lakekeeper.Warehouse {
	t.Helper()

	profile, err := profile.NewS3StorageSettings("testacc", "local-01",
		profile.WithEndpoint("http://minio:9000/"),
		profile.WithPathStyleAccess(),
		profile.WithS3KeyPrefix(keyPrefix),
	)
	if err != nil {
		t.Fatalf("error creating storage profile, %v", err)
	}

	opts := lakekeeper.CreateWarehouseOptions{
		Name:              acctest.RandString(8),
		ProjectID:         projectID,
		StorageProfile:    *profile.AsProfile(),
		StorageCredential: storage.StorageCredentialWrapper{StorageCredential: storage.NewS3CredentialAccessKey("minio-root-user", "minio-root-password", "")},
		DeleteProfile:     &lakekeeper.HardDeleteProfile{Type: "hard"},
	}

	warehouse, _, err := TestLakekeeperClient.Warehouse.CreateWarehouse(&opts)
	if err != nil {
		t.Fatalf("could not create test warehouse: %v", err)
	}

	t.Cleanup(func() {
		opts := lakekeeper.DeleteWarehouseOptions{
			Force:     true,
			ProjectID: &projectID,
		}
		if _, err := TestLakekeeperClient.Warehouse.DeleteWarehouse(warehouse.ID, &opts); err != nil {
			t.Fatalf("could not cleanup test warehouse: %v", err)
		}
	})

	return warehouse
}

// CreateRole is a test helper for creating a role.
func CreateRole(t *testing.T, projectID string) *lakekeeper.Role {
	t.Helper()

	description := acctest.RandString(32)

	opts := lakekeeper.CreateRoleOptions{
		Name:        acctest.RandString(8),
		Description: &description,
	}

	if projectID != "" {
		opts.ProjectID = &projectID
	}

	role, _, err := TestLakekeeperClient.Role.CreateRole(&opts)
	if err != nil {
		t.Fatalf("could not create test role: %v", err)
	}

	t.Cleanup(func() {
		if _, err := TestLakekeeperClient.Role.DeleteRole(role.ID, projectID); err != nil {
			t.Fatalf("could not cleanup test role: %v", err)
		}
	})

	return role
}

// CreateUser is a test helper for creating a user.
func CreateUser(t *testing.T, id string) *lakekeeper.User {
	t.Helper()

	name := acctest.RandomWithPrefix("acctest")
	email := fmt.Sprintf("%s@test.com", name)
	userType := lakekeeper.HumanUserType

	opts := lakekeeper.ProvisionUserOptions{
		ID:       &id,
		Name:     &name,
		Email:    &email,
		UserType: &userType,
	}

	user, _, err := TestLakekeeperClient.User.ProvisionUser(&opts)
	if err != nil {
		t.Fatalf("could not create test user: %v", err)
	}

	t.Cleanup(func() {
		if _, err := TestLakekeeperClient.User.DeleteUser(user.ID); err != nil {
			t.Fatalf("could not cleanup test user: %v", err)
		}
	})

	return user
}
