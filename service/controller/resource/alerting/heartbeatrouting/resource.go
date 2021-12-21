package heartbeatrouting

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const Name = "heartbeatrouting"

type Config struct {
	Installation string
	K8sClient    k8sclient.Interface
	Logger       micrologger.Logger
	OpsgenieKey  string
}

type Resource struct {
	installation string
	k8sClient    k8sclient.Interface
	logger       micrologger.Logger
	opsgenieKey  string
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}
	if config.OpsgenieKey == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.OpsgenieKey must not be empty", config)
	}

	r := &Resource{
		installation: config.Installation,
		k8sClient:    config.K8sClient,
		logger:       config.Logger,
		opsgenieKey:  config.OpsgenieKey,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
