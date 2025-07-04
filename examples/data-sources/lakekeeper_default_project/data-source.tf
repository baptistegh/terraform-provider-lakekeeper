data "lakekeeper_default_project" "default" {}

resource "lakekeeper_role" "admin" {
  name       = "admin"
  project_id = data.lakekeeper_default_project.default.id
}