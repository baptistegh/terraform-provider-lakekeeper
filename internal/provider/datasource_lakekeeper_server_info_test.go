//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperServerInfo_basic(t *testing.T) {
	server, _, err := testutil.TestLakekeeperClient.ServerV1().Info(t.Context())
	if err != nil {
		t.Fatalf("could not get server info, %s", err.Error())
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "lakekeeper_server_info" "foo" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_server_info.foo", "bootstrapped", strconv.FormatBool(server.Bootstrapped)),
					resource.TestCheckResourceAttr("data.lakekeeper_server_info.foo", "authz_backend", server.AuthzBackend),
					resource.TestCheckResourceAttr("data.lakekeeper_server_info.foo", "aws_system_identities_enabled", strconv.FormatBool(server.AWSSystemIdentitiesEnabled)),
					resource.TestCheckResourceAttr("data.lakekeeper_server_info.foo", "azure_system_identities_enabled", strconv.FormatBool(server.AzureSystemIdentitiesEnabled)),
					resource.TestCheckResourceAttr("data.lakekeeper_server_info.foo", "gcp_system_identities_enabled", strconv.FormatBool(server.GCPSystemIdentitiesEnabled)),
					resource.TestCheckResourceAttr("data.lakekeeper_server_info.foo", "server_id", server.ServerID),
					resource.TestCheckResourceAttr("data.lakekeeper_server_info.foo", "version", server.Version),
					resource.TestCheckResourceAttr("data.lakekeeper_server_info.foo", "queues.#", fmt.Sprint(rune(len(server.Queues)))),
					resource.TestCheckResourceAttr("data.lakekeeper_server_info.foo", "default_project_id", server.DefaultProjectID),
				),
			},
		},
	})
}
