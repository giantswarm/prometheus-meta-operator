package alert

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	rules, err := r.GetRules(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "ensuring the prometheus rules exists")
	for _, desiredRule := range rules {
		ruleName := desiredRule.GetName()
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("checking if rule %s needs to be created or updated", ruleName))
		currentRule, err := r.prometheusClient.MonitoringV1().PrometheusRules(key.Namespace(cluster)).Get(ctx, ruleName, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("rule %s needs to be created", ruleName))
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating rule %s", ruleName))
			currentRule, err = r.prometheusClient.MonitoringV1().PrometheusRules(key.Namespace(cluster)).Create(ctx, desiredRule, metav1.CreateOptions{})
			if err != nil {
				r.logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("could not create rule %s", ruleName))
				return microerror.Mask(err)
			}
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("created rule %s", ruleName))
		} else if err != nil {
			return microerror.Mask(err)
		}

		if hasChanged(currentRule, desiredRule) {
			updateMeta(currentRule, desiredRule)
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("rule %s needs to be updated", ruleName))
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updating rule %s", ruleName))

			_, err := r.prometheusClient.MonitoringV1().PrometheusRules(key.Namespace(cluster)).Update(ctx, desiredRule, metav1.UpdateOptions{})

			if err != nil {
				r.logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("could not update rule %s", ruleName))
				return microerror.Mask(err)
			}

			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("updated rule %s", ruleName))
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("rule %s does not need to be updated", ruleName))
		}

	}
	r.logger.LogCtx(ctx, "level", "debug", "message", "ensured the prometheus rules")

	return nil
}

func updateMeta(c, d *v1.PrometheusRule) {
	d.SetGenerateName(c.GetGenerateName())
	d.SetUID(c.GetUID())
	d.SetResourceVersion(c.GetResourceVersion())
	d.SetGeneration(c.GetGeneration())
	d.SetSelfLink(c.GetSelfLink())
	labels := c.GetLabels()
	for k, v := range d.GetLabels() {
		labels[k] = v
	}
	d.SetLabels(labels)
	annotations := c.GetAnnotations()
	for k, v := range d.GetAnnotations() {
		annotations[k] = v
	}
	d.SetAnnotations(annotations)
	d.SetCreationTimestamp(c.GetCreationTimestamp())
	d.SetDeletionTimestamp(c.GetDeletionTimestamp())
	d.SetDeletionGracePeriodSeconds(c.GetDeletionGracePeriodSeconds())
	d.SetFinalizers(c.GetFinalizers())
	d.SetOwnerReferences(c.GetOwnerReferences())
	d.SetClusterName(c.GetClusterName())
	d.SetManagedFields(c.GetManagedFields())
}
