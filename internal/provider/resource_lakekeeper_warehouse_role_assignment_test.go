//go:build acceptance
// +build acceptance

package provider

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperWarehouseRoleAssignment_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	keyPrefix := fmt.Sprintf("key-prefix-%d", rand.Int())
	warehouse := testutil.CreateWarehouse(t, project.ID, keyPrefix)

	role := testutil.CreateRole(t, project.ID)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_warehouse_role_assignment" "test" {
						warehouse_id = "%s"
						role_id = "%s"
						assignments = ["ownership"]
					}
				`, warehouse.ID, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "id", warehouse.ID+":"+role.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "warehouse_id", warehouse.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.0", "ownership"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_warehouse_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_warehouse_role_assignment" "test" {
						warehouse_id = "%s"
						role_id = "%s"
						assignments = ["modify"]
					}
				`, warehouse.ID, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "id", warehouse.ID+":"+role.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "warehouse_id", warehouse.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.0", "modify"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_warehouse_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_warehouse_role_assignment" "test" {
						warehouse_id = "%s"
						role_id = "%s"
						assignments = ["ownership", "manage_grants", "create"]
					}
				`, warehouse.ID, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "id", warehouse.ID+":"+role.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "warehouse_id", warehouse.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.#", "3"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.0", "create"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.1", "manage_grants"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.2", "ownership"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_warehouse_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// delete all assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_warehouse_role_assignment" "test" {
						warehouse_id = "%s"
						role_id = "%s"
						assignments = []
					}
				`, warehouse.ID, role.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "id", warehouse.ID+":"+role.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "warehouse_id", warehouse.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "role_id", role.ID),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.#", "0"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_warehouse_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLakekeeperWarehouseRoleAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_warehouse_role_assignment" {
			continue
		}

		warehouseID, roleID := splitInternalID(types.StringValue(rs.Primary.ID))

		assignments, _, err := testutil.TestLakekeeperClient.PermissionV1().WarehousePermission().GetAssignments(context.Background(), warehouseID, nil)
		if err != nil {
			return fmt.Errorf("could not list project assignments to check destroy, %w", err)
		}

		for _, v := range assignments.Assignments {
			if v.Assignee.Value == roleID && v.Assignee.Type == permissionv1.RoleType {
				return fmt.Errorf("project assignment still exists")
			}
		}
	}

	return nil
}
