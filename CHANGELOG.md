# Changelog

## [0.4.1](https://github.com/nobbs/kubectl-mapr-ticket/compare/v0.4.0...v0.4.1) (2024-02-27)


### Bug Fixes

* **deps:** update golang.org/x/exp digest to 2c58cdc ([#58](https://github.com/nobbs/kubectl-mapr-ticket/issues/58)) ([2fcef6d](https://github.com/nobbs/kubectl-mapr-ticket/commit/2fcef6d0fe90146dd0d8f31a8b7503fd90608137))
* **deps:** update golang.org/x/exp digest to 814bf88 ([#63](https://github.com/nobbs/kubectl-mapr-ticket/issues/63)) ([1aaf728](https://github.com/nobbs/kubectl-mapr-ticket/commit/1aaf728f56ef8c9bbe5797b8d24768ba4e3e8737))
* **deps:** update golang.org/x/exp digest to ec58324 ([#60](https://github.com/nobbs/kubectl-mapr-ticket/issues/60)) ([702c0d0](https://github.com/nobbs/kubectl-mapr-ticket/commit/702c0d06ee36af965b04904afd938012fbaf613f))
* **deps:** update kubernetes packages to v0.29.2 ([#61](https://github.com/nobbs/kubectl-mapr-ticket/issues/61)) ([6023551](https://github.com/nobbs/kubectl-mapr-ticket/commit/60235512ccd920d393b4b29a46ddb86bd49f0198))


### Tests

* add more tests for various code paths ([#55](https://github.com/nobbs/kubectl-mapr-ticket/issues/55)) ([4931b3d](https://github.com/nobbs/kubectl-mapr-ticket/commit/4931b3db67d17fcd318fc88ddbf86cd556cf5899))


### Miscellaneous Chores

* **deps:** update codecov/codecov-action action to v4 ([#56](https://github.com/nobbs/kubectl-mapr-ticket/issues/56)) ([d05da5e](https://github.com/nobbs/kubectl-mapr-ticket/commit/d05da5ef238dc225d4961aa23ea1cd697679ccc8))
* **deps:** update dependency golang to v1.22.0 ([#62](https://github.com/nobbs/kubectl-mapr-ticket/issues/62)) ([629f544](https://github.com/nobbs/kubectl-mapr-ticket/commit/629f544f16fed058c2d853d5c76ce235bbf8be00))
* **deps:** update golangci/golangci-lint-action action to v4 ([#59](https://github.com/nobbs/kubectl-mapr-ticket/issues/59)) ([da53337](https://github.com/nobbs/kubectl-mapr-ticket/commit/da533373138702e25d4486460787d62119ca406d))

## [0.4.0](https://github.com/nobbs/kubectl-mapr-ticket/compare/v0.3.0...v0.4.0) (2024-01-27)


### Features

* add `inspect` command ([#52](https://github.com/nobbs/kubectl-mapr-ticket/issues/52)) ([f4bd4bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/f4bd4bcbc1d03f5464afac3f578a0262bba0f234))


### Bug Fixes

* ci not properly updating version ([c066e8a](https://github.com/nobbs/kubectl-mapr-ticket/commit/c066e8a950890bc285080c8cbd590421c590df0b))
* complete multiple values for sort-by flags ([8f4bf49](https://github.com/nobbs/kubectl-mapr-ticket/commit/8f4bf496b24515bc8f7d02c47a6dd457a8b06683))
* **deps:** update module github.com/nobbs/mapr-ticket-parser to v0.1.3 ([f4bd4bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/f4bd4bcbc1d03f5464afac3f578a0262bba0f234))


### Continuous Integration

* use gotestsum for tests ([f4bd4bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/f4bd4bcbc1d03f5464afac3f578a0262bba0f234))


### Documentation

* add godoc comments to cmd packages ([f4bd4bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/f4bd4bcbc1d03f5464afac3f578a0262bba0f234))
* add inspect subcommand to README ([f4bd4bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/f4bd4bcbc1d03f5464afac3f578a0262bba0f234))
* add some more go doc strings ([f4bd4bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/f4bd4bcbc1d03f5464afac3f578a0262bba0f234))


### Miscellaneous Chores

* add license header to all files ([f4bd4bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/f4bd4bcbc1d03f5464afac3f578a0262bba0f234))
* **deps:** upgrade all dependencies ([b376b29](https://github.com/nobbs/kubectl-mapr-ticket/commit/b376b291e45fd0247ff3addce5c27fbcdbf118c6))
* switch to direnv for managing local dev env ([f4bd4bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/f4bd4bcbc1d03f5464afac3f578a0262bba0f234))
* switch to prerelease version of mapr-ticket-parser ([f4bd4bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/f4bd4bcbc1d03f5464afac3f578a0262bba0f234))
* typo in SPDX license header ([9120486](https://github.com/nobbs/kubectl-mapr-ticket/commit/912048669a42b2196a59a87bc9f380d87074ac78))
* update pre-commit config ([03df59a](https://github.com/nobbs/kubectl-mapr-ticket/commit/03df59aaf0bcccfe6af2338ba9bef621e2c1ebb8))

## [0.3.0](https://github.com/nobbs/kubectl-mapr-ticket/compare/v0.2.1...v0.3.0) (2024-01-25)


### Features

* add --all-namespaces flag to used-by command, also tests ([#40](https://github.com/nobbs/kubectl-mapr-ticket/issues/40)) ([7c1af4c](https://github.com/nobbs/kubectl-mapr-ticket/commit/7c1af4c37e0e9afc9be92d2e4f740a2037fd225d))
* add claim command to list all PVCs using tickets ([#44](https://github.com/nobbs/kubectl-mapr-ticket/issues/44)) ([0e2eefd](https://github.com/nobbs/kubectl-mapr-ticket/commit/0e2eefd82e3e4b5acee4d8699422ad480ecc39b6))
* add sort option to `volume` command ([a8f52ee](https://github.com/nobbs/kubectl-mapr-ticket/commit/a8f52eeb252edc6204a54a67278b42a31a4f7b6a))
* add ticket status to volume command ([07f6ce0](https://github.com/nobbs/kubectl-mapr-ticket/commit/07f6ce057e0cc74cb69988b6475d62dd13f8d98f))
* implement claim sorting ([#50](https://github.com/nobbs/kubectl-mapr-ticket/issues/50)) ([3499c07](https://github.com/nobbs/kubectl-mapr-ticket/commit/3499c078562bf7b1c5b0aa72f89a078a13c845f5))


### Bug Fixes

* add pvc alias to claim command ([cbeb8c9](https://github.com/nobbs/kubectl-mapr-ticket/commit/cbeb8c9fbca3b5a9dc6887f67ca96c920dd13890))
* **deps:** update kubernetes packages to v0.29.1 ([#46](https://github.com/nobbs/kubectl-mapr-ticket/issues/46)) ([d757733](https://github.com/nobbs/kubectl-mapr-ticket/commit/d757733e1de9a52f114eb40728f4cc117736e640))
* rename commands, `list` to `secret` and `usedby` to `volume` ([868b96a](https://github.com/nobbs/kubectl-mapr-ticket/commit/868b96a87b1c481e0cc66123cfe19db980ddfbbb))
* streamline sort options for all commands ([3499c07](https://github.com/nobbs/kubectl-mapr-ticket/commit/3499c078562bf7b1c5b0aa72f89a078a13c845f5))
* ticket status not properly parsed in `claim` command ([#48](https://github.com/nobbs/kubectl-mapr-ticket/issues/48)) ([901b82a](https://github.com/nobbs/kubectl-mapr-ticket/commit/901b82a10250c9bda945a0e29124b4b3178a7046))
* volume sort options ([7c0288b](https://github.com/nobbs/kubectl-mapr-ticket/commit/7c0288bee85eba9364aac16849eae5fef0955538))


### Tests

* add test for cli completion functions ([7875ffd](https://github.com/nobbs/kubectl-mapr-ticket/commit/7875ffd145ee3ada1a1c9dbc68c779b461680976))
* add tests for `version.String()` ([5253528](https://github.com/nobbs/kubectl-mapr-ticket/commit/525352894f7da34713bd5734067b894a8cae541c))
* add tests for duration pflag type ([7875ffd](https://github.com/nobbs/kubectl-mapr-ticket/commit/7875ffd145ee3ada1a1c9dbc68c779b461680976))
* add tests for types ([3499c07](https://github.com/nobbs/kubectl-mapr-ticket/commit/3499c078562bf7b1c5b0aa72f89a078a13c845f5))
* add tests for util/cli.go ([7875ffd](https://github.com/nobbs/kubectl-mapr-ticket/commit/7875ffd145ee3ada1a1c9dbc68c779b461680976))
* run in parallel ([3499c07](https://github.com/nobbs/kubectl-mapr-ticket/commit/3499c078562bf7b1c5b0aa72f89a078a13c845f5))


### Continuous Integration

* add pre-commit config ([7875ffd](https://github.com/nobbs/kubectl-mapr-ticket/commit/7875ffd145ee3ada1a1c9dbc68c779b461680976))
* exclude test files from funlen check ([7875ffd](https://github.com/nobbs/kubectl-mapr-ticket/commit/7875ffd145ee3ada1a1c9dbc68c779b461680976))
* **lint:** tune golangci-lint, add gci linter ([69d7edd](https://github.com/nobbs/kubectl-mapr-ticket/commit/69d7eddadf334cb0d21eb7d66bf4b5aa3c1c1ce6))
* only build if test and lint pass ([3499c07](https://github.com/nobbs/kubectl-mapr-ticket/commit/3499c078562bf7b1c5b0aa72f89a078a13c845f5))


### Documentation

* add badges to README.md ([fba5b28](https://github.com/nobbs/kubectl-mapr-ticket/commit/fba5b283e6c14c97d56bd03fddd951b0e8815f5f))
* add documentation for completion functions ([5253528](https://github.com/nobbs/kubectl-mapr-ticket/commit/525352894f7da34713bd5734067b894a8cae541c))
* update README.md ([3499c07](https://github.com/nobbs/kubectl-mapr-ticket/commit/3499c078562bf7b1c5b0aa72f89a078a13c845f5))


### Code Refactoring

* cleaning up the codebase even more ([#49](https://github.com/nobbs/kubectl-mapr-ticket/issues/49)) ([a39fb61](https://github.com/nobbs/kubectl-mapr-ticket/commit/a39fb617a2a4f8b277865f1d01c74c678679e441))
* major code reorganization ([#47](https://github.com/nobbs/kubectl-mapr-ticket/issues/47)) ([03258b5](https://github.com/nobbs/kubectl-mapr-ticket/commit/03258b5deb7c7506af361594fc7d0891542c891f))
* move cli to internal package ([5253528](https://github.com/nobbs/kubectl-mapr-ticket/commit/525352894f7da34713bd5734067b894a8cae541c))
* move duration functions util functions ([7875ffd](https://github.com/nobbs/kubectl-mapr-ticket/commit/7875ffd145ee3ada1a1c9dbc68c779b461680976))
* remove --all flag from volume command ([04d7a82](https://github.com/nobbs/kubectl-mapr-ticket/commit/04d7a82ba0ea9a5d42e27dfdf2cac3fa60e5208a))
* remove duplicate util definitions ([7875ffd](https://github.com/nobbs/kubectl-mapr-ticket/commit/7875ffd145ee3ada1a1c9dbc68c779b461680976))
* rename ListItem to TicketSecret ([5253528](https://github.com/nobbs/kubectl-mapr-ticket/commit/525352894f7da34713bd5734067b894a8cae541c))
* simplify namespace handling ([5253528](https://github.com/nobbs/kubectl-mapr-ticket/commit/525352894f7da34713bd5734067b894a8cae541c))
* update secret filter implementation ([5253528](https://github.com/nobbs/kubectl-mapr-ticket/commit/525352894f7da34713bd5734067b894a8cae541c))


### Miscellaneous Chores

* major refactoring and test coverage improvements ([#39](https://github.com/nobbs/kubectl-mapr-ticket/issues/39)) ([5253528](https://github.com/nobbs/kubectl-mapr-ticket/commit/525352894f7da34713bd5734067b894a8cae541c))
* prepare debug logging ([#37](https://github.com/nobbs/kubectl-mapr-ticket/issues/37)) ([7875ffd](https://github.com/nobbs/kubectl-mapr-ticket/commit/7875ffd145ee3ada1a1c9dbc68c779b461680976))
* use charmbracelet/log for debug logging ([7875ffd](https://github.com/nobbs/kubectl-mapr-ticket/commit/7875ffd145ee3ada1a1c9dbc68c779b461680976))

## [0.2.1](https://github.com/nobbs/kubectl-mapr-ticket/compare/v0.2.0...v0.2.1) (2024-01-08)


### Bug Fixes

* add option to sort secrets by number of PVCs using them ([#36](https://github.com/nobbs/kubectl-mapr-ticket/issues/36)) ([cc43942](https://github.com/nobbs/kubectl-mapr-ticket/commit/cc439424f037a503b83a92985dc2cb4b436e9b8c))
* set namespace correctly for `used-by` command ([#33](https://github.com/nobbs/kubectl-mapr-ticket/issues/33)) ([b6163fa](https://github.com/nobbs/kubectl-mapr-ticket/commit/b6163fa064c274078ed4e171e0fa369191ed69c0)), closes [#31](https://github.com/nobbs/kubectl-mapr-ticket/issues/31)


### Documentation

* add shell completion instructions ([#34](https://github.com/nobbs/kubectl-mapr-ticket/issues/34)) ([d7e47fb](https://github.com/nobbs/kubectl-mapr-ticket/commit/d7e47fb7d4a0ab5ede5d32a937735e74dc146fb9)), closes [#32](https://github.com/nobbs/kubectl-mapr-ticket/issues/32)
* use curl instead of wget ([d7e47fb](https://github.com/nobbs/kubectl-mapr-ticket/commit/d7e47fb7d4a0ab5ede5d32a937735e74dc146fb9))

## [0.2.0](https://github.com/nobbs/kubectl-mapr-ticket/compare/v0.1.2...v0.2.0) (2024-01-07)


### Features

* add shell completions for cli flags and args ([6264bfd](https://github.com/nobbs/kubectl-mapr-ticket/commit/6264bfdd7a55ec228206d620073e97aa33f26a6a))
* clean up cli implementation, add completions ([#30](https://github.com/nobbs/kubectl-mapr-ticket/issues/30)) ([6264bfd](https://github.com/nobbs/kubectl-mapr-ticket/commit/6264bfdd7a55ec228206d620073e97aa33f26a6a))


### Bug Fixes

* remove `--all-namespaces` flag from `usedby` ([6264bfd](https://github.com/nobbs/kubectl-mapr-ticket/commit/6264bfdd7a55ec228206d620073e97aa33f26a6a))


### Continuous Integration

* **release:** add `krew` update action ([ea72720](https://github.com/nobbs/kubectl-mapr-ticket/commit/ea72720eec131e2b4fbf0f1654c65216c8c8487f))


### Documentation

* add `krew` installation method ([ea72720](https://github.com/nobbs/kubectl-mapr-ticket/commit/ea72720eec131e2b4fbf0f1654c65216c8c8487f))


### Code Refactoring

* move cli help text to constants ([6264bfd](https://github.com/nobbs/kubectl-mapr-ticket/commit/6264bfdd7a55ec228206d620073e97aa33f26a6a))


### Miscellaneous Chores

* add `krew` release manifest template ([#5](https://github.com/nobbs/kubectl-mapr-ticket/issues/5)) ([ea72720](https://github.com/nobbs/kubectl-mapr-ticket/commit/ea72720eec131e2b4fbf0f1654c65216c8c8487f))
* add helper script to enable kubectl completion ([6264bfd](https://github.com/nobbs/kubectl-mapr-ticket/commit/6264bfdd7a55ec228206d620073e97aa33f26a6a))
* fix changelog section order ([#27](https://github.com/nobbs/kubectl-mapr-ticket/issues/27)) ([7d27551](https://github.com/nobbs/kubectl-mapr-ticket/commit/7d275513fa83993c056eb5372685eab88dae0b23))

## [0.1.2](https://github.com/nobbs/kubectl-mapr-ticket/compare/v0.1.1...v0.1.2) (2024-01-02)


### Features

* add `--expires-before &lt;duration&gt;` to `list` ([#17](https://github.com/nobbs/kubectl-mapr-ticket/issues/17)) ([3b802ff](https://github.com/nobbs/kubectl-mapr-ticket/commit/3b802ffd5370c34c2df99ee4785b75c89489a685)), closes [#16](https://github.com/nobbs/kubectl-mapr-ticket/issues/16)
* add `json` and `yaml` output options to list ([f72175e](https://github.com/nobbs/kubectl-mapr-ticket/commit/f72175e98ac459a5985a9f8e27c6e9daad4456e6))
* add `used-by` command ([#13](https://github.com/nobbs/kubectl-mapr-ticket/issues/13)) ([30bacf4](https://github.com/nobbs/kubectl-mapr-ticket/commit/30bacf4f1997b46253060931b9749c18c3b20159))
* add `used-by` option to `list` command ([30bacf4](https://github.com/nobbs/kubectl-mapr-ticket/commit/30bacf4f1997b46253060931b9749c18c3b20159))
* add `wide` output option ([f72175e](https://github.com/nobbs/kubectl-mapr-ticket/commit/f72175e98ac459a5985a9f8e27c6e9daad4456e6))
* add filtering options to list command ([#10](https://github.com/nobbs/kubectl-mapr-ticket/issues/10)) ([98728bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/98728bcfa2e0141d5976ff37c2c2560e447105db))
* add sort option to `list` command ([#25](https://github.com/nobbs/kubectl-mapr-ticket/issues/25)) ([0ea6727](https://github.com/nobbs/kubectl-mapr-ticket/commit/0ea672751a8d6e719624e522bcf2e8dc5401656e)), closes [#19](https://github.com/nobbs/kubectl-mapr-ticket/issues/19)
* add uid and gid filters to list command ([f72175e](https://github.com/nobbs/kubectl-mapr-ticket/commit/f72175e98ac459a5985a9f8e27c6e9daad4456e6))
* implement `used-by` command to find ticket using pvs ([30bacf4](https://github.com/nobbs/kubectl-mapr-ticket/commit/30bacf4f1997b46253060931b9749c18c3b20159))
* improve `list` output ([30bacf4](https://github.com/nobbs/kubectl-mapr-ticket/commit/30bacf4f1997b46253060931b9749c18c3b20159))
* more output options for list command ([#12](https://github.com/nobbs/kubectl-mapr-ticket/issues/12)) ([f72175e](https://github.com/nobbs/kubectl-mapr-ticket/commit/f72175e98ac459a5985a9f8e27c6e9daad4456e6))
* show number of PVs using a ticket on `list -i` ([#24](https://github.com/nobbs/kubectl-mapr-ticket/issues/24)) ([48cd830](https://github.com/nobbs/kubectl-mapr-ticket/commit/48cd83035fbe2071814efe22a63c79cbd445bf54)), closes [#20](https://github.com/nobbs/kubectl-mapr-ticket/issues/20)


### Bug Fixes

* **deps:** update k8s.io/utils digest to e7106e6 ([#18](https://github.com/nobbs/kubectl-mapr-ticket/issues/18)) ([599fb42](https://github.com/nobbs/kubectl-mapr-ticket/commit/599fb428be2e191e25a8e8ba679b1df37e3bc62b))
* **deps:** update module sigs.k8s.io/yaml to v1.4.0 ([#14](https://github.com/nobbs/kubectl-mapr-ticket/issues/14)) ([d0ab4f5](https://github.com/nobbs/kubectl-mapr-ticket/commit/d0ab4f56466ba5edde8b012817f2fc653ab14bf5))
* global flags now properly global ([98728bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/98728bcfa2e0141d5976ff37c2c2560e447105db))
* handle duration overflow (&gt; 292 years) ([#26](https://github.com/nobbs/kubectl-mapr-ticket/issues/26)) ([c6bc824](https://github.com/nobbs/kubectl-mapr-ticket/commit/c6bc82469e04002964119c34d8358472c0675fdc)), closes [#23](https://github.com/nobbs/kubectl-mapr-ticket/issues/23)
* list dereferenced from secret in loop ([ad3b089](https://github.com/nobbs/kubectl-mapr-ticket/commit/ad3b08989c5116a8e186f013bbe7294d5017ec6d))
* ticket expiry check inverted ([98728bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/98728bcfa2e0141d5976ff37c2c2560e447105db))


### Continuous Integration

* **release:** also update versions in README ([30bacf4](https://github.com/nobbs/kubectl-mapr-ticket/commit/30bacf4f1997b46253060931b9749c18c3b20159))
* **renovate:** run `go mod tidy` before PRs ([951d01e](https://github.com/nobbs/kubectl-mapr-ticket/commit/951d01e3c1f1b2e8dfc28d36857cc4e0dd6eb247))


### Documentation

* add `used-by` to README ([30bacf4](https://github.com/nobbs/kubectl-mapr-ticket/commit/30bacf4f1997b46253060931b9749c18c3b20159))


### Code Refactoring

* rename some packages ([4df6597](https://github.com/nobbs/kubectl-mapr-ticket/commit/4df65970b84fd12c69004912236d0a3bf3b9bd65))
* restructure internal code to support filtering ([98728bc](https://github.com/nobbs/kubectl-mapr-ticket/commit/98728bcfa2e0141d5976ff37c2c2560e447105db))


### Tests

* add empty tests to fix coverage reporting ([#22](https://github.com/nobbs/kubectl-mapr-ticket/issues/22)) ([47bce7a](https://github.com/nobbs/kubectl-mapr-ticket/commit/47bce7a39d83fba4ee4597f4eb1364c094478366)), closes [#21](https://github.com/nobbs/kubectl-mapr-ticket/issues/21)


### Miscellaneous Chores

* **deps:** update go.mod ([30bacf4](https://github.com/nobbs/kubectl-mapr-ticket/commit/30bacf4f1997b46253060931b9749c18c3b20159))
* **deps:** use v0.1.1 of mapr-ticket-parser ([844f00e](https://github.com/nobbs/kubectl-mapr-ticket/commit/844f00efb79f58588d5a4ff1a4af7043a3077990))
* **lint:** disable complexity linters for now ([30bacf4](https://github.com/nobbs/kubectl-mapr-ticket/commit/30bacf4f1997b46253060931b9749c18c3b20159))

## [0.1.1](https://github.com/nobbs/kubectl-mapr-ticket/compare/v0.1.0...v0.1.1) (2023-12-30)


### Continuous Integration

* add codecov reporting ([#8](https://github.com/nobbs/kubectl-mapr-ticket/issues/8)) ([05f0d65](https://github.com/nobbs/kubectl-mapr-ticket/commit/05f0d65a897637f85b496eb7ce7f8975b1295f7d))
* **build:** also package LICENSE and README into release assets ([#6](https://github.com/nobbs/kubectl-mapr-ticket/issues/6)) ([030e4c1](https://github.com/nobbs/kubectl-mapr-ticket/commit/030e4c1e9904b7d1bfbdddef188d0af96273e464))
* **tests:** use CGO as race won't work otherwise ([05f0d65](https://github.com/nobbs/kubectl-mapr-ticket/commit/05f0d65a897637f85b496eb7ce7f8975b1295f7d))


### Documentation

* add readme ([#9](https://github.com/nobbs/kubectl-mapr-ticket/issues/9)) ([dca7660](https://github.com/nobbs/kubectl-mapr-ticket/commit/dca766070bf207e97e4a4e51506723a0405b1cc7))


### Code Refactoring

* add proper usage string to cli ([dca7660](https://github.com/nobbs/kubectl-mapr-ticket/commit/dca766070bf207e97e4a4e51506723a0405b1cc7))

## 0.1.0 (2023-12-30)


### Features

* initial implementation of ticket list command ([8043689](https://github.com/nobbs/kubectl-mapr-ticket/commit/80436895cf07b160faf3ad6a3a4fda3999ec7425))


### Build System

* fix go version in go.mod ([471b8db](https://github.com/nobbs/kubectl-mapr-ticket/commit/471b8db26887563169542a822e7b2968e0248ba8))


### Miscellaneous Chores

* configure Renovate ([#2](https://github.com/nobbs/kubectl-mapr-ticket/issues/2)) ([a3002e5](https://github.com/nobbs/kubectl-mapr-ticket/commit/a3002e51e14a7819669532477322b88bfb6ffff8))


### Continuous Integration

* add build on release ([763f003](https://github.com/nobbs/kubectl-mapr-ticket/commit/763f00300cb484d4429d2d4f3bd982720ffbff38))
* **build:** disable CGO ([1608ca2](https://github.com/nobbs/kubectl-mapr-ticket/commit/1608ca2d3a2a8160d0726e5d4402b613bfd9c7fb))
* **build:** multi-arch and multi-os build ([2ff900b](https://github.com/nobbs/kubectl-mapr-ticket/commit/2ff900bfc99dad05e7fd1b3fae362a58ec2e6353))
* **build:** overwrite release assets ([de50d0c](https://github.com/nobbs/kubectl-mapr-ticket/commit/de50d0c88d977e612bc51e26121dfa9b0fc75781))
