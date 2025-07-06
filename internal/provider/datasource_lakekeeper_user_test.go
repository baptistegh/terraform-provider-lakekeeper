//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperUser_basic(t *testing.T) {
	rID := fmt.Sprintf("oidc~%s", uuid.New().String())
	user := testutil.CreateUser(t, rID)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_user" "foo" {
					id = "%s"
				}`, user.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_user.foo", "id", rID),
					resource.TestCheckResourceAttr("data.lakekeeper_user.foo", "name", user.Name),
					resource.TestCheckResourceAttr("data.lakekeeper_user.foo", "email", *user.Email),
					resource.TestCheckResourceAttr("data.lakekeeper_user.foo", "created_at", user.CreatedAt),
					resource.TestCheckResourceAttr("data.lakekeeper_user.foo", "last_updated_with", user.LastUpdatedWith),
					resource.TestCheckResourceAttr("data.lakekeeper_user.foo", "user_type", string(user.UserType)),
				),
			},
		},
	})
}
