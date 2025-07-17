resource "lakekeeper_role_role_assignment" "data_analysts" {
  role_id     = "a4653498-1dd9-4f12-a2e4-1cc7d4023226"
  assignee_id = "cb6ee351-68ff-4299-87f2-876964f6d8dd"
  assignments = ["ownership"]
}
