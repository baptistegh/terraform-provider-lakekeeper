//go:build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperWarehouses_basic(t *testing.T) {
	project := testutil.CreateProject(t)
	warehouse := testutil.CreateWarehouse(t, project.ID, acctest.RandString(8))

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_warehouses" "this" {
				  project_id = "%s"
				}
				`, project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_warehouses.this", "warehouses.#", "1"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouses.this", "warehouses.0.id", warehouse.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouses.this", "warehouses.0.name", warehouse.Name),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouses.this", "warehouses.0.active", "true"),
				),
			},
		},
	})
}
