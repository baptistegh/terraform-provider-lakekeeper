//go:build acceptance

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperDefaultProject_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "lakekeeper_default_project" "foo" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_default_project.foo", "id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.lakekeeper_default_project.foo", "name", "Default Project"),
				),
			},
		},
	})
}
