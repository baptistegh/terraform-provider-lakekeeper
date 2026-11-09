data "lakekeeper_projects" "all" {}

# Projects by name, for example to look up each organization's project.
locals {
  project_ids = { for p in data.lakekeeper_projects.all.projects : p.name => p.id }
}
