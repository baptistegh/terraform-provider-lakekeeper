# id must be the same as in the JWT token.
# prefixed by the IDP (oidc~ or kubernetes~)

resource "lakekeeper_user" "john_doe" {
  id        = "oidc~91d18c8-1da4-471e-89f1-6e43eb4dcb38"
  name      = "John Doe"
  email     = "john.doe@lakekeeper.io"
  user_type = "human"
}
