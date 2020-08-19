package servicemonitor

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	serviceMonitors, err := toServiceMonitors(obj, r.provider)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "deleting servicemonitor")
	for _, serviceMonitor := range serviceMonitors {
		err := r.prometheusClient.MonitoringV1().ServiceMonitors(serviceMonitor.GetNamespace()).Delete(ctx, serviceMonitor.GetName(), metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}
	r.logger.LogCtx(ctx, "level", "debug", "message", "deleted servicemonitor")

	return nil
}
