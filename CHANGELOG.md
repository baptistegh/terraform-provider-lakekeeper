## 0.2.3 (2025-07-17)


IMPROVEMENTS:

* feat: add lakekeeper_server_user_assignment resource by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/101
* feat: add lakekeeper_server_role_assignment resource by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/102
* feat: add lakekeeper_project_user/role_assignment resources by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/103
* feat: add role (role/user) assignments by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/104

MISCELLANEOUS CHORES:

* chore: disable dependabot on docker compose by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/100


## [0.2.4](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/main-v0.2.3...main-v0.2.4) (2025-07-19)


### ⚠ BREAKING CHANGES

* refactor lakekeeper go-client ([#44](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/44))

### Features

* add lakekeeper_project_user/role_assignment resources ([#103](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/103)) ([52272ac](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/52272aca1570015b9a94a99c02f3080156f6c485))
* add lakekeeper_server_role_assignment resource ([#102](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/102)) ([3a6cf95](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/3a6cf953e4212ef4ed58fd86a5026f7376ab816f))
* add lakekeeper_server_user_assignment resource ([#101](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/101)) ([e981a8e](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/e981a8e7fa75bd502369bcc86f7640bef72fb2f4))
* add project datasource/resource ([#7](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/7)) ([0f105e0](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/0f105e074f9ef047289a83dc83e41e1c8ee5f45b))
* add role (role/user) assignments ([#104](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/104)) ([1729f59](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/1729f598f342e53942e1c164b587c6aab8f8adba))
* add user datasource and resource ([#19](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/19)) ([862b21f](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/862b21ff351cc60295bd225c49000761c4ff389e))
* add warehouse resource (create/delete/read) ([#22](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/22)) ([c91d8a8](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/c91d8a8f3a5a035863803316238948210965e30d))
* add whoami datasource ([#18](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/18)) ([58a0e7d](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/58a0e7db489258a97d950c449270ef14e37534b4))
* **ci:** bump golangci-lint to version 2 ([#3](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/3)) ([9ccb729](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/9ccb729963691f6c3879944325469a60f4599a57))
* **client:** add rename warehouse method ([#73](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/73)) ([6ac483f](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/6ac483f996e25704d4f97241997549d2c89b8e61))
* **client:** add set warehouse protection method ([#75](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/75)) ([5c05190](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/5c05190c16b884a9eb649ad0d115b9fcac9b7b96))
* **client:** add static token based authentication ([#50](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/50)) ([c80b61d](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/c80b61d9e522add7c32e6961d95649c1cd5018d9))
* **client:** add update delete profile for warehouses ([#86](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/86)) ([8fe2bac](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/8fe2bacf71e98abbaa784c0de16a6f1983d54de9))
* **client:** add update storage credential for warehouses ([#85](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/85)) ([021a054](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/021a054162025a9e96619c271d5605bc3c9f411c))
* **client:** add update storage profile method for warehouses ([#84](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/84)) ([819d21f](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/819d21f7ad56fb571912ff148380ac5736990141))
* **client:** add warehouse activate/deactive methods ([#74](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/74)) ([cad7360](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/cad7360bdf7ed44444cfb9e0fff5c605797aea86))
* **client:** better error handling ([#31](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/31)) ([45a08b1](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/45a08b1b2df10f4321c42bc31b8bf7463464ae39))
* **client:** new datastructure design ([#58](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/58)) ([8806937](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/88069371595ac54addcbb8b8d2e8e2e174e342ea))
* implement project renaming ([#54](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/54)) ([b0171ec](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/b0171ec8635feb45e1eea045a5f0cd70fb5d038a))
* implements update method for warehouse resource ([#96](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/96)) ([ebd7bd9](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/ebd7bd9477fd4d49d1a88e3283c317405c3df76c))
* **permission:** add lakekeeper_warehouse_user/role_assignment resources ([#116](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/116)) ([7065514](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/7065514de26e2b1e0be13dac6f0e501823e65941))
* **provider:** add role resource and datasource ([#32](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/32)) ([9071f74](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/9071f74cfc9a42159837c6266543a2b59068fdc9))
* **provider:** add warehouse datasource on default project ([#36](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/36)) ([886bd88](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/886bd881bf958902f98e72393fbbf4651a13bfec))
* refactor lakekeeper go-client ([#44](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/44)) ([85af60b](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/85af60b1079478005c80a1785be5c2aedade7bfd))
* **test:** add acceptance test on lakekeeper_default_project datasource ([#12](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/12)) ([9a4473e](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/9a4473e29f708dd8bd3b3534ed124b8e0fd781d8))


### Bug Fixes

* **ci:** error on label removal on PR ([#70](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/70)) ([4c7a735](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/4c7a735008df09f6dd548746ba50c97ab5ce4481))
* **ci:** issue-comment-triage set correct permissions ([#52](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/52)) ([7e8aa87](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/7e8aa879b86cf6eb34531d920cb95a8994f6d098))
* **provider:** update on user not applying ([#35](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/35)) ([d810ab0](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/d810ab051bdb89d4b77a629e498c999a63e9f10e))


### Miscellaneous Chores

* prepare next release ([6772777](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/6772777eb3e14e41076589065d6c2dc5e41efdec))

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
