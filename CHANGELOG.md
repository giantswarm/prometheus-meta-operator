# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [Unreleased]

### Added

- Add service monitor for nginx-ingress-controller

### Removed

- Removed service resource as this is no longer needed (it was used for cortex frontend)

### Fixed

- Fix an error during alert update: metadata.resourceVersion: Invalid value

## [0.1.1] - 2020-05-27

### Added

- Change chart namespace from giantswarm to monitoring

## [0.1.0] - 2020-05-27

### Added

- First release.


[Unreleased]: https://github.com/giantswarm/aws-operator/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/giantswarm/aws-operator/releases/tag/v0.1.1
[0.1.0]: https://github.com/giantswarm/aws-operator/releases/tag/v0.1.0
