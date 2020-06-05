package generic

import (
	"context"
	"reflect"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cr, err := r.toCR(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "creating")
	current, err := r.client.Get(ct.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = r.client.Create(cr)
	} else if apierrors.IsAlreadyExists(err) {
		current.ObjectMeta = nil
		if !reflect.DeepEqual(current, cr) {
			_, err = r.client.Update(cr)
			if err != nil {
				return microerror.Mask(err)
			}

		}
	} else if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "level", "debug", "message", "created")

	return nil
}
