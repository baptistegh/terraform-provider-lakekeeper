---
page_title: "Overview"
subcategory: "Manage Warehouses"
description: |-

---

# Manage Warehouses

This guide explains how to manage warehouses in Lakekeeper using Terraform.

3 storage profiles are [supported by Lakekeeper](https://github.com/lakekeeper/lakekeeper?tab=readme-ov-file#storage-profile-support):

* S3
* Google Cloud Storage
* Azure Data Lake Storage

All the warehouse resource object have the following attributes:

* `name` (String) Name of the warehouse to create. Must be unique within a project and may not contain "/"
* `project_id` (String) The project ID to which the warehouse belongs.
* `active` (Boolean) Whether the warehouse is active. Default is `true`.
* `managed_access` (Boolean) Whether the managed access is configured on this warehouse. Default is `false`.
* `protected` (Boolean) Whether the warehouse is protected from being deleted. Default is `false`.

The selected storage profile must be indicated in the `storage_profile` attribute.

You also can [set up a soft delete profile](#soft-delete-profile) with `delete_profile`.

## Storage Profiles

You can have more informations on all the storage profiles in the dedicated pages:

* [S3](./2-warehouse-s3.md.tmpl) 
* [Google Cloud Storage](./3-warehouse-gcs.md.tmpl) 
* [Azure Data Lake Storage Gen2](./4-warehouse-adls.md.tmpl) 

<a id="soft-delete-profile"></a>
## Enabling Soft Delete

You can enable the soft delete profile with the field `delete_profile`. By default, `hard` delete profile is used.

You control the expiration with `delete_profile.expiration_seconds`.

```terraform
resource "lakekeeper_warehouse" "cloudflare" {
  project_id = data.lakekeeper_default_project.id
  name       = "aws"
  storage_profile = {
    s3 = {
      region       = "us-east-1"
      bucket       = "mybucket"
      sts_enabled  = false
      credential = {
        cloudflare_r2 = {
          access_key_id = "AKIAEXAMPLE1234567890"
          secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
          account_id = "your-account-id"
          token = "your-token"
        }
      }
    }
  }
  delete_profile = {
    type               = "soft"
    expiration_seconds = 3600 // default is 3600
  }
}
```