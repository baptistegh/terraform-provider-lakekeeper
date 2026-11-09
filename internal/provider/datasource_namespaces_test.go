//go:build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperNamespaces_basic(t *testing.T) {
	project := testutil.CreateProject(t)
	warehouse := testutil.CreateWarehouse(t, project.ID, acctest.RandString(8))
	prefix := acctest.RandString(8)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "lakekeeper_namespace" "a" {
				  project_id     = "%[1]s"
				  warehouse_name = "%[2]s"
				  name           = "%[3]sa"
				}
				resource "lakekeeper_namespace" "b" {
				  project_id     = "%[1]s"
				  warehouse_name = "%[2]s"
				  name           = "%[3]sb"
				}
				resource "lakekeeper_namespace" "child" {
				  project_id     = "%[1]s"
				  warehouse_name = "%[2]s"
				  name           = "${lakekeeper_namespace.a.name}.c"
				}
				data "lakekeeper_namespaces" "top" {
				  warehouse_id = "%[4]s"
				  depends_on   = [lakekeeper_namespace.a, lakekeeper_namespace.b, lakekeeper_namespace.child]
				}
				data "lakekeeper_namespaces" "children" {
				  warehouse_id = "%[4]s"
				  parent       = ["%[3]sa"]
				  depends_on   = [lakekeeper_namespace.child]
				}
				`, project.ID, warehouse.Name, prefix, warehouse.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lakekeeper_namespaces.top", "namespaces.#", "2"),
					resource.TestCheckResourceAttr("data.lakekeeper_namespaces.top", "namespaces.0.parts.#", "1"),
					resource.TestCheckResourceAttr("data.lakekeeper_namespaces.top", "namespaces.0.parts.0", prefix+"a"),
					resource.TestCheckResourceAttrPair("data.lakekeeper_namespaces.top", "namespaces.0.id", "lakekeeper_namespace.a", "namespace_id"),
					resource.TestCheckResourceAttr("data.lakekeeper_namespaces.top", "namespaces.1.parts.0", prefix+"b"),
					resource.TestCheckResourceAttr("data.lakekeeper_namespaces.children", "namespaces.#", "1"),
					resource.TestCheckResourceAttr("data.lakekeeper_namespaces.children", "namespaces.0.parts.#", "2"),
					resource.TestCheckResourceAttr("data.lakekeeper_namespaces.children", "namespaces.0.parts.0", prefix+"a"),
					resource.TestCheckResourceAttr("data.lakekeeper_namespaces.children", "namespaces.0.parts.1", "c"),
					resource.TestCheckResourceAttrPair("data.lakekeeper_namespaces.children", "namespaces.0.id", "lakekeeper_namespace.child", "namespace_id"),
				),
			},
		},
	})
}
