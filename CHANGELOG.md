# Changelog

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
