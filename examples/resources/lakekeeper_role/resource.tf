resource "lakekeeper_project" "awesome" {
  name = "awesome"
}

resource "lakekeeper_user" "john_doe" {
  project_id  = lakekeeper_project.awesome.id
  name        = "toto-editors"
  description = "Here I can describe the role."
}
