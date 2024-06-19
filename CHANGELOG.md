# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Remove line-breaks in alerting links which suppress links in notifications.

## [4.77.1] - 2024-06-19

### Fixed

- Reverse ingress removal condition to remove the ingress when mimir is enabled.

## [4.77.0] - 2024-06-19

### Removed

- Remove AlertManager link in opsgenie and slack templating when mimir is enabled.

### Changed

- Remove unused scrape_timeout inhibition.
- Some improvements towards Mimir:
  - Internal rework to remove the use generic resource to ease out the migration to Mimir.
  - Update generic resource so we can delete resources when mimir is enabled.
  - Remove legacy prometheus resources when Mimir is enabled.
  - Remove alertmanager ingress when Mimir is enabled.
  - Ignore the prometheus-to-grafana-cloud prometheus in the remove write controller.
- Change Alert link to point to Mimir alerting UI when Mimir is enabled.
  - Rename Prometheus link to Source

## [4.76.0] - 2024-06-03

### Changed

- Delete per cluster heartbeats when Mimir is enabled.
- Delete per cluster heartbeat alertmanager wiring when Mimir is enabled.

## [4.75.1] - 2024-05-23

### Changed

- Remove prometheus remote write agent configuration when mimir is enabled.
- Remove unnecessary prometheus control-plane affinity.

## [4.75.0] - 2024-05-13

### Added

- Add `cluster_control_plane_unhealthy` inhibition.
- Allow Prometheus Agent Sharding strategy to be overridden per cluster.

### Removed

- Removed `apiserver_down` inhibition.

### Fixed

- Use `kubernetes.io/tls` type for TLS secrets.

## [4.74.0] - 2024-05-02

### Changed

- Expose prometheus agent sharding strategies as prometheus-meta-operator configuration parameters so we can experiment with the scaling strategies.

## [4.73.1] - 2024-05-01

### Fixed

- Ensure proxy url is set when needed within slack_configs

## [4.73.0] - 2024-04-30

### Changed

