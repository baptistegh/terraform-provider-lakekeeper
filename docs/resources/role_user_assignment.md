---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lakekeeper_role_user_assignment Resource - terraform-provider-lakekeeper"
subcategory: ""
description: |-
  The lakekeeper_role_role_assignment resource allows to manage the lifecycle of a user assignement to a role.
  Upstream API: Lakekeeper REST API docs https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_role_assignments
---

# lakekeeper_role_user_assignment (Resource)

The `lakekeeper_role_role_assignment` resource allows to manage the lifecycle of a user assignement to a role.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/update_role_assignments)

## Example Usage

```terraform
resource "lakekeeper_role_user_assignment" "data_analysts" {
  role_id     = "a4653498-1dd9-4f12-a2e4-1cc7d4023226"
  user_id     = "cb6ee351-68ff-4299-87f2-876964f6d8dd"
  assignments = ["assignee"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `assignments` (Set of String) List of role assignments for this role. values can be `ownership` or `assignee`
- `role_id` (String) The ID of the role.
- `user_id` (String) The ID of the user to assign to this role.

### Read-Only

- `id` (String) The internal ID of this resource. In the form: <role_id>:<user_id>

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
# you can import user assignments to a role by the id in the form <role_id>/<user_id>
terraform import lakekeeper_role_role_assignment.data_analysts "a4653498-1dd9-4f12-a2e4-1cc7d4023226/cb6ee351-68ff-4299-87f2-876964f6d8dd"
```
