//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperProject_basic(t *testing.T) {

	rName := acctest.RandString(8)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
				resource "lakekeeper_project" "foo" {
				  name = "%s"
				}
				`, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project.foo", "name", rName),
					resource.TestCheckResourceAttrSet("lakekeeper_project.foo", "id"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLakekeeperProjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_project" {
			continue
		}

		_, _, err := testutil.TestLakekeeperClient.Project.GetProject(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Project with id %s still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}
