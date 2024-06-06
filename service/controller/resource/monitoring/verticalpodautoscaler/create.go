package verticalpodautoscaler

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/resourceutils"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
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
		resourceutils.UpdateMeta(current, desired)
		_, err = r.vpaClient.AutoscalingV1().VerticalPodAutoscalers(desired.GetNamespace()).Update(ctx, desired, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.Debugf(ctx, "updated")
	}

	return nil
}
