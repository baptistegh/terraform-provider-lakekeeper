//go:build acceptance

package provider

import (
	"fmt"
	"slices"
	"testing"

	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"

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
	if err := grantAssignments(t.Context(), projectAssignmentScope(testutil.TestLakekeeperClient, project.ProjectId),
		testAssignment{permissions.PrincipalUser, user.Id, string(managementv1.PROJECTRELATION_PROJECT_ADMIN)},
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
				}`, project.ProjectId, user.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_project_user_access.foo", "id", fmt.Sprintf("%s/%s", project.ProjectId, user.Id)),
					resource.TestCheckResourceAttr("data.lakekeeper_project_user_access.foo", "project_id", project.ProjectId),
					resource.TestCheckResourceAttr("data.lakekeeper_project_user_access.foo", "user_id", user.Id),
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
