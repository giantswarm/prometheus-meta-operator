package servicemonitor

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	serviceMonitors, err := toServiceMonitors(obj, r.provider)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "creating servicemonitor")
	for _, desired := range serviceMonitors {
		current, err := r.prometheusClient.MonitoringV1().ServiceMonitors(desired.GetNamespace()).Get(ctx, desired.GetName(), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			current, err = r.prometheusClient.MonitoringV1().ServiceMonitors(desired.GetNamespace()).Create(ctx, desired, metav1.CreateOptions{})
		}
		if err != nil {
			return microerror.Mask(err)
		}

		if hasChanged(current, desired) {
			desired.ObjectMeta = current.ObjectMeta
			_, err = r.prometheusClient.MonitoringV1().ServiceMonitors(desired.GetNamespace()).Update(ctx, desired, metav1.UpdateOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}
	r.logger.LogCtx(ctx, "level", "debug", "message", "created servicemonitor")

	return nil
}
