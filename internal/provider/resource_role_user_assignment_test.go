//go:build acceptance

package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperRoleUserAssignment_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	userID := fmt.Sprintf("oidc~%s", acctest.RandString(32))
	user := testutil.CreateUser(t, userID)

	role := testutil.CreateRole(t, project.ProjectId)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperRoleUserAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_role_user_assignment" "test" {
						role_id = "%s"
						user_id = "%s"
						assignments = ["ownership"]
					}
				`, role.Id, user.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "id", role.Id+"/"+user.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "role_id", role.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "user_id", user.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "assignments.0", "ownership"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_role_user_assignment" "test" {
						role_id = "%s"
						user_id = "%s"
						assignments = ["assignee"]
					}
				`, role.Id, user.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "id", role.Id+"/"+user.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "role_id", role.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "user_id", user.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "assignments.0", "assignee"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_role_user_assignment" "test" {
						role_id = "%s"
						user_id = "%s"
						assignments = ["ownership", "assignee"]
					}
				`, role.Id, user.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "id", role.Id+"/"+user.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "role_id", role.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "user_id", user.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "assignments.#", "2"),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "assignments.0", "assignee"),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "assignments.1", "ownership"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Check some basic attribute validation
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_role_user_assignment" "test" {
						role_id = "%s/%s"
						user_id = "%s"
						assignments = ["ownership"]
					}
				`, project.ProjectId, role.Id, user.Id),
				ExpectError: regexp.MustCompile("Attribute role_id must be a role UUID and NOT include the project UUID"),
			},
			// delete all assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_role_user_assignment" "test" {
						role_id = "%s"
						user_id = "%s"
						assignments = []
					}
				`, role.Id, user.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "id", role.Id+"/"+user.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "role_id", role.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "user_id", user.Id),
					resource.TestCheckResourceAttr("lakekeeper_role_user_assignment.test", "assignments.#", "0"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLakekeeperRoleUserAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_role_user_assignment" {
			continue
		}

		roleID, userID := splitInternalID(types.StringValue(rs.Primary.ID))

		relations, err := roleAssignmentScope(testutil.TestLakekeeperClient, roleID).relations(context.Background(), permissions.PrincipalUser, userID)
		if err != nil {
			return fmt.Errorf("could not list project assignments to check destroy, %w", err)
		}

		if len(relations) > 0 {
			return fmt.Errorf("project assignment still exists")
		}
	}

	return nil
}
