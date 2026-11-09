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

func TestAccDataLakekeeperWarehouseRoleAccess_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	role := testutil.CreateRole(t, project.ProjectId)

	keyPrefix := acctest.RandString(8)
	warehouse := testutil.CreateWarehouse(t, project.ProjectId, keyPrefix)

	// assignment
	if err := grantAssignments(t.Context(), warehouseAssignmentScope(testutil.TestLakekeeperClient, warehouse.Id),
		testAssignment{permissions.PrincipalRole, role.Id, string(managementv1.WAREHOUSERELATION_CREATE)},
	); err != nil {
		t.Fatalf("could not create warehouse access, %v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_warehouse_role_access" "foo" {
					warehouse_id = "%s"
					role_id = "%s"
				}`, warehouse.Id, role.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "id", fmt.Sprintf("%s/%s", warehouse.Id, role.Id)),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "warehouse_id", warehouse.Id),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "role_id", role.Id),
					resource.TestCheckResourceAttrSet("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.#"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.lakekeeper_warehouse_role_access.foo",
						tfjsonpath.Path(
							tfjsonpath.New("allowed_actions"),
						),
						knownvalue.SetPartial([]knownvalue.Check{
							knownvalue.StringFunc(func(v string) error {
								if !slices.Contains([]string{
									"create_namespace",
									"get_config",
									"get_metadata",
									"list_namespaces",
									"include_in_list",
									"list_deleted_tabulars",
									"get_all_tasks",
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
