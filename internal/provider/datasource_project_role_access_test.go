//go:build acceptance

package provider

import (
	"fmt"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDataLakekeeperProjectRoleAccess_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	role := testutil.CreateRole(t, project.ID)

	// assignment
	if _, err := testutil.TestLakekeeperClient.PermissionV1().ProjectPermission().Update(
		t.Context(),
		project.ID,
		&permissionv1.UpdateProjectPermissionsOptions{
			Writes: []*permissionv1.ProjectAssignment{
				{
					Assignee: permissionv1.UserOrRole{
						Type:  permissionv1.RoleType,
						Value: role.ID,
					},
					Assignment: permissionv1.AdminProjectAssignment,
				},
			},
		},
	); err != nil {
		t.Fatalf("could not create project access, %v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_project_role_access" "foo" {
					project_id = "%s"
					role_id = "%s"
				}`, project.ID, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "id", fmt.Sprintf("%s/%s", project.ID, role.ID)),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "project_id", project.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "role_id", role.ID),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.lakekeeper_project_role_access.foo",
						tfjsonpath.Path(
							tfjsonpath.New("allowed_actions"),
						),
						knownvalue.SetPartial([]knownvalue.Check{
							knownvalue.StringExact(string(permissionv1.CreateRole)),
							knownvalue.StringExact(string(permissionv1.CreateWarehouse)),
							knownvalue.StringExact(string(permissionv1.DeleteProject)),
							knownvalue.StringExact(string(permissionv1.RenameProject)),
							knownvalue.StringExact(string(permissionv1.ListWarehouses)),
							knownvalue.StringExact(string(permissionv1.ListRoles)),
							knownvalue.StringExact(string(permissionv1.SearchRoles)),
							knownvalue.StringExact(string(permissionv1.ReadProjectAssignments)),
							knownvalue.StringExact(string(permissionv1.GrantProjectRoleCreator)),
							knownvalue.StringExact(string(permissionv1.GrantProjectCreate)),
							knownvalue.StringExact(string(permissionv1.GrantProjectDescribe)),
							knownvalue.StringExact(string(permissionv1.GrantProjectModify)),
							knownvalue.StringExact(string(permissionv1.GrantProjectSelet)),
							knownvalue.StringExact(string(permissionv1.GrantProjectAdmin)),
							knownvalue.StringExact(string(permissionv1.GrantSecurityAdmin)),
							knownvalue.StringExact(string(permissionv1.GrantDataAdmin)),
							knownvalue.StringExact(string(permissionv1.GetProjectEndpointStatistics))}),
					),
				},
			},
		},
	})
}
