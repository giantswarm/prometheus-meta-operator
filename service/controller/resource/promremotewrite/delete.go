package promremotewrite

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/resourcecanceledcontext"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "deleting")
	{
		//get remotewrite
		remoteWrite, err := ToRemoteWrite(obj)
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.Debugf(ctx, "remotewrite obj,", remoteWrite.Spec.ClusterSelector)

		labelMap, err := metav1.LabelSelectorAsMap(&remoteWrite.Spec.ClusterSelector)
		if err != nil {
			return microerror.Mask(err)
		}

		// fetch current prometheus
		prometheusList, err := r.prometheusClient.
			MonitoringV1().
			Prometheuses(metav1.NamespaceAll).
			List(ctx, metav1.ListOptions{LabelSelector: labels.SelectorFromSet(labelMap).String()})
		if err != nil {
			return microerror.Maskf(errorFetchingPrometheus, "Could not fetch Prometheus with label selector '%#q'", remoteWrite.Spec.ClusterSelector.String())
		}
		if prometheusList == nil && len(prometheusList.Items) == 0 {
			r.logger.Debugf(ctx, "no prometheus found, cancel reconciliation")
			resourcecanceledcontext.SetCanceled(ctx)
			return nil
		}

		// loop
		for _, current := range prometheusList.Items {

			// omit remotewrite from Prometheus once RemoteWrite CR is deleted
			if desired, ok := omitPrometheusRemoteWrite(*remoteWrite, *current); ok {
				r.logger.Debugf(ctx, fmt.Sprintf("updating Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))
				updateMeta(current, desired)
				_, err = r.prometheusClient.MonitoringV1().
					Prometheuses(current.GetNamespace()).
					Update(ctx, desired, metav1.UpdateOptions{})
				if err != nil {
					return microerror.Mask(err)
				}
			} else {
				r.logger.Debugf(ctx, fmt.Sprintf("no update required for Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))
			}

		}

	}
	r.logger.Debugf(ctx, "deleted")

	return nil
}
