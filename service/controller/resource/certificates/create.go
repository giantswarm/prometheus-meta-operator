package certificates

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/resourceutils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	desired, err := r.getDesiredObject(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Debugf(ctx, "creating")
	c := r.k8sClient.K8sClient().CoreV1().Secrets(desired.GetNamespace())
	current, err := c.Get(ctx, desired.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		current, err = c.Create(ctx, desired, metav1.CreateOptions{})
	}
	if err != nil {
		return microerror.Mask(err)
	}

	if r.hasChanged(current, desired) {
		resourceutils.UpdateMeta(current, desired)
		_, err = c.Update(ctx, desired, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
	}
	r.logger.Debugf(ctx, "created")

	return nil
}
