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

### How to update upstream code

We store modified upstream code for our own usage.

- pkg/alertmanager/config
- pkg/prometheus/common/config

Example on how to update pkg/alertmanager/config :

```
$ git checkout -b upstream-code
$ git tag -d $(git tag -l)
$ git remote add -f alertmanager https://github.com/prometheus/alertmanager.git
$ git checkout v0.22.2
$ git subtree split -P config/ -b alertmanager-config
$ git checkout upstream-code
$ git subtree merge --squash -P pkg/alertmanager/config alertmanager-config
# fix conflicts if any and commit
# push for review
$ git push -u origin HEAD

# restore local tags
$ git tag -d $(git tag -l)
$ git fetch origin
```

[operatorkit]: https://github.com/giantswarm/operatorkit
[prometheus-operator]: https://github.com/prometheus-operator/prometheus-operator
