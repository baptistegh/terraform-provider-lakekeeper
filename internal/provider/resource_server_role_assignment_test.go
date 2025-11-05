//go:build acceptance

package provider

import (
	"context"
	"fmt"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperServerRoleAssignment_basic(t *testing.T) {

	project := testutil.CreateProject(t)
	role := testutil.CreateRole(t, project.ID)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperServerRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_server_role_assignment" "test" {
						role_id = "%s"
						assignments = ["operator"]
					}
				`, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "assignments.0", "operator"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_server_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_server_role_assignment" "test" {
						role_id = "%s"
						assignments = ["admin"]
					}
				`, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "assignments.0", "admin"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_server_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_server_role_assignment" "test" {
						role_id = "%s"
						assignments = ["admin", "operator"]
					}
				`, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "assignments.#", "2"),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "assignments.0", "admin"),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "assignments.1", "operator"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_server_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// delete all assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_server_role_assignment" "test" {
						role_id = "%s"
						assignments = []
					}
				`, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_server_role_assignment.test", "assignments.#", "0"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_server_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLakekeeperServerRoleAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_server_role_assignment" {
			continue
		}

		assignments, _, err := testutil.TestLakekeeperClient.PermissionV1().ServerPermission().GetAssignments(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("could not list server assignments to check destroy, %w", err)
		}

		for _, v := range assignments.Assignments {
			if v.Assignee.Value == rs.Primary.ID && v.Assignee.Type == permissionv1.RoleType {
				return fmt.Errorf("server assignment still exists")
			}
		}
	}

	return nil
}
