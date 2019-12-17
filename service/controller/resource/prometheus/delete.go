package prometheus

import (
	"context"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	prometheus, err := toPrometheus(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.prometheusClient.MonitoringV1().Prometheuses(prometheus.GetNamespace()).Delete(prometheus.GetName(), &metav1.DeleteOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
