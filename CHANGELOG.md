# Changelog


## [0.4.2](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.4.1...v0.4.2) (2025-08-06)


### Features

* use the string constants of the client ([#164](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/164)) ([bf5bbf5](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/bf5bbf53397c9bba43ec5e24bf73d54ad58ded9f))


### Miscellaneous Chores

* **deps:** bump github.com/baptistegh/go-lakekeeper from 0.0.14 to 0.0.16 ([#161](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/161)) ([a359ecb](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/a359ecbe75f0285b399883bc1ee4d52f771269bc))
* **deps:** bump github.com/baptistegh/go-lakekeeper from 0.0.16 to 0.0.17 ([#165](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/165)) ([3d4817b](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/3d4817b45ada6e0ac3c17833a3ed3d3fff88b864))

## [0.4.1](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.4.0...v0.4.1) (2025-08-03)


### Documentation

* fix broken navigation ([#159](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/159)) ([3289050](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/3289050c91450a859e29e841b115ef4800855c34))

## [0.4.0](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.3.2...v0.4.0) (2025-08-03)


### ⚠ BREAKING CHANGES

* **warehouse:** use the same design on datasource ([#156](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/156))
* **warehouse:** wrap storage profiles and credential in nested object ([#153](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/153))
* **warehouse:** redesign storage profile and credential ([#152](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/152))

### Features

* **warehouse:** redesign storage profile and credential ([#152](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/152)) ([9f33757](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/9f337572312720ee4498bbf536b16a7199eda371))
* **warehouse:** use the same design on datasource ([#156](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/156)) ([17f3e6d](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/17f3e6d21f48680e7e48268c2abb2a25501b3f2c))
* **warehouse:** wrap storage profiles and credential in nested object ([#153](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/153)) ([568909f](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/568909fe9e3cf5858073b738b1658ca0f0035640))


### Documentation

* add manage warehouses guides ([#158](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/158)) ([281886b](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/281886b736227e39ad1617e173803cfb2ffd240f))
* Fix Terraform docs link in README.md ([#157](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/157)) ([b370e1c](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/b370e1c40877b6dd6c8052ab087f500faff5d7ea))
* Improve playground scenario ([#153](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/153)) ([568909f](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/568909fe9e3cf5858073b738b1658ca0f0035640))


### Miscellaneous Chores

* add pr title checker ([#155](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/155)) ([5a10e3e](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/5a10e3ea3a09b461d7303420086d8d71bf3af31a))
* **ci:** improve tests workflow ([#154](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/154)) ([3043335](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/3043335f8f63c3a18e71a3c6f4d95d1ca721f4a1))
* **ci:** rename archive.formats in goreleaser config ([#150](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/150)) ([eabcefa](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/eabcefa54da411efdc2815f05b3cc386d4396e0c))
* **ci:** set up changelog sections for release please ([#148](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/148)) ([2c5b72f](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/2c5b72f0e5037358e5353295f7a663b000cf1100))
* **deps:** bump github.com/baptistegh/go-lakekeeper ([953273f](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/953273fd76599de700398953b81ca31ff8cfee96))
* **deps:** bump github.com/baptistegh/go-lakekeeper from 0.0.11 to 0.0.14 ([#147](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/147)) ([953273f](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/953273fd76599de700398953b81ca31ff8cfee96))
* **deps:** bump github.com/hashicorp/terraform-plugin-framework from 1.15.0 to 1.15.1 in the terraform-plugin group ([#146](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/146)) ([9431d6b](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/9431d6be8ca57f05bb97079932c94257c27062d3))
* rename datasource and resource files ([#151](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/151)) ([1c908f7](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/1c908f7ed062a0f269a8acef875c7d85d712979b))

## [0.3.2](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.3.1...v0.3.2) (2025-07-30)


### Bug Fixes

* oauth2 scope was not correctly sent ([#144](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/144)) ([b5f70c6](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/b5f70c6db3cff9827eef300c4bddc417fb5b14b0))

## [0.3.1](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.3.0...v0.3.1) (2025-07-24)


### Bug Fixes

* **build:** conflicts on go tags with uniseg@v0.4.7 ([#139](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/139)) ([2ed43bb](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/2ed43bb23de6d8a0a02616eae32763bd605661ef))
* **user:** regex matching if ID is unknown ([#143](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/143)) ([d63e70d](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/d63e70d769c742444a8f2b87e4cc042f92b5fa76))

## [0.3.0](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.2.5...v0.3.0) (2025-07-23)


### ⚠ BREAKING CHANGES

* replace : by / in composite ids ([#138](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/138))

### Features

* allow setting warehouse parameters as unknown variables ([#153](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/153)) ([79536a0](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/79536a0dd621c7b1796bd1112b3b436b36b12917))
* **docs:** add missing parameters documentation ([#137](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/137)) ([f77cf13](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/f77cf13d967cb34ddc68d12f654ce893871d20bc))
* replace : by / in composite ids ([#138](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/138)) ([ff53e7d](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/ff53e7d2156706a82c32854d2a910c60ccaabbd6))

## [0.2.5](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.2.4...v0.2.5) (2025-07-21)


### Features

* add lakekeeeper_warehouse_user/role_access data sources ([0853f22](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/0853f2208fc814796b8a0a55d112d8c9831952f3))
* add lakekeeper_warehouse_assignments datasource ([0853f22](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/0853f2208fc814796b8a0a55d112d8c9831952f3))
* add lakekeeper_warehouse_assignments datasource ([#128](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/128)) ([0853f22](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/0853f2208fc814796b8a0a55d112d8c9831952f3))
* add playground examples ([#120](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/120)) ([f22dc38](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/f22dc386a4c546b30c46d47feec2e017606623dc))
* **project:** add permission datasources ([#129](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/129)) ([08af17a](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/08af17a778a82afaa202f11abda2a1695f3d2264))
* **warehouse:** add managed access property ([#127](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/127)) ([a06c619](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/a06c6191fe6ed852f7fa68f15f79115e73d9be02))

## [0.2.4](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.2.3...v0.2.4) (2025-07-19)


### Features

* **permission:** add lakekeeper_warehouse_user/role_assignment resources ([#116](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/116)) ([7065514](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/7065514de26e2b1e0be13dac6f0e501823e65941))


### Miscellaneous Chores

* prepare next release ([6772777](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/6772777eb3e14e41076589065d6c2dc5e41efdec))

## 0.2.3 (2025-07-17)

IMPROVEMENTS:

* feat: add lakekeeper_server_user_assignment resource by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/101
* feat: add lakekeeper_server_role_assignment resource by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/102
* feat: add lakekeeper_project_user/role_assignment resources by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/103
* feat: add role (role/user) assignments by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/104

MISCELLANEOUS CHORES:

* chore: disable dependabot on docker compose by @baptistegh in https://github.com/baptistegh/terraform-provider-lakekeeper/pull/100

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
