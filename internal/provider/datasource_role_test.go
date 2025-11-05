//go:build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataLakekeeperRole_basic(t *testing.T) {

	project := testutil.CreateProject(t)
	defaultProject, _, err := testutil.TestLakekeeperClient.ProjectV1().GetDefault(t.Context())
	if err != nil {
		t.Fatalf("could not get default project, %v", err)
	}

	roleDefaultProject := testutil.CreateRole(t, defaultProject.ID)
	roleNewProject := testutil.CreateRole(t, project.ID)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "lakekeeper_default_project" "default" {}
				data "lakekeeper_role" "default" {
					project_id = data.lakekeeper_default_project.default.id
					role_id = "%s"
				}
				data "lakekeeper_role" "new" {
					role_id = "%s"
					project_id = "%s"
				}`, roleDefaultProject.ID, roleNewProject.ID, project.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// role in default project
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "id", "00000000-0000-0000-0000-000000000000/"+roleDefaultProject.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "role_id", roleDefaultProject.ID),
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "project_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "name", roleDefaultProject.Name),
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "description", *roleDefaultProject.Description),
					resource.TestCheckResourceAttr("data.lakekeeper_role.default", "created_at", roleDefaultProject.CreatedAt),

					// role in created project
					resource.TestCheckResourceAttr("data.lakekeeper_role.new", "id", fmt.Sprintf("%s/%s", project.ID, roleNewProject.ID)),
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
