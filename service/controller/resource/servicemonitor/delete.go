package servicemonitor

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	serviceMonitors, err := toServiceMonitors(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "deleting servicemonitor")
	for _, serviceMonitor := range serviceMonitors {
		err := r.prometheusClient.MonitoringV1().ServiceMonitors(serviceMonitor.GetNamespace()).Delete(serviceMonitor.GetName(), &metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}
	r.logger.LogCtx(ctx, "deleted servicemonitor")

	return nil
}
