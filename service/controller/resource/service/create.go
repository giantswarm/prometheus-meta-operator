package service

import (
	"context"

	"github.com/giantswarm/microerror"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	// This resource is being removed because we do not need it currently.
	// see: https://github.com/giantswarm/giantswarm/issues/11479
	err := r.EnsureDeleted(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
