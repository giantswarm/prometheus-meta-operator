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
$ git subtree merge --squash -P pkg/alertmanager/config alertmanager-config
# fix conflicts (the usual way) if any

# restore local tags
$ git tag -d $(git tag -l)
$ git fetch

# push for review
$ git push -u origin HEAD

/!\ Do not merge with squash, once approved merge to master manually.
/!\ We need to preserve commit history otherwise following git subtree commands won't work.
$ git checkout master
$ git merge alertmanager-0.22.2
$ git push
```


[operatorkit]: https://github.com/giantswarm/operatorkit
[prometheus-operator]: https://github.com/prometheus-operator/prometheus-operator
