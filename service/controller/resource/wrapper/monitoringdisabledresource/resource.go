package monitoringdisabledresource

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v8/pkg/resource"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

type Config struct {
	Resource resource.Interface
	Logger   micrologger.Logger
}

type monitoringDisabledWrapper struct {
	resource resource.Interface
	logger   micrologger.Logger
}

// New returns a new monitoring disabled wrapper according to the configured resource's
// implementation, which might be resource.Interface or crud.Interface. This has
// then different implications on how to measure metrics for the different
// methods of the interfaces.
func New(config Config) (resource.Interface, error) {
	if config.Resource == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}

	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &monitoringDisabledWrapper{
		resource: config.Resource,
		logger:   config.Logger,
	}

	return r, nil
}

func (r *monitoringDisabledWrapper) EnsureCreated(ctx context.Context, obj interface{}) error {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	if key.IsMonitoringDisabled(cluster) {
		r.logger.Debugf(ctx, "monitoring disabled, cleaning up existing monitoring")
		return r.resource.EnsureDeleted(ctx, obj)
	}

	return r.resource.EnsureCreated(ctx, obj)
}

func (r *monitoringDisabledWrapper) EnsureDeleted(ctx context.Context, obj interface{}) error {
	return r.resource.EnsureDeleted(ctx, obj)
}

func (r *monitoringDisabledWrapper) Name() string {
	return r.resource.Name()
}
