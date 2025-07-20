//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperWarehouse_basic(t *testing.T) {

	rName := acctest.RandString(8)

	project := testutil.CreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseDestroy,
		Steps: []resource.TestStep{
			// Create a warehouse with S3 storage profile
			{
				Config: fmt.Sprintf(`		
				resource "lakekeeper_warehouse" "s3" {
					name = "%s"
					project_id = "%s"
					storage_profile = {
						type = "s3"
						bucket = "testacc"
						endpoint = "http://minio:9000/"
						region = "eu-west-1"
						sts_enabled = false
					}
					storage_credential = {
						type = "s3_access_key",
						access_key_id = "minio-root-user"
						secret_access_key = "minio-root-password"
					}
				}
				`, rName, project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.s3", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.s3", "warehouse_id"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "protected", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "active", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "managed_access", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.type", "s3"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.bucket", "testacc"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.endpoint", "http://minio:9000/"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.region", "eu-west-1"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.sts_enabled", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_credential.type", "s3_access_key"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_credential.access_key_id", "minio-root-user"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_credential.secret_access_key", "minio-root-password"),
				),
			},
			// Import is not configured
			{
				ResourceName:      "lakekeeper_warehouse.s3",
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile("Import Not Implemented"),
			},
			{
				Config: fmt.Sprintf(`		
				resource "lakekeeper_warehouse" "s3" {
					name = "%s"
					project_id = "%s"
					managed_access = true
					storage_profile = {
						type = "s3"
						bucket = "testacc"
						endpoint = "http://minio:9000/"
						region = "eu-west-1"
						sts_enabled = true
						assume_role_arn = "arn:aws:iam::123456789012:role/AssumeRole"
					}
					storage_credential = {
						type = "s3_access_key",
						access_key_id = "minio-root-user"
						secret_access_key = "minio-root-password-1"
					}
				}
				`, rName, project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.s3", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.s3", "warehouse_id"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "protected", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "active", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "managed_access", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.type", "s3"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.bucket", "testacc"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.endpoint", "http://minio:9000/"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.region", "eu-west-1"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.sts_enabled", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_credential.type", "s3_access_key"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_credential.access_key_id", "minio-root-user"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_credential.secret_access_key", "minio-root-password-1"),
				),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_GCS(t *testing.T) {

	rName := acctest.RandString(8)

	project := testutil.CreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseDestroy,
		Steps: []resource.TestStep{
			// Create a warehouse with S3 storage profile
			{
				Config: fmt.Sprintf(`		
				resource "lakekeeper_warehouse" "gcs" {
					name = "%s"
					project_id = "%s"
					storage_profile = {
						type = "gcs"
						bucket = "testacc"
					}
					storage_credential = {
						type = "gcs_gcp_system_identity",
					}
				}
				`, rName, project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.gcs", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.gcs", "warehouse_id"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs", "protected", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs", "active", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs", "managed_access", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs", "storage_profile.type", "gcs"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs", "storage_profile.bucket", "testacc"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs", "storage_credential.type", "gcs_gcp_system_identity"),
				),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_ADLS(t *testing.T) {

	rName := acctest.RandString(8)

	project := testutil.CreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseDestroy,
		Steps: []resource.TestStep{
			// Create a warehouse with S3 storage profile
			{
				Config: fmt.Sprintf(`		
				resource "lakekeeper_warehouse" "adls" {
					name = "%s"
					project_id = "%s"
					storage_profile = {
						type = "adls"
						account_name = "testacc"
						filesystem = "testfs"
					}
					storage_credential = {
						type = "az_shared_access_key",
						az_key = "test-key"
					}
				}
				`, rName, project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.adls", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.adls", "warehouse_id"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "protected", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "active", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "managed_access", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.type", "adls"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.account_name", "testacc"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.filesystem", "testfs"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_credential.type", "az_shared_access_key"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_credential.az_key", "test-key"),
				),
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
		if _, _, err := testutil.TestLakekeeperClient.WarehouseV1(projectID).Get(warehouseID); err == nil {
			return fmt.Errorf("Warehouse with id %s still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}
