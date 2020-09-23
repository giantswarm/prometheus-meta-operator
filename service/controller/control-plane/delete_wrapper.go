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

func DeleteWrap(resources []resource.Interface, config resourcesConfig) ([]resource.Interface, error) {
	var wrapped []resource.Interface

	for _, r := range resources {
		c := deleteResourceConfig{
			Resource: r,
		}

		resource, err := newDeleteResource(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrapped = append(wrapped, resource)
	}

	return wrapped, nil
}

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
