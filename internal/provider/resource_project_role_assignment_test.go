//go:build acceptance

package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperProjectRoleAssignment_basic(t *testing.T) {

	project := testutil.CreateProject(t)
	role := testutil.CreateRole(t, project.ID)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperProjectRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_project_role_assignment" "test" {
						project_id = "%s"
						role_id = "%s"
						assignments = ["data_admin"]
					}
				`, project.ID, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "id", project.ID+"/"+role.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "assignments.0", "data_admin"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_project_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_project_role_assignment" "test" {
						project_id = "%s"
						role_id = "%s"
						assignments = ["security_admin"]
					}
				`, project.ID, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "id", project.ID+"/"+role.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "assignments.0", "security_admin"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_project_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_project_role_assignment" "test" {
						project_id = "%s"
						role_id = "%s"
						assignments = ["security_admin", "describe"]
					}
				`, project.ID, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "id", project.ID+"/"+role.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "assignments.#", "2"),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "assignments.0", "describe"),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "assignments.1", "security_admin"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_project_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Check some basic attribute validation
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_project_role_assignment" "test" {
						project_id = "%s"
						role_id = "%s/%s"
						assignments = ["describe"]
					}
				`, project.ID, project.ID, role.ID),
				ExpectError: regexp.MustCompile("Attribute role_id must be a role UUID and NOT include the project UUID"),
			},
			// delete all assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_project_role_assignment" "test" {
						project_id = "%s"
						role_id = "%s"
						assignments = []
					}
				`, project.ID, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "id", project.ID+"/"+role.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "project_id", project.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "assignments.#", "0"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_project_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLakekeeperProjectRoleAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_project_role_assignment" {
			continue
		}

		projectID, roleID := splitInternalID(types.StringValue(rs.Primary.ID))

		assignments, _, err := testutil.TestLakekeeperClient.PermissionV1().ProjectPermission().GetAssignments(context.Background(), projectID, nil)
		if err != nil {
			return fmt.Errorf("could not list project assignments to check destroy, %w", err)
		}

		for _, v := range assignments.Assignments {
			if v.Assignee.Value == roleID && v.Assignee.Type == permissionv1.RoleType {
				return fmt.Errorf("project assignment still exists")
			}
		}
	}

	return nil
}
