# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.41.1] - 2021-06-22

### Added

- Add term to not count api-server errors for clusters in transitioning state.
- Business-hours alert for azure clusters not updating in time.

### Changed

- Increase `ManagementClusterWebhookDurationExceedsTimeout` duration from 5m to 15m.

### Fixed

- Fix CoreDNSMaxHPAReplicasReached alert to not fire in case max and min are equal.
- Business-hours alert for azure clusters not creating in time.

### Removed

- Remove AlertManager ingress to avoid conflicts with the existing one, until the new AlertManager is ready to replace the one from _g8s-prometheus_

## [1.41.0] - 2021-06-17

### Added

- Add `AppPendingUpdate` alert.
- Add scrapeconfig for `falco-exporter` on management clusters.
- Add Alertmanager managed by Prometheus Operator.
- Add Alertmanager ingress.
- Add `WorkloadClusterDeploymentNotSatisfiedLudacris` to monitor `metrics-server` in workload clusters.
- Add `CoreDNSMaxHPAReplicasReached` business hours alert for when CoreDNS has been scaled to its maximum for too long.

### Changed

- Lower Prometheus disk space alert from 10% to 5%.
- Change severity of `ChartOperatorDown` alert to notify.
- Merge all provider certificate.management-cluster.rules into one prometheus rule.

### Fixed

- Fix service name in ingress.

## [1.40.0] - 2021-06-14

### Changed

- Lower `kubelet` SLO from 99.9% to 99%.

## [1.39.0] - 2021-06-11

### Added

- Add ServiceLevelBurnRateTicket alert.
- Add the prometheus log level option
- Add high and low burn rates as recording rules.

### Changed

- Move managed apps SLO alerts to the service-level format.
- Set `HighNumberOfAllocatedSockets` to notify not page
- Extract `kubelet` and `api-server` SLO targets to their own recording rules.
- Extract `kubelet` and `api-server` alerting thresholds to their own recording rules.
- Change `ServiceLevelBurnRateTooHigh` to use new created values.

### Fixed

- Fixed the way VPA `maxAllowed` parameter for memory is calculated so that we
  avoid going over node memory capacity with the memory limit (`maxAllowed` is
  used for request and limit is that multiplied by 1.2).

## [1.38.0] - 2021-05-28

### Changed

- Increased alert duration of `PrometheusCantCommunicateWithKubernetesAPI`.
- Refactor resources to namespace monitoring and alerting code.
- Add cluster-autoscaler to `WorkloadClusterContainerIsRestartingTooFrequentlyFirecracker`

### Removed

- Remove `tlscleanup` and `volumeresizehack` resources as they are not needed anymore.

## [1.37.0] - 2021-05-26

### Added

- Add HTTP proxy support to Prometheus Remote Write.

## [1.36.0] - 2021-05-25

### Added

- Added alert `HighNumberOfAllocatedSockets` for High number of allocated sockets
- Added alert `HighNumberOfOrphanedSockets` for High number of orphaned sockets
- Added alert `HighNumberOfTimeWaitSockets` for High number of time wait sockets
- Added alert `AWSWorkloadClusterNodeTooManyAutoTermination` for terminate unhealthy feature.
- Preserve and merge global HTTP client config when generating heartbeat
  receivers in AlertManager config; this allows it to be used in environments
  where internet access is only allowed through a proxy.

### Changed

- Include `cluster-api-core-unique-webhook` into `DeploymentNotSatisfiedFirecracker` and `DeploymentNotSatisfiedChinaFirecracker`.
- Increased duration for `PrometheusPersistentVolumeSpaceTooLow` alert
- Increased duration for `WorkloadClusterEtcdDBSizeTooLarge` alert.
- Increased duration for `WorkloadClusterEtcdHasNoLeader` alert.
- Silence `OperatorkitErrorRateTooHighCelestial` and `OperatorkitCRNotDeletedCelestial` outside working hours.
- Update Prometheus to 2.27.1
- Add atlas, and installation tag onto Heartbeats.

### Fixed

- Fix `PrometheusFailsToCommunicateWithRemoteStorageAPI` alert not firing on china clusters.

## [1.35.0] - 2021-05-12

### Added

- Add alert `alertmanager-dashboard` not satisfied.

## [1.34.1] - 2021-05-10

### Fixed

- inhibit KubeStateMetricsDown and KubeStateMetricsMissing

## [1.34.0] - 2021-05-06

### Changed

- Lower the severity to notify for managed app's error budget alerts

### Fixed

- Fix ManagedApp alert
- Fix `InhibitionKubeStateMetricsDown` not firing long enough

## [1.33.0] - 2021-04-27

### Changed

- Raise prometheus cpu limit to 150%.

### Removed

- Remove `PodLimitAlmostReachedAWS` and `EBSVolumeMountErrors` alerts as they were not used.

