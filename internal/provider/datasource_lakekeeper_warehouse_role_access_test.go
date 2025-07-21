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
					Assignment: permissionv1.ModifyWarehouseAssignment,
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
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "id", fmt.Sprintf("%s:%s", warehouse.ID, role.ID)),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "warehouse_id", warehouse.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "role_id", role.ID),

					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.#", "11"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.0", "activate"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.1", "deactivate"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.2", "delete"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.3", "get_config"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.4", "get_metadata"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.5", "include_in_list"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.6", "list_deleted_tabulars"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.7", "list_namespaces"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.8", "modify_storage"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.9", "modify_storage_credential"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_role_access.foo", "allowed_actions.10", "rename"),
				),
			},
		},
	})
}
