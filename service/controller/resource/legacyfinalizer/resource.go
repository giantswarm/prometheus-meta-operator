package legacyfinalizer

import (
	"context"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
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

// EnsureCreated ensures that reconciled AzureConfig CR gets orphaned finalizer
// deleted.
func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cr, err := key.ToClusterMR(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Debugf(ctx, "ensuring resource doesn't have orphaned prometheus-meta-operator-control-plane-controller finalizer present")

	{
		// Refresh the CR object.
		err := r.ctrlClient.Get(ctx, client.ObjectKey{Name: cr.GetName(), Namespace: cr.GetNamespace()}, cr)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	var exists bool
	finalizers := cr.GetFinalizers()
	for i, v := range finalizers {
		if strings.TrimSpace(v) == legacyFinalizer {
			exists = true

			// Drop it.
			cr.SetFinalizers(append(finalizers[:i], finalizers[i+1:]...))
			break
		}
	}

	if exists {
		r.logger.Debugf(ctx, "deleting legacy finalizer from resource")

		err := r.ctrlClient.Update(ctx, cr)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.Debugf(ctx, "deleted legacy finalizer from resource")
		return nil
	}

	r.logger.Debugf(ctx, "ensured resource doesn't have orphaned prometheus-meta-operator-control-plane-controller finalizer present")

	return nil
}

// EnsureDeleted is no-op.
func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	return nil
}

// Name returns the resource name.
func (r *Resource) Name() string {
	return Name
}
