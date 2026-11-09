//go:build acceptance

package provider

import (
	"fmt"
	"testing"

	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperProjectAssignments_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	user := testutil.CreateUser(t, fmt.Sprintf("oidc~%s", acctest.RandString(8)))
	role := testutil.CreateRole(t, project.ProjectId)

	// assignment 2
	if err := grantAssignments(t.Context(), projectAssignmentScope(testutil.TestLakekeeperClient, project.ProjectId),
		testAssignment{permissions.PrincipalUser, user.Id, string(managementv1.PROJECTRELATION_SELECT)},
		testAssignment{permissions.PrincipalRole, role.Id, string(managementv1.PROJECTRELATION_MODIFY)},
	); err != nil {
		t.Fatalf("could not create project assignments, %v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_project_assignments" "foo" {
					id = "%s"
				}`, project.ProjectId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "id", project.ProjectId),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.#", "3"), // also have the admin user

					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.0.assignee_id", testutil.DefaultUserID),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.0.assignee_type", "user"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.0.assignment", string(managementv1.PROJECTRELATION_PROJECT_ADMIN)),

					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.1.assignee_id", user.Id),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.1.assignee_type", "user"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.1.assignment", string(managementv1.PROJECTRELATION_SELECT)),

					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.2.assignee_id", role.Id),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.2.assignee_type", "role"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.2.assignment", string(managementv1.PROJECTRELATION_MODIFY)),
				),
			},
		},
	})
}
