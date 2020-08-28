package service

import (
	"github.com/giantswarm/operatorkit/v2/pkg/flag/service/kubernetes"

	"github.com/giantswarm/prometheus-meta-operator/flag/service/etcd"
	"github.com/giantswarm/prometheus-meta-operator/flag/service/installation"
	"github.com/giantswarm/prometheus-meta-operator/flag/service/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/flag/service/provider"
	"github.com/giantswarm/prometheus-meta-operator/flag/service/vault"
)

// Service is an intermediate data structure for command line configuration flags.
type Service struct {
	Kubernetes   kubernetes.Kubernetes
	Prometheus   prometheus.Prometheus
	Provider     provider.Provider
	Installation installation.Installation
	Etcd         etcd.Etcd
	Vault        vault.Vault
}
