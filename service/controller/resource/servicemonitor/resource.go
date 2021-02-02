package servicemonitor

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/servicemonitor/service"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "servicemonitor"
)

type Config struct {
	PrometheusClient promclient.Interface
	Logger           micrologger.Logger
	Installation     string
	Provider         string
}

// TODO: remove this resource in the next release.
type Resource struct {
	prometheusClient promclient.Interface
	logger           micrologger.Logger
	installation     string
	provider         string
}

func New(config Config) (*Resource, error) {
	if config.PrometheusClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.PrometheusClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Provider == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Provider must not be empty", config)
	}

	r := &Resource{
		prometheusClient: config.PrometheusClient,
		logger:           config.Logger,
		installation:     config.Installation,
		provider:         config.Provider,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func toServiceMonitors(cluster metav1.Object, provider string, installation string) ([]*promv1.ServiceMonitor, error) {
	serviceMonitors := []*promv1.ServiceMonitor{
		service.APIServer(cluster, provider, installation),
		service.NginxIngressController(cluster, provider, installation),
	}

	if (provider == "aws" || provider == "azure") && key.ClusterType(cluster) == "workload_cluster" {
		serviceMonitors = append(serviceMonitors, service.ClusterAutoscaler(cluster, provider, installation))
	}

	return serviceMonitors, nil
}
