//go:build acceptance

package provider

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperWarehouseRoleAssignment_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	keyPrefix := fmt.Sprintf("key-prefix-%d", rand.Int())
	warehouse := testutil.CreateWarehouse(t, project.ProjectId, keyPrefix)

	role := testutil.CreateRole(t, project.ProjectId)

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
				`, warehouse.Id, role.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "id", warehouse.Id+"/"+role.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "warehouse_id", warehouse.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "role_id", role.Id),
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
				`, warehouse.Id, role.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "id", warehouse.Id+"/"+role.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "warehouse_id", warehouse.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "role_id", role.Id),
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
						assignments = ["ownership", "manage_grants", "create", "manage_tags"]
					}
				`, warehouse.Id, role.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "id", warehouse.Id+"/"+role.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "warehouse_id", warehouse.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "role_id", role.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.#", "4"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.0", "create"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.1", "manage_grants"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.2", "manage_tags"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "assignments.3", "ownership"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_warehouse_role_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Check some basic attribute validation
			{
				Config: fmt.Sprintf(`
					resource "lakekeeper_warehouse_role_assignment" "test" {
						warehouse_id = "%s/%s"
						role_id = "%s"
						assignments = ["ownership"]
					}
				`, project.ProjectId, warehouse.Id, role.Id),
				ExpectError: regexp.MustCompile("Attribute warehouse_id must be a warehouse UUID and NOT include the project\nUUID"),
			},
			{
				Config: fmt.Sprintf(`
					resource "lakekeeper_warehouse_role_assignment" "test" {
						warehouse_id = "%s"
						role_id = "%s/%s"
						assignments = ["ownership"]
					}
				`, warehouse.Id, project.ProjectId, role.Id),
				ExpectError: regexp.MustCompile("Attribute role_id must be a role UUID and NOT include the project UUID"),
			},
			// delete all assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_warehouse_role_assignment" "test" {
						warehouse_id = "%s"
						role_id = "%s"
						assignments = []
					}
				`, warehouse.Id, role.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "id", warehouse.Id+"/"+role.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "warehouse_id", warehouse.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_role_assignment.test", "role_id", role.Id),
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

		relations, err := warehouseAssignmentScope(testutil.TestLakekeeperClient, warehouseID).relations(context.Background(), permissions.PrincipalRole, roleID)
		if err != nil {
			return fmt.Errorf("could not list project assignments to check destroy, %w", err)
		}

		if len(relations) > 0 {
			return fmt.Errorf("project assignment still exists")
		}
	}

	return nil
}
