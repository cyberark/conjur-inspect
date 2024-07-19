# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed
- Nothing should go in this section, please add to the latest unreleased version
  (and update the corresponding date), or add a new version.

## [0.4.0] - 2024-04-24

### Added
- Added `--container-id` argument, which takes either a container ID or name
  to run container specific checks against. CNJR-4620
- Added container specific check to record the `podman inspect` or
  `docker inspect` output for the specified container. CNJR-2399
- The JSON inspection report is now included in the output archive file. This
  allows the single file to include all relevant data.
  CNJR-1649
- Added Conjur health and info checks for the specified container. CNJR-2397
- Added container logs to the output archive file, when a container ID is
  specified. CNJR-2395

## [0.3.0] - 2023-01-26

### Changed
- Renamed `conjur-preflight` to `conjur-inspect`.
  [cyberark/conjur-inspect#30](https://github.com/cyberark/conjur-inspect/pull/30)

### Added
- `conjur-preflight` now has a follower port connectivity report where it checks
  if the required ports that the follower needs to be deployed are open.
  [conjurinc/conjur-preflight#18](https://github.com/conjurinc/conjur-preflight/pull/18)
- `conjur-preflight` now has a cli option flag that can be used to output
  JSON formatted reports
  [conjurinc/conjur-preflight#18](https://github.com/conjurinc/conjur-preflight/pull/23)
- Added a progress indicator for the running checks.
  [conjurinc/conjur-preflight#24](https://github.com/conjurinc/conjur-preflight/pull/24)
- Raw report data is now collected during an inspection and saved as a report
  result archive. The name and location of this archived are managed with the
  new `--report-id` and `--raw-data-dir` command arguments.
  [conjurinc/conjur-preflight#26](https://github.com/conjurinc/conjur-preflight/pull/26)
- Added container runtime checks to report on the available Docker and/or
  Podman version and configuration.
  [cyberark/conjur-inspect#15](https://github.com/cyberark/conjur-inspect/pull/15)
- Added more debug logs for shell command execution to improve error
  troubleshooting experience.
  [cyberark/conjur-inspect#31](https://github.com/cyberark/conjur-inspect/pull/31)

### Fixed
- Previously, if the user running `conjur-preflight` has insufficient permission
  to access a partition mountpoint, then the check would result in a segfault.
  Now, this logs a warning and skips the mountpoint in the report.
  [conjurinc/conjur-preflight#28](https://github.com/conjurinc/conjur-preflight/pull/28)

### Security
- Forced golang.org/x/net to use v0.7.0
  [cyberark/conjur-inspect#34](https://github.com/cyberark/conjur-inspect/pull/34)

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
