---
page_title: "Warehouse on S3"
subcategory: "Manage Warehouses"
description: |-

---


# Set Up a Warehouse on S3

The S3 storage support enables the warehouse to be stored on [AWS S3](https://aws.amazon.com/s3/) or any other S3 compatible storage (eg. [MinIO](https://min.io/)).

[Cloudflare R2](https://developers.cloudflare.com/r2/) is also supported with this storage profile.

The s3 storage profile has the following required attributes:

* `bucket` (String) The bucket name for the storage profile.
* `region` (String) Region to use for S3 requests.
* `sts_enabled` (Boolean) Whether to enable STS for S3 storage profile. Required if the storage type is `s3`. If enabled, the `sts_role_arn` or `assume_role_arn` must be provided.

And the following optional attributes:

* `endpoint` (String) Optional endpoint to use for S3 requests, if not provided the region will be used to determine the endpoint. If both region and endpoint are provided, the endpoint will be used. Example: `http://s3-de.my-domain.com:9000`
* `path_style_access` (Boolean) Path style access for S3 requests. If the underlying S3 supports both, we recommend to not set path_style_access.
* `flavor` (String) S3 flavor to use. Defaults to `aws`.
* `key_prefix` (String) Subpath in the filesystem to use.
* `allow_alternative_protocols` (Boolean) Allow `s3a://`, `s3n://`, `wasbs://` in locations. This is disabled by default. We do not recommend to use this setting except for migration.
* `assume_role_arn` (String) Optional ARN to assume when accessing the bucket from Lakekeeper for S3 storage profile
* `aws_kms_key_arn` (String) ARN of the KMS key used to encrypt the S3 bucket, if any.  
* `push_s3_delete_disabled` (Boolean) Controls whether the `s3.delete-enabled=false` flag is sent to clients.
* `remote_signing_url_style` (String) S3 URL style detection mode for remote signing. One of `auto`, `path-style`, `virtual-host`. Default: `auto`.
* `sts_role_arn` (String)
* `sts_token_validity_seconds` (Number) The validity of the STS tokens in seconds. Default is `3600`.

## Credentials

You must set credentials to access the storage. The `s3` storage profile support 3 different credential types:

* `access_key` 
* `aws_system_identity` 
* `cloudflare_r2`

## Examples

### Access Key 

```terraform
resource "lakekeeper_warehouse" "aws" {
  project_id = data.lakekeeper_default_project.id
  name       = "aws"
  storage_profile = {
    s3 = {
      region       = "us-east-1"
      bucket       = "mybucket"
      sts_enabled  = true
      sts_role_arn = "arn:aws:iam::123456789012:role/AssumeRole"
      credential = {
        access_key = {
          access_key_id     = "AKIAEXAMPLE1234567890"
          secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
        }
      }
    }
  }
}
```

### Amazon Web Services System Identity

```terraform
resource "lakekeeper_warehouse" "aws" {
  project_id = data.lakekeeper_default_project.id
  name       = "aws"
  storage_profile = {
    s3 = {
      region       = "us-east-1"
      bucket       = "mybucket"
      sts_enabled  = true
      sts_role_arn = "arn:aws:iam::123456789012:role/AssumeRole"
      credential = {
        aws_system_identity = {
          external_id = "1234567890"
        }
      }
    }
  }
}
```

### Cloudflare R2

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
}
```

