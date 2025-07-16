resource "lakekeeper_server_user_assignment" "john_doe" {
  user_id     = "oidc~91d18c8-1da4-471e-89f1-6e43eb4dcb38"
  assignments = ["operator"]
}
