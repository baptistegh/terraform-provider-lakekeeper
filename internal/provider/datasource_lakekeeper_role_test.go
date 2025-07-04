//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperRole_basic(t *testing.T) {

	project := testutil.CreateProject(t)

	roleDefaultProject := testutil.CreateRole(t, "")
	roleNewProject := testutil.CreateRole(t, project.ID)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_role" "default" {
					name = "%s"
				}
				data "lakekeeper_role" "new" {
					name = "%s"
					project_id = "%s"
				}`, roleDefaultProject.Name, roleNewProject.Name, project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// role in default project
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "id", "00000000-0000-0000-0000-000000000000:"+roleDefaultProject.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "role_id", roleDefaultProject.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "project_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "name", roleDefaultProject.Name),
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "description", *roleDefaultProject.Description),
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "created_at", roleDefaultProject.CreatedAt),

					// role in created project
					resource.TestCheckResourceAttr("data.lakekeeper_role.new", "id", fmt.Sprintf("%s:%s", project.ID, roleNewProject.ID)),
					resource.TestCheckResourceAttr("data.lakekeeper_role.new", "role_id", roleNewProject.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_role.new", "project_id", project.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_role.new", "name", roleNewProject.Name),
					resource.TestCheckResourceAttr("data.lakekeeper_role.new", "description", *roleNewProject.Description),
					resource.TestCheckResourceAttr("data.lakekeeper_role.new", "created_at", roleNewProject.CreatedAt),
				),
			},
		},
	})
}
