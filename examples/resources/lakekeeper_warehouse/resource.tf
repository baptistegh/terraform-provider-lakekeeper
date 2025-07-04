data "lakekeeper_project" "my_project" {
  name = "my_project"
}

resource "lakekeeper_warehouse" "aws" {
  project_id = data.lakekeeper_project.my_project.id
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
