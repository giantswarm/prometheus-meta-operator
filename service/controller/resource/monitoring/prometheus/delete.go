package prometheus

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	object, err := r.getObjectMeta(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.config.Logger.Debugf(ctx, "deleting")
	err = r.config.PrometheusClient.MonitoringV1().Prometheuses(object.GetNamespace()).Delete(ctx, object.GetName(), metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		// fall through
	} else if err != nil {
		return microerror.Mask(err)
	}
	r.config.Logger.Debugf(ctx, "deleted")

	return nil
}
