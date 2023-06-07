package service

import (
	"github.com/giantswarm/operatorkit/v7/pkg/flag/service/kubernetes"

	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/ingress"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/installation"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/opsgenie"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/provider"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/security"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag/service/vault"
)

// Service is an intermediate data structure for command line configuration flags.
type Service struct {
	Ingress      ingress.Ingress
	Installation installation.Installation
	Kubernetes   kubernetes.Kubernetes
	Opsgenie     opsgenie.Opsgenie
	Prometheus   prometheus.Prometheus
	Provider     provider.Provider
	Security     security.Security
	Vault        vault.Vault
}