## [1.32.1] - 2021-04-22

### Fixed

- Adjust container restarting too often firecracker.

## [1.32.0] - 2021-04-19

### Added

- Add alert for `kube-state-metrics` missing.
- Tune remote write configuration to avoid loss of data.

### Changed

- Only fire `KubeStateMetricsDown` if `kube-state-metrics` is down.

## [1.31.0] - 2021-04-16

### Added

- Page firecracker for failed cluster transitions.
- Page Firecracker in working hours for restarting containers.
- Add recording rules for kube-mixins
- `MatchingNumberOfPrometheusAndCluster` now has a runbook, link added to alert.

### Changed

- Keep the `container_network.*` metrics as they are needed for the [kubernetes mixins dashboards](https://github.com/giantswarm/g8s-grafana/tree/master/helm/g8s-grafana/dashboards/mixin)

## [1.30.0] - 2021-04-12

### Removed

- Remove Gatekeeper alerts and targets.

## [1.29.1] - 2021-04-09

### Fixed

- Fix inhibition for `MatchingNumberOfPrometheusAndCluster` alert by matching it with
  source from Management Cluster instead of the cluster the alert is firing for.

## [1.29.0] - 2021-04-09

### Added

- Add `PrometheusCantCommunicateWithRemoteStorageAPI` to alert when Prometheus fails to send samples to Cortex.
- Add workload type and name labels for `ManagedAppBasicError*` alerts
- Add alert for master node in HA setup down for too long.
- Add aggregation for docker actions.

### Fixed

- Fix prometheus storage alert

### Removed

- Removed unnecessary whitespace in additional scrape configs.

## [1.28.0] - 2021-04-01

### Added

- Add support to calculate maximum CPU.
- Include cadvisor metrics from the pod in `draughtsman` namespace.
- Add `PrometheusPersistentVolumeSpaceTooLow` alert for prometheus storage going over 90 percent.

### Changed

- Split `ManagementClusterCertificateWillExpireInLessThanTwoWeeks` alert per provider.
- Increased duration time for flapping `WorkloadClusterWebhookDurationExceedsTimeout` alert

### Fixed

- Changed prometheus volume space alert ownership to atlas:
  - `PersistentVolumeSpaceTooLow` -> `PrometheusPersistentVolumeSpaceTooLow`

### Removed

- Do not monitor docker for CAPI clusters

### Removed

- Remove promxy resource.

## [1.27.4] - 2021-03-26

- Add recording rules for dex activity, creating the metrics
  - `aggregation:dex_requests_status_ok`
  - `aggregation:dex_requests_status_4xx`
  - `aggregation:dex_requests_status_5xx`

## [1.27.3] - 2021-03-25

- Fix prometheus/common secret token in imported code.

## [1.27.2] - 2021-03-25

### Fixed

- Fix alertmanager secretToken in imported alertmanager code.

## [1.27.1] - 2021-03-25

### Fixed

- Remove follow_redirects from alertmanager config
  - Update prometheus/alertmanger@v0.21.0
  - Update prometheus/common@v0.17.0

## [1.27.0] - 2021-03-24

### Changed

- Update architect to 2.4.2

### Removed

- Removed memory-intensive notify only systemd alerts.

## [1.26.0] - 2021-03-24

### Changed

- Push to `shared-app-collection`
- Rename `EtcdWorkloadClusterDown` to `WorkloadClusterEtcdDown`
- Increased memory limits by 1.2 factor

### Fixed

- Support vmware for `WorkloadClusterEtcdDown`
- Add vmware to the list of valid providers

## [1.25.2] - 2021-03-23

### Fixed

- Disable follow redirect for alertmanager

## [1.25.1] - 2021-03-22

### Fixed

- Set prometheus minimum CPU to 100m

## [1.25.0] - 2021-03-22

### Added

- Add support for monitoring vmware clusters
- Add support to get the API Server URL for both legacy and CAPI clusters

### Changed

- Upgrade ingress version to networking.k8s.io/v1beta1

### Fix

- Fix typo in `MatchingNumberOfPrometheusAndCluster` alert
- Fix scrapeconfig to use secured ports for kubernetes control plane components for CAPI clusters
- Fix scrapeconfig to proxy all calls through the API Server for CAPI clusters

## [1.24.8] - 2021-03-18

### Fix

- Avoid alerting for `MatchingNumberOfPrometheusAndCluster` when a cluster is
  being deleted.


## [1.24.7] - 2021-03-18

### Added

- Add support to copy CAPI cluster's certificates
- Add aggregation `aggregation:giantswarm:api_auth_giantswarm_successful_attempts_total`.

## [1.24.6] - 2021-03-02

### Fixed

- Fix equality check on the VPA CR to prevent it being overriden and losing it's status information on every prometheus-meta-operator deployment.
- Inhibit `MatchingNumberOfPrometheusAndCluster` when kube-state-metrics is down
  to prevent bogus pages when `kube_pod_container_status_running` metric
  isn't available

## [1.24.5] - 2021-03-02

### Added

- Set the prometheus UI Web page title.
- Add 'app' label to metrics pushed from `app-exporter` to cortex

## [1.24.4] - 2021-02-26

### Changed

- Avoid alerting for ETCD backups outside business hours.

## [1.24.3] - 2021-02-24

### Changed

- Use `resident_memory` when calculating docker memory usage.

## [1.24.2] - 2021-02-24

### Added

- Add 'catalog' label to metrics pushed from `app-exporter` to cortex

## [1.24.1] - 2021-02-23

### Fixed

- Fixed syntax error in expressions of `ManagementClusterPodPending*` alerts

## [1.24.0] - 2021-02-23

### Added

- Add Alert for missing prometheus for a workload cluster
- Add `ManagementClusterPodStuckFirecracker` and `WorkloadClusterPodStuckFirecracker` alerts for Firecracker.
- Add `ManagementClusterPodStuckCelestial` alert for Celestial.
- Send samples per second to cortex

### Changed

- Move Cluster Autoscaller app installation/upgrade related alerts to team Batman.

## [1.23.1] - 2021-02-22

### Added

- Add `TestClusterTooOld` for testing installations
- Added Mayu as a scrape target as well as puma's pods

### Changed

- Apply prometheus rule group (which includes
- Discover ETCD targets through the LoadBalancer using the `giantswarm.io/etcd-domain` annotation

### Fixed

- Remove `PersistentVolumeSpaceTooLow` from Workload Clusters.

## [1.23.0] - 2021-02-17

### Added

- Add the sig-customer alerts:
  - `WorkloadClusterCertificateWillExpireInLessThanAMonth`
  - `WorkloadClusterCertificateWillExpireMetricMissing`
- Add the ludacris alerts:
  - `CadvisorDown`
  - `CalicoRestartRateTooHigh`
  - `CertOperatorVaultTokenAlmostExpiredMissing`
  - `CertOperatorVaultTokenAlmostExpired`
  - `ClusterServiceVaultTokenAlmostExpiredMissing`
  - `ClusterServiceVaultTokenAlmostExpired`
  - `CollidingOperatorsLudacris`
  - `CoreDNSCPUUsageTooHigh`
  - `CoreDNSDeploymentNotSatisfied`
  - `CoreDNSLatencyTooHigh`
  - `DeploymentNotSatisfiedLudacris` and assign it to rocket `DeploymentNotSatisfiedRocket`
  - `DockerMemoryUsageTooHigh` for both Ludacris and Biscuit
  - `DockerVolumeSpaceTooLow` for both Ludacris and Biscuit
  - `EtcdVolumeSpaceTooLow` for both Ludacris and Biscuit
  - `JobFailed` renamed to `ManagementClusterJobFailed`
  - `KubeConfigMapCreatedMetricMissing`
  - `KubeDaemonSetCreatedMetricMissing`
  - `KubeDeploymentCreatedMetricMissing`
  - `KubeEndpointCreatedMetricMissing`
  - `KubeNamespaceCreatedMetricMissing`
  - `KubeNodeCreatedMetricMissing`
  - `KubePodCreatedMetricMissing`
  - `KubeReplicaSetCreatedMetricMissing`
  - `KubeSecretCreatedMetricMissing`
  - `KubeServiceCreatedMetricMissing`
  - `KubeStateMetricsDown`
  - `KubeletConditionBad`
  - `KubeletDockerOperationsErrorsTooHigh`
  - `KubeletDockerOperationsLatencyTooHigh`
  - `KubeletPLEGLatencyTooHigh`
  - `KubeletVolumeSpaceTooLow` for both Ludacris and Biscuit
  - `LogVolumeSpaceTooLow` for both Ludacris and Biscuit
  - `MachineAllocatedFileDescriptorsTooHigh`
  - `MachineEntropyTooLow`
  - `MachineLoadTooHigh` and moved it to biscuit
  - `MachineMemoryUsageTooHigh` and moved it to biscuit
  - `ManagementClusterAPIServerAdmissionWebhookErrors`
  - `ManagementClusterAPIServerLatencyTooHigh`
  - `ManagementClusterContainerIsRestartingTooFrequently`
  - `ManagementClusterCriticalSystemdUnitFailed`
  - `ManagementClusterDaemonSetNotSatisfiedLudacris`
  - `ManagementClusterDaemonSetNotSatisfiedLudacris`
  - `ManagementClusterDisabledSystemdUnitActive`
  - `ManagementClusterHighNumberSystemdUnits`
  - `ManagementClusterNetExporterCPUUsageTooHigh`
  - `ManagementClusterSystemdUnitFailed`
  - `ManagementClusterWebhookDurationExceedsTimeout`
  - `Network95thPercentileLatencyTooHigh`
  - `NetworkCheckErrorRateTooHigh`
  - `NodeConnTrackAlmostExhausted`
  - `NodeExporterCollectorFailed`
  - `NodeExporterDeviceError`
  - `NodeExporterDown`
  - `NodeExporterMissing`
  - `NodeHasConstantOOMKills`
  - `NodeStateFlappingUnderLoad`
  - `OperatorNotReconcilingLudacris`
  - `OperatorkitErrorRateTooHighLudacris`
  - `PersistentVolumeSpaceTooLow` for both Ludacris and Biscuit
  - `ReleaseNotReady`
  - `RootVolumeSpaceTooLow` for both Ludacris and Biscuit
  - `SYNRetransmissionRateTooHigh`
  - `ServiceLevelBurnRateTooHigh`
  - `WorkloadClusterAPIServerAdmissionWebhookErrors`
  - `WorkloadClusterAPIServerLatencyTooHigh`
  - `WorkloadClusterCriticalSystemdUnitFailed`
  - `WorkloadClusterDaemonSetNotSatisfiedLudacris`
  - `WorkloadClusterDisabledSystemdUnitActive`
  - `WorkloadClusterHighNumberSystemdUnits`
  - `WorkloadClusterNetExporterCPUUsageTooHigh`
  - `WorkloadClusterSystemdUnitFailed`
  - `WorkloadClusterWebhookDurationExceedsTimeout`

### Changed

- Migrate and rename `EBSVolumeMountErrors` to `ManagementClusterEBSVolumeMountErrors` and `WorkloadClusterEBSVolumeMountErrors`

### Removed

- Removing legacy finalizers resource used to remove old custom resource finalizers

## [1.22.0] - 2021-02-16

### Changed

- Improved inhibition alert `InhibitionClusterStatusUpdating` to  inhibit alerts 10 minutes after the update has finished to avoid unecessery pages.

## [1.21.0] - 2021-02-16

### Changed

- Split `ManagementClusterAppFailed` per team

### Added

- Add the solution engineer alerts:
  - `AzureQuotaUsageApproachingLimit`
  - `NATGatewaysPerVPCApproachingLimit`
  - `ServiceUsageApproachingLimit`

## [1.20.0] - 2021-02-16

### Added

- Add the rocket alerts:
  - `BackendServerUP`
  - `ClockOutOfSyncKVM`
  - `CollidingOperatorsRocket`
  - `DNSCheckErrorRateTooHighKVM`
  - `DNSErrorRateTooHighKVM`
  - `EtcdWorkloadClusterDownKVM`
  - `IngressExporterDown`
  - `KVMManagementClusterDeploymentScaledDownToZero`
  - `KVMNetworkErrorRateTooHigh`
  - `ManagementClusterCriticalPodMetricMissingKVM`
  - `ManagementClusterCriticalPodNotRunningKVM`
  - `ManagementClusterMasterNodeMissingRocket`
  - `ManagementClusterPodLimitAlmostReachedKVM`
  - `ManagementClusterPodPendingFor15Min`
  - `MayuSystemdUnitIsNotRunning`
  - `NetworkInterfaceLeftoverWithoutCluster`
  - `OnpremManagementClusterMissingNodes`
  - `OperatorNotReconcilingRocket`
  - `OperatorkitCRNotDeletedRocket`
  - `OperatorkitErrorRateTooHighRocket`
  - `WorkloadClusterCriticalPodMetricMissingKVM`
  - `WorkloadClusterCriticalPodNotRunningKVM`
  - `WorkloadClusterEndpointIPDown`
  - `WorkloadClusterEtcdCommitDurationTooHighKVM`
  - `WorkloadClusterEtcdDBSizeTooLargeKVM`
  - `WorkloadClusterEtcdHasNoLeaderKVM`
  - `WorkloadClusterEtcdNumberOfLeaderChangesTooHighKVM`
  - `WorkloadClusterMasterNodeMissingRocket`
  - `WorkloadClusterPodLimitAlmostReachedKVM`
- Added the firecracker rules to PMO:
  - `AWSClusterCreationFailed`
  - `AWSClusterUpdateFailed`
  - `AWSManagementClusterDeploymentScaledDownToZero`
  - `AWSManagementClusterMissingNodes`
  - `AWSNetworkErrorRateTooHigh`
  - `ClockOutOfSyncAWS`
  - `CloudFormationStackFailed`
  - `CloudFormationStackRollback`
  - `ClusterAutoscalerAppFailedAWS`
  - `ClusterAutoscalerAppNotInstalledAWS`
  - `ClusterAutoscalerAppPendingInstallAWS`
  - `ClusterAutoscalerAppPendingUpgradeAWS`
  - `CollidingOperatorsFirecracker`
  - `ContainerIsRestartingTooFrequentlyFirecracker`
  - `CredentialdCantReachKubernetes`
  - `DNSCheckErrorRateTooHighAWS`
  - `DNSErrorRateTooHighAWS`
  - `DefaultCredentialsMissing`
  - `DeploymentNotSatisfiedChinaFirecracker`
  - `DeploymentNotSatisfiedFirecracker`
  - `ELBHostsOutOfService`
  - `EtcdWorkloadClusterDownAWS`
  - `FluentdMemoryHighUtilization`
  - `JobHasNotBeenScheduledForTooLong`
  - `KiamMetadataFindRoleErrors`
  - `ManagementClusterDaemonSetNotSatisfiedChinaFirecracker`
  - `ManagementClusterDaemonSetNotSatisfiedFirecracker`
  - `OperatorNotReconcilingFirecracker`
  - `OperatorkitCRNotDeletedFirecracker`
  - `OperatorkitErrorRateTooHighFirecracker`
  - `TooManyCredentialsForOrganization`
  - `TrustedAdvisorErroring`
  - `WorkloadClusterCriticalPodNotRunningAWS`
  - `WorkloadClusterCriticalPodMetricMissingAWS`
  - `WorkloadClusterDaemonSetNotSatisfiedFirecracker`
  - `WorkloadClusterEtcdCommitDurationTooHighAWS`
  - `WorkloadClusterEtcdDBSizeTooLargeAWS`
  - `WorkloadClusterEtcdHasNoLeaderAWS`
  - `WorkloadClusterEtcdNumberOfLeaderChangesTooHighAWS`
  - `WorkloadClusterMasterNodeMissingFirecracker`
  - `WorkloadClusterPodLimitAlmostReachedAWS`
- Splitting `NodeIsUnschedulable` per team
- Split `ContainerIsRestartingTooFrequentlyFirecracker` into `WorkloadClusterContainerIsRestartingTooFrequentlyFirecracker` and `ManagementClusterContainerIsRestartingTooFrequentlyFirecracker`
- Add the following biscuit alerts to split alerts between workload and management cluster:
  - `ManagementClusterCriticalPodNotRunning`
  - `ManagementClusterCriticalPodMetricMissing`
  - `ManagementClusterPodLimitAlmostReached`

### Changed

- Move `AzureManagementClusterMissingNodes` and `AWSManagementClusterMissingNodes` to team biscuit `ManagementClusterMissingNodes`
- Move `ManagementClusterPodStuckAzure` and `ManagementClusterPodStuckAWS` to team biscuit `ManagementClusterPodPendingFor15Min`
- Renamed the following alerts:
  - `AzureClusterAutoscalerIsRestartingFrequently` -> `WorkloadClusterAutoscalerIsRestartingFrequentlyAzure`
  - `CriticalPodNotRunningAzure` -> `WorkloadClusterCriticalPodNotRunningAzure`
  - `CriticalPodMetricMissingAzure` -> `WorkloadClusterCriticalPodMetricMissingAzure`
  - `MasterNodeMissingCelestial` -> `WorkloadClusterMasterNodeMissingCelestial`
  - `NodeUnexpectedTaintNodeWithImpairedVolumes` -> `WorkloadClusterNodeUnexpectedTaintNodeWithImpairedVolumes`
  - `PodLimitAlmostReachedAzure` -> `WorkloadClusterPodLimitAlmostReachedAzure`

### Fixed

- Do not page biscuit for a failing prometheus

## [1.19.2] - 2021-02-12

### Fixed

- Fix incorrect prometheus memory usage recording rules after we migrated to the new monitoring setup

### Changed

- Use azure-collector instead of azure-operator in `AzureClusterCreationFailed` alert

### Removed

- Removing service monitor resource used to clean up unused service monitor CR

## [1.19.1] - 2021-02-10

### Fixed

- Fix empty prometheus rules in helm template issues for aws and kvm installations

## [1.19.0] - 2021-02-10

### Added

- Added the celestial rules to PMO:
  - `AzureClusterAutoscalerIsRestartingFrequently`
  - `AzureClusterCreationFailed`
  - `AzureDeploymentIsRunningForTooLong`
  - `AzureDeploymentStatusFailed`
  - `AzureManagementClusterDeploymentScaledDownToZero`
  - `AzureManagementClusterMissingNodes`
  - `AzureNetworkErrorRateTooHigh`
  - `AzureServicePrincipalExpirationDateUnknown`
  - `AzureServicePrincipalExpiresInOneMonth`
  - `AzureServicePrincipalExpiresInOneWeek`
  - `AzureVMSSRateLimit30MinutesAlmostReached`
  - `AzureVMSSRateLimit30MinutesReached`
  - `AzureVMSSRateLimit3MinutesAlmostReached`
  - `AzureVMSSRateLimit3MinutesReached`
  - `ClockOutOfSyncAzure`
  - `ClusterAutoscalerAppFailedAzure`
  - `ClusterAutoscalerAppNotInstalledAzure`
  - `ClusterAutoscalerAppPendingInstallAzure`
  - `ClusterAutoscalerAppPendingUpgradeAzure`
  - `ClusterWithNoResourceGroup`
  - `CollidingOperatorsCelestial`
  - `CriticalPodMetricMissingAzure`
  - `CriticalPodNotRunningAzure`
  - `DNSCheckErrorRateTooHighAzure`
  - `DNSErrorRateTooHighAzure`
  - `DeploymentNotSatisfiedCelestial`
  - `EtcdWorkloadClusterDownAzure`
  - `LatestETCDBackup1DayOld`
  - `LatestETCDBackup2DaysOld`
  - `ManagementClusterNotBackedUp24h`
  - `MasterNodeMissingCelestial`
  - `OperatorNotReconcilingCelestial`
  - `OperatorkitCRNotDeletedCelestial`
  - `OperatorkitErrorRateTooHighCelestial`
  - `PodLimitAlmostReachedAzure`
  - `ManagementClusterPodStuckAzure` (renamed from `PodStuckAzure`)
  - `ReadsRateLimitAlmostReached`
  - `VPNConnectionProvisioningStateBad`
  - `VPNConnectionStatusBad`
  - `WorkloadClusterEtcdCommitDurationTooHighAzure`
  - `WorkloadClusterEtcdDBSizeTooLargeAzure`
  - `WorkloadClusterEtcdHasNoLeaderAzure`
  - `WorkloadClusterEtcdNumberOfLeaderChangesTooHighAzure`
  - `WritesRateLimitAlmostReached`
  - `ETCDBackupJobFailedOrStuck` (renamed from `BackupJobFailedOrStuck`)
- Added node `role` label to `kubelet` metrics as it's needed by `MasterNodeMissingCelestial` alert

### Removed

- Removed axolotl from Chinese rules as the installation has been decommissioned

## [1.18.0] - 2021-02-08

### Removed

- Added the batman alerts to PMO:
  - `AppExporterDown`
  - `AppOperatorNotReady`
  - `AppWithoutTeamLabel`
  - `CertManagerPodHighMemoryUsage`
  - `CertificateSecretWillExpireInLessThanTwoWeeks`
  - `ChartOperatorDown`
  - `ChartOrphanConfigMap`
  - `ChartOrphanSecret`
  - `CollidingOperatorsBatman`
  - `CordonedAppExpired`
  - `DeploymentNotSatisfiedBatman`
  - `DeploymentNotSatisfiedChinaBatman`
  - `ElasticsearchClusterHealthStatusRed`
  - `ElasticsearchClusterHealthStatusYellow`
  - `ElasticsearchDataVolumeSpaceTooLow`
  - `ElasticsearchHeapUsageWarning`
  - `ElasticsearchPendingTasksTooHigh`
  - `ExternalDNSCantAccessRegistry`
  - `ExternalDNSCantAccessSource`
  - `HelmHistorySecretCountTooHigh`
  - `IngressControllerDeploymentNotSatisfied`
  - `IngressControllerMemoryUsageTooHigh`
  - `IngressControllerReplicaSetNumberTooHigh`
  - `IngressControllerSSLCertificateWillExpireSoon`
  - `IngressControllerServiceHasNoEndpoints`
  - `ManagedAppBasicErrorBudgetBurnRateAboveSafeLevel`
  - `ManagedAppBasicErrorBudgetBurnRateInLast10mTooHigh`
  - `ManagedAppBasicErrorBudgetEstimationWarning`
  - `ManagedLoggingElasticsearchClusterDown`
  - `ManagedLoggingElasticsearchDataNodesNotSatisfied`
  - `ManagementClusterAppFailed`
  - `OperatorNotReconcilingBatman`
  - `OperatorkitErrorRateTooHighBatman`
  - `RepeatedHelmOperation`
  - `TillerHistoryConfigMapCountTooHigh`
  - `TillerRunningPods`
  - `TillerUnreachable`
  - `WorkloadClusterAppFailed`
  - `WorkloadClusterDeploymentNotSatisfied`
  - `WorkloadClusterDeploymentScaledDownToZero`
  - `WorkloadClusterManagedDeploymentNotSatisfied`

## [1.17.2] - 2021-02-04

### Changed

- (internal) Rely on `Ingress` for OAuth2 proxy to configure TLS for Prometheus
  domain, as it also configures management of the certificates, instead of
  creating copies which could break access in case they became out of date.

### Fixed

- Fix incorrect prometheus memory usage recording rule

## [1.17.1] - 2021-02-02

### Fixed

- Fixed incorrect label in GatekeeperDown alert.

## [1.17.0] - 2021-02-02

### Added

- Added the `NoHealthyJumphost` alert
- Added the biscuit alerts to PMO:
  - `AppCollectionDeploymentFailed`
  - `CalicoNodeMemoryHighUtilization`
  - `CrsyncDeploymentNotSatisfied`
  - `CrsyncTooManyTagsMissing`
  - `DeploymentNotSatisfiedBiscuit`
  - `DeploymentNotSatisfiedChinaBiscuit`
  - `DraughtsmanRateLimitAlmostReached`
  - `EtcdDown`
  - `GatekeeperDown`
  - `GatekeeperWebhookMissing`
  - `KeyPairStorageAlmostFull`
  - `ManagementClusterHasLessThanThreeNodes`
  - `ManagementClusterCriticalSystemdUnitFailed`
  - `ManagementClusterDisabledSystemdUnitActive`
  - `ManagementClusterEtcdCommitDurationTooHigh`
  - `ManagementClusterEtcdDBSizeTooLarge`
  - `ManagementClusterEtcdHasNoLeader`
  - `ManagementClusterEtcdNumberOfLeaderChangesTooHigh`
  - `ManagementClusterHighNumberSystemdUnits`
  - `ManagementClusterPodPending`
  - `ManagementClusterSystemdUnitFailed`
  - `VaultIsDown`
  - `VaultIsSealed`

### Changed

- Renamed control plane and tenant cluster respectively to management cluster and
  workload cluster. Renamed some alerts:
  - ControlPlaneCertificateWillExpireInLessThanTwoWeeks > ManagementClusterCertificateWillExpireInLessThanTwoWeeks
  - ControlPlaneDaemonSetNotSatisfiedAtlas > ManagementClusterDaemonSetNotSatisfiedAtlas
  - ControlPlaneDaemonSetNotSatisfiedChinaAtlas > ManagementClusterDaemonSetNotSatisfiedChinaAtlas
  - PrometheusCantCommunicateWithTenantAPI > PrometheusCantCommunicateWithKubernetesAPI
- Rename ETCDDown alert to ManagementClusterEtcdDown
- Enable alerts only on the corresponding providers

### Fixed

- Fix missing app label on kube-apiserver target
- Fix missing app label on nginx-ingress-controller target

## [1.16.1] - 2021-01-28

### Fixed

- Fix recording rules to apply them to all prometheuses

## [1.16.0] - 2021-01-28

### Changed

- Reenable `Remote Write` to Cortex

### Added

- Trigger final heartbeat before deleting the cluster to clean up opened heartbeat alerts

### Removed

- Remove webhook from `AlertManagerNotificationsFailing` alert.

## [1.15.0] - 2021-01-22

### Changed

- Use giantswarm/prometheus image

### Fixed

- Fix recording rules creation
- Fix prometheus container image tag to not use latest
- Fix prometheus minimal memory in VPA

## [1.14.0] - 2021-01-13

### Added

- Add inhibition rules.
- Set Prometheus pod max memory usage (via vpa) to 90% of lowest node allocatable memory
- Prometheus monitors itself

### Changed

- Ignore missing unhealthy prometheus instances in promxy to avoid it from crash looping
- Added the biscuit alerts to PMO:
  - `ControlPlaneCertificateWillExpireInLessThanTwoWeeks`
- Add topologySpreadConstraint to evenly spread prometheus pods
- Ignore Slack in `AlertManagerNotificationsFailing` alert.
- Set heartbeat alert to up for 10mn

### Removed

- Removed g8s-prometheus target
- Removed alert resource

## [1.13.0] - 2021-01-05

### Added

- Add priority class `prometheus` and use it for all managed Prometheus pods in
  order to allow scheduler to evict other pods with lower priority to make
  space for Prometheus

## [1.12.0] - 2020-12-02

### Changed

- Change PrometheusCantCommunicateWithTenantAPI to ignore promxy
- Set prometheus default resources to 100m of CPU and 1Gi of memory
- Reduced number of metrics ingested from nginx-ingress-controller in order to
  reduce memory requirements of Prometheus.

## [1.11.0] - 2020-12-01

### Added

- Create `VerticalPodAutoscaler` resource for each Prometheus configuring the
  VPA to manage Prometheus pod requests and limits to allow dynamic scaling but
  prevent scheduling and OOM issues.

### Changed

- Change prometheus affinity from "Prefer" to "Required".

## [1.10.3] - 2020-11-25

### Fixed

- Fix initial heartbeat ping so that it only triggers on creation.

## [1.10.2] - 2020-11-25

### Fixed

- Set prometheus cpu requests and limits to 0.25 CPU.

## [1.10.1] - 2020-11-24

### Fixed

- Set prometheus cpu requests and limits to 1 CPU.
- Set prometheus memory requests and limits to 5Gi.

## [1.10.0] - 2020-11-20

### Added

- Add team atlas alerts (in helm chart).

### Changed

- Set heartbeat client log level to fatal to avoid polluting our logs.
- Set prometheus to select rules from monitoring namespace.

### Removed

- Set alert resource to delete PrometheusRules in cluster namespace.

### Fixed

- Fix prometheus targets.
- Fix duplicated scrapping of nginx-ingress-controller.

## [1.9.0] - 2020-11-11

### Added

- Add support for `Remote Write` to Cortex
- Added recording rules
- Add node affinity to prefer not scheduling on master nodes
- Added `pipeline` tag to _Hearbeat_ alert to be able to see if it affects
  a stable or testing installation at first glance

### Changed

- Increase memory request from 100Mi to 5Gi

### Fixed

- Fix kube-state-metrics scraping port on Control Planes.
- Fixed creating of alerts, it was failing due to a typo in template path

## [1.8.0] - 2020-10-21

### Added

- Add pod, container, node and node role labels
- Allow ignoring clusters using the `giantswarm.io/monitoring: false` label on cluster CRs
- Add monitoring of control plane bastions
- Add heartbeat alert to prometheus
- Create heartbeat in opsgenie
- Route heartbeat alerts to corresponding opsgenie heartbeat

## [1.7.0] - 2020-10-14

### Added

- Add alertmanager config

### Fixed

- Fix a bug where promxy configmap keep growing and lead to OutOfMemory issues.
- Fix an issue where prometheus fails to be created due to resource order.

## [1.6.0] - 2020-10-12

### Changed

- Set retention size to 90Gi and duration to 2w
- Increased storage to 100Gi

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


[Unreleased]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.41.1...HEAD
[1.41.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.41.0...v1.41.1
[1.41.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.40.0...v1.41.0
[1.40.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.39.0...v1.40.0
[1.39.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.38.0...v1.39.0
[1.38.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.37.0...v1.38.0
[1.37.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.36.0...v1.37.0
[1.36.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.35.0...v1.36.0
[1.35.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.34.1...v1.35.0
[1.34.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.34.0...v1.34.1
[1.34.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.33.0...v1.34.0
[1.33.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.32.1...v1.33.0
[1.32.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.32.0...v1.32.1
[1.32.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.31.0...v1.32.0
[1.31.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.30.0...v1.31.0
[1.30.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.29.1...v1.30.0
[1.29.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.29.0...v1.29.1
[1.29.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.28.0...v1.29.0
[1.28.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.27.4...v1.28.0
[1.27.4]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.27.3...v1.27.4
[1.27.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.27.2...v1.27.3
[1.27.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.27.1...v1.27.2
[1.27.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.27.0...v1.27.1
[1.27.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.26.0...v1.27.0
[1.26.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.25.2...v1.26.0
[1.25.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.25.1...v1.25.2
[1.25.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.25.0...v1.25.1
[1.25.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.24.8...v1.25.0
[1.24.8]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.24.6...v1.24.8
[1.24.6]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.24.7...v1.24.6
[1.24.7]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.24.6...v1.24.7
[1.24.6]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.24.5...v1.24.6
[1.24.5]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.24.4...v1.24.5
[1.24.4]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.24.3...v1.24.4
[1.24.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.24.2...v1.24.3
[1.24.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.24.1...v1.24.2
[1.24.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.24.0...v1.24.1
[1.24.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.23.1...v1.24.0
[1.23.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.23.0...v1.23.1
[1.23.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.22.0...v1.23.0
[1.22.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.21.0...v1.22.0
[1.21.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.20.0...v1.21.0
[1.20.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.19.2...v1.20.0
[1.19.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.19.1...v1.19.2
[1.19.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.19.0...v1.19.1
[1.19.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.18.0...v1.19.0
[1.18.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.17.2...v1.18.0
[1.17.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.17.1...v1.17.2
[1.17.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.17.0...v1.17.1
[1.17.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.16.1...v1.17.0
[1.16.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.16.0...v1.16.1
[1.16.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.15.0...v1.16.0
[1.15.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.14.0...v1.15.0
[1.14.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.13.0...v1.14.0
[1.13.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.10.3...v1.11.0
[1.10.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.10.2...v1.10.3
[1.10.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.10.1...v1.10.2
[1.10.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.10.0...v1.10.1
[1.10.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.9.0...v1.10.0
[1.9.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.8.0...v1.9.0
[1.8.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.7.0...v1.8.0
[1.7.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.6.0...v1.7.0
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
