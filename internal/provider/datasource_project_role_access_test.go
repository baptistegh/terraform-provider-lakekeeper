//go:build acceptance

package provider

import (
	"fmt"
	"slices"
	"testing"

	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDataLakekeeperProjectRoleAccess_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	role := testutil.CreateRole(t, project.ProjectId)

	// assignment
	if err := grantAssignments(t.Context(), projectAssignmentScope(testutil.TestLakekeeperClient, project.ProjectId),
		testAssignment{permissions.PrincipalRole, role.Id, string(managementv1.PROJECTRELATION_PROJECT_ADMIN)},
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
				}`, project.ProjectId, role.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "id", fmt.Sprintf("%s/%s", project.ProjectId, role.Id)),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "project_id", project.ProjectId),
					resource.TestCheckResourceAttr("data.lakekeeper_project_role_access.foo", "role_id", role.Id),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.lakekeeper_project_role_access.foo",
						tfjsonpath.Path(
							tfjsonpath.New("allowed_actions"),
						),
						knownvalue.SetPartial([]knownvalue.Check{
							knownvalue.StringFunc(func(v string) error {
								if !slices.Contains([]string{
									"create_role",
									"create_warehouse",
									"delete",
									"rename",
									"list_warehouses",
									"list_roles",
									"search_roles",
									"read_assignments",
									"grant_role_creator",
									"grant_create",
									"grant_describe",
									"grant_modify",
									"grant_select",
									"grant_project_admin",
									"grant_security_admin",
									"grant_data_admin",
									"get_endpoint_statistics",
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
