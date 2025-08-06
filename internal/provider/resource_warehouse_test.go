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
						s3 = {
							bucket = "testacc"
							endpoint = "http://minio:9000/"
							region = "eu-west-1"
							sts_enabled = false
							remote_signing_url_style = "path"
							credential = {
								access_key = {
									access_key_id = "minio-root-user"
									secret_access_key = "minio-root-password"
								}
							}
						}
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
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.bucket", "testacc"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.endpoint", "http://minio:9000/"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.region", "eu-west-1"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.sts_enabled", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.remote_signing_url_style", "path"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.credential.access_key.access_key_id", "minio-root-user"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.credential.access_key.secret_access_key", "minio-root-password"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.credential.aws_system_identity"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.credential.cloudflare_r2"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "delete_profile.type", "hard"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.s3", "storage_profile.adls"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.s3", "storage_profile.gcs"),
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
						s3 = {
							bucket = "testacc"
							endpoint = "http://minio:9000/"
							region = "eu-west-1"
							sts_enabled = true
							remote_signing_url_style = "virtual_host"
							assume_role_arn = "arn:aws:iam::123456789012:role/AssumeRole"
							credential = {
								access_key = {
									access_key_id = "minio-root-user"
									secret_access_key = "minio-root-password-1"
								}
							}	
						}
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
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.bucket", "testacc"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.endpoint", "http://minio:9000/"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.region", "eu-west-1"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.sts_enabled", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.remote_signing_url_style", "virtual_host"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.credential.access_key.access_key_id", "minio-root-user"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.credential.access_key.secret_access_key", "minio-root-password-1"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.credential.aws_system_identity"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.s3", "storage_profile.s3.credential.cloudflare_r2"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.s3", "storage_profile.adls"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.s3", "storage_profile.gcs"),
				),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_GCS_SystemIdentity(t *testing.T) {

	rName := acctest.RandString(8)
	rPrefix := acctest.RandString(8)
	rBucket := acctest.RandString(8)

	project := testutil.CreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseDestroy,
		Steps: []resource.TestStep{
			// Create a warehouse with S3 storage profile
			{
				Config: fmt.Sprintf(`		
				resource "lakekeeper_warehouse" "gcs_system_identity" {
					name = "%s"
					project_id = "%s"
					storage_profile = {
						gcs  = {
							bucket = "%s"
							key_prefix = "%s"
							credential = {
								gcp_system_identity = {}
							}
						}
					}
				}
				`, rName, project.ID, rBucket, rPrefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.gcs_system_identity", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.gcs_system_identity", "warehouse_id"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_system_identity", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_system_identity", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_system_identity", "protected", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_system_identity", "active", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_system_identity", "managed_access", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_system_identity", "storage_profile.gcs.bucket", rBucket),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_system_identity", "storage_profile.gcs.key_prefix", rPrefix),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_system_identity", "storage_profile.gcs.credential.gcp_system_identity.enabled", "true"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.gcs_system_identity", "storage_profile.gcs.credential.service_account_key"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.gcs_system_identity", "storage_profile.adls"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.gcs_system_identity", "storage_profile.s3"),
				),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_GCS_ServiceAccountKey(t *testing.T) {

	rName := acctest.RandString(8)
	rPrefix := acctest.RandString(8)
	rBucket := acctest.RandString(8)

	key := `{"type":"service_account","project_id":"project-id","private_key_id":"some_key_id","private_key":"-----BEGIN PRIVATE KEY-----\nPRIVATE KEY DATA\n-----END PRIVATE KEY-----\n","client_email":"my-service-account@project-id.iam.gserviceaccount.com","client_id":"123456789012345678901","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url":"https://www.googleapis.com/robot/v1/metadata/x509/my-service-account%40project-id.iam.gserviceaccount.com"}`

	project := testutil.CreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseDestroy,
		Steps: []resource.TestStep{
			// Create a warehouse with S3 storage profile
			{
				Config: fmt.Sprintf(`		
					resource "lakekeeper_warehouse" "gcs_service_account_key" {
						name = "%s"
						project_id = "%s"
						storage_profile = {
							gcs  = {
								bucket = "%s"
								key_prefix = "%s"
								credential = {
									service_account_key = {
										key = file("testdata/key.json")
									}
								}
							}
						}
					}
				`, rName, project.ID, rBucket, rPrefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.gcs_service_account_key", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_warehouse.gcs_service_account_key", "warehouse_id"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "protected", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "active", "true"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "managed_access", "false"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "storage_profile.gcs.bucket", rBucket),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "storage_profile.gcs.key_prefix", rPrefix),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "storage_profile.gcs.credential.service_account_key.key", key),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "storage_profile.gcs.credential.gcp_system_identity"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "storage_profile.adls"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.gcs_service_account_key", "storage_profile.s3"),
				),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_ADLS_SharedAccessKey(t *testing.T) {

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
						adls = {
							account_name = "testacc"
							filesystem = "testfs"
							credential = {
								shared_access_key = {
									key = "test-key"
								}						
							}
						}
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
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.account_name", "testacc"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.filesystem", "testfs"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.shared_access_key.key", "test-key"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.client_credentials"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.azure_system_identity"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.s3"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.gcs"),
				),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_ADLS_ClientCredentials(t *testing.T) {

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
						adls = {
							account_name = "testacc"
							filesystem = "testfs"
							credential = {
								client_credentials = {
									client_id = "test-client-id"
									client_secret = "test-client-secret"
									tenant_id = "test-tenant-id"
								}
							}
						}
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
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.account_name", "testacc"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.filesystem", "testfs"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.client_credentials.client_id", "test-client-id"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.client_credentials.client_secret", "test-client-secret"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.client_credentials.tenant_id", "test-tenant-id"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.shared_access_key"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.azure_system_identity"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.s3"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.gcs"),
				),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_ADLS_SystemIdentity(t *testing.T) {

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
						adls = {
							account_name = "testacc"
							filesystem = "testfs"
							credential = {
								azure_system_identity = {}
							}
			}
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
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.account_name", "testacc"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.filesystem", "testfs"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.azure_system_identity.enabled", "true"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.shared_access_key"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.adls.credential.client_credentials"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.s3"),
					resource.TestCheckNoResourceAttr("lakekeeper_warehouse.adls", "storage_profile.gcs"),
				),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_MultipleStorage(t *testing.T) {
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
						adls = {
							account_name = "testacc"
							filesystem = "testfs"
							credential = {
								azure_system_identity = {}
							}
						}
						s3 = {
							bucket = "testacc"
							endpoint = "http://minio:9000/"
							region = "eu-west-1"
							sts_enabled = false
							credential = {
								access_key = {
									access_key_id = "minio-root-user"
									secret_access_key = "minio-root-password"
								}
							}
						}
					}
				}
				`, rName, project.ID),
				ExpectError: regexp.MustCompile("Incorrect Warehouse creation request"),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_MultipleCreds(t *testing.T) {
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
						s3 = {
							bucket = "testacc"
							endpoint = "http://minio:9000/"
							region = "eu-west-1"
							sts_enabled = false
							credential = {
								access_key = {
									access_key_id = "minio-root-user"
									secret_access_key = "minio-root-password"
								}
								aws_system_identity = {
									external_id = "test"
								}
							}
						}
					}
				}
				`, rName, project.ID),
				ExpectError: regexp.MustCompile("Incorrect Warehouse creation request"),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_IncorrectStorageFamily(t *testing.T) {
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
						s3 = {
							bucket = "testacc"
							endpoint = "http://minio:9000/"
							region = "eu-west-1"
							sts_enabled = false
							credential = {
								azure_system_identity = {}
							}
						}
					}
				}
				`, rName, project.ID),
				ExpectError: regexp.MustCompile("Incorrect Warehouse creation request"),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_NoStorageProfile(t *testing.T) {
	rName := acctest.RandString(8)

	project := testutil.CreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`		
				resource "lakekeeper_warehouse" "no_profile" {
					name = "%s"
					project_id = "%s"
				}
				`, rName, project.ID),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config: fmt.Sprintf(`		
				resource "lakekeeper_warehouse" "no_profile_inside" {
					name = "%s"
					project_id = "%s"
					storage_profile = {}
				}
				`, rName, project.ID),
				ExpectError: regexp.MustCompile("Incorrect Warehouse creation request"),
			},
		},
	})
}

func TestAccLakekeeperWarehouse_NoCreds(t *testing.T) {
	rName := acctest.RandString(8)

	project := testutil.CreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`		
				resource "lakekeeper_warehouse" "s3_no_creds" {
					name = "%s"
					project_id = "%s"
					storage_profile = {
						s3 = {
							bucket = "testacc"
							region = "eu-west-1"
							sts_enabled = false
						}
					}
				}
				`, rName, project.ID),
				ExpectError: regexp.MustCompile("Incorrect attribute value type"),
			},
			{
				Config: fmt.Sprintf(`		
				resource "lakekeeper_warehouse" "s3_no_creds_inside" {
					name = "%s"
					project_id = "%s"
					storage_profile = {
						s3 = {
							bucket = "testacc"
							region = "eu-west-1"
							sts_enabled = false
							credential = {}
						}
					}
				}
				`, rName, project.ID),
				ExpectError: regexp.MustCompile("Incorrect Warehouse creation request"),
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
		if _, _, err := testutil.TestLakekeeperClient.WarehouseV1(projectID).Get(context.Background(), warehouseID); err == nil {
			return fmt.Errorf("Warehouse with id %s still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}
