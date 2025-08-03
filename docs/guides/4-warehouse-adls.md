---
page_title: "Warehouse on Azure Data Lake Storage Gen2"
subcategory: "Manage Warehouses"
description: |-

---

# Set Up a Warehouse on ADLS

The ADLS storage support enables the warehouse to be stored on [Azure Data Lake Storage Gen2](https://docs.microsoft.com/en-us/azure/storage/common/storage-introduction).

The `adls` storage profile has the following required attributes:

* `account_name` (String) Name of the azure storage account.
* `filesystem` (String) Name of the adls filesystem, in blobstorage also known as container.

And the following optional attributes:

* `allow_alternative_protocols` (Boolean) Allow alternative protocols such as wasbs:// in locations. This is disabled by default. We do not recommend to use this setting except for migration.
* `authority_host` (String) The authority host to use for authentication. Defaults to `https://login.microsoftonline.com`.
* `host` (String) The host to use for the storage account. Defaults to `dfs.core.windows.net`.
* `key_prefix` (String) Subpath in the filesystem to use.
* `sas_token_validity_seconds` (Number) The validity of the sas token in seconds. Default is `3600`.

## Credentials

You must set credentials to access the storage. The `adls` storage profile support 3 different credentials:

* `client_credential` 
* `shared_access_key` 
* `azure_system_identity`

## Examples

### Client Credentials 

```terraform
resource "lakekeeper_warehouse" "aws" {
  project_id = data.lakekeeper_default_project.id
  name       = "aws"
  storage_profile = {
    adls = {
      account_name = "myaccount"
      filesystem   = "myfilesystem"
      credential = {
        client_credentials = {
          client_id     = "your-client-id"
          client_secret = "your-client-secret"
          tenant_id     = "your-tenant-id"
        }
      }
    }
  }
}
```

### Shared Access Key

```terraform
resource "lakekeeper_warehouse" "aws" {
  project_id = data.lakekeeper_default_project.id
  name       = "aws"
  storage_profile = {
    als = {
      account_name = "myaccount"
      filesystem   = "myfilesystem"
      credential = {
        shared_access_key = {
          key = "your-shared-access-key"
        }
      }
    }
  }
}
```

### Azure System Identity

```terraform
resource "lakekeeper_warehouse" "cloudflare" {
  project_id = data.lakekeeper_default_project.id
  name       = "aws"
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
```

