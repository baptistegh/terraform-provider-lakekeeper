---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lakekeeper_warehouse_role_access Data Source - terraform-provider-lakekeeper"
subcategory: ""
description: |-
  The lakekeeper_warehouse_role_access data source retrieves the accesses a role can have on a warehouse.
  Upstream API: Lakekeeper REST API docs https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_access_by_id
---

# lakekeeper_warehouse_role_access (Data Source)

The `lakekeeper_warehouse_role_access` data source retrieves the accesses a role can have on a warehouse.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_access_by_id)

## Example Usage

```terraform
data "lakekeeper_warehouse_role_access" "foo" {
  warehouse_id = "116d3ba8-1c38-4548-b39c-aaed6c325406"
  role_id      = "9d8fd24b-55a7-473f-ad7f-66d4b8e8d6ae"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `role_id` (String) ID of the role.
- `warehouse_id` (String) ID of the warehouse.

### Read-Only

- `allowed_actions` (Set of String) List of the role's allowed actions on the warehouse. The possible values are `create_namespace` `delete` `modify_storage` `modify_storage_credential` `get_config` `get_metadata` `list_namespaces` `include_in_list` `deactivate` `activate` `rename` `list_deleted_tabulars` `read_assignments` `grant_create` `grant_describe` `grant_modify` `grant_select` `grant_pass_grants` `grant_manage_grants` `change_ownership`
- `id` (String) The internal ID of this data source, in the form <warehouse_id>:<role_id>.
