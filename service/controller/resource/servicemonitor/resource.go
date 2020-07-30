package servicemonitor

import (
	"reflect"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/servicemonitor/service"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "servicemonitor"
)

type Config struct {
	PrometheusClient promclient.Interface
	Logger           micrologger.Logger
	Provider         string
}

type Resource struct {
	prometheusClient promclient.Interface
	logger           micrologger.Logger
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
		provider:         config.Provider,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func toServiceMonitors(obj interface{}, provider string) ([]*promv1.ServiceMonitor, error) {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	serviceMonitors := []*promv1.ServiceMonitor{
		service.APIServer(cluster, provider),
		service.NginxIngressController(cluster, provider),
	}

	if provider == "aws" || provider == "azure" {
		serviceMonitors = append(serviceMonitors, service.ClusterAutoscaler(cluster, provider))
	}

	return serviceMonitors, nil
}

func hasChanged(current, desired *promv1.ServiceMonitor) bool {
	return !reflect.DeepEqual(current.Spec, desired.Spec)
}
