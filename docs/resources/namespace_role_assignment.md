---
page_title: "lakekeeper_namespace_role_assignment Resource"
subcategory: ""
description: Manage a role's direct namespace permissions.
---

# lakekeeper_namespace_role_assignment

Manage all direct permissions for one namespace and role. Other principals are preserved.
Use one resource per namespace-role pair. Import existing permissions before managing them.
Requires the OpenFGA backend. Namespace grants also apply to descendants according to Lakekeeper's model.

```hcl
resource "lakekeeper_namespace_role_assignment" "readers" {
  project_id   = lakekeeper_project.example.id
  namespace_id = lakekeeper_namespace.example.namespace_id
  role_id      = lakekeeper_role.readers.role_id
  assignments  = ["select"]
}
```

Required attributes: `project_id`, `namespace_id`, `role_id` (UUID strings) and `assignments` (set of strings).
Assignments accept `ownership`, `pass_grants`, `manage_grants`, `describe`, `select`, `create`, and `modify`.
The computed `id` is `project_id/namespace_id/role_id` and is also the import identifier.

Removing a resource revokes its direct grants. An empty set also revokes them while retaining the Terraform resource.
Updates revoke obsolete permissions before adding new ones. A failed write can leave partial progress; refresh and reapply.
Failed reads and deletes return errors and preserve state. A missing permissions endpoint is an error, including with allow-all.

The namespace resource exposes its UUID as `namespace_id`. Set its optional `managed_access = true` to centralize descendant grants.
Leaving `managed_access` unset leaves that setting unchanged and avoids the OpenFGA-only managed-access endpoint.
