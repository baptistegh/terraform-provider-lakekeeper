terraform {
  required_providers {
    lakekeeper = {
      source = "baptistegh/lakekeeper"
    }
  }
}

provider "lakekeeper" {
  endpoint      = "http://localhost:8181"
  auth_url      = "http://localhost:30080/realms/iceberg/protocol/openid-connect/token"
  client_id     = "lakekeeper-admin"
  client_secret = "KNjaj1saNq5yRidVEMdf1vI09Hm0pQaL"
}

data "lakekeeper_server_info" "this" {}