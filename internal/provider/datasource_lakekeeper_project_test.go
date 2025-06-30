//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperProject_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
				data "lakekeeper_project" "foo" {
				  name = "%s"
				}
				`, project.Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_project.foo", "id", project.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_project.foo", "name", project.Name),
				),
			},
		},
	})
}
