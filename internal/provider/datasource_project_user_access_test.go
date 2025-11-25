//go:build acceptance

package provider

import (
	"fmt"
	"slices"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDataLakekeeperProjectUserAccess_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	user := testutil.CreateUser(t, fmt.Sprintf("oidc~%s", acctest.RandString(8)))

	// assignment
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
				data "lakekeeper_project_user_access" "foo" {
					project_id = "%s"
					user_id = "%s"
				}`, project.ID, user.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_project_user_access.foo", "id", fmt.Sprintf("%s/%s", project.ID, user.ID)),
					resource.TestCheckResourceAttr("data.lakekeeper_project_user_access.foo", "project_id", project.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_project_user_access.foo", "user_id", user.ID),
					resource.TestCheckResourceAttrSet("data.lakekeeper_project_user_access.foo", "allowed_actions.#"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.lakekeeper_project_user_access.foo",
						tfjsonpath.Path(
							tfjsonpath.New("allowed_actions"),
						),
						knownvalue.SetPartial([]knownvalue.Check{
							knownvalue.StringFunc(func(v string) error {
								if !slices.Contains([]string{
									string(permissionv1.CreateRole),
									string(permissionv1.CreateWarehouse),
									string(permissionv1.DeleteProject),
									string(permissionv1.RenameProject),
									string(permissionv1.ListWarehouses),
									string(permissionv1.ListRoles),
									string(permissionv1.SearchRoles),
									string(permissionv1.ReadProjectAssignments),
									string(permissionv1.GrantProjectRoleCreator),
									string(permissionv1.GrantProjectCreate),
									string(permissionv1.GrantProjectDescribe),
									string(permissionv1.GrantProjectModify),
									string(permissionv1.GrantProjectSelet),
									string(permissionv1.GrantProjectAdmin),
									string(permissionv1.GrantSecurityAdmin),
									string(permissionv1.GrantDataAdmin),
									string(permissionv1.GetProjectEndpointStatistics),
								}, v) {
									return fmt.Errorf("%s is not an allowed action", v)
								}
								return nil
							}),
						}),
					),
				},
			},
		},
	})
}
