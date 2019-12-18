package servicemonitor

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	services, err := toServiceMonitors(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, service := range services {
		_, err = r.prometheusClient.MonitoringV1().ServiceMonitors(service.GetNamespace()).Create(service)
		if apierrors.IsAlreadyExists(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
