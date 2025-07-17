//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperProjectUserAssignment_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	userID := fmt.Sprintf("oidc~%s", acctest.RandString(32))
	user := testutil.CreateUser(t, userID)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperProjectUserAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_project_user_assignment" "test" {
						project_id = "%s"
						user_id = "%s"
						assignments = ["project_admin"]
					}
				`, project.ID, user.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "id", project.ID+":"+user.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "user_id", user.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "assignments.0", "project_admin"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_project_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_project_user_assignment" "test" {
						project_id = "%s"
						user_id = "%s"
						assignments = ["role_creator"]
					}
				`, project.ID, user.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "id", project.ID+":"+user.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "user_id", user.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "assignments.0", "role_creator"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_project_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_project_user_assignment" "test" {
						project_id = "%s"
						user_id = "%s"
						assignments = ["role_creator", "select"]
					}
				`, project.ID, user.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "id", project.ID+":"+user.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "user_id", user.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "assignments.#", "2"),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "assignments.0", "role_creator"),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "assignments.1", "select"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_project_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// delete all assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_project_user_assignment" "test" {
						project_id = "%s"
						user_id = "%s"
						assignments = []
					}
				`, project.ID, user.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "id", project.ID+":"+user.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "user_id", user.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_user_assignment.test", "assignments.#", "0"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_project_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLakekeeperProjectUserAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_project_user_assignment" {
			continue
		}

		projectID, userID := splitInternalID(types.StringValue(rs.Primary.ID))

		assignments, _, err := testutil.TestLakekeeperClient.PermissionV1().ProjectPermission().GetAssignments(projectID, nil)
		if err != nil {
			return fmt.Errorf("could not list project assignments to check destroy, %w", err)
		}

		for _, v := range assignments.Assignments {
			if v.Assignee.Value == userID && v.Assignee.Type == permissionv1.UserType {
				return fmt.Errorf("project assignment still exists")
			}
		}
	}

	return nil
}
