resource "lakekeeper_project_user_assignment" "john_doe" {
  project_id  = "5653bd71-1f1c-4a2c-913b-fbd92d6c1157"
  user_id     = "oidc~91d18c8-1da4-471e-89f1-6e43eb4dcb38"
  assignments = ["select", "modify"]
}
