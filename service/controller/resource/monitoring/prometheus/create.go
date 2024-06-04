package prometheus

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

	desired, err := r.toPrometheus(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Debugf(ctx, "creating")
	current, err := r.prometheusClient.MonitoringV1().Prometheuses(desired.GetNamespace()).Get(ctx, desired.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		current, err = r.prometheusClient.MonitoringV1().Prometheuses(desired.GetNamespace()).Create(ctx, desired, metav1.CreateOptions{})
	}

	if err != nil {
		return microerror.Mask(err)
	}

	if r.hasChanged(current, desired) {
		updateMeta(current, desired)
		_, err = r.prometheusClient.MonitoringV1().Prometheuses(desired.GetNamespace()).Update(ctx, desired, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
	}
	r.logger.Debugf(ctx, "created")

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
	// without this, it's impossible to change labels on resources
	if len(d.GetLabels()) == 0 {
		d.SetLabels(c.GetLabels())
	}
	// without this, it's impossible to change annotations on resources
	if len(d.GetAnnotations()) == 0 {
		d.SetAnnotations(c.GetAnnotations())
	}
	d.SetFinalizers(c.GetFinalizers())
	d.SetOwnerReferences(c.GetOwnerReferences())
	d.SetManagedFields(c.GetManagedFields())
}
