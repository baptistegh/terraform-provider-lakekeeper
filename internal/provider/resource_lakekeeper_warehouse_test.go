//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperWarehouse_basic(t *testing.T) {

	rName := acctest.RandString(8)
	rPrefix := acctest.RandString(12)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseDestroy,
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
						bucket = "testacc"
						endpoint = "http://minio:9000/"
						region ="local-01"
						sts_enabled = false
						path_style_access = true
						key_prefix = "%s"
					}
					storage_credential = {
						type = "s3_access_key",
						access_key_id = "minio-root-user"
						secret_access_key = "minio-root-password"
					}
				}
				`, rName, rPrefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "protected", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "active", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.type", "s3"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.bucket", "testacc"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.endpoint", "http://minio:9000/"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.region", "local-01"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.sts_enabled", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.path_style_access", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_credential.type", "s3_access_key"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_credential.access_key_id", "minio-root-user"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_credential.secret_access_key", "minio-root-password"),
				),
			},
			// Update must throw an error for now
			{
				Config: fmt.Sprintf(`				
				resource "lakekeeper_warehouse" "s3" {
					name = "%s"
					protected = false
					active = true
					storage_profile = {
						type = "s3"
						bucket = "testacc"
						endpoint = "http://minio:9000/"
						region ="local-01"
						sts_enabled = false
						path_style_access = true
						key_prefix = "%s"
					}
					storage_credential = {
						type = "s3_access_key",
						access_key_id = "minio-root-user-2"
						secret_access_key = "minio-root-password"
					}
				}
				`, rName, rPrefix),
				ExpectError: regexp.MustCompile("requested to perform an in-place upgrade"),
			},
			// Import is not configured
			{
				ResourceName:      "lakekeeper_warehouse.s3",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile("Import Not Implemented"),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_NonImplemented(t *testing.T) {

	rName := acctest.RandString(8)
	rPrefix := acctest.RandString(12)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseDestroy,
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
						bucket = "testacc"
						endpoint = "http://minio:9000/"
						region ="local-01"
						sts_enabled = false
						path_style_access = true
						key_prefix = "%s"
					}
					storage_credential = {
						type = "s3_access_key",
						access_key_id = "minio-root-user"
						secret_access_key = "minio-root-password"
					}
				}
				`, rName, rPrefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.type", "s3"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_credential.type", "s3_access_key"),
				),
			},
			// Update must throw an error for now
			{
				Config: fmt.Sprintf(`				
				resource "lakekeeper_warehouse" "s3" {
					name = "%s"
					protected = false
					active = true
					storage_profile = {
						type = "s3"
						bucket = "testacc"
						endpoint = "http://minio:9000/"
						region ="local-01"
						sts_enabled = false
						path_style_access = true
						key_prefix = "%s"
					}
					storage_credential = {
						type = "s3_access_key",
						access_key_id = "minio-root-user-2"
						secret_access_key = "minio-root-password"
					}
				}
				`, rName, rPrefix),
				ExpectError: regexp.MustCompile("requested to perform an in-place upgrade"),
			},
			// Import is not configured
			{
				ResourceName:      "lakekeeper_warehouse.s3",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile("Import Not Implemented"),
			},
		},
	})
}

func testAccCheckLakekeeperWarehouseDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_warehouse" {
			continue
		}

		projectID, warehouseID := splitInternalID(types.StringValue(rs.Primary.ID))
		_, err := testutil.TestLakekeeperClient.GetWarehouseByID(context.Background(), projectID, warehouseID)
		if err == nil {
			return fmt.Errorf("Warehouse with id %s still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}
