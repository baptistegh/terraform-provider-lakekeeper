//go:build acceptance

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

func TestAccLakekeeperRole_basic(t *testing.T) {

	rName := acctest.RandString(8)
	rDescription := acctest.RandString(32)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_default_project" "default" {}			
				resource "lakekeeper_role" "foo" {
				  name = "%s"
				  project_id = data.lakekeeper_default_project.default.id
				}
				`, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role.foo", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_role.foo", "project_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckNoResourceAttr("lakekeeper_role.foo", "description"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "role_id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "created_at"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update resource
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_default_project" "default" {}		
				resource "lakekeeper_role" "foo" {
				  name = "%s"
				  project_id = data.lakekeeper_default_project.default.id
				  description = "%s"
				}
				`, rName, rDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role.foo", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_role.foo", "description", rDescription),
					resource.TestCheckResourceAttr("lakekeeper_role.foo", "project_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "role_id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "created_at"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "updated_at"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Rename resource
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_default_project" "default" {}		
				resource "lakekeeper_role" "foo" {
				  name = "new-name"
				  project_id = data.lakekeeper_default_project.default.id
				  description = "%s"
				}
				`, rDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role.foo", "name", "new-name"),
					resource.TestCheckResourceAttr("lakekeeper_role.foo", "description", rDescription),
					resource.TestCheckResourceAttr("lakekeeper_role.foo", "project_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "role_id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "created_at"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "updated_at"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLakekeeperRole_duplicate(t *testing.T) {

	rName := acctest.RandString(8)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`	
				data "lakekeeper_default_project" "default" {}				
				resource "lakekeeper_role" "foo" {
				  project_id = data.lakekeeper_default_project.default.id
				  name = "%s"
				}				
				resource "lakekeeper_role" "toto" {
				  project_id = data.lakekeeper_default_project.default.id
				  name = "%s"
				}
				`, rName, rName),
				ExpectError: regexp.MustCompile("Role(Name)?AlreadyExists"),
			},
		},
	})
}

func TestAccLakekeeperRole_project(t *testing.T) {

	rName := acctest.RandString(8)
	project := testutil.CreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`	
				data "lakekeeper_default_project" "default" {}				
				resource "lakekeeper_role" "foo" {
				  name = "%s"
				  project_id = data.lakekeeper_default_project.default.id
				}				
				resource "lakekeeper_role" "toto" {
				  name = "%s"
				  project_id = "%s"
				}
				`, rName, rName, project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role.foo", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_role.foo", "project_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckNoResourceAttr("lakekeeper_role.foo", "description"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "role_id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.foo", "created_at"),
					resource.TestCheckResourceAttr("lakekeeper_role.toto", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_role.toto", "project_id", project.ID),
					resource.TestCheckNoResourceAttr("lakekeeper_role.toto", "description"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.toto", "id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.toto", "role_id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.toto", "project_id"),
					resource.TestCheckResourceAttrSet("lakekeeper_role.toto", "created_at"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role.toto",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLakekeeperRoleDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_role" {
			continue
		}

		projectID, roleID := splitInternalID(types.StringValue(rs.Primary.ID))

		_, _, err := testutil.TestLakekeeperClient.RoleV1(projectID).Get(context.Background(), roleID)
		if err == nil {
			return fmt.Errorf("Role with id %s still exists", rs.Primary.ID)
		}
		return nil
	}
	return nil
}
