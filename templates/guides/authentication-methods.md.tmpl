---
page_title: "Authentication Methods"
subcategory: ""
description: |-

---

# Authentication Methods

This guide explains how to configure authentication for the lakekeeper provider in Terraform.
Authentication is performed via an OpenID Connect (OIDC)-compatible identity provider.

## Prerequisites

* Access to a Lakekeeper API endpoint
* A configured OIDC client with:
    * `client_id`
    * `client_secret`
    * Token URL (`auth_url`)
    * Required scope(s), e.g., `lakekeeper`

## Example Provider Configuration

```hcl
provider "lakekeeper" {
  endpoint      = "<lakekeeper_base_url>"
  auth_url      = "<token_endpoint>"
  client_id     = "<client_id>"
  client_secret = var.client_secret
  scopes        = ["lakekeeper"]
}

variable "client_secret" {
  type      = string
  sensitive = true
}
```

You can then set Terraform variables like:
- If a variable does not have any value set, you will be prompted by Terraform to provide the value.
- Use Terraform VAR environment variables: `TF_VAR_client_secret="<secret>" terraform plan`
- Use Terraform flags: `terraform plan -var="client_secret=<secret>"`
- Use Lakekeeper Terraform Provider environment variables: `LAKEKEEPER_CLIENT_SECRET="<secret>" terraform plan`

If you prefer to set up the provider with environment variables, you can use:
- `LAKEKEEPER_ENDPOINT`
- `LAKEKEEPER_AUTH_URL`
- `LAKEKEEPER_CLIENT_ID`
- `LAKEKEEPER_CLIENT_SECRET`

### Parameter Details

| Parameter       | Description                                                                                             | Default Value    |
| --------------- | ------------------------------------------------------------------------------------------------------- | ---------------- |
| `endpoint`      | The URL of the Lakekeeper API                                                                           |
| `auth_url`      | The OIDC token endpoint (typically includes /realms/<realm>/protocol/openid-connect/token for Keycloak) |
| `client_id`     | The OIDC client ID configured in your identity provider                                                 |
| `client_secret` | The corresponding secret for the client                                                                 |
| `scopes`        | A list of scopes to request                                                                             | `["lakekeeper"]` |

Other parameters are available, you can find them on the provider home page documentation.


## Authentication Flow

When you run Terraform commands, the provider:

1. Sends a `POST` request to the `auth_url` with the OIDC client credentials.
2. Retrieves an access token (JWT).
3. Uses the token to authenticate requests to the Lakekeeper API.

## Security Tips

- Never commit your client_secret to version control.
- Use environment variables or secret management tools to inject secrets securely.
- Consider using Terraform Cloud Environment Variables or a tool like HashiCorp Vault.