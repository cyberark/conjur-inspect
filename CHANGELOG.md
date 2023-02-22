# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed
- Nothing should go in this section, please add to the latest unreleased version
  (and update the corresponding date), or add a new version.

## [0.3.0] - 2023-01-26

### Added
- `conjur-preflight` now has a follower port connectivity report where it checks
  if the required ports that the follower needs to be deployed are open.
  [conjurinc/conjur-preflight#18](https://github.com/conjurinc/conjur-preflight/pull/18)
- `conjur-preflight` now has a cli option flag that can be used to output
  JSON formatted reports
  [conjurinc/conjur-preflight#18](https://github.com/conjurinc/conjur-preflight/pull/23)

## [0.2.0] - 2023-01-20

### Added
- `conjur-preflight` now detects whether it is writing standard output to a
  terminal or to a file/pipe and formats with rich or plain text accordingly.
  [conjurinc/conjur-preflight#19](https://github.com/conjurinc/conjur-preflight/pull/19)
- A new CLI flag `--debug` causes `conjur-preflight` to log more verbose
  information about the execution of the application and its checks.
  [conjurinc/conjur-preflight#19](https://github.com/conjurinc/conjur-preflight/pull/19)
- `conjur-preflight` now includes disk related checks for read, write, and sync
  latency, as well as read and write operations per second (IOPs). These require
  `fio` and `libaio` to be present as a prerequesite for these checks.
  [conjurinc/conjur-preflight#19](https://github.com/conjurinc/conjur-preflight/pull/19)

### Fixed
- Previously, the application version was not properly embedded in the final
  executable. Now the application and the reports it produces reflect the
  correct version number.
  [conjurinc/conjur-preflight#19](https://github.com/conjurinc/conjur-preflight/pull/19)

### Security
- Added replace statements to go.mod to prune dependencies with known vulnerabilities from
  the dependency tree.
  [conjurinc/conjur-preflight#21](https://github.com/conjurinc/conjur-preflight/pull/21)
  [conjuring/conjur-preflight#22](https://github.com/conjurinc/conjur-preflight/pull/22)

## [0.1.0] - 2023-01-12

### Added
- Initial reports for CPU, memory, disk space, and OS version.
  [conjurinc/conjur-preflight#14](https://github.com/conjurinc/conjur-preflight/pull/14)

## [0.0.0] - 2022-12-20

### Changed
- Reset repository from bash-based CLI to golang CLI.

[Unreleased]: https://github.com/conjurinc/conjur-preflight/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/conjurinc/conjur-preflight/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/conjurinc/conjur-preflight/compare/v0.0.0...v0.1.0
