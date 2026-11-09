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
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperWarehouseUserAssignment_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	keyPrefix := fmt.Sprintf("key-prefix-%d", rand.Int())
	warehouse := testutil.CreateWarehouse(t, project.ProjectId, keyPrefix)

	userID := fmt.Sprintf("oidc~%s", acctest.RandString(32))
	user := testutil.CreateUser(t, userID)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperWarehouseUserAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_warehouse_user_assignment" "test" {
						warehouse_id = "%s"
						user_id = "%s"
						assignments = ["ownership"]
					}
				`, warehouse.Id, user.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "id", warehouse.Id+"/"+user.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "warehouse_id", warehouse.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "user_id", user.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "assignments.0", "ownership"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_warehouse_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// change one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_warehouse_user_assignment" "test" {
						warehouse_id = "%s"
						user_id = "%s"
						assignments = ["select"]
					}
				`, warehouse.Id, user.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "id", warehouse.Id+"/"+user.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "warehouse_id", warehouse.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "user_id", user.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "assignments.0", "select"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_warehouse_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// add one assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_warehouse_user_assignment" "test" {
						warehouse_id = "%s"
						user_id = "%s"
						assignments = ["select", "describe", "write"]
					}
				`, warehouse.Id, user.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "id", warehouse.Id+"/"+user.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "warehouse_id", warehouse.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "user_id", user.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "assignments.#", "3"),
					resource.TestCheckTypeSetElemAttr("lakekeeper_warehouse_user_assignment.test", "assignments.*", "write"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "assignments.0", "describe"),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "assignments.1", "select"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_warehouse_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Check some basic attribute validation
			{
				Config: fmt.Sprintf(`
					resource "lakekeeper_warehouse_user_assignment" "test" {
						warehouse_id = "%s/%s"
						user_id = "%s"
						assignments = ["ownership"]
					}
				`, project.ProjectId, warehouse.Id, user.Id),
				ExpectError: regexp.MustCompile("Attribute warehouse_id must be a warehouse UUID and NOT include the project\nUUID"),
			},
			// delete all assignments
			{
				Config: fmt.Sprintf(`				
					resource "lakekeeper_warehouse_user_assignment" "test" {
						warehouse_id = "%s"
						user_id = "%s"
						assignments = []
					}
				`, warehouse.Id, user.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "id", warehouse.Id+"/"+user.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "warehouse_id", warehouse.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "user_id", user.Id),
					resource.TestCheckResourceAttr("lakekeeper_warehouse_user_assignment.test", "assignments.#", "0"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_warehouse_user_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLakekeeperWarehouseUserAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_warehouse_user_assignment" {
			continue
		}

		warehouseID, userID := splitInternalID(types.StringValue(rs.Primary.ID))

		relations, err := warehouseAssignmentScope(testutil.TestLakekeeperClient, warehouseID).relations(context.Background(), permissions.PrincipalUser, userID)
		if err != nil {
			return fmt.Errorf("could not list project assignments to check destroy, %w", err)
		}

		if len(relations) > 0 {
			return fmt.Errorf("project assignment still exists")
		}
	}

	return nil
}
