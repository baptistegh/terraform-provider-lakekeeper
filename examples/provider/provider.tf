provider "lakekeeper" {
  endpoint      = "http://localhost:8181"
  auth_url      = "http://localhost:30080/realms/iceberg/protocol/openid-connect/token"
  client_id     = "lakekeeper-admin"
  client_secret = "<redacted>"
  scopes        = ["lakekeeper"]
}
