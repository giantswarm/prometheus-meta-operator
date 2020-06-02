package prometheus

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	prometheus, err := toPrometheus(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "creating prometheus")
	_, err = r.prometheusClient.MonitoringV1().Prometheuses(prometheus.GetNamespace()).Create(prometheus)
	if apierrors.IsAlreadyExists(err) {
		// fall through
	} else if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "created prometheus")

	return nil
}
