package legacyfinalizer

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	Name = "finalizer"

	// Finalizer of old operator's controller.
	legacyFinalizer = "operatorkit.giantswarm.io/prometheus-meta-operator-control-plane-controller"
)

type Config struct {
	CtrlClient client.Client
	Logger     micrologger.Logger
}

// Resource does garbage collection on the AzureConfig CR finalizers.
type Resource struct {
	ctrlClient client.Client
	logger     micrologger.Logger
}

func New(config Config) (*Resource, error) {
	if config.CtrlClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.CtrlClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &Resource{
		ctrlClient: config.CtrlClient,
		logger:     config.Logger,
	}

	return r, nil
}

// Name returns the resource name.
func (r *Resource) Name() string {
	return Name
}
