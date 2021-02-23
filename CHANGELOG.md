# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [Unreleased]

### Added

- Add `TestClusterTooOld` for testing installations
- Add `ManagementClusterPodStuckFirecracker` and `WorkloadClusterPodStuckFirecracker` alerts for Firecracker.
- Add `ManagementClusterPodStuckCelestial` alert for Celestial.

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


[Unreleased]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.23.0...HEAD
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
