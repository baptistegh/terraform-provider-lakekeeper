resource "lakekeeper_project" "bi" {
  name = "bi"
}

# create a warehouse: S3 storage with Access Key
resource "lakekeeper_warehouse" "aws" {
  project_id     = lakekeeper_project.bi.id
  name           = "aws"
  protected      = false
  active         = true
  managed_access = true
  storage_profile = {
    s3 = {
      region = "us-east-1"
      bucket = "mybucket"
      credential = {
        access_key = {
          access_key_id     = "AKIAEXAMPLE1234567890"
          secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
        }
      }
    }
  }
  delete_profile = {
    type               = "soft"
    expiration_seconds = 3600
  }
}

# Create a warehouse: GCS with Service Account Key
resource "lakekeeper_warehouse" "gcs" {
  project_id     = lakekeeper_project.bi.id
  name           = "gcs"
  protected      = false
  active         = true
  managed_access = false
  storage_profile = {
    gcs = {
      bucket = "mybucket"
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

# Create a warehouse: ADLS with Azure System Identity
resource "lakekeeper_warehouse" "adls" {
  project_id     = lakekeeper_project.bi.id
  name           = "adls"
  protected      = false
  active         = true
  managed_access = false
  storage_profile = {
    adls = {
      account_name = "myaccount"
      filesystem   = "fs"
      credential = {
        azure_system_identity = {}
      }
    }
  }
}
