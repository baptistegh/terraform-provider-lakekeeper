resource "lakekeeper_project_role_assignment" "data_analysts" {
  project_id  = "a4653498-1dd9-4f12-a2e4-1cc7d4023226"
  role_id     = "cb6ee351-68ff-4299-87f2-876964f6d8dd"
  assignments = ["project_admin"]
}
