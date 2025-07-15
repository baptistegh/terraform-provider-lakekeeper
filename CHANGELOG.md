## 0.2.2 (2025-07-15)

IMPROVEMENTS:

* feat: implements update method for warehouse resource by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/96

MISCELLANEOUS CHORES:

* chore: skip storage validation and bump go client to v0.0.6 by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/95
* ci: fix labeler to run on PR by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/92
* chore(deps): lakekeeper go client is now a dependency by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/91

## 0.2.1 (2025-07-09)

MISCELLANEOUS CHORES:

* docs: fix some inconsistencies by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/89

## 0.2.0 (2025-07-09)

BREAKING CHANGES:

* feat!: refactor lakekeeper go-client by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/44

IMPROVEMENTS:

* feat(client): add static token based authentication by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/50
* feat: Implement project renaming in Terraform provider by @IDerr in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/54
* feat(client): add rename warehouse method by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/73
* feat(client): add warehouse activate/deactive methods by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/74
* feat(client): add set warehouse protection method by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/75
* feat(client): proposal new datastructure design by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/58

BUG FIXES:

* fix(ci): error on label removal on PR by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/70
* ci: fix pull request triage on forks by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/72

DEPENDENCIES:

* chore(deps): bump mvdan.cc/sh/v3 from 3.11.0 to 3.12.0 by @dependabot in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/56
* chore(deps): bump bitnami/minio from 2025.4.22 to 2025.5.24 in /run by @dependabot in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/78
* chore(deps): bump keycloak/keycloak from 26.0.7 to 26.3.0 in /run by @dependabot in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/81
* chore(deps): bump bitnami/postgresql from 16.3.0 to 17.5.0 in /run by @dependabot in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/80
* chore(deps): bump openfga/openfga from v1.8 to v1.9 in /run by @dependabot in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/79

MISCELLANEOUS CHORES:

* chore(ci): fix missing permissions preventing label assignment by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/52
* chore: fix pr request triage on fork by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/55
* ci: add workflow auto merge on review by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/59
* ci: add labeler by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/69
* ci: fix dependabot on docker compose by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/76
* ci: remove ci cache by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/83
* docs: fix some typo errors by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/53
* ci: run acceptance tests on 3 last lakekeeper versions by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/82

## 0.1.1 (2025-07-06)

IMPROVEMENTS:

* docs: add contribution guidelines by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/43
* docs: add code of conduct by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/47
  
DEPENDENCIES:

* chore(deps): bump github.com/hashicorp/terraform-plugin-testing from 1.13.0 to 1.13.2 in the terraform-plugin group by @dependabot in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/48

MISCELLANEOUS CHORES:

* chore(ci): refactor linters by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/41
* chore: force coventional commits on pr by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/45
* ci: remove feature label by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/46
* ci: add dependabot on docker compose images by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/40

## 0.1.0 (2025-07-05)

FEATURES:

* feat: add whoami datasource by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/18
* feat: add user datasource and resource by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/19
* feat: add warehouse resource (create/delete/read) by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/22
* feat(provider): add role resource and datasource by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/32
* feat(provider): add warehouse datasource on default project by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/36

IMPROVEMENTS:

* chore(test): add validation on user IDs by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/23
* feat(client): better error handling by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/31
* docs: refine documentation in README and examples by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/37

BUG FIXES:

* fix(provider): update on user not applying by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/35
  
DEPENDENCIES:

* chore(deps): bump github.com/hashicorp/terraform-plugin-docs from 0.21.0 to 0.22.0 in /tools by @dependabot in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/25

MISCELLANEOUS CHORES:

* chore(ci): only run tests on PR and push on main by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/13
* feat(test): add acceptance test on lakekeeper_default_project datasource by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/12
* chore(ci): add release note categories by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/20
* chore(docs): add non-production notice by @IDerr in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/21
* chore(ci): add terraform and sh linters by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/26
* chore: add improvements category in release notes by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/34
* chore: add MAINTAINERS.md by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/38
