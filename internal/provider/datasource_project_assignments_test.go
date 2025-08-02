//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperProjectAssignments_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	user := testutil.CreateUser(t, fmt.Sprintf("oidc~%s", acctest.RandString(8)))
	role := testutil.CreateRole(t, project.ID)

	// assignment 2
	if _, err := testutil.TestLakekeeperClient.PermissionV1().ProjectPermission().Update(
		t.Context(),
		project.ID,
		&permissionv1.UpdateProjectPermissionsOptions{
			Writes: []*permissionv1.ProjectAssignment{
				{
					Assignee: permissionv1.UserOrRole{
						Type:  permissionv1.UserType,
						Value: user.ID,
					},
					Assignment: permissionv1.SelectProjectAssignment,
				}, {
					Assignee: permissionv1.UserOrRole{
						Type:  permissionv1.RoleType,
						Value: role.ID,
					},
					Assignment: permissionv1.ModifyProjectAssignment,
				},
			},
		},
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
				}`, project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "id", project.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.#", "3"), // also have the admin user

					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.0.assignee_id", testutil.DefaultUserID),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.0.assignee_type", string(permissionv1.UserType)),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.0.assignment", string(permissionv1.AdminProjectAssignment)),

					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.1.assignee_id", user.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.1.assignee_type", string(permissionv1.UserType)),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.1.assignment", string(permissionv1.SelectProjectAssignment)),

					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.2.assignee_id", role.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.2.assignee_type", string(permissionv1.RoleType)),
					resource.TestCheckResourceAttr("data.lakekeeper_project_assignments.foo", "assignments.2.assignment", string(permissionv1.ModifyProjectAssignment)),
				),
			},
		},
	})
}
