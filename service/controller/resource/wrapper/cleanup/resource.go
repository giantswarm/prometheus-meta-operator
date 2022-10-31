package cleanup

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
)

type Config struct {
	Resource resource.Interface
	Logger   micrologger.Logger
}

type cleanupWrapper struct {
	resource resource.Interface
	logger   micrologger.Logger
}

// New returns a new cleanup Wrapper to always call EnsureDeleted
func New(config Config) (resource.Interface, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &cleanupWrapper{
		resource: config.Resource,
		logger:   config.Logger,
	}

	return r, nil
}

func (r *cleanupWrapper) EnsureCreated(ctx context.Context, obj interface{}) error {

	return r.resource.EnsureDeleted(ctx, obj)
}

func (r *cleanupWrapper) EnsureDeleted(ctx context.Context, obj interface{}) error {

	return r.resource.EnsureDeleted(ctx, obj)
}

func (r *cleanupWrapper) Name() string {
	return r.resource.Name()
}
