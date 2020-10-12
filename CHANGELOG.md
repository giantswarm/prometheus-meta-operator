# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [Unreleased]

## [1.6.0] - 2020-10-12

### Changed

- Set retention size to 90Gi and duration to 2w
- Increased storage to 100Gi
- Increased memory limit to 300Mi

## [1.5.1] - 2020-10-07

### Fixed

- Fix promxy config marshaling
- Fix promxy config not being updated

## [1.5.0] - 2020-10-07

### Added

- Support for managing Promxy configuration

### Removed

- Old namespace deleter resource

## [1.4.0] - 2020-09-25

### Added

- Add oauth ingress
- Add tls certificate for ingress
- Add ingress for individual prometheuses

## [1.3.0] - 2020-09-24

### Added

- Scraping of tenant cluster prometheus
- Scraping of control plane prometheus
- Add installation label
- Add labelling schema alert

### Changed

- Set honor labels to true
- Change control plane namespace to reflect the installation name instead of 'kubernetes'

## [1.2.0] - 2020-09-03

### Added

- Add monitoring label
- Add etcd target for control planes
- Add vault target
- Add gatekeeper target
- Add managed-app target
- Add cert-operator target
- Add bridge-operator target
- Add flannel-operator target
- Add ingress-exporter target
- Add coreDNS target
- Add azure-collector target

### Removed

- frontend, ingress, and service resources.

### Fixed

- prevented data loss in `Cluster` resources by always using the correct
  version of the type as configured in CRDs storage version (#101)
- avoids trying to read dependant objects from the cluster when processing
  deletion, as they may be gone already and errors here were disrupting cleanup
  and preventing the finalizer from being removed (#115)

## [1.1.0] - 2020-08-27

### Added

- Scraping of the control plane operators
    - aws-operator
    - azure-operator
    - kvm-operator
    - app-operator
    - chart-operator
    - cluster-operator
    - etcd-backup-operator
    - node-operator
    - release-operator
    - organization-operator
    - prometheus-meta-operator
    - rbac-operator
    - draughtsman
- Scraping of the monitoring targets
    - app-exporter
    - cert-exporter
    - vault-exporter
    - node-exporter
    - net-exporter
    - kube-state-metrics
    - alertmanager
    - grafana
    - prometheus
    - prometheus-config-controller
    - fluentbit
- Scraping of the control plane apis
    - tokend
    - companyd
    - userd
    - api
    - kubernetesd
    - credentiald
    - cluster-service
- New control-plane controller, reconciling kubernetes api service (#92)

## [1.0.1] - 2020-08-25

### Changed

- Rename controller name and finalizers

## [1.0.0] - 2020-08-20

### Added

- Scraping of kube-proxy (#88)
- Scraping of kube-scheduler (#87)
- Scraping of kube-controller-manager (#85)
- Scraping of etcd (#81)
- Scraping of kubelet (#82)
- Scraping of legacy docker, calico-node, cluster-autoscaler, aws-node and cadvisor (#78)

### Changed

- Moved prometheus storage from `emptyDir` to a `persistentVolumeClaim`
- Remove tenant cluster prometheus limits
- Updated backward incompatible Kubernetes dependencies to v1.18.5.

## [0.3.2] - 2020-07-24

### Changed

- Set TC prometheus memory limit to 1Gi (#73)

## [0.3.1] - 2020-07-17

### Changed

- Set TC prometheus memory limit to 200Mi

## [0.3.0] - 2020-07-15

### Changed

- Scale prometheus-meta-operator replicas back to one.

### Added

- Set prometheus request/limits (cpu: 100m, memory: 100Mi)

## [0.2.1] - 2020-07-01

### Fixed

- Fixed release process

## [0.2.0] - 2020-06-29

### Added

- Add service monitor for nginx-ingress-controller
- Reconcile CAPI (Cluster) and legacy cluster CRs (AWSConfig, AzureConfig, KVMConfig)

### Changed

- Reduced prometheus server replicas to one (#45)
- Reduced default prometheus-meta-operator replicas to zero as having both this and previous (g8s-prometheus) solutions on at the same time is overloading some control planes

### Removed

- Removed cortex frontend as it's an optimisation that's not currently needed
- Removed service and ingress resources as they are no longer needed (they were used for the cortex frontend)

### Fixed

- Fix an error during alert update: metadata.resourceVersion: Invalid value

## [0.1.1] - 2020-05-27

### Added

- Change chart namespace from giantswarm to monitoring

## [0.1.0] - 2020-05-27

### Added

- First release.


[Unreleased]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.6.0...HEAD
[1.6.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.5.1...v1.6.0
[1.5.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.5.0...v1.5.1
[1.5.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v0.3.2...v1.0.0
[0.3.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v0.2.1...v0.3.0
[0.2.1]: https://github.com/giantswarm/prometheus-meta-operator/releases/tag/v0.2.1
[0.2.0]: https://github.com/giantswarm/prometheus-meta-operator/releases/tag/v0.2.0
[0.1.1]: https://github.com/giantswarm/prometheus-meta-operator/releases/tag/v0.1.1
[0.1.0]: https://github.com/giantswarm/prometheus-meta-operator/releases/tag/v0.1.0
