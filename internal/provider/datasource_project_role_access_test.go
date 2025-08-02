//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.#", "16"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.0", "create_role"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.1", "create_warehouse"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.2", "delete"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.3", "grant_create"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.4", "grant_data_admin"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.5", "grant_describe"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.6", "grant_modify"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.7", "grant_project_admin"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.8", "grant_role_creator"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.9", "grant_security_admin"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.10", "grant_select"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.11", "list_roles"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.12", "list_warehouses"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.13", "read_assignments"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.14", "rename"),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "allowed_actions.15", "search_roles"),
				),
			},
		},
	})
}
