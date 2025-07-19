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

resource "lakekeeper_project" "main" {
  name = "Main Project"
}

resource "lakekeeper_project" "second" {
  name = "Main Project 2"
}

resource "lakekeeper_role" "access" {
  project_id = lakekeeper_project.second.id
  name       = "access-warehouse"
}

resource "lakekeeper_project_user_assignment" "main_peter" {
  project_id  = lakekeeper_project.main.id
  user_id     = lakekeeper_user.peter.id
  assignments = ["project_admin"]
}

resource "lakekeeper_warehouse" "warehouse" {
  name       = "test-warehouse"
  project_id = lakekeeper_project.main.id
  storage_profile = {
    type   = "gcs"
    bucket = "testbucket"
  }
  storage_credential = {
    type = "gcs_gcp_system_identity"
  }
}

resource "lakekeeper_project_role_assignment" "test" {
  project_id  = lakekeeper_project.second.id
  role_id     = lakekeeper_role.access.role_id
  assignments = ["project_admin"]
}

resource "lakekeeper_server_user_assignment" "admin" {
  user_id     = lakekeeper_user.peter.id
  assignments = ["admin"]
}