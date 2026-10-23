resource "lakekeeper_namespace_role_assignment" "readers" {
  project_id   = lakekeeper_project.example.id
  namespace_id = lakekeeper_namespace.example.namespace_id
  role_id      = lakekeeper_role.readers.role_id
  assignments  = ["select"]
}
