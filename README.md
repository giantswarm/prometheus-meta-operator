[![CircleCI](https://circleci.com/gh/giantswarm/prometheus-meta-operator.svg?&style=shield)](https://circleci.com/gh/giantswarm/prometheus-meta-operator) [![Docker Repository on Quay](https://quay.io/repository/giantswarm/prometheus-meta-operator/status "Docker Repository on Quay")](https://quay.io/repository/giantswarm/prometheus-meta-operator)

# prometheus-meta-operator

The prometheus-meta-operator watches Cluster CRs and creates [prometheus-operator] CRs. It is implemented
using [operatorkit].

## Getting Project

Clone the git repository: https://github.com/giantswarm/prometheus-meta-operator.git

### How to build

Build it using the standard `go build` command.

```
go build github.com/giantswarm/prometheus-meta-operator
```

You may want to regenerate the unit test files with:
```
go test -v ./... -update
```

### How to update upstream code

We store modified upstream code for our own usage.

- pkg/alertmanager/config
- pkg/prometheus/common/config

#### Initial upstream setup

Add the upstream git repository:

```
$ git remote add alertmanager https://github.com/prometheus/alertmanager.git
```

On first run commands are the same as for Upgrade except for `git subtree merge` which has to be replaced with:

```
$ git subtree add --squash -P pkg/alertmanager/config alertmanager-config
```


#### Upgrade

```
# add upstream tags
$ git tag -d $(git tag -l)
$ git fetch alertmanager

$ git checkout v0.22.2
$ git subtree split -P config/ -b alertmanager-config
$ git checkout -b alertmanager-0.22.2 origin/master
$ git subtree merge --message "Upgrade alertmanager/config to v0.22.2" --squash -P pkg/alertmanager/config alertmanager-config
# fix conflicts (the usual way) if any

# restore local tags
$ git tag -d $(git tag -l)
$ git fetch

# push for review
$ git push -u origin HEAD

/!\ Do not merge with squash, once approved merge to master manually.
/!\ We need to preserve commit history otherwise following git subtree commands won't work.
$ git checkout master
$ git merge --ff-only alertmanager-0.22.2
$ git push
```

# remoteWrite CRs

Prometheus-meta-operator also manages remoteWrite custom resources.


## remoteWrite CRDs

Code for remoteWrite CRDs is in the `api/v1alpha1/` directory.

The actual CRDs are in `config/crd/monitoring.giantswarm.io_remotewrites.yaml`

To generate the CRDs from code, just use `make generate`.

## Deployment

CRDs deployment is managed within the helm chart.
The remoteWrite CRD is located under the chart's templates directory as a symbolic link to the generated yaml file. 

[operatorkit]: https://github.com/giantswarm/operatorkit
[prometheus-operator]: https://github.com/prometheus-operator/prometheus-operator

# Custom Prometheus volume size

Prometheus-meta-operator provides a way of setting custom Prometheus volume size.

The Prometheus volume size can be set on the cluster CR using the dedicated annotation `monitoring.giantswarm.io/prometheus-volume-size`

Three values are possible:

* `small` = 30 Gi
* `medium` = 100 Gi
* `large` = 200 Gi

while `medium` is the default value.


The retention size of prometheis will be set according to the volume size: we apply a ratio of 90%:

* `small` (30 Gi) => retentionSize = 27Gi
* `medium` (100 Gi) => retentionSize = 90Gi
* `large` (200 Gi) => retentionSize = 180Gi

Check [Prometheus Volume Sizing](https://docs.giantswarm.io/getting-started/observability/monitoring/prometheus/volume-size/) for more details.

# Prometheus Agent Sharding

Prometheus Meta Operator configures the Prometheus Agent instances running in workload clusters (pre-mimir setup cf. observability-operator).

To be able to ingest metrics without disrupting the workload running in the clusters, Prometheus Meta Operator can shard the number of running Prometheus Agents.

The default configuration is defined in PMO itself PMO add a new shard every 1M time series present in the WC prometheus running on the management cluster. To avoid scaling down too abruptly, we defined a scale down threshold of 20%.

As this default value was not enough to avoid workload disruptions, we added 2 ways to be able to override the scale up series count target and the scale down percentage.

1. Those values can be configured at the installation level by overriding the following values:

```yaml
prometheusAgent:
  shardScaleUpSeriesCount: 1000000
  shardScaleDownPercentage: 0.20
```

2. Those values can also be set per cluster using the following cluster annotations:

```yaml
monitoring.giantswarm.io/prometheus-agent-scale-up-series-count: 1000000
monitoring.giantswarm.io/prometheus-agent-scale-down-percentage: 0.20
```
