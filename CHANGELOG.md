# Changelog

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


### Build System

* fix go version in go.mod ([471b8db](https://github.com/nobbs/kubectl-mapr-ticket/commit/471b8db26887563169542a822e7b2968e0248ba8))


### Miscellaneous Chores

* configure Renovate ([#2](https://github.com/nobbs/kubectl-mapr-ticket/issues/2)) ([a3002e5](https://github.com/nobbs/kubectl-mapr-ticket/commit/a3002e51e14a7819669532477322b88bfb6ffff8))


### Continuous Integration

* add build on release ([763f003](https://github.com/nobbs/kubectl-mapr-ticket/commit/763f00300cb484d4429d2d4f3bd982720ffbff38))
* **build:** disable CGO ([1608ca2](https://github.com/nobbs/kubectl-mapr-ticket/commit/1608ca2d3a2a8160d0726e5d4402b613bfd9c7fb))
* **build:** multi-arch and multi-os build ([2ff900b](https://github.com/nobbs/kubectl-mapr-ticket/commit/2ff900bfc99dad05e7fd1b3fae362a58ec2e6353))
* **build:** overwrite release assets ([de50d0c](https://github.com/nobbs/kubectl-mapr-ticket/commit/de50d0c88d977e612bc51e26121dfa9b0fc75781))


### Features

* initial implementation of ticket list command ([8043689](https://github.com/nobbs/kubectl-mapr-ticket/commit/80436895cf07b160faf3ad6a3a4fda3999ec7425))
