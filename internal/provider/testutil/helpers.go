//go:build acceptance

package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	"github.com/lakekeeper/go-lakekeeper/pkg/client"
	"github.com/lakekeeper/go-lakekeeper/pkg/storage/credential"
	"github.com/lakekeeper/go-lakekeeper/pkg/storage/profile"
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
}

var TestLakekeeperClient *client.Client

var DefaultUserID string

func init() {
	client, err := testLakekeeperConfig.NewLakekeeperClient(context.Background())
	if err != nil {
		panic("failed to create test client: " + err.Error())
	}

	TestLakekeeperClient = client

	user, err := TestLakekeeperClient.Users.Whoami(context.Background())
	if err != nil {
		panic("failed to get current user: " + err.Error())
	}

	DefaultUserID = user.Id
}

// CreateProject is a test helper for creating a project.
func CreateProject(t *testing.T) *managementv1.GetProjectResponse {
	t.Helper()

	resp, err := TestLakekeeperClient.Projects.Create(t.Context(), managementv1.NewCreateProjectRequest(acctest.RandomWithPrefix("acctest")))
	if err != nil {
		t.Fatalf("could not create test project: %v", err)
	}

	t.Cleanup(func() {
		if err := TestLakekeeperClient.Projects.Delete(context.Background(), resp.ProjectId); err != nil {
			t.Fatalf("could not cleanup test project: %v", err)
		}
	})

	project, err := TestLakekeeperClient.Projects.Get(t.Context(), resp.ProjectId)
	if err != nil {
		t.Fatalf("could not get test project: %v", err)
	}

	return project
}

// CreateWarehouse is a test helper for creating a warehouse.
func CreateWarehouse(t *testing.T, projectID, keyPrefix string) *managementv1.GetWarehouseResponse {
	t.Helper()

	req := managementv1.NewCreateWarehouseRequest(
		profile.NewS3Profile("testacc", "local-01",
			profile.WithS3Endpoint("http://minio:9000/"),
			profile.WithS3PathStyleAccess(),
			profile.WithS3KeyPrefix(keyPrefix),
		),
		acctest.RandString(8),
	)
	req.SetProjectId(projectID)
	req.SetStorageCredential(credential.NewS3AccessKey("minio-root-user", "minio-root-password"))
	hard := managementv1.TabularDeleteProfileHardAsTabularDeleteProfile(managementv1.NewTabularDeleteProfileHard("hard"))
	req.SetDeleteProfile(hard)

	w, err := TestLakekeeperClient.Warehouses.Create(t.Context(), req)
	if err != nil {
		t.Fatalf("could not create test warehouse: %v", err)
	}

	t.Cleanup(func() {
		if _, err := TestLakekeeperClient.WarehouseAPI.DeleteWarehouse(context.Background(), w.Id).Force(true).Execute(); err != nil {
			t.Fatalf("could not cleanup test warehouse: %v", err)
		}
	})

	warehouse, err := TestLakekeeperClient.Warehouses.Get(t.Context(), w.Id)
	if err != nil {
		t.Fatalf("could not create test warehouse: %v", err)
	}

	return warehouse
}

// CreateRole is a test helper for creating a role.
func CreateRole(t *testing.T, projectID string) *managementv1.Role {
	t.Helper()

	req := managementv1.NewCreateRoleRequest(acctest.RandString(8))
	req.SetDescription(acctest.RandString(32))

	role, err := TestLakekeeperClient.Roles.Create(t.Context(), projectID, req)
	if err != nil {
		t.Fatalf("could not create test role: %v", err)
	}

	t.Cleanup(func() {
		if err := TestLakekeeperClient.Roles.Delete(context.Background(), projectID, role.Id); err != nil {
			t.Fatalf("could not cleanup test role: %v", err)
		}
	})

	return role
}

// CreateUser is a test helper for creating a user.
func CreateUser(t *testing.T, id string) *managementv1.User {
	t.Helper()

	name := acctest.RandomWithPrefix("acctest")

	req := managementv1.NewCreateUserRequest()
	req.SetId(id)
	req.SetName(name)
	req.SetEmail(fmt.Sprintf("%s@test.com", name))
	req.SetUserType(managementv1.USERTYPE_HUMAN)

	user, err := TestLakekeeperClient.Users.Create(t.Context(), req)
	if err != nil {
		t.Fatalf("could not create test user: %v", err)
	}

	t.Cleanup(func() {
		if err := TestLakekeeperClient.Users.Delete(context.Background(), user.Id); err != nil {
			t.Fatalf("could not cleanup test user: %v", err)
		}
	})

	return user
}
