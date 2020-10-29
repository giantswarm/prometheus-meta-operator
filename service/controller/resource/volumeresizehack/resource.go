package volumeresizehack

import (
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const Name = "volumeresizehack"

type Config struct {
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
}

type Resource struct {
	k8sClient        k8sclient.Interface
	logger           micrologger.Logger
	prometheusClient promclient.Interface
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.PrometheusClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.PrometheusClient must not be empty", config)
	}

	r := &Resource{
		logger:           config.Logger,
		k8sClient:        config.K8sClient,
		prometheusClient: config.PrometheusClient,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
