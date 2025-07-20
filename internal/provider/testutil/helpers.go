//go:build acceptance
// +build acceptance

package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/credential"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
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

var DefaultUserID string

func init() {
	client, err := testLakekeeperConfig.NewLakekeeperClient(context.Background())
	if err != nil {
		panic("failed to create test client: " + err.Error())
	}

	TestLakekeeperClient = client

	user, _, err := TestLakekeeperClient.UserV1().Whoami()
	if err != nil {
		panic("failed to get current user: " + err.Error())
	}

	DefaultUserID = user.ID
}

// CreateProject is a test helper for creating a project.
func CreateProject(t *testing.T) *managementv1.Project {
	t.Helper()

	opts := managementv1.CreateProjectOptions{
		Name: acctest.RandomWithPrefix("acctest"),
	}

	resp, _, err := TestLakekeeperClient.ProjectV1().Create(&opts)
	if err != nil {
		t.Fatalf("could not create test project: %v", err)
	}

	t.Cleanup(func() {
		if _, err := TestLakekeeperClient.ProjectV1().Delete(resp.ID); err != nil {
			t.Fatalf("could not cleanup test project: %v", err)
		}
	})

	project, _, err := TestLakekeeperClient.ProjectV1().Get(resp.ID)
	if err != nil {
		t.Fatalf("could not get test project: %v", err)
	}

	return project
}

// CreateWarehouse is a test helper for creating a warehouse.
func CreateWarehouse(t *testing.T, projectID, keyPrefix string) *managementv1.Warehouse {
	t.Helper()

	storage := profile.NewS3StorageSettings("testacc", "local-01",
		profile.WithEndpoint("http://minio:9000/"),
		profile.WithPathStyleAccess(),
		profile.WithS3KeyPrefix(keyPrefix),
	)

	creds := credential.NewS3CredentialAccessKey("minio-root-user", "minio-root-password")

	opts := managementv1.CreateWarehouseOptions{
		Name:              acctest.RandString(8),
		StorageProfile:    storage.AsProfile(),
		StorageCredential: creds.AsCredential(),
		DeleteProfile:     profile.NewTabularDeleteProfileHard().AsProfile(),
	}

	w, _, err := TestLakekeeperClient.WarehouseV1(projectID).Create(&opts)
	if err != nil {
		t.Fatalf("could not create test warehouse: %v", err)
	}

	t.Cleanup(func() {
		opts := managementv1.DeleteWarehouseOptions{
			Force: core.Ptr(true),
		}
		if _, err := TestLakekeeperClient.WarehouseV1(projectID).Delete(w.ID, &opts); err != nil {
			t.Fatalf("could not cleanup test warehouse: %v", err)
		}
	})

	warehouse, _, err := TestLakekeeperClient.WarehouseV1(projectID).Get(w.ID)
	if err != nil {
		t.Fatalf("could not create test warehouse: %v", err)
	}

	return warehouse
}

// CreateRole is a test helper for creating a role.
func CreateRole(t *testing.T, projectID string) *managementv1.Role {
	t.Helper()

	description := acctest.RandString(32)

	opts := managementv1.CreateRoleOptions{
		Name:        acctest.RandString(8),
		Description: &description,
	}

	role, _, err := TestLakekeeperClient.RoleV1(projectID).Create(&opts)
	if err != nil {
		t.Fatalf("could not create test role: %v", err)
	}

	t.Cleanup(func() {
		if _, err := TestLakekeeperClient.RoleV1(projectID).Delete(role.ID); err != nil {
			t.Fatalf("could not cleanup test role: %v", err)
		}
	})

	return role
}

// CreateUser is a test helper for creating a user.
func CreateUser(t *testing.T, id string) *managementv1.User {
	t.Helper()

	name := acctest.RandomWithPrefix("acctest")
	email := fmt.Sprintf("%s@test.com", name)
	userType := managementv1.HumanUserType

	opts := managementv1.ProvisionUserOptions{
		ID:       &id,
		Name:     &name,
		Email:    &email,
		UserType: &userType,
	}

	user, _, err := TestLakekeeperClient.UserV1().Provision(&opts)
	if err != nil {
		t.Fatalf("could not create test user: %v", err)
	}

	t.Cleanup(func() {
		if _, err := TestLakekeeperClient.UserV1().Delete(user.ID); err != nil {
			t.Fatalf("could not cleanup test user: %v", err)
		}
	})

	return user
}
