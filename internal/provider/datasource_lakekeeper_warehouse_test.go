//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperWarehouse_basic(t *testing.T) {

	keyPrefix := acctest.RandString(8)
	warehouse := testutil.CreateWarehouse(t, "", keyPrefix)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_warehouse" "default" {
					name = "%s"
				}`, warehouse.Name),
				Check: resource.ComposeAggregateTestCheckFunc(
					// warehouse in default project
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse.default", "id", "00000000-0000-0000-0000-000000000000:"+warehouse.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse.default", "warehouse_id", warehouse.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse.default", "project_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse.default", "name", warehouse.Name),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse.default", "storage_profile.type", "s3"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse.default", "storage_profile.path_style_access", "true"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse.default", "storage_profile.endpoint", "http://minio:9000/"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse.default", "storage_profile.allow_alternative_protocols", "true"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse.default", "storage_profile.key_prefix", keyPrefix),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse.default", "delete_profile.type", "hard"),
					resource.TestCheckNoResourceAttr("data.lakekeeper_warehouse.default", "delete_profile.expiration_seconds"),
				),
			},
		},
	})
}
