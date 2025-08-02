//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"fmt"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperRoleRoleAssignment_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	assignee := testutil.CreateRole(t, project.ID)
	role := testutil.CreateRole(t, project.ID)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperRoleRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_role_role_assignment" "test" {
						role_id = "%s"
						assignee_id = "%s"
						assignments = ["ownership"]
					}
				`, role.ID, assignee.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "id", role.ID+"/"+assignee.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignee_id", assignee.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignments.0", "ownership"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_role_role_assignment" "test" {
						role_id = "%s"
						assignee_id = "%s"
						assignments = ["assignee"]
					}
				`, role.ID, assignee.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "id", role.ID+"/"+assignee.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignee_id", assignee.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignments.0", "assignee"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_role_role_assignment" "test" {
						role_id = "%s"
						assignee_id = "%s"
						assignments = ["ownership", "assignee"]
					}
				`, role.ID, assignee.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "id", role.ID+"/"+assignee.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignee_id", assignee.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignments.#", "2"),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignments.0", "assignee"),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignments.1", "ownership"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// delete all assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_role_role_assignment" "test" {
						role_id = "%s"
						assignee_id = "%s"
						assignments = []
					}
				`, role.ID, assignee.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "id", role.ID+"/"+assignee.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignee_id", assignee.ID),
					resource.TestCheckResourceAttr("lakekeeper_role_role_assignment.test", "assignments.#", "0"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_role_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLakekeeperRoleRoleAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_role_role_assignment" {
			continue
		}

		roleID, assigneeID := splitInternalID(types.StringValue(rs.Primary.ID))

		assignments, _, err := testutil.TestLakekeeperClient.PermissionV1().RolePermission().GetAssignments(context.Background(), roleID, nil)
		if err != nil {
			return fmt.Errorf("could not list project assignments to check destroy, %w", err)
		}

		for _, v := range assignments.Assignments {
			if v.Assignee.Value == assigneeID && v.Assignee.Type == permissionv1.RoleType {
				return fmt.Errorf("project assignment still exists")
			}
		}
	}

	return nil
}
