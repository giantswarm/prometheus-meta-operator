package alert

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	serviceMonitors, err := toServiceMonitors(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, serviceMonitor := range serviceMonitors {
		_, err = r.prometheusClient.MonitoringV1().ServiceMonitors(serviceMonitor.GetNamespace()).Create(serviceMonitor)
		if apierrors.IsAlreadyExists(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
