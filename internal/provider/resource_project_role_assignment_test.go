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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperProjectRoleAssignment_basic(t *testing.T) {

	project := testutil.CreateProject(t)
	role := testutil.CreateRole(t, project.ProjectId)

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
				`, project.ProjectId, role.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "id", project.ProjectId+"/"+role.Id),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "project_id", project.ProjectId),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "role_id", role.Id),
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
				`, project.ProjectId, role.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "id", project.ProjectId+"/"+role.Id),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "project_id", project.ProjectId),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "role_id", role.Id),
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
						assignments = ["security_admin", "describe", "write"]
					}
				`, project.ProjectId, role.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "id", project.ProjectId+"/"+role.Id),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "project_id", project.ProjectId),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "role_id", role.Id),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "assignments.#", "3"),
					resource.TestCheckTypeSetElemAttr("lakekeeper_project_role_assignment.test", "assignments.*", "write"),
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
				`, project.ProjectId, project.ProjectId, role.Id),
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
				`, project.ProjectId, role.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "id", project.ProjectId+"/"+role.Id),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "project_id", project.ProjectId),
					resource.TestCheckResourceAttr("lakekeeper_project_role_assignment.test", "role_id", role.Id),
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

		relations, err := projectAssignmentScope(testutil.TestLakekeeperClient, projectID).relations(context.Background(), permissions.PrincipalRole, roleID)
		if err != nil {
			return fmt.Errorf("could not list project assignments to check destroy, %w", err)
		}

		if len(relations) > 0 {
			return fmt.Errorf("project assignment still exists")
		}
	}

	return nil
}
