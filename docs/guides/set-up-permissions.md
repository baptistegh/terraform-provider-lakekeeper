---
page_title: "Set Up Permissions"
subcategory: ""
description: |-

---

# Set Up Permissions

This guide explains how to manage users, projects, warehouses, roles, and permission assignments in Lakekeeper using Terraform.

## Load the Default Project (or create a new one)

First, retrieve the default project used as a base for assigning users and roles:

```hcl
data "lakekeeper_default_project" "default" {}
```

This allows other resources to reference the `project_id` consistently.

You can also create a seperate project in case of large company setup.

```terraform
resource "lakekeeper_project" "project" {
    name = "my-project"
}
```

## Define Users

Each user must be declared using their OIDC subject identifier, along with their name, email, and user type:

```terraform
resource "lakekeeper_user" "anna" {
  id        = "oidc~d223d88c-85b6-4859-b5c5-27f3825e47f6"
  name      = "Anna Cold"
  email     = "anna@example.com"
  user_type = "human"
}

resource "lakekeeper_user" "peter" {
  id        = "oidc~cfb55bf6-fcbb-4a1e-bfec-30c6649b52f8"
  name      = "Peter Cold"
  email     = "peter@example.com"
  user_type = "human"
}
```

## Assign Users Directly to Objects

Users can be explicitly assigned to a project. You can grant specific project-level permissions such as `project_admin`.

```terraform
resource "lakekeeper_project_user_assignment" "default_anna" {
  project_id  = data.lakekeeper_default_project.default.id
  user_id     = lakekeeper_user.anna.id
  assignments = ["project_admin"]
}
```

> In this example, Anna is granted admin rights on the default project.

_You can assign users to other objects with the appropriate resources._

## Assign Role to Objects

Grant specific permissions to a role for a warehouse using a `lakekeeper_warehouse_role_assignment`:

```terraform
resource "lakekeeper_role" "select" {
  project_id  = data.lakekeeper_default_project.default.id
  name        = "test-role"
  description = "this role gives select permissions on test-warehouse"
}

resource "lakekeeper_warehouse_role_assignment" "wh_select" {
  warehouse_id = lakekeeper_warehouse.gcs.warehouse_id
  role_id      = lakekeeper_role.select.role_id
  assignments  = ["select", "describe"]
}
```

> A `test-role` is created.
> it grants read-only access to `test-warehouse` warehouse.

_You can assign roles to other objects with the appropriate resources._

## Assign the Role to a User

Finally, link the role to a specific user by assigning the `assignee` permission to them:

```terraform
resource "lakekeeper_role_user_assignment" "select_peter" {
  role_id     = lakekeeper_role.select.role_id
  user_id     = lakekeeper_user.peter.id
  assignments = ["assignee"]
}
```

> Peter now has read-only access to the warehouse via the `test-role`.