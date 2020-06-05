package service

import (
	"github.com/giantswarm/operatorkit/flag/service/kubernetes"

	"github.com/giantswarm/prometheus-meta-operator/flag/service/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/flag/service/provider"
)

// Service is an intermediate data structure for command line configuration flags.
type Service struct {
	Kubernetes kubernetes.Kubernetes
	Prometheus prometheus.Prometheus
	Provider   provider.Provider
}
