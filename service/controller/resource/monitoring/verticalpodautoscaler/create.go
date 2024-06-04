package verticalpodautoscaler

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	if r.mimirEnabled {
		r.logger.Debugf(ctx, "mimir is enabled, deleting")
		return r.EnsureDeleted(ctx, obj)
	}

	desired, err := r.getObject(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Debugf(ctx, "checking if vpa cr already exists")
	current, err := r.vpaClient.AutoscalingV1().VerticalPodAutoscalers(desired.GetNamespace()).Get(ctx, desired.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		r.logger.Debugf(ctx, "creating")
		_, err = r.vpaClient.AutoscalingV1().VerticalPodAutoscalers(desired.GetNamespace()).Create(ctx, desired, metav1.CreateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.Debugf(ctx, "created")
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Debugf(ctx, "checking if vpa cr needs to be updated")
	if hasChanged(current, desired) {
		r.logger.Debugf(ctx, "updating")
		updateMeta(current, desired)
		_, err = r.vpaClient.AutoscalingV1().VerticalPodAutoscalers(desired.GetNamespace()).Update(ctx, desired, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.Debugf(ctx, "updated")
	}

	return nil
}

func updateMeta(c, d metav1.Object) {
	d.SetGenerateName(c.GetGenerateName())
	d.SetUID(c.GetUID())
	d.SetResourceVersion(c.GetResourceVersion())
	d.SetGeneration(c.GetGeneration())
	d.SetSelfLink(c.GetSelfLink())
	d.SetCreationTimestamp(c.GetCreationTimestamp())
	d.SetDeletionTimestamp(c.GetDeletionTimestamp())
	d.SetDeletionGracePeriodSeconds(c.GetDeletionGracePeriodSeconds())
	d.SetLabels(c.GetLabels())
	d.SetAnnotations(c.GetAnnotations())
	d.SetFinalizers(c.GetFinalizers())
	d.SetOwnerReferences(c.GetOwnerReferences())
	d.SetManagedFields(c.GetManagedFields())
}
