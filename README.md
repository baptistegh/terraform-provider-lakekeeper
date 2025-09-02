# Terraform Provider for Lakekeeper

[![docs](https://img.shields.io/static/v1?label=Docs&message=terraform-provider-lakekeeper&color=5c4ee5)](https://registry.terraform.io/providers/baptistegh/lakekeeper/latest/docs)
[![docs](https://img.shields.io/badge/status-preview-orange)](https://github.com/baptistegh/terraform-provider-lakekeeper/releases)
[![Nightly Acceptance Tests](https://github.com/baptistegh/terraform-provider-lakekeeper/actions/workflows/nightly.yml/badge.svg)](https://github.com/baptistegh/terraform-provider-lakekeeper/actions/workflows/nightly.yml)
[![Tests](https://github.com/baptistegh/terraform-provider-lakekeeper/actions/workflows/test.yml/badge.svg)](https://github.com/baptistegh/terraform-provider-lakekeeper/actions/workflows/test.yml)

Terraform provider for [Lakekeeper](https://docs.lakekeeper.io/).

> [!IMPORTANT]  
> **⚠️ Preview Status Notice**
>
> - ⚠️ **Breaking changes** may occur without notice  
> - 🔄 APIs and behavior may change significantly between versions  
> - 🧪 Use at your own risk for development and testing purposes only
>
> **🚧 This project is currently in _preview_ and under active development. It is _not production-ready_ and should not be used in production environments.**
>
> 💬 We welcome [feedback, bug reports, and contributions](https://github.com/baptistegh/terraform-provider-lakekeeper/issues) during this preview phase.
> Please report issues or share your experience to help us improve the provider before its stable release.

## Table of Contents

- [Terraform Provider for Lakekeeper](#terraform-provider-for-lakekeeper)
  - [Table of Contents](#table-of-contents)
  - [Documentation](#documentation)
  - [Installation](#installation)
  - [Supported Versions](#supported-versions)
  - [Playground](#playground)
    - [Overview](#overview)
    - [Setup Instructions](#setup-instructions)
  - [Development](#development)
    - [Local Environment](#local-environment)
    - [Tests](#tests)
    - [Releases](#releases)

## Documentation

All documentation for this provider can be found on the Terraform Registry: <https://registry.terraform.io/providers/baptistegh/lakekeeper/latest/docs>.

The provider is also available from [OpenTofu registry](https://search.opentofu.org/provider/baptistegh/lakekeeper/latest).

## Installation

This provider can be installed automatically by using the `terraform` configuration block:

```terraform
terraform {
  required_providers {
    lakekeeper = {
      source = "baptistegh/lakekeeper"
    }
  }
}
```

## Supported Versions

This provider will officially support the latest three last versions of Lakekeeper starting v0.9.3, although older versions may still work.

The following Lakekeeper versions are used when running acceptance tests in CI:

- Unreleased (latest-main)
- v0.9.5
- v0.9.4
- v0.9.3

_The provider should be compatible with Lakekeeper >= v0.9.3 thanks to the introduction of skip storage validation._
_See: [lakekeeper/lakekeeper#1239](https://github.com/lakekeeper/lakekeeper/pull/1239)_

_Acceptance tests are executed using Terraform v1.9.8._

## Playground

A sample playground project is available in the `playground/` directory. It is based on the [Official Lakekeeper Examples](https://github.com/lakekeeper/lakekeeper/tree/main/examples/access-control-simple).

### Overview

This playground will set up the following structure:

- **Warehouses:**
  - `test-warehouse-gcs`: Configured using a `gcs` (Google Cloud Storage) storage profile.
  - `test-warehouse-s3`: Configured using a `s3` (AWS S3) storage profile.
  - `test-warehouse-adls`: Configured using a `adls` (Azure Data Lake Storage) storage profile.
- **Roles:**
  - `read-role`: with `select` and `describe` permissions on all the warehouses.
  - `write-role`: with `create` and `modify` permissions on the `test-warehouse-s3` warehouse.
- **Users:**
  - **Anna:**
    - username: `anna`
    - password: `iceberg`
    - `project_admin` assignment on the default project
  - **Peter:**
    - username: `peter`
    - password: `iceberg`
    - assignee to the `read-role` and `write-role`

_Anna is a `project_admin` and Peter has read access on all the warehouses and can write on the `test-warehouse-s3` warehouse._ 

### Setup Instructions

To create and launch the playground:

```sh
make testacc-up   # Sets up required services via Docker Compose (Lakekeeper, Keycloak, OpenFGA, PostgreSQL)
make playground   # Creates the playground structure
```

To tear down and clean up the playground:

```sh
make playground-destroy
```

You can connect to the web interface at <http://localhost:8181> using one of the user credentials listed above to explore the configured resources.

Feel free to modify `playground/main.tf` to customize the structure according to your needs.

## Development

This project requires Go 1.24 and Terraform >= 1.9.8.

After cloning the repository, you can build the project by running `make build`

### Local Environment

You can spin up a local developer environment via [Docker Compose](https://docs.docker.com/compose/) by running `make testacc-up`. This will spin up a few containers for Lakekeeper, Keycloak, PostgreSQL and OpenFGA, which can be used for testing the provider.

To stop the environment you can use the `make clean`.

### Tests

Every resource supported by this provider will have a reasonable amount of acceptance test coverage.

You can run acceptance tests against a Lakekeeper instance by running `make testacc`.

### Releases

This provider uses [GoReleaser](https://goreleaser.com/]) to build and publish releases. Each release published to GitHub contains binary files for Linux, macOS (darwin), and Windows, as configured within the [`.goreleaser.yml`](https://github.com/baptistegh/terraform-provider-lakekeeper/blob/main/.goreleaser.yml) file.

Each release also contains a `terraform-provider-lakekeeper_${RELEASE_VERSION}_SHA256SUMS` file that can be used to check integrity.

You can find the [list of releases](https://github.com/baptistegh/terraform-provider-lakekeeper/releases) and the [changelog](https://github.com/baptistegh/terraform-provider-lakekeeper/blob/main/CHANGELOG.md) for each version.
