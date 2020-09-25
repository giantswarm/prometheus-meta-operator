package service

import (
	"github.com/giantswarm/operatorkit/v2/pkg/flag/service/kubernetes"

	"github.com/giantswarm/prometheus-meta-operator/flag/service/installation"
	"github.com/giantswarm/prometheus-meta-operator/flag/service/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/flag/service/provider"
	"github.com/giantswarm/prometheus-meta-operator/flag/service/security"
	"github.com/giantswarm/prometheus-meta-operator/flag/service/vault"
)

// Service is an intermediate data structure for command line configuration flags.
type Service struct {
	Kubernetes   kubernetes.Kubernetes
	Prometheus   prometheus.Prometheus
	Provider     provider.Provider
	Security     security.Security
	Installation installation.Installation
	Vault        vault.Vault
}
