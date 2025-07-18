---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lakekeeper_server_info Data Source - terraform-provider-lakekeeper"
subcategory: ""
description: |-
  The lakekeeper_server_info data source retrieves information about a lakekeeper instance.
  Upstream API: Lakekeeper REST API docs https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/get_server_info
---

# lakekeeper_server_info (Data Source)

The `lakekeeper_server_info` data source retrieves information about a lakekeeper instance.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/get_server_info)

## Example Usage

```terraform
data "lakekeeper_server_info" "infos" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `authz_backend` (String) Authorization backend configured
- `aws_system_identities_enabled` (Boolean)
- `azure_system_identities_enabled` (Boolean)
- `bootstrapped` (Boolean) True if the server has been bootstrapped
- `default_project_id` (String) The default project ID
- `gcp_system_identities_enabled` (Boolean)
- `queues` (List of String)
- `server_id` (String) The ID of the server
- `version` (String) The current version of the running server
