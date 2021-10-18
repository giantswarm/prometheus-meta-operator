# Targets

- alertmanager
- api
- app-exporter
- app-operator
- cadvisor
- calico
- cert-operator
- chart-operator
- cluster-autoscaler
- cluster-operator
- cluster-service
- companyd
- coredns
- credentiald
- crsync
- docker
- draughtsman
- etcd
- etcd-backup-operator
- fluentbit
- g8s-grafana
- g8s-prometheus
- gatekeeper
- kube-proxy
- kube-state-metrics
- kubelet
- kubernetes-apiserver
- kubernetes-controller-manager
- kubernetes-scheduler
- kubernetesd
- net-exporter
- nginx-ingress-controller
- node-exporter
- node-operator
- organization-operator
- prometheus-meta-operator
- release-operator
- tokend
- userd
- vault
- vault-exporter

### AWS specific

- aws-node
- aws-operator
- cluster-autoscaler

### Azure specific

- azure-operator
- cluster-autoscaler

### KVM specific

- kvm-operator
- bridge-operator
- flannel-operator
- ingress-exporter

## Missing

- prometheus-operator-app
- CP prometheus
- TC prometheus
- azure-collector

## Not scraped

- cert-manager
- default-http-backend
- dex
- dex-k8s-authenticator
- external-dns
- oauth2-proxy
- passage
- passage-redis
- tiller

### AWS specific

- admission-controller-unique
- calico-typha
- kiam
- opa-mutator-app

### Azure specific

- azure-admission-controller
- calico-typha

### KVM specific

- calico-kube-controllers
- heapster
- monitoring-influxdb

## Unscrapeable

* happa: not exposing any metrics

* metrics-server: not scrapping because according to official documentation it should not be used [as a source of monitoring solution metrics](https://github.com/kubernetes-sigs/metrics-server#kubernetes-metrics-server).

## Known issues

* docker: fails to be scrapped, due to metrics port not being exposed on giantswarm releases below 11.0.0, see [commit](https://github.com/giantswarm/k8scloudconfig/commit/6ecc07e665c3e854dfa8be102a8c6446d1d9dc3c#diff-be6122463e3fe598d118a80e09254d3d)

* nginx-ingress-controller: fails to be scrapped, due to network policy not allowing metrics port (10254). This is fix from [giantswarm release v10.0.0](https://github.com/giantswarm/releases/tree/master/aws/archived/v10.1.0)
