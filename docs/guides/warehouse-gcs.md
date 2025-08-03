---
page_title: "Warehouse on Google Cloud Storage"
subcategory: "Manage Warehouses"
description: |-

---


# Set Up a Warehouse on GCS

The GCS storage support enables the warehouse to be stored on [Google Cloud Storage](https://cloud.google.com/storage).

The `gcs` storage profile has the following required attributes:

* `bucket` (String) The bucket name.

And the following optional attributes:

* `key_prefix` (String) Subpath in the filesystem to use.


## Credentials

You must set credentials to access the storage. The `gcs` storage profile support 2 different credential:

* `service_account_key` 
* `gcp_system_identity`

## Examples

### Service Account Key

You can provide the key by a file or a json encoded string.

```terraform
resource "lakekeeper_warehouse" "aws" {
  project_id = data.lakekeeper_default_project.id
  name       = "aws"
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
}
```

### Google Cloud Platform System Identity

```terraform
resource "lakekeeper_warehouse" "aws" {
  project_id = data.lakekeeper_default_project.id
  name       = "aws"
  storage_profile = {
    gcs = {
      bucket = "mybucket"
      credential = {
        gcp_system_identity = {}
      }
    }
  }
}
```
