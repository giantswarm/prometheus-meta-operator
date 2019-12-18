package servicemonitor

import (
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const (
	Name = "servicemonitor"
)

type Config struct {
	PrometheusClient promclient.Interface
	Logger           micrologger.Logger
}

type Resource struct {
	prometheusClient promclient.Interface
	logger           micrologger.Logger
}

func New(config Config) (*Resource, error) {
	if config.PrometheusClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.PrometheusClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &Resource{
		prometheusClient: config.PrometheusClient,
		logger:           config.Logger,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
