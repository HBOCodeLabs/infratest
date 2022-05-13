# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/2.0.0.html).

## Unreleased

### Maintenance
* A CHANGELOG file has been added, and the CI process now ensures it has been updated on pull
  requests.
* The workflow that runs on commits to the `main` branch has been updated to not generate releases
  any more. Instead, a newly created `Release` workflow will be used, which automatically updates
  the CHANGELOG file with the version input at run time (the workflow is manually invoked whenever
  we want to generate a release).

### Added
* A new method, `vault.AssertSecretExits`, for asserting that secrets exist in Hashicorp
  [Vault](https://vaultproject.io).

## [v0.9.0] - 2022-05-20

## [v0.8.0] - 2022-05-13

## [v0.7.2] - 2022-03-11

## [v0.7.1] - 2022-03-01

## [v0.7.0] - 2022-02-23

## [v0.6.0] - 2022-02-01

## [v0.5.0] - 2022-02-01

## [v0.4.0] - 2022-01-21

## [v0.3.0] - 2021-12-21

## [v0.2.0] - 2021-12-21

## [v0.1.2] - 2021-11-15

## [v0.1.0] - 2021-11-11


[v0.9.0]: https://github.com/hbocodelabs/infratest/compare/v0.8.v0...v0.9.0
[v0.8.0]: https://github.com/hbocodelabs/infratest/compare/v0.7.2...v0.8.0
[v0.7.2]: https://github.com/hbocodelabs/infratest/compare/v0.7.1...v0.7.2
[v0.7.1]: https://github.com/hbocodelabs/infratest/compare/v0.7.v0...v0.7.1
[v0.7.0]: https://github.com/hbocodelabs/infratest/compare/v0.6.v0...v0.7.0
[v0.6.0]: https://github.com/hbocodelabs/infratest/compare/v0.5.v0...v0.6.0
[v0.5.0]: https://github.com/hbocodelabs/infratest/compare/v0.4.v0...v0.5.0
[v0.4.0]: https://github.com/hbocodelabs/infratest/compare/v0.3.v0...v0.4.0
[v0.3.0]: https://github.com/hbocodelabs/infratest/compare/v0.2.v0...v0.3.0
[v0.2.0]: https://github.com/hbocodelabs/infratest/compare/v0.1.2...v0.2.0
[v0.1.2]: https://github.com/hbocodelabs/infratest/compare/v0.1.v0...v0.1.2
[v0.1.0]: https://github.com/hbocodelabs/infratest/releases/tag/v0.1.0
