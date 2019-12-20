package alert

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	prometheusRules, err := toPrometheusRules(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, prometheusRule := range prometheusRules {
		_, err = r.prometheusClient.MonitoringV1().PrometheusRules(prometheusRule.GetNamespace()).Create(prometheusRule)
		if apierrors.IsAlreadyExists(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}
