package managementcluster

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v8/pkg/resource"
	v1 "k8s.io/api/core/v1"
)

type cpResourceConfig struct {
	Resource     resource.Interface
	Installation string
}

type cpResource struct {
	resource     resource.Interface
	installation string
}

// Wrap wrap resource and replace the input object with a Service which is named after the installation.
func Wrap(resources []resource.Interface, config resourcesConfig) ([]resource.Interface, error) {
	var wrapped []resource.Interface

	for _, r := range resources {
		c := cpResourceConfig{
			Resource:     r,
			Installation: config.Installation,
		}

		cpResource, err := newCPResource(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrapped = append(wrapped, cpResource)
	}

	return wrapped, nil
}

func newCPResource(config cpResourceConfig) (*cpResource, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}

	r := &cpResource{
		resource:     config.Resource,
		installation: config.Installation,
	}

	return r, nil
}

func (r *cpResource) EnsureCreated(ctx context.Context, obj interface{}) error {
	err := r.resource.EnsureCreated(ctx, r.cpObject(obj))
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *cpResource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	err := r.resource.EnsureDeleted(ctx, r.cpObject(obj))
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *cpResource) Name() string {
	return r.resource.Name()
}

func (r *cpResource) cpObject(obj interface{}) *v1.Service {
	svc := obj.(*v1.Service)
	svc.SetName(r.installation)

	return svc
}
