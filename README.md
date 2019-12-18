[![CircleCI](https://circleci.com/gh/giantswarm/prometheus-meta-operator.svg?&style=shield)](https://circleci.com/gh/giantswarm/prometheus-meta-operator) [![Docker Repository on Quay](https://quay.io/repository/giantswarm/prometheus-meta-operator/status "Docker Repository on Quay")](https://quay.io/repository/giantswarm/prometheus-meta-operator)

# prometheus-meta-operator

The prometheus-meta-operator watches Cluster CR and creates [prometheus-operator] CR. It is implemented
using [operatorkit].

## Getting Project

Clone the git repository: https://github.com/giantswarm/prometheus-meta-operator.git

### How to build

Build it using the standard `go build` command.

```
go build github.com/giantswarm/prometheus-meta-operator
```


[operatorkit]: https://github.com/giantswarm/operatorkit
[prometheus-operator]: https://github.com/coreos/prometheus-operator
