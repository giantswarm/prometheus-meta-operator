package controlplane

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v2/pkg/resource"
)

type deleteResourceConfig struct {
	Resource resource.Interface
}

type deleteResource struct {
	resource resource.Interface
}

// newDeleteResource convert a resource to be a delete only resource.
// EnsureDeleted method will be called in every cases (creation and deletion).
func newDeleteResource(config deleteResourceConfig) (*deleteResource, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	r := &deleteResource{
		resource: config.Resource,
	}

	return r, nil
}

func (r *deleteResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	err := r.resource.EnsureDeleted(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *deleteResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	err := r.resource.EnsureDeleted(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *deleteResource) Name() string {
	return r.resource.Name()
}
