package servicemonitor

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	services, err := toServiceMonitors(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, service := range services {
		err := r.prometheusClient.MonitoringV1().ServiceMonitors(service.GetNamespace()).Delete(service.GetName(), &metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
