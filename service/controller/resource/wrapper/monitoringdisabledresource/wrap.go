package monitoringdisabledresource

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/operatorkit/v8/pkg/resource"
)

// WrapConfig is the configuration used to wrap resources with disabled monitoring resource.
type WrapConfig struct {
	Logger micrologger.Logger
}

// Wrap wraps each given resource with a disabled monitoring resource and returns the list of
// wrapped resources.
func Wrap(resources []resource.Interface, config WrapConfig) ([]resource.Interface, error) {
	var wrapped []resource.Interface

	for _, r := range resources {
		c := Config{
			Resource: r,
			Logger:   config.Logger,
		}

		monitoringDisabledResource, err := New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrapped = append(wrapped, monitoringDisabledResource)
	}

	return wrapped, nil
}
