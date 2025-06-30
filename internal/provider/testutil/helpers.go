//go:build acceptance || flakey || settings
// +build acceptance flakey settings

package testutil

import (
	"context"
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
	InitialBootstrap: true,
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
