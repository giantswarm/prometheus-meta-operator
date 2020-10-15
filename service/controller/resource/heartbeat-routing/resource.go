package heartbeatrouting

import (
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const Name = "heartbeatrouting"

type Config struct {
	K8sClient    k8sclient.Interface
	Logger       micrologger.Logger
	Installation string
	Provider     string
}

type Resource struct {
	k8sClient    k8sclient.Interface
	logger       micrologger.Logger
	installation string
	provider     string
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

	r := &Resource{
		logger:       config.Logger,
		k8sClient:    config.K8sClient,
		installation: config.Installation,
		provider:     config.Provider,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}
