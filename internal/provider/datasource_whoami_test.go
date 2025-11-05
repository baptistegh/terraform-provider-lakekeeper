//go:build acceptance

package provider

import (
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperWhoami_basic(t *testing.T) {
	user, _, err := testutil.TestLakekeeperClient.UserV1().Whoami(t.Context())
	if err != nil {
		t.Fatalf("could not get current user, %s", err.Error())
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "lakekeeper_whoami" "foo" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_whoami.foo", "id", user.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_whoami.foo", "name", user.Name),
					resource.TestCheckResourceAttr("data.lakekeeper_whoami.foo", "created_at", user.CreatedAt),
					resource.TestCheckResourceAttr("data.lakekeeper_whoami.foo", "last_updated_with", user.LastUpdatedWith),
					resource.TestCheckResourceAttr("data.lakekeeper_whoami.foo", "user_type", string(user.UserType)),
				),
			},
		},
	})
}
