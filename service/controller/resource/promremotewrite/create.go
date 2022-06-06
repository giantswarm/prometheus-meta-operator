package promremotewrite

import (
	"context"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "reconcile remotewrite")
	{
		// get remotewrite
		remoteWrite, err := ToRemoteWrite(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		// fetch current prometheus
		prometheusList, err := r.prometheusClient.
			MonitoringV1().
			Prometheuses(metav1.NamespaceAll).
			List(ctx, metav1.ListOptions{LabelSelector: remoteWrite.Spec.ClusterSelector.String()})
		if err != nil {
			return microerror.Maskf(errorFetchingPrometheus, "Could not fetch Prometheus with label selector '%T'", remoteWrite.Spec.ClusterSelector.String())
		}
		if prometheusList == nil && len(prometheusList.Items) == 0 {
			return microerror.Maskf(noSuchPrometheusForLabel, "No Such Prometheus found with Label '%T'", remoteWrite.Spec.ClusterSelector.String())
		}
		currentPrometheus := prometheusList.Items[0]

		desired, err := toPrometheusRemoteWrite(remoteWrite, *currentPrometheus)
		if err != nil {
			return microerror.Mask(err)
		}

		updateMeta(currentPrometheus, desired)
		_, err = r.prometheusClient.MonitoringV1().
			Prometheuses(currentPrometheus.GetNamespace()).
			Update(ctx, desired, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

	}

	r.logger.Debugf(ctx, "updated")

	return nil
}

func updateMeta(c, d metav1.Object) {
	d.SetGenerateName(c.GetGenerateName())
	d.SetUID(c.GetUID())
	d.SetResourceVersion(c.GetResourceVersion())
	d.SetGeneration(c.GetGeneration())
	d.SetSelfLink(c.GetSelfLink())
	d.SetCreationTimestamp(c.GetCreationTimestamp())
	d.SetDeletionTimestamp(c.GetDeletionTimestamp())
	d.SetDeletionGracePeriodSeconds(c.GetDeletionGracePeriodSeconds())
	d.SetManagedFields(c.GetManagedFields())
}
