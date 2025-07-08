resource "lakekeeper_project" "bi" {
  name = "bi"
}

resource "lakekeeper_warehouse" "aws" {
  project_id = lakekeeper_project.bi.id
  name       = "aws"
  protected  = false
  active     = true
  storage_profile = {
    type   = "s3"
    region = "us-east-1"
  }
  delete_profile = {
    type               = "soft"
    expiration_seconds = 3600
  }
  storage_credential = {
    type              = "s3_access_key"
    access_key_id     = "AKIAEXAMPLE1234567890"
    secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  }
}
