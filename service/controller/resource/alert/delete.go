package alert

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	rules, err := getRules(obj, r.installation)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, rule := range rules {
		r.logger.LogCtx(ctx, "level", "debug", "message", "deleting rule", "rule", rule.GetName())
		err = r.prometheusClient.MonitoringV1().PrometheusRules(key.Namespace(cluster)).Delete(ctx, rule.GetName(), metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// We ignore not found errors here
			return nil
		} else if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "deleted rule", "rule", rule.GetName())
	}

	return nil
}
