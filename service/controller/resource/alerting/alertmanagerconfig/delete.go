package alertmanagerconfig

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	object := r.getObjectMeta()

	r.logger.Debugf(ctx, "deleting")
	err := r.k8sClient.K8sClient().CoreV1().Secrets(object.GetNamespace()).Delete(ctx, object.GetName(), metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		// fall through
	} else if err != nil {
		return microerror.Mask(err)
	}
	r.logger.Debugf(ctx, "deleted")

	return nil
}
