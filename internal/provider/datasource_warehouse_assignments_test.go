//go:build acceptance

package provider

import (
	"fmt"
	"testing"

	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperWarehouseAssignments_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	user := testutil.CreateUser(t, fmt.Sprintf("oidc~%s", acctest.RandString(8)))
	role := testutil.CreateRole(t, project.ProjectId)

	keyPrefix := acctest.RandString(8)
	warehouse := testutil.CreateWarehouse(t, project.ProjectId, keyPrefix)

	// assignment 2
	if err := grantAssignments(t.Context(), warehouseAssignmentScope(testutil.TestLakekeeperClient, warehouse.Id),
		testAssignment{permissions.PrincipalUser, user.Id, string(managementv1.WAREHOUSERELATION_OWNERSHIP)},
		testAssignment{permissions.PrincipalRole, role.Id, string(managementv1.WAREHOUSERELATION_DESCRIBE)},
	); err != nil {
		t.Fatalf("could not create warehouse assignments, %v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_warehouse_assignments" "foo" {
					id = "%s"
				}`, warehouse.Id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "id", warehouse.Id),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "assignments.#", "3"), // also have the admin user

					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "assignments.0.assignee_id", testutil.DefaultUserID),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "assignments.0.assignee_type", "user"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "assignments.0.assignment", string(managementv1.WAREHOUSERELATION_OWNERSHIP)),

					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "assignments.1.assignee_id", user.Id),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "assignments.1.assignee_type", "user"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "assignments.1.assignment", string(managementv1.WAREHOUSERELATION_OWNERSHIP)),

					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "assignments.2.assignee_id", role.Id),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "assignments.2.assignee_type", "role"),
					resource.TestCheckResourceAttr("data.lakekeeper_warehouse_assignments.foo", "assignments.2.assignment", string(managementv1.WAREHOUSERELATION_DESCRIBE)),
				),
			},
		},
	})
}
