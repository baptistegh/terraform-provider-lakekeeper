//go:build acceptance

package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/testutil"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLakekeeperUser_basic(t *testing.T) {

	rID := acctest.RandomWithPrefix("oidc~")
	rName := acctest.RandString(8)
	rUpdatedName := acctest.RandString(12)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLakekeeperUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`				
				resource "lakekeeper_user" "foo" {
				  id = "%s"
				  name = "%s"
				  email = "%s@local.local"
				  user_type = "human"
				}
				`, rID, rName, rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "id", rID),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "name", rName),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "email", rName+"@local.local"),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "user_type", "human"),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "last_updated_with", "create-endpoint"),
					resource.TestCheckResourceAttrSet("lakekeeper_user.foo", "created_at"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_user.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update User
			{
				Config: fmt.Sprintf(`				
				resource "lakekeeper_user" "foo" {
				  id = "%s"
				  name = "%s"
				  email = "%s@local.local"
				  user_type = "application"
				}
				`, rID, rUpdatedName, rUpdatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "id", rID),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "name", rUpdatedName),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "email", rUpdatedName+"@local.local"),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "user_type", "application"),
					resource.TestCheckResourceAttr("lakekeeper_user.foo", "last_updated_with", "create-endpoint"),
					resource.TestCheckResourceAttrSet("lakekeeper_user.foo", "updated_at"),
					resource.TestCheckResourceAttrSet("lakekeeper_user.foo", "created_at"),
				),
			},
			// Verify import
			{
				ResourceName:      "lakekeeper_user.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLakekeeperUser_invalidID(t *testing.T) {

	rName := acctest.RandString(8)
	rID := acctest.RandString(12)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create user with wrong ID
			{
				Config: fmt.Sprintf(`				
				resource "lakekeeper_user" "failed" {
				  id = "%s"
				  name = "%s"
				  email = "%s@local.local"
				  user_type = "human"
				}
				`, rID, rName, rName),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
		},
	})
}

func TestAccLakekeeperUser_invalidType(t *testing.T) {

	rID := acctest.RandomWithPrefix("oidc~")
	rName := acctest.RandString(8)
	rType := acctest.RandString(8)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create user with wrong type
			{
				Config: fmt.Sprintf(`				
				resource "lakekeeper_user" "failed" {
				  id = "%s"
				  name = "%s"
				  email = "%s@local.local"
				  user_type = "%s"
				}
				`, rID, rName, rName, rType),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			// Create user without type
			{
				Config: fmt.Sprintf(`			
				resource "lakekeeper_user" "failed" {
				  id = "%s"
				  name = "%s"
				  email = "%s@local.local"
				}
				`, rID, rName, rName),
				ExpectError: regexp.MustCompile("required"),
			},
		},
	})
}

func testAccCheckLakekeeperUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lakekeeper_user" {
			continue
		}

		_, _, err := testutil.TestLakekeeperClient.UserV1().Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("User with id %s still exists", rs.Primary.ID)
		}
	}
	return nil
}
