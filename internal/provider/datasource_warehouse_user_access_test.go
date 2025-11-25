//go:build acceptance

package provider

import (
	"fmt"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccDataLakekeeperWarehouseUserAccess_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	user := testutil.CreateUser(t, fmt.Sprintf("oidc~%s", acctest.RandString(8)))

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
						Type:  permissionv1.UserType,
						Value: user.ID,
					},
					Assignment: permissionv1.DescribeWarehouseAssignment,
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
				data "lakekeeper_warehouse_user_access" "foo" {
					warehouse_id = "%s"
					user_id = "%s"
				}`, warehouse.ID, user.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_user_access.foo", "id", fmt.Sprintf("%s/%s", warehouse.ID, user.ID)),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_user_access.foo", "warehouse_id", warehouse.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_user_access.foo", "user_id", user.ID),
					resource.TestCheckResourceAttrSet("data.lakekeeper_warehouse_user_access.foo", "allowed_actions.#"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.lakekeeper_warehouse_user_access.foo",
						tfjsonpath.Path(
							tfjsonpath.New("allowed_actions"),
						),
						knownvalue.SetPartial([]knownvalue.Check{
							knownvalue.StringExact(string(permissionv1.GetConfig)),
							knownvalue.StringExact(string(permissionv1.GetMetadata)),
							knownvalue.StringExact(string(permissionv1.ListNamespaces)),
							knownvalue.StringExact(string(permissionv1.IncludeInList)),
							knownvalue.StringExact(string(permissionv1.ListDeletedTabulars)),
							knownvalue.StringExact(string(permissionv1.GetAllTasks)),
							knownvalue.StringExact(string(permissionv1.GetWarehouseEndpointStatistics)),
						}),
					),
				},
			},
		},
	})
}
