package podmonitor

import (
	"reflect"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "podmonitor"
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

func toPodMonitors(obj interface{}, provider string) ([]*promv1.PodMonitor, error) {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	podMonitors := []*promv1.PodMonitor{}

	return podMonitors, nil
}

func hasChanged(current, desired *promv1.PodMonitor) bool {
	return !reflect.DeepEqual(current.Spec, desired.Spec)
}
