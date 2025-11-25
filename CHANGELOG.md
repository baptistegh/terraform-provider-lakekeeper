# Changelog


## [0.4.7](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.4.6...v0.4.7) (2025-11-25)


### Features

* add namespace resource ([#216](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/216)) ([59a9654](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/59a96549ff91ecbff3d86b717cea3d1a7cf2da3b))
* **warehouse:** add `get_endpoint_statistics` action in permissions ([#217](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/217)) ([f5474a6](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/f5474a6981d60e04ec1326d0e9b93ab9dbae8ad2))


### Bug Fixes

* **deps:** update all non-major dependencies (minor) ([#208](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/208)) ([c67e202](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/c67e20247c42b2a9d58c6a477ab90d894dc82365))
* **deps:** update module github.com/baptistegh/go-lakekeeper to v0.0.20 ([#218](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/218)) ([10d9f5a](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/10d9f5a4ef052465900e12a5dfdeef730ae12ada))
* **deps:** update module github.com/baptistegh/go-lakekeeper to v0.0.22 ([#220](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/220)) ([1b52449](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/1b5244928646bb5e7c8264e3ed59a8ed3d73813d))


### Miscellaneous Chores

* **ci:** configure renovate ([#204](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/204)) ([50791b1](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/50791b157186515b0da498f784f6b151eb57fbde))
* **ci:** Remove daily schedule from renovate.json ([#205](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/205)) ([892d8f9](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/892d8f990739f4eabd4acdc8879fae97ab4f44e0))
* **ci:** Update renovate.json for Go dependency management ([#211](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/211)) ([bd075fb](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/bd075fb973486f72ba9ecc319d85dccb1863ccba))
* Configure Renovate ([#200](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/200)) ([dc99d0d](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/dc99d0d99e60ed6b8f989d027e5629e025cd9399))
* **deps:** bump golang.org/x/crypto ([2c07839](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/2c07839b19293ab30fb23cd9be355f031e62087f))
* **deps:** bump golang.org/x/crypto from 0.42.0 to 0.45.0 in the go_modules group across 1 directory ([#212](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/212)) ([2c07839](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/2c07839b19293ab30fb23cd9be355f031e62087f))
* **deps:** bump golang.org/x/oauth2 from 0.32.0 to 0.33.0 ([#197](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/197)) ([17587ca](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/17587caac39a50f468e820df0b42bfac2a2830ea))
* **deps:** bump golangci/golangci-lint-action from 8.0.0 to 9.0.0 in the github-actions group ([#198](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/198)) ([b6d170a](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/b6d170a4b096837f09a1b5b0d98b5c616f733fb5))
* **deps:** update actions/checkout action to v5.0.1 ([#206](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/206)) ([32238bb](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/32238bb724d3fb04e59b517f5bc0a37f6d01ca01))
* **deps:** update actions/checkout action to v6 ([#214](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/214)) ([ddb6779](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/ddb67799dd606ad0e168872077777a6a47c2060a))
* **deps:** update all non-major dependencies ([bd0ef35](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/bd0ef35c6ce548e1f27d75c6e2da32934fde7c93))
* **deps:** update all non-major dependencies (minor) ([#213](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/213)) ([bd0ef35](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/bd0ef35c6ce548e1f27d75c6e2da32934fde7c93))
* **deps:** update postgres docker tag to v18 ([#209](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/209)) ([d1fb8e5](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/d1fb8e5ae725c8268fe535e9190fee816794cbd2))
* remove use of deprecated GetDefault() ([#219](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/219)) ([44388e8](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/44388e82e9d9dffc8893fcf573b9252f6495014e))
* **renovate:** remove grouping dependencies ([#215](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/215)) ([5cc2dbc](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/5cc2dbcbab6cc70dc520cae1a1c3dd753e853023))
* **renovate:** Update configuration ([#207](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/207)) ([b83b4d2](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/b83b4d24c2d0a908cfb936260dd3842275a8eb54))
* **tests:** use slices contains for allowed actions checks ([#221](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/221)) ([89d2fc4](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/89d2fc43dad79318fb62b27f2cfd05d5821a1f16))

## [0.4.6](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.4.5...v0.4.6) (2025-11-05)


### Bug Fixes

* **authn:** use dedicated http client for OAuth2 Token Source ([#193](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/193)) ([60d0225](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/60d02251a646c044bce0fdde0d34a12520a0c0f1))


### Miscellaneous Chores

* use go1.25 ([#194](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/194)) ([1ce0adf](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/1ce0adf31501d8eb115106869efefd972b6888dd))

## [0.4.5](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.4.4...v0.4.5) (2025-10-29)


### Bug Fixes

* Add some input validation to role/warehouse IDs ([#189](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/189)) ([8606cad](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/8606cad6e271226869ad3bbc190cf56b557545e1))
* Use PlanModifiers to reduce the amount of unnecessary "(known after apply)" entires in a plan ([#188](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/188)) ([563c3cd](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/563c3cdbfdbf72e326ed93896e4008d3c00f12e8))
* **warehouse:** use actions instead of assignments ([#184](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/184)) ([c4a069f](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/c4a069f77169ada495d27e4563d23ebc3faa0618))


### Miscellaneous Chores

* add compatibility for lakekeeper 0.10.0 ([#183](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/183)) ([aa5fbb8](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/aa5fbb89bf9333e6435cc6290dfac146444bee9e))
* **ci:** use latest Lakekeeper versions in CI ([#190](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/190)) ([72277b5](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/72277b55f89039316a93057f6579dfb2c3dd0d0a))
* **deps:** bump github.com/baptistegh/go-lakekeeper from 0.0.17 to 0.0.18 ([#181](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/181)) ([33a03be](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/33a03be58431e8a31495c6d5034e6e4bd0533340))
* **deps:** bump github.com/hashicorp/terraform-plugin-docs from 0.23.0 to 0.24.0 in the terraform-plugin group ([#187](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/187)) ([faf81cc](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/faf81cc992debb173ae7bd78da11b95d7a0c3247))
* **deps:** bump github.com/hashicorp/terraform-plugin-framework from 1.16.0 to 1.16.1 in the terraform-plugin group ([#182](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/182)) ([8dd680d](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/8dd680da35c65d6973e119b89666e8dda451db62))
* **deps:** bump github.com/hashicorp/terraform-plugin-framework-validators from 0.18.0 to 0.19.0 in the terraform-plugin group ([#185](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/185)) ([4a6a174](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/4a6a1747e3450e1c66050f17210ff55ff8715eb9))
* **deps:** bump golang.org/x/oauth2 from 0.30.0 to 0.31.0 ([#178](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/178)) ([01db534](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/01db53495d540cd959fa501299f774ffbd7fefc8))
* **deps:** bump golang.org/x/oauth2 from 0.31.0 to 0.32.0 ([#186](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/186)) ([3bba8e8](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/3bba8e833a3cd95488d55eac15abf388c2c9203e))
* **deps:** bump the github-actions group with 2 updates ([#176](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/176)) ([1a0e83a](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/1a0e83aca0eb4867d0f0cb16a93f520e59ed0053))
* **deps:** bump the terraform-plugin group with 3 updates ([#179](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/179)) ([0dd04e8](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/0dd04e8adf9f9a173f4b27c69a3597ebdae9b658))
* fix typo in `README.md` ([#191](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/191)) ([f4794dc](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/f4794dc639d995a308fb4e90f23e17d6c61138bf))
* remove bitnami postgresql image ([#180](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/180)) ([35756b2](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/35756b22441f1a8ef161d51450cf8004b7664212))

## [0.4.4](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.4.3...v0.4.4) (2025-09-02)


### Bug Fixes

* handle role renaming ([#174](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/174)) ([89acec1](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/89acec1d002aba4ba160582f98c350b47df0982f))


### Miscellaneous Chores

* add tests on lakekeeper v0.9.4 / v0.9.5 ([#175](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/175)) ([0b0fa42](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/0b0fa4279e9eee5d9d48dfce4b442871ba0ae9cb))
* **deps:** bump mfinelli/setup-shfmt from 3 to 4 in the github-actions group ([#172](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/172)) ([15afd7b](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/15afd7b5d05c5a56bc59af05c6d92c10eb6557f8))

## [0.4.3](https://github.com/baptistegh/terraform-provider-lakekeeper/compare/v0.4.2...v0.4.3) (2025-08-20)


### Features

* **warehouse:** add import on warehouse resources ([#170](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/170)) ([b9b761c](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/b9b761cec1ce61883e380e6dde09fa698616cb57))


### Bug Fixes

* **roles:** typo error on import state project_id ([#169](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/169)) ([23e0c16](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/23e0c163b12d1fb5de48382dfe8ac02d6a44d413))


### Documentation

* add link to OpenTofu registry ([#166](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/166)) ([c31678c](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/c31678cfbf173c6d6ddd7f9ff4f3cd3655b34a9a))


### Miscellaneous Chores

* **deps:** bump actions/checkout from 4 to 5 in the github-actions group ([#168](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/168)) ([b9c678f](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/b9c678f2a31759188667f94da6911ad6cee9d9b3))
* **deps:** bump github.com/hashicorp/terraform-plugin-testing from 1.13.2 to 1.13.3 in the terraform-plugin group ([#171](https://github.com/baptistegh/terraform-provider-lakekeeper/issues/171)) ([83da6b1](https://github.com/baptistegh/terraform-provider-lakekeeper/commit/83da6b159820b0f83152ca3392b51f2c1e0f1240))

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
