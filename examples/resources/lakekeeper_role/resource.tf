resource "lakekeeper_project" "awesome" {
  name = "awesome"
}

resource "lakekeeper_role" "editors" {
  project_id  = lakekeeper_project.awesome.id
  name        = "editors"
  description = "Here I can describe the role."
}
