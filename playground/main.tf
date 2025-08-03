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


resource "lakekeeper_user" "anna" {
  id        = "oidc~d223d88c-85b6-4859-b5c5-27f3825e47f6"
  name      = "Anna Cold"
  email     = "anna@example.com"
  user_type = "human"
}

resource "lakekeeper_user" "peter" {
  id        = "oidc~cfb55bf6-fcbb-4a1e-bfec-30c6649b52f8"
  name      = "Peter Cold"
  email     = "peter@example.com"
  user_type = "human"
}

resource "lakekeeper_project_user_assignment" "default_anna" {
  project_id  = data.lakekeeper_default_project.default.id
  user_id     = lakekeeper_user.anna.id
  assignments = ["project_admin"]
}

// To add perter as a Super Admin, uncomment this section
// resource "lakekeeper_server_user_assignment" "default_peter" {
//   user_id     = lakekeeper_user.peter.id
//   assignments = ["operator"]
// }

resource "lakekeeper_warehouse" "s3" {
  name       = "test-warehouse-s3"
  project_id = data.lakekeeper_default_project.default.id
  storage_profile = {
    s3 = {
      bucket          = "testbucket"
      region          = "us-west-1"
      sts_enabled     = true
      assume_role_arn = "arn:aws:iam::123456789012:role/MyDeploymentRole"
      credential = {
        access_key = {
          access_key_id     = "AKIAEXAMPLE1234567890"
          secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
        }
      }
    }
  }
}

resource "lakekeeper_warehouse" "gcs" {
  name       = "test-warehouse-gcs"
  project_id = data.lakekeeper_default_project.default.id
  storage_profile = {
    gcs = {
      bucket = "testmybucket"
      credential = {
        service_account_key = {
          key = file("key.json")
        }
      }
    }
  }
  delete_profile = {
    type               = "soft"
    expiration_seconds = 3600
  }
}


resource "lakekeeper_warehouse" "adls" {
  name       = "test-warehouse-adls"
  project_id = data.lakekeeper_default_project.default.id
  storage_profile = {
    adls = {
      account_name = "myaccount"
      filesystem   = "myfilesystem"
      credential = {
        azure_system_identity = {}
      }
    }
  }
}

resource "lakekeeper_role" "read" {
  project_id  = data.lakekeeper_default_project.default.id
  name        = "read-role"
  description = "this role gives select permissions on all the warehouses"
}

resource "lakekeeper_role_user_assignment" "read_peter" {
  role_id     = lakekeeper_role.read.role_id
  user_id     = lakekeeper_user.peter.id
  assignments = ["assignee"]
}

resource "lakekeeper_role" "write" {
  project_id  = data.lakekeeper_default_project.default.id
  name        = "write-role"
  description = "this role gives create and modify permissions on test-warehouse-s3"
}

resource "lakekeeper_role_user_assignment" "write_peter" {
  role_id     = lakekeeper_role.write.role_id
  user_id     = lakekeeper_user.peter.id
  assignments = ["assignee"]
}

# Adding select and describe permissions for the role to all the warehouses
resource "lakekeeper_warehouse_role_assignment" "all_read" {
  for_each = tomap({
    s3   = lakekeeper_warehouse.s3.warehouse_id,
    gcs  = lakekeeper_warehouse.gcs.warehouse_id,
    adls = lakekeeper_warehouse.adls.warehouse_id
  })

  warehouse_id = each.value
  role_id      = lakekeeper_role.read.role_id
  assignments  = ["select", "describe"]
}

resource "lakekeeper_warehouse_role_assignment" "s3_write" {
  warehouse_id = lakekeeper_warehouse.s3.warehouse_id
  role_id      = lakekeeper_role.write.role_id
  assignments  = ["create", "modify"]
}