package alert

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	prometheusRules, err := toPrometheusRules(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "deleting alert rules")
	for _, prometheusRule := range prometheusRules {
		err := r.prometheusClient.MonitoringV1().ServiceMonitors(prometheusRule.GetNamespace()).Delete(prometheusRule.GetName(), &metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}
	r.logger.LogCtx(ctx, "deleted alert rules")

	return nil
}
