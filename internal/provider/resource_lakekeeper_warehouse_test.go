//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLakekeeperWarehouse_basic(t *testing.T) {

	rName := acctest.RandString(8)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a warehouse with S3 storage profile
			{
				Config: fmt.Sprintf(`				
				resource "lakekeeper_warehouse" "s3" {
					name = "%s"
					protected = false
					active = true
					storage_profile = {  
						type = "s3"
						region = "us-west-2"
						bucket = "test-bucket"
						sts_enabled = false
					}
					delete_profile = {
						type = "soft"
						expiration_seconds = 3600
					}
				}
				`, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse.foo", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.foo", "storage_profile.type", "s3"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_warehouse.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// func testAccCheckLakekeeperWarehouseDestroy(s *terraform.State) error {
// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "lakekeeper_warehouse" {
// 			continue
// 		}
//
// 		_, err := testutil.TestLakekeeperClient.GetWarehouseByID(context.Background(), rs.Primary.ID) TODO: check how to accept projectID here
// 		if err == nil {
// 			return fmt.Errorf("Warehouse with id %s still exists", rs.Primary.ID)
// 		}
// 		return nil
// 	}
// 	return nil
// }
