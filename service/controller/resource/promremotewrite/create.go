package promremotewrite

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/resourcecanceledcontext"
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
		r.logger.Debugf(ctx, "remotewrite obj,", remoteWrite.Spec.ClusterSelector)

		// fetch current prometheus
		prometheusList, err := r.prometheusClient.
			MonitoringV1().
			Prometheuses(metav1.NamespaceAll).
			List(ctx, metav1.ListOptions{LabelSelector: remoteWrite.Spec.ClusterSelector.String()})
		if err != nil {
			return microerror.Maskf(errorFetchingPrometheus, "Could not fetch Prometheus with label selector '%T'", remoteWrite.Spec.ClusterSelector.String())
		}
		if prometheusList == nil && len(prometheusList.Items) == 0 {
			r.logger.Debugf(ctx, "no prometheus found, cancel reconciliation")
			resourcecanceledcontext.SetCanceled(ctx)
			return nil
		}

		// loop
		for _, current := range prometheusList.Items {

			desired, ok := toPrometheusRemoteWrite(*remoteWrite, *current)
			if ok {
				r.logger.Debugf(ctx, fmt.Sprintf("updating Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))
			} else {
				r.logger.Debugf(ctx, fmt.Sprintf("no update required for Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))
				continue
			}

			updateMeta(current, desired)
			_, err = r.prometheusClient.MonitoringV1().
				Prometheuses(current.GetNamespace()).
				Update(ctx, desired, metav1.UpdateOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
		}

	}

	r.logger.Debugf(ctx, "updated")

	return nil
}

func updateMeta(c, d metav1.Object) {
	d.SetGenerateName(c.GetGenerateName())
	d.SetResourceVersion(c.GetResourceVersion())
	d.SetGeneration(c.GetGeneration())
	d.SetSelfLink(c.GetSelfLink())
	d.SetCreationTimestamp(c.GetCreationTimestamp())
	d.SetDeletionTimestamp(c.GetDeletionTimestamp())
	d.SetDeletionGracePeriodSeconds(c.GetDeletionGracePeriodSeconds())
	d.SetManagedFields(c.GetManagedFields())
}
