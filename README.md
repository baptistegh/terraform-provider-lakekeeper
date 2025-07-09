<div align="center">
    <h1>Terraform Provider Lakekeeper</h1>
    <p>
Terraform provider for <a href="https://docs.lakekeeper.io/">Lakekeeper</a>.<p>
    <a href="https://registry.terraform.io/providers/baptistegh/lakekeeper/latest/docs"><img src="https://img.shields.io/static/v1?label=Docs&message=terraform-provider-lakekeeper&color=000000&style=for-the-badge" /></a>
  <a href="https://github.com/baptistegh/terraform-provider-lakekeeper/releases"><img src="https://img.shields.io/badge/status-preview-orange?style=for-the-badge" /></a>
</div>

### ⚠️ Preview Status Notice

**🚧 This project is currently in _preview_ and under active development. It is _not production-ready_ and should not be used in production environments.**

- ⚠️ **Breaking changes** may occur without notice  
- 🔄 APIs and behavior may change significantly between versions  
- 🧪 Use at your own risk for development and testing purposes only  

💬 We welcome [feedback, bug reports, and contributions](https://github.com/baptistegh/terraform-provider-lakekeeper/issues) during this preview phase. Please report issues or share your experience to help us improve the provider before its stable release.



## Docs

All documentation for this provider can be found on the Terraform Registry: https://registry.terraform.io/providers/baptistegh/lakekeeper/latest/docs.

## Installation

This provider can be installed automatically using Terraform >=0.13 by using the `terraform` configuration block:

```terraform
terraform {
  required_providers {
    lakekeeper = {
      source = "baptistegh/lakekeeper"
      version = "0.2.1"
    }
  }
}
```

## Supported Versions

This provider will officially support the latest three last versions of Lakekeeper, although older versions may still work.

The following Lakekeeper versions are used when running acceptance tests in CI:

- 0.9.2 (latest)
- 0.9.1
- 0.9.0

_Acceptance tests are executed using Terraform v1.12.2._

## Releases

This provider uses [GoReleaser](https://goreleaser.com/]) to build and publish releases. Each release published to GitHub contains binary files for Linux, macOS (darwin), and Windows, as configured within the [`.goreleaser.yml`](https://github.com/baptistegh/terraform-provider-lakekeeper/blob/main/.goreleaser.yml) file.

Each release also contains a `terraform-provider-lakekeeper_${RELEASE_VERSION}_SHA256SUMS` file that can be used to check integrity.

You can find the list of releases [here](https://github.com/baptistegh/terraform-provider-lakekeeper/releases). You can find the changelog for each version [here](https://github.com/baptistegh/terraform-provider-lakekeeper/blob/main/CHANGELOG.md).

## Development

This project requires Go 1.24 and Terraform >= 1.9.8.

After cloning the repository, you can build the project by running `make build`

### Local Environment

You can spin up a local developer environment via [Docker Compose](https://docs.docker.com/compose/) by running `make testacc-up`. This will spin up a few containers for Lakekeeper, Keycloak, PostgreSQL and OpenFGA, which can be used for testing the provider.

To stop the environment you can use the `make clean`.

### Tests

Every resource supported by this provider will have a reasonable amount of acceptance test coverage.

You can run acceptance tests against a Lakekeeper instance by running `make testacc`.