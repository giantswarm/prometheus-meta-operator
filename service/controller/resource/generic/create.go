package generic

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cr, err := r.toCR(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "creating")
	_, err = r.client.Create(cr)
	if apierrors.IsAlreadyExists(err) {
		_, err = r.client.Update(cr)
	}

	if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "level", "debug", "message", "created")

	return nil
}
