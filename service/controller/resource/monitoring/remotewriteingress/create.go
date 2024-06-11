package remotewriteingress

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/resourceutils"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	if r.config.MimirEnabled {
		r.config.Logger.Debugf(ctx, "mimir is enabled, deleting heartbeat if it exists")
		return r.EnsureDeleted(ctx, obj)
	}

	desired, err := r.toIngress(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.config.Logger.Debugf(ctx, "creating")
	current, err := r.config.K8sClient.K8sClient().NetworkingV1().Ingresses(desired.GetNamespace()).Get(ctx, desired.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		current, err = r.config.K8sClient.K8sClient().NetworkingV1().Ingresses(desired.GetNamespace()).Create(ctx, desired, metav1.CreateOptions{})
	}

	if err != nil {
		return microerror.Mask(err)
	}

	if r.hasChanged(current, desired) {
		resourceutils.UpdateMeta(current, desired)
		_, err = r.config.K8sClient.K8sClient().NetworkingV1().Ingresses(desired.GetNamespace()).Update(ctx, desired, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
	}
	r.config.Logger.Debugf(ctx, "created")

	return nil
}