- To ensure that customers can define their own AlertmanagerConfig CRs, we need to remove the default alertmanager matcher injection (cf. upstream https://github.com/prometheus-operator/prometheus-operator/issues/4033)

### Added

- Add `SlackApiToken` configuration directive.

## [4.72.0] - 2024-04-03

### Changed

- This PR adds a receiver and a route for the mimir heartbeat. We need to add them here until we use mimir's alertmanager.

## [4.71.0] - 2024-03-19

### Added

- Add team `honeybadger` slack router and receiver.

### Removed

- Remove the `azure` provider.

## [4.70.3] - 2024-03-18

### Fixed

- Fix missing data in labelling schema to add missing labels to avoid issues with aggregations of data coming from prometheus agents that have some extra labels set as opposed to the existing prometheus scrape config.

## [4.70.2] - 2024-03-13

### Fixed

- Fix missing prometheus link in notification template.

## [4.70.1] - 2024-03-12

### Fixed

- Fix noop resource creation code.

## [4.70.0] - 2024-03-12

### Removed

- Remove alerting from Prometheus if mimir is enabled.

## [4.69.0] - 2024-03-12

### Changed

- Disable rule evaluation in Prometheus when Mimir is enabled.

### Removed

- Remove `prometheus` and `prometheus_replica` external labels when Mimir is enabled.

## [4.68.4] - 2024-03-06

### Removed

- Remove `falco-exporter` from static scrapeconfig as they are now monitored via servicemonitors.

## [4.68.3] - 2024-02-19

### Fixed

- Fix alertmanager ciliumnetworkpolicy to allow access to coredns.

## [4.68.2] - 2024-02-19

### Fixed

- Add missing ciliumnetworkpolicy for alertmanager.

## [4.68.1] - 2024-02-15

### Added

- Add proxy port to CiliumNetworkPolicy if needed.

## [4.68.0] - 2024-02-14

### Added

- Add CNP for prometheus-meta-operator to be able to talk to the api-server in locked-down clusters.

## [4.67.3] - 2024-02-13

### Added

- Add `update` method to cilium netpol resource.

## [4.67.2] - 2024-02-13

### Fixed

- Add grafana-cloud squid proxy port to prometheus CNP.

## [4.67.1] - 2024-02-13

### Fixed

- Fix error for already existing `ciliumNetworkPolicy`.

## [4.67.0] - 2024-02-12

### Added

- Add `ciliumNetworkPolicy` for all Prometheus instances on the MC.

## [4.66.1] - 2024-02-07

### Fixed

- Fix VPA to support latest Prometheus-operator version (based on observability-bundle 1.2.0) as the latest version of the Prometheus CR now supports the `scale` subresource which causes issues with VPA.

## [4.66.0] - 2024-02-06

### Changed

- Support multi-provider Management clusters.

### Fixed

- Fix how we enable `remote-write-receiver` to avoid deprecated warnings.
- Fix test generation to split capi and vintage tests generated files.

### Removed

- Free retention duration property of it's 2 weeks limitation if the free storage allows it.

## [4.65.0] - 2024-01-29

### Added

- Add `CiliumNetworkPolicy` for all created Prometheuses.
- Always set `shards` to 1 for all created Prometheuses.

## [4.64.0] - 2024-01-19

### Changed

- Improved the blackhole routing for `stable-testing` MCs to silence more alerts related to test WCs

## [4.63.1] - 2023-12-12

### Fixed

- Fix Tinkerers slack receiver repeat_interval config.

## [4.63.0] - 2023-12-06

### Changed

- Configure `gsoci.azurecr.io` as the default container image registry.
- Change Tinkerers slack receiver repeat_interval to 2 weeks.

## [4.62.0] - 2023-11-27

### Changed

- Group alerts by teams.

## [4.61.0] - 2023-11-22

### Changed

- Upgrade to go 1.21
- Upgrade internal dependencies.
- Increased `group_wait` in alertmanager config from 1 to 5m.

## [4.60.0] - 2023-11-02

### Changed

- Silence `ManagementClusterAppFailed` for WCs of `stable-testing` MCs.

## [4.59.0] - 2023-10-11

### Changed

- Upgrade Prometheus to 2.47.1 and configure keepDroppedTargets to 5.

### Fixed

- Alert template: fix newlines / whitespace trimming if opsrecipe is *not* specified or a dashboard *is* specified.

## [4.58.0] - 2023-10-09

### Removed

- Remove custom SLO handling in alertmanager config.

## [4.57.0] - 2023-10-04

### Fixed

- Fix Prometheus PSP by adding seccomp profile to RuntimeDefault.

### Added

- Handle `remoteTimeout` in RemoteWrite secret and set it to 60s (hardcoded to 30s with `prometheus-agent < 0.6.4`).

### Removed

- Remove the temporary code in pmo to avoid RemoteWriteSecret update on anteater/deu01 and anteater/seu01.

## [4.56.0] - 2023-10-03

### Fixed

- Set Prometheus seccomp profile to RuntimeDefault.

## [4.55.0] - 2023-10-02

### Changed

- Add condition for PSP installation in helm chart

## [4.54.1] - 2023-09-28

### Added

- Routing rule for `ClusterUnhealthyPhase` and test clusters on stable-testing MCs to route to blackhole

## [4.54.0] - 2023-09-28

### Added

- Add cert-manager clusterIssuer configuration option for Ingresses.
- Add support for EKS as a provider.

## [4.53.0] - 2023-09-27

### Changed

- Temporary avoid RemoteWriteSecret update on anteater/deu01 and anteater/seu01.

### Removed

- Remove KVM related things that are not used anymore.
- Revert `prometheus-agent` max shards to 10 to prevent incessant paging.

## [4.52.0] - 2023-09-26

### Changed

- Support absolute Grafana dashboard URLs.
- Increase `prometheus-agent` max shards to 50 to improve agent stability.

## [4.51.0] - 2023-09-25

### Changed

- Ignore `PrometheusMetaOperatorReconcileErrors` alerts on `stable-testing`.
- Increase `group_wait` from AlertManager config to let more time to inhibition alerts to be executed.

## [4.50.0] - 2023-09-25

### Changed

- Only send silenced page-level SLOTH alerts to `phoenix`'s slack alert channel, rather than all alerts.

## [4.49.2] - 2023-09-22

### Fixed

- Reverted support absolute Grafana dashboard URLs.

## [4.49.1] - 2023-09-21

### Changed

- Ignore kube-proxy target on EKS or clusters with observability bundle >= 0.8.3 (where the kube-proxy service monitor is enabled).

## [4.49.0] - 2023-09-21

### Changed

- Adapt scrape targets to EKS clusters.
- computation of number of shards: rely on max number of series over the last 6h.

### Fixed

- Support absolute Grafana dashboard URLs.
- Fix api server url in case the CAPI provider sets https prefix in the CAPI CR status.

## [4.48.0] - 2023-09-19

### Changed

- Support flux-managed clusters.

## [4.47.0] - 2023-09-14

### Changed

- Enable Opsgenie alerts for Shield.

### Fixed

- Change source for the organization label.

## [4.46.0] - 2023-08-21

### Added

- Add team `tinkerers` slack router and receiver.

### Changed

- Apply Kyverno policy exception to the PMO replicaset as well.

### Fixed

- Fix null receiver name.
- Remove `aws-load-balancer-controller` from the list of ignored targets.

## [4.45.1] - 2023-07-18

## Changed

- Skip alerts named `WorkloadClusterApp*` in `stable-testing` installations.

## [4.45.0] - 2023-07-18

### Changed

- When the cluster pipeline is set to stable-testing, only route management cluster alerts to opsgenie.

### Removed

- Clean up unused targets (moved to service monitors).
- Remove #harbor-implementation slack integration.

## [4.44.0] - 2023-07-04

### Removed

- Clean up some of the vintage targets.

## [4.43.2] - 2023-07-03

### Fixed

- Set number of shards to existing value if Prometheus is not reachable to avoid a race condition on cluster creation.

## [4.43.1] - 2023-06-29

### Fixed

- Fix some security concerns.

## [4.43.0] - 2023-06-29

### Changed

- Change shard computation to be based on number of head series.

## [4.42.0] - 2023-06-26

### Added

- ReRoute clippy alerts to `phoenix slack` until all team labels are changed

## [4.41.0] - 2023-06-22

### Added

- Added scrape for `vault-etcd-backups-exporter` towards legacy vault VMs.
- Add Kyverno Policy Exceptions.

## [4.40.0] - 2023-06-19

### Added

- Add back Prometheus CPU limits.
- Add alert routing for `team-turtles`

### Changed

- Move `gcp`, `capa` and `capz` alerts to team phoenix.

### Fixed

- Update dropped labels in KSM metrics to avoid duplicate samples.
- Drop unused greedy KSM metrics.
- Remove imagePullSecrets towards https://github.com/giantswarm/giantswarm/issues/27267.

## [4.39.0] - 2023-06-07

### Removed

- Prometheus CPU limits

## [4.38.2] - 2023-06-02

### Added

- Add static node-exporter target to the list of ignored targets because this target is needed for releases still using node-exporter < 1.14 (< aws 18.2.0).

## [4.38.1] - 2023-05-31

### Added

- Add alert routing for team bigmac

## [4.38.0] - 2023-05-22

### Changed

- Dynamically compute agent shards number according to cluster size.

## [4.37.0] - 2023-05-10

### Added

- Add sharding capabilities to the Prometheus Agent.
- Create new remote-write-secret Secret and remote-write-config ConfigMap per cluster to not have a bad workaround in the observability bundle.

### Fixed

- Fix prometheus control plane node toleration.

### Removed

- Stop pushing to `openstack-app-collection`.

## [4.36.4] - 2023-05-02

### Fixed

- Fix forgotten kube-state-metrics down source labels.

## [4.36.3] - 2023-05-02

### Fixed

- Add missing node label to kubelet.

## [4.36.2] - 2023-05-01

### Changed

- Increased heartbeat delay before alert from 25min to 60min
- Updated alertmanager heartbeat config for faster response

## [4.36.1] - 2023-04-28

### Fixed

- Keep accepting 'master' as `role` label value for etcd scraping.

## [4.36.0] - 2023-04-27

### Changed

- Deprecate `role=master` in favor of `role=control-plane`.

## [4.35.4] - 2023-04-25

### Changed

- Add Inhibition rule when prometheus-agent is down.

## [4.35.3] - 2023-04-18

### Changed

- Change Atlas slack alert router to only route alerts with page and/or notify severity matcher.

### Fixed

- Fix list of ignored targets for Vintage WCs.

## [4.35.2] - 2023-04-13

### Fixed

- Allow PMO to patch secrets so it can remove finalizers.

## [4.35.1] - 2023-04-12

### Added

- Add finalizer to remote-write-config to block cluster deletion until PMO deleted the secret.

## [4.35.0] - 2023-04-11

### Changed

- Handle prometheus scrape target removal based on the observability bundle version.

## [4.34.0] - 2023-04-06

### Added

- Add `loki` namespace in cAdvisor scrape config for MC.

### Fixed

- Fix proxy configuration as no_proxy was not respected.

## [4.33.0] - 2023-04-04

### Changed

- Add more flexibility in the configuration so prometheus image, pvc size and so on can be overwritten by configuration.

## [4.32.0] - 2023-03-30

### Removed

- Drop node-exporter metrics (`node_filesystem_files` `node_filesystem_readonly` `node_nfs_requests_total` `node_network_carrier` `node_network_transmit_colls_total` `node_network_carrier_changes_total` `node_network_transmit_packets_total` `node_network_carrier_down_changes_total` `node_network_carrier_up_changes_total` `node_network_iface_id` `node_xfs_.+` `node_ethtool_.+`)
- Drop kong metrics (`kong_latency_count` `kong_latency_sum`)
- Drop kube-state-metrics metrics (`kube_.+_metadata_resource_version`)
- Drop nginx-ingress-controller metrics (`nginx_ingress_controller_bytes_sent_sum` `nginx_ingress_controller_request_size_count` `nginx_ingress_controller_response_size_count` `nginx_ingress_controller_response_duration_seconds_sum` `nginx_ingress_controller_response_duration_seconds_count` `nginx_ingress_controller_ingress_upstream_latency_seconds` `nginx_ingress_controller_ingress_upstream_latency_seconds_sum` `nginx_ingress_controller_ingress_upstream_latency_seconds_count`)

## [4.31.1] - 2023-03-28

### Changed

- Prometheus-agent tuning: revert maxSamplesPerSend to 150000

## [4.31.0] - 2023-03-28

### Removed

- Drop `rest_client_rate_limiter_duration_seconds_bucket` `rest_client_request_size_bytes_bucket` `rest_client_response_size_bytes_bucket` from Kubernetes component metrics.
- Drop `coredns_dns_response_size_bytes_bucket` and `coredns_dns_request_size_bytes_bucket` from coredns metrics.
- Drop `nginx_ingress_controller_connect_duration_seconds_bucket` `nginx_ingress_controller_header_duration_seconds_bucket` `nginx_ingress_controller_bytes_sent_count` `nginx_ingress_controller_request_duration_seconds_sum` from nginx-ingress-controller metrics.
- Drop `kong_upstream_target_health` and `kong_latency_bucket` Kong metrics.
- Drop `kube_pod_tolerations` `kube_pod_status_scheduled` `kube_replicaset_metadata_generation` `kube_replicaset_status_observed_generation` `kube_replicaset_annotations` and`kube_replicaset_status_fully_labeled_replicas` kube-state-metrics metrics.
- Drop `promtail_request_duration_seconds_bucket` and `loki_request_duration_seconds_bucket` metrics from promtail and loki.

## [4.30.0] - 2023-03-28

### Removed

- Remove immutable secret deletion not needed after 4.27.0.
- Remove alertmanager ownership job.

## [4.29.2] - 2023-03-28

### Changed

- VPA settings: set memory limit to 80% node size
- Drop `awscni_assigned_ip_per_cidr` metric from aws cni.
- Drop `uid` label from kubelet.
- Drop `image_id` label from kube-state-metrics.

## [4.29.1] - 2023-03-27

### Changed

- Prometheus remotewrite endpoints for agents: increase max body size from 10m to 50m

### Removed

- Removed pod_id relabelling as it's not needed anymore.

## [4.29.0] - 2023-03-27

### Changed

- Bump Prometheus default image to `v2.43.0`
- Prometheus-agent tuning: increase maxSamplesPerSend from 150000 to 300000

### Removed

- Drop some unused metrics from cAdvisor.
- Remove draughtsman references.

## [4.28.0] - 2023-03-23

### Remove

- Drop `id` and `name` label from cAdvisor metrics.

## [4.27.0] - 2023-03-22

### Changed

- Allow changes in the remote write api endpoint secret.

### Fixed

- The region as external Label for capa,gcp and capz

### Removed

- Drop `uid` label from kube-state-metrics metrics.
- Drop `container_id` label from kube-state-metrics metrics.

## [4.26.0] - 2023-03-20

### Changed

- Prometheus resources: set requests=limits. Still allowing prometheus up to 90% of node capacity.
- Prometheus TSDB size: reduce it to 85% of disk space, to keep space for WAL before alerts fire.
- Prometheus-agent tuning: increase maxSamplesPerSend from 50000 to 150000

## [4.25.3] - 2023-03-15

### Fixed

- Fix ownership job

## [4.25.2] - 2023-03-15

### Changed

- Updated `RetentionSize` property in Prometheus CR according to Volume Storage Size (90%)
- Allow ownership job patch `alertmanagerConfigSelector` to fail in case the label has been already removed.

## [4.25.1] - 2023-03-14

### Changed

- Followup on Alertmanager resource to Helm
  - Set Alertmanager enabled by default.
  - remove the label managed-by: pmo from alertmanagerConfigSelector.

## [4.25.0] - 2023-03-14

### Changed

- Move Alertmanager resource to Helm
  - Delete `controller resource alerting/alertmanager`.
  - Create alertmanager template in helm.
  - Delete the obsolete static scraping configs for alertmanager.
  - Add a hook job that change the ownership labels for alertmanager resource.

## [4.24.1] - 2023-03-09

### Changed

- VPA settings: changes in 4.24.0 were wrong, resulting in too low limits.
    - Previous logic (4.23.0) was right, and limits were 90% node size.
    - Comments have been updated for better understanding
    - limit has been reverted to 90% node size
    - code for CPU limits has been updated to do the same kind of calculations
    - tests have been updated for more meaningful results

## [4.24.0] - 2023-03-02

### Changed

- Un-drop `nginx_ingress_controller_request_duration_seconds_bucket` for workload clusters
- Add additional annotations on all `ingress` objects to support DNS record creation via `external-dns`
- VPA settings: changed max memory requests from 90% to 80% of node RAM, so that memory limit is 96% node RAM (avoids crashing node with big prometheis)
- VPA settings: remove useless config for `prometheus-config-reloader` and `rules-configmap-reloader`: now it's only 1 container called `config-reloader`, and default config scales it down just nice!

## [4.23.0] - 2023-02-28

### Changed

- Look up PrometheusRules CR in the whole MC only labelled with `application.giantswarm.io/team`

## [4.22.0] - 2023-02-27

### Fixed

- Removed Prometheus readinessProbe 5mn delay; since there is already a 15mn startupProbe

### Changed

- Increased ScrapeInterval and EvaluationInterval from 30s to 60s
    - pros: twice less CPU usage, less disk usage
    - cons: up to 30s more delay in alerts, and very short usage peaks get smoothed over 1 minute
- Addep `kyverno` namespace to WC & MC default scrape config for cadvisor metrics

## [4.21.0] - 2023-02-23

### Changed

- Use resource.Quantity.AsApproximateFloat64() instead of AsInt64(), in order to avoid conversion issue when multiply cpu, e.g. 3880m
- Use label selector to selects only worker nodes for vpa resource to get maxCPU and maxMemory
- List nodes from API only once in VPA resource
- Improve VPA maxAllowedCPU, use 70% of the node allocatable CPU.
- Prevent 'notify' severity alerts from being sent to '#alert-phoenix' Slack channel (too noisy).
- Update getMaxCPU use 50% of the node allocatable CPU.

### Added

- Send SLO (sloth based) notify level alerts to '#alert-phoenix' Slack channel.

## [4.20.6] - 2023-02-13

### Changed

- Fix list of targets to scrape or ignore.

## [4.20.5] - 2023-02-09

### Fixed

- Manage etcd certificates differently between CAPI/Vintage. On Vintage, etcd certificates are binded via a volume. On CAPI, certificates are binded via a secret.
- Pass the missing Provider property to `etcdcertificates.Config`

### Added

- Add `.provider.flavor` property in Helm chart.

## [4.20.4] - 2023-02-07

### Fixed

- Fix certificates created as directories rathen than files

## [4.20.3] - 2023-02-02

### Fixed

- Fix heartbeat update mechanism to prevent leftover labels in OpsGenie.

## [4.20.2] - 2023-01-18

### Fixed

- Remove proxy support to remote write endpoint consumers.

### Added

- Add alertmanagerservicemonitor resource, to scrape metrics from alertmanager.
- Added target and source matchers for stack_failed label.

## [4.20.1] - 2023-01-17

### Fixed

- Enable `remote-write-receiver` via `EnableFeatures` field added in `CommonPrometheusFields` (schema 0.62.0)

## [4.20.0] - 2023-01-17

### Changed

- Upgrade `prometheus` from 2.39.1 to 2.41.0 and `alertmanager` from 0.23.0 to 0.25.0.

## [4.19.2] - 2023-01-12

### Fixed

- Fix getDefaultAppVersion org namespace

### Changed

- Bump alpine from 3.17.0 to 3.17.1

## [4.19.1] - 2023-01-11

### Added

- Add proxy support to remote write endpoint consumers.

### Fixed

- Fix node-exporter target discovery

## [4.19.0] - 2023-01-02

### Changed

- remotewrite ingress allows bigger requests
- prometheus-agent: increase max samples per send. ⚠️ Warning: updates an immutable secret, will require manual actions at deployment.

## [4.18.0] - 2022-12-19

### Changed

- Allow remote write over insecure endpoint certificate.
- Ignore remotewrite feature in kube-system namespace.

## [4.17.0] - 2022-12-14

### Changed

- Deploy needed resources for the agent to run on Vintage MCs.

### Fixed

* opsgenie alert templating: list of firing instances
* slack alert templating: list of firing instances
* fix dashboard url

## [4.16.0] - 2022-12-07

### Changed

- Change HasPrometheusAgent function to ignore prometheus-agent scraping targets on CAPA and CAPVCD.
- Do not reconcile service monitors in kube-system for CAPA and CAPVCD MCs.
- Change label selector used to discover `PodMonitors` and `ServiceMonitors`
  to avoid a duplicate scrape introduced in https://github.com/giantswarm/observability-bundle/pull/18.
- README: how to generate test

## [4.15.0] - 2022-12-05

### Changed

- Send less alerts into Atlas alert slack channels (filtering out heartbeats and inhibitions)
- Opsgenie messages: revert to markdown

### Added

- Add capz provider

## [4.14.0] - 2022-11-30

### Changed

- Change HasPrometheusAgent function to ignore prometheus-agent scraping targets on gcp.
- Do not reconcile service monitors in kube-system for gcp MCs .

## [4.13.0] - 2022-11-30

### Changed

- Improve HasPrometheusAgent function to ignore prometheus-agent scraping targets.
- Bump alpine from 3.16.3 to 3.17.0
- Do not reconcile service monitors in CAPO MCs.

### Removed

- Remove option to disable PVC creations only used on KVM installations.
- Remove deprecated ingress v1beta1 (only used on kvm).

## [4.12.0] - 2022-11-25

### Changed

- Improve opsgenie notification template.

## [4.11.2] - 2022-11-24

### Fixed

- Remove `vault` targets for CAPI clusters.

## [4.11.1] - 2022-11-22

### Fixed

- Fix reconciliation issues on vintage MCs.

## [4.11.0] - 2022-11-18

### Removed

- Remove the `CLUSTER-prometheus/app-exporter-CLUSTER/0` job in favor of Service Monitor provided by the app.

### Added

- Ensure the remote write endpoint configuration is enabled for MCs
- Add Inhibition rule for prometheus-agent to ignore clusters that doesn't deploy the agent.

### Changed

- Send non-paging alert to Atlas slack channels.

### Fixed

- Fix a reconciliation bug on CAPI MC that were reconciled twice.

## [4.10.0] - 2022-11-15

### Fixed

- Fix scraping of controller-manager and kube-scheduler for vintage MCs.

### Removed

- Old remotewrite secret
- Removed targets for clusters using the prometheus agent.

## [4.9.2] - 2022-11-03

### Changed

- prometheus PSP: allow "projected" volumes

## [4.9.1] - 2022-10-31

### Changed

- Change `remotewritesecret` to always delete the secret, as it's not needed anymore in favor of `RemoteWrite.spec.secrets`.

## [4.9.0] - 2022-10-28

### Added

- Added cadvisor scraping for `flux-*` namespaces.

## [4.8.1] - 2022-10-26

### Fixed

- Reduce label selector to find Prometheus PVC

## [4.8.0] - 2022-10-26

## [4.7.1] - 2022-10-21

### Fixed

- Fix alertmanager psp (add projected and downardAPI)

## [4.7.0] - 2022-10-20

### Added

- Add scraping of Cilium on Management Clusters.
- Add externalLabels for the remote write endpoint configuration.

### Changed

- Configure working queue config for the remote write receiver (reducing max number of shards but increasing sample capacity).

### Fixed

- Fix remote write endpoint ingress buffer size to avoid the use of temporary buffer file in the ingress controller.

## [4.6.4] - 2022-10-17

### Changed

- Move shield route so alerts for shield don't go to opsgenie at all, only to their slack.

## [4.6.3] - 2022-10-17

### Changed

- Customize Prometheus volume size based via the `monitoring.giantswarm.io/prometheus-volume-size` annotation
- Change remotewrite endpoint secrets namespace to clusterID ns.
- Add `.svc` suffix to the alertmanager target to make PMO work behind a corporate proxy.
- Upgrade to go 1.19
- Bump prometheus-operator to v0.54.0

### Added

- Enable remote write receiver.
- Generate prometheus remote write agent secret and config.
- Configure prometheus remote write agent ingress.
- Add Slack channel for Team Shield.

## [4.6.2] - 2022-09-13

### Fixed

- Fix controller manager port to use default or a value from annotation.
- Fix scheduler port to use default or a value from annotation.
- Bump github.com/labstack/echo to v4.9.0 to fix sonatype-2022-5436 CVE.

## [4.6.1] - 2022-09-12

### Fixed

- Drop original `label_topology_kubernetes_io_region` & `label_topology_kubernetes_io_zone` labels.

## [4.6.0] - 2022-09-12

### Added

- Relabeling for labels `label_topology_kubernetes_io_region` & `label_topology_kubernetes_io_zone` to `region` & `zone`.

## [4.5.1] - 2022-08-24

### Fixed

- Fix CAPI MCs being seen as workload cluster.

## [4.5.0] - 2022-08-24

### Changed

- Change CAPI version from v1alpha3 to v1beta1.

## [4.4.1] - 2022-08-19

### Fixed

- Fix Team hydra config.

## [4.4.0] - 2022-08-17

### Added

- Add service priority as a tag in opsgenie alerts.
- Add Team Hydra receiver and route.

### Fixed

- Upgrade go-kit/kit to fix CVE-2022-24450 and CVE-2022-29946.
- Upgrade getsentry/sentry-go to fix CVE-2021-23772, CVE-2021-42576, CVE-2020-26892, and CVE-2021-3127.

## [4.3.0] - 2022-08-02

### Fixed

- Fix psp names for prometheus and alertmanager.

## [4.2.0] - 2022-07-28

### Changed

- Set node-exporter namespace to `kube-system` for CAPI MCs and all WC, and to `monitoring` for vintage MCs.
- Set cert-exporter namespace to `kube-system` for CAPI MCs and all WC, and to `monitoring` for vintage MCs.

### Fixed

- Added `pod_name` as a label to distinguish between multiple etcd pods when running in-cluster (e.g. CAPI).

### Added

- Push to `gcp-app-collection`.

### Changed

- Bump alpine from 3.16.0 to 3.16.1

## [4.1.0] - 2022-07-20

### Changed

- Upgrade operatorkit from v7.0.1 to v7.1.0.
- Upgrade github.com/sirupsen/logrus from 1.8.1 to 1.9.0.

### Added

- errors_total metric for each controller (comes with operatorkit upgrade).

### Fixed

- Cleanup of RemoteWrite Status (configuredPrometheuses, syncedSecrets) in case a cluster gets deleted.

## [4.0.1] - 2022-07-14

### Fixed

- Fix creation of new prometheus instance once a cluster has been created

## [4.0.0] - 2022-07-13

### Added

- Implement remotewrite CR logic, in order to configure Prometheus remotewrite config.
- Add HTTP_PROXY in remotewrite config
- Add unit tests for remotewrite resource
- Add Secrets field in the RemoteWrite CR
- Implement sync RemoteWrite Secrets logic
- Adding RemoteWrite.status field to ensure cleanup
- Add psp and service account for prometheus and alertmanager

### Changed

- Rename `vcd` to `cloud-director`
- Monitor using a podmonitor.

### Fixed

- Fix API server discovery.

### Removed

- Remove duplicate scrape config targets.

### Fixed

- Fix API server discovery.
- Add `patch` verb for `remoteWrite` resources.

## [3.8.0] - 2022-06-30

### Added

- Add Secrets field in the RemoteWrite CR

## [3.7.0] - 2022-06-20

This release was created on release-v3.5.x branch to fix release 3.6.0 see PR#992

### Changed

- Change remote write name to grafana-cloud.

## [3.6.0] - 2022-06-08

### Added

- Add remotewrite controller.
- Deployment of remoteWrite CRD in Helm chart
- Ignore remotewrite field when updating prometheus CR.
- Add `PodMonitor` support for workload cluster Prometheus.

### Fixed

- dependencies updates
- fix build by ignoring CVEs we can't fix for the moment
- Upgrade docker image from Alpine 3.15.1 to Alpine 3.16.0

### Added

- remoteWrite CustomResourceDefinition

## [3.5.0] - 2022-05-17

### Added

- Add Cluster Service Priority label.
- Add customer and organization label to metrics.
- Add VCD provider.

## [3.4.3] - 2022-05-10

### Fixed

- Add 5mn initial delay before performing readiness checks.

## [3.4.2] - 2022-05-09

### Fixed

- Use 'ip' node label as target to scrape etcd on MCs.

## [3.4.1] - 2022-05-05

### Fixed

- Fix CAPI cluster detection for legacy Management Clusters.

## [3.4.0] - 2022-05-04

### Added

- Add `PodMonitor` support on management clusters.

## [3.3.0] - 2022-05-04

### Changed

- Add `nodepool` label to `kube-state-metrics` metrics.
- Improve CAPI cluster detection.

## [3.2.0] - 2022-04-13

### Changed

- Change how MC managed with CAPI are reconciled in PMO (using the cluster CR instead of the Kubernetes Service)

### Fixed

- Fix etcd service discovery for CAPI clusters.

### Removed

- Remove skip resource

## [3.1.0] - 2022-04-08

### Added

- Add support for etcd-certificates on OpenStack.
- Add context to generic resources.

### Fixed

- Add skip resource, to fix MC duplicated handling.

## [3.0.0] - 2022-03-28

### Added

- Add alertmanager ingress.
- Configure alertmanager and wire prometheus to both legacy and new alertmanagers.

### Changed

- Remove deprecated matcher types from alertmanager config.
- Changed scrape_interval to 180s and scrape_timeout to 60s for azure-collector.

### Removed

- Remove old teams from alertmanager template.
- Remove code to manage legacy alertmanager.

## [2.4.0] - 2022-03-16

### Changed

- Migrate to rbac/v1 from rbac/v1beta1.
- Change additional scraping config to keep cadvisor metrics for `kong.*` named namespaces

### Fixed

- Do not trail right whitespaces in config.

## [2.3.0] - 2022-03-04

### Changed

- Support ingress v1 by default.
- Scrape node-exporter trough apiserver proxy.

### Fixed

- Old references to Firecracker and Celestial replaced with Phoenix

## [2.2.1] - 2022-02-24

### Fixed

- Fix failing `aggregation:prometheus:memory_percentage` due to duplicated series from node exporter.

## [2.2.0] - 2022-01-20

### Changed

- Allow overriding the scraping protocol

### Fixed

- Set ingress class name in ingress spec instead of annotation to prepare supporting ingress v1.

## [2.1.1] - 2022-01-12

### Fixed

- Prevent panic when encountering a different user in the CAPI kubeconfig.

## [2.1.0] - 2022-01-10

## Added

- Added support for OpenStack provider

## [2.0.0] - 2022-01-03

### Changed

- Disable cluster-api controller on KVM installations.
- Disable legacy controller on AWS and Azure installations.
- Upgrade to Go 1.17
- Upgrade github.com/giantswarm/microkit v0.2.2 to v1.0.0
- Upgrade github.com/giantswarm/versionbundle v0.2.0 to v1.0.0
- Upgrade github.com/giantswarm/microendpoint v0.2.0 to v1.0.0
- Upgrade github.com/giantswarm/microerror v0.3.0 to v0.4.0
- Upgrade github.com/giantswarm/micrologger v0.5.0 to v0.6.0
- Upgrade github.com/spf13/viper v1.9.0 to v1.10.0
- Upgrade github.com/giantswarm/k8sclient v5.12.0 to v7.0.1
- Upgrade k8s.io/api v0.19.4 to v0.21.4
- Upgrade k8s.io/apiextensions-apiserver v0.19.4 to v0.21.4
- Upgrade sigs.k8s.io/controller-runtime v0.6.4 to v0.8.3
- Upgrade k8s.io/client-go v0.19.4 to v0.21.4
- Upgrade github.com/giantswarm/operatorkit v4.3.1 to v7.0.0
- Upgrade sigs.k8s.io/cluster-api v0.3.19 to v0.4.5
- Upgrade sigs.k8s.io/controller-runtime v0.8.3 to v0.9.7
- Upgrade github.com/prometheus-operator v0.50.0 to v0.52.1

### Removed

- Remove k8sclient.G8sClient
- Remove versionbundle.Changelog
- Remove github.com/giantswarm/cluster-api v0.3.13-gs

## [1.53.0] - 2021-12-17

### Changed

- Renamed `cancel_if_has_no_workers` inhibition to `cancel_if_cluster_has_no_workers` to make it explicit it's about clusters and not node pools.

## [1.52.1] - 2021-12-14

### Fixed

- Fix relabeling for `__meta_kubernetes_service_annotation_giantswarm_io_monitoring_app_label`

## [1.52.0] - 2021-12-13

### Added

- Add new inhibition for clusters without workers.
- Add relabeling for `__meta_kubernetes_service_annotation_giantswarm_io_monitoring_app_label`

### Changed

- Upgrade alertmanager to v0.23.0
- Upgrade prometheus-operator v0.49.0 to v0.50.0

### Fixed

- Avoid defaulting of `role` label (containing the role of the k8s node). If data is missing we can't reliably default it.

## [1.51.2] - 2021-10-28

### Fixed

- Fix finding certificates in organization namespaces.

### Removed

- Remove cloud limit alerts from customer channel.

## [1.51.1] - 2021-09-10

### Fixed

- Re-introduce `v1alpha2` scheme.

## [1.51.0] - 2021-09-09

### Changed

- Drop `v1alpha2` scheme.
- Reconcile `v1alpha3` cluster.

### Fixed

- Do not create the legacy controller on new installations.

## [1.50.0] - 2021-08-16

### Changed

- Upgrade prometheus-operator to v0.49.0

### Fixed

- Fix an issue where prometheus config is empty, due to missing serviceMonitorSelector.

## [1.49.0] - 2021-08-11

### Added

- Add `additionalScrapeConfigs` flag which accepts a string which will be appended to the management cluster scrape config
  template for installation specific configuration.

## [1.48.0] - 2021-08-09

### Added

- Add receiver and route for `#noise-falco` Slack channel.

## [1.47.0] - 2021-08-05

### Changed

- Add the service label in the alert templates for the `ServiceLevelBurnRateTooHigh` alert.
- Update Prometheus to 2.28.1.
- Allow the use of Prometheus Operator Service Monitor for management clusters.

## [1.46.0] - 2021-07-14

### Changed

- Use `giantswarm/config` to generate managed configuration.

## [1.45.0] - 2021-06-28

### Changed

- Use Grafana Cloud remote-write URL from config instead of hardcoding it, to
  allow overriding the URL in installations which can't access Grafana Cloud
  directly.

## [1.44.2] - 2021-06-24

## [1.44.1] - 2021-06-24

## [1.44.0] - 2021-06-23

### Removed

- Migrate existing rules to https://github.com/giantswarm/prometheus-rules.

## [1.43.0] - 2021-06-22

### Changed

- Removed `ServiceLevelBurnRateTicket` alert.

## [1.42.0] - 2021-06-22

### Changed

- Removed `NodeExporterDown` alert and use SLO framework to monitor node-exporters.
- Change `ServiceLevelBurnRateTooHigh` and `ServiceLevelBurnRateTooHighTicket` to opt-out for services.

## [1.41.2] - 2021-06-22

### Fixed

- Fix typo in `AzureClusterCreationFailed` and `AzureClusterUpgradeFailed`

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

[Unreleased]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.77.1...HEAD
[4.77.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.77.0...v4.77.1
[4.77.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.76.0...v4.77.0
[4.76.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.75.1...v4.76.0
[4.75.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.75.0...v4.75.1
[4.75.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.74.0...v4.75.0
[4.74.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.73.1...v4.74.0
[4.73.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.73.0...v4.73.1
[4.73.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.72.0...v4.73.0
[4.72.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.71.0...v4.72.0
[4.71.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.70.3...v4.71.0
[4.70.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.70.2...v4.70.3
[4.70.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.70.1...v4.70.2
[4.70.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.70.0...v4.70.1
[4.70.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.69.0...v4.70.0
[4.69.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.68.4...v4.69.0
[4.68.4]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.68.3...v4.68.4
[4.68.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.68.2...v4.68.3
[4.68.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.68.1...v4.68.2
[4.68.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.68.0...v4.68.1
[4.68.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.67.3...v4.68.0
[4.67.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.67.2...v4.67.3
[4.67.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.67.1...v4.67.2
[4.67.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.67.0...v4.67.1
[4.67.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.66.1...v4.67.0
[4.66.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.66.0...v4.66.1
[4.66.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.65.0...v4.66.0
[4.65.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.64.0...v4.65.0
[4.64.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.63.1...v4.64.0
[4.63.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.63.0...v4.63.1
[4.63.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.62.0...v4.63.0
[4.62.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.61.0...v4.62.0
[4.61.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.60.0...v4.61.0
[4.60.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.59.0...v4.60.0
[4.59.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.58.0...v4.59.0
[4.58.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.57.0...v4.58.0
[4.57.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.56.0...v4.57.0
[4.56.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.55.0...v4.56.0
[4.55.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.54.1...v4.55.0
[4.54.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.54.0...v4.54.1
[4.54.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.53.0...v4.54.0
[4.53.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.52.0...v4.53.0
[4.52.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.51.0...v4.52.0
[4.51.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.50.0...v4.51.0
[4.50.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.49.2...v4.50.0
[4.49.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.49.1...v4.49.2
[4.49.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.49.0...v4.49.1
[4.49.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.48.0...v4.49.0
[4.48.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.47.0...v4.48.0
[4.47.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.46.0...v4.47.0
[4.46.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.45.1...v4.46.0
[4.45.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.45.0...v4.45.1
[4.45.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.44.0...v4.45.0
[4.44.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.43.2...v4.44.0
[4.43.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.43.1...v4.43.2
[4.43.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.43.0...v4.43.1
[4.43.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.42.0...v4.43.0
[4.42.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.41.0...v4.42.0
[4.41.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.40.0...v4.41.0
[4.40.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.39.0...v4.40.0
[4.39.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.38.2...v4.39.0
[4.38.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.38.1...v4.38.2
[4.38.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.38.0...v4.38.1
[4.38.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.37.0...v4.38.0
[4.37.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.36.4...v4.37.0
[4.36.4]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.36.3...v4.36.4
[4.36.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.36.2...v4.36.3
[4.36.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.36.1...v4.36.2
[4.36.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.36.0...v4.36.1
[4.36.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.35.4...v4.36.0
[4.35.4]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.35.3...v4.35.4
[4.35.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.35.2...v4.35.3
[4.35.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.35.1...v4.35.2
[4.35.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.35.0...v4.35.1
[4.35.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.34.0...v4.35.0
[4.34.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.33.0...v4.34.0
[4.33.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.32.0...v4.33.0
[4.32.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.31.1...v4.32.0
[4.31.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.31.0...v4.31.1
[4.31.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.30.0...v4.31.0
[4.30.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.29.2...v4.30.0
[4.29.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.29.1...v4.29.2
[4.29.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.29.0...v4.29.1
[4.29.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.28.0...v4.29.0
[4.28.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.27.0...v4.28.0
[4.27.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.26.0...v4.27.0
[4.26.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.25.3...v4.26.0
[4.25.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.25.2...v4.25.3
[4.25.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.25.1...v4.25.2
[4.25.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.25.0...v4.25.1
[4.25.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.24.1...v4.25.0
[4.24.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.24.0...v4.24.1
[4.24.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.23.0...v4.24.0
[4.23.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.22.0...v4.23.0
[4.22.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.21.0...v4.22.0
[4.21.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.20.6...v4.21.0
[4.20.6]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.20.5...v4.20.6
[4.20.5]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.20.4...v4.20.5
[4.20.4]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.20.3...v4.20.4
[4.20.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.20.2...v4.20.3
[4.20.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.20.1...v4.20.2
[4.20.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.20.0...v4.20.1
[4.20.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.19.2...v4.20.0
[4.19.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.19.1...v4.19.2
[4.19.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.19.0...v4.19.1
[4.19.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.18.0...v4.19.0
[4.18.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.17.0...v4.18.0
[4.17.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.16.0...v4.17.0
[4.16.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.15.0...v4.16.0
[4.15.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.14.0...v4.15.0
[4.14.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.13.0...v4.14.0
[4.13.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.12.0...v4.13.0
[4.12.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.11.2...v4.12.0
[4.11.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.11.1...v4.11.2
[4.11.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.11.0...v4.11.1
[4.11.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.10.0...v4.11.0
[4.10.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.9.2...v4.10.0
[4.9.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.9.1...v4.9.2
[4.9.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.9.0...v4.9.1
[4.9.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.8.1...v4.9.0
[4.8.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.8.0...v4.8.1
[4.8.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.7.1...v4.8.0
[4.7.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.7.0...v4.7.1
[4.7.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.6.4...v4.7.0
[4.6.4]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.6.3...v4.6.4
[4.6.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.6.2...v4.6.3
[4.6.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.6.1...v4.6.2
[4.6.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.6.0...v4.6.1
[4.6.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.5.1...v4.6.0
[4.5.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.5.0...v4.5.1
[4.5.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.4.1...v4.5.0
[4.4.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.4.0...v4.4.1
[4.4.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.3.0...v4.4.0
[4.3.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.2.0...v4.3.0
[4.2.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.1.0...v4.2.0
[4.1.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.0.1...v4.1.0
[4.0.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v4.0.0...v4.0.1
[4.0.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.8.0...v4.0.0
[3.8.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.7.0...v3.8.0
[3.7.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.6.0...v3.7.0
[3.6.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.5.0...v3.6.0
[3.5.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.4.3...v3.5.0
[3.4.3]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.4.2...v3.4.3
[3.4.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.4.1...v3.4.2
[3.4.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.4.0...v3.4.1
[3.4.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.3.0...v3.4.0
[3.3.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.2.0...v3.3.0
[3.2.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.1.0...v3.2.0
[3.1.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v3.0.0...v3.1.0
[3.0.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v2.4.0...v3.0.0
[2.4.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v2.3.0...v2.4.0
[2.3.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v2.2.1...v2.3.0
[2.2.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v2.2.0...v2.2.1
[2.2.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v2.1.1...v2.2.0
[2.1.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v2.1.0...v2.1.1
[2.1.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.53.0...v2.0.0
[1.53.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.52.1...v1.53.0
[1.52.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.52.0...v1.52.1
[1.52.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.51.2...v1.52.0
[1.51.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.51.1...v1.51.2
[1.51.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.51.0...v1.51.1
[1.51.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.50.0...v1.51.0
[1.50.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.49.0...v1.50.0
[1.49.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.48.0...v1.49.0
[1.48.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.47.0...v1.48.0
[1.47.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.46.0...v1.47.0
[1.46.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.45.0...v1.46.0
[1.45.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.44.2...v1.45.0
[1.44.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.44.1...v1.44.2
[1.44.1]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.44.0...v1.44.1
[1.44.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.43.0...v1.44.0
[1.43.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.42.0...v1.43.0
[1.42.0]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.41.2...v1.42.0
[1.41.2]: https://github.com/giantswarm/prometheus-meta-operator/compare/v1.41.1...v1.41.2
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
