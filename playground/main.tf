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

data "lakekeeper_default_project" "default" {}

resource "lakekeeper_user" "peter" {
  id        = "oidc~cfb55bf6-fcbb-4a1e-bfec-30c6649b52f8"
  name      = "Peter Cold"
  email     = "peter@example.com"
  user_type = "human"
}

resource "lakekeeper_user" "anna" {
  id        = "oidc~d223d88c-85b6-4859-b5c5-27f3825e47f6"
  name      = "Anna Cold"
  email     = "anna@example.com"
  user_type = "human"
}

resource "lakekeeper_project_user_assignment" "default_anna" {
  project_id  = data.lakekeeper_default_project.default.id
  user_id     = lakekeeper_user.anna.id
  assignments = ["project_admin"]
}

resource "lakekeeper_role" "select" {
  project_id  = data.lakekeeper_default_project.default.id
  name        = "test-role"
  description = "this role gives select permissions on test-warehouse"
}

resource "lakekeeper_warehouse" "gcs" {
  name       = "test-warehouse"
  project_id = data.lakekeeper_default_project.default.id
  storage_profile = {
    type   = "gcs"
    bucket = "testbucket"
  }
  storage_credential = {
    type = "gcs_gcp_system_identity"
  }
}

resource "lakekeeper_warehouse_role_assignment" "wh_select" {
  warehouse_id = lakekeeper_warehouse.gcs.warehouse_id
  role_id      = lakekeeper_role.select.role_id
  assignments  = ["select", "describe"]
}

resource "lakekeeper_role_user_assignment" "select_peter" {
  role_id     = lakekeeper_role.select.role_id
  user_id     = lakekeeper_user.peter.id
  assignments = ["assignee"]
}