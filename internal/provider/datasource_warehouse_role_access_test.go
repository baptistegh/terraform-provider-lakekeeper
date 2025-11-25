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

func TestAccDataLakekeeperWarehouseRoleAccess_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	role := testutil.CreateRole(t, project.ID)

	keyPrefix := acctest.RandString(8)
	warehouse := testutil.CreateWarehouse(t, project.ID, keyPrefix)

	// assignment
	if _, err := testutil.TestLakekeeperClient.PermissionV1().WarehousePermission().Update(
		t.Context(),
		warehouse.ID,
		&permissionv1.UpdateWarehousePermissionsOptions{
			Writes: []*permissionv1.WarehouseAssignment{
				{
					Assignee: permissionv1.UserOrRole{
						Type:  permissionv1.RoleType,
						Value: role.ID,
					},
					Assignment: permissionv1.CreateWarehouseAssignment,
				},
			},
		},
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
				}`, warehouse.ID, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "id", fmt.Sprintf("%s/%s", warehouse.ID, role.ID)),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "warehouse_id", warehouse.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "role_id", role.ID),
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
									string(permissionv1.CreateNamespace),
									string(permissionv1.GetConfig),
									string(permissionv1.GetMetadata),
									string(permissionv1.ListNamespaces),
									string(permissionv1.IncludeInList),
									string(permissionv1.ListDeletedTabulars),
									string(permissionv1.GetAllTasks),
									string(permissionv1.GetWarehouseEndpointStatistics),
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
