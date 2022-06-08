package promremotewrite

import (
	"context"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "deleting")
	{ //// get remotewrite
		//remoteWrite, err := ToRemoteWrite(obj)
		//if err != nil {
		//	return microerror.Mask(err)
		//}
		//
		//// fetch current prometheus
		//prometheusList, err := r.prometheusClient.
		//	MonitoringV1().
		//	Prometheuses(metav1.NamespaceAll).
		//	List(ctx, metav1.ListOptions{LabelSelector: remoteWrite.Spec.ClusterSelector.String()})
		//if err != nil {
		//	return microerror.Maskf(errorFetchingPrometheus, "Could not fetch Prometheus with label selector '%T'", remoteWrite.Spec.ClusterSelector.String())
		//}
		//if prometheusList == nil && len(prometheusList.Items) == 0 {
		//	return microerror.Maskf(noSuchPrometheusForLabel, "No Such Prometheus found with Label '%T'", remoteWrite.Spec.ClusterSelector.String())
		//}
		//currentPrometheus := prometheusList.Items[0]
		//
		//// omit remotewrite from Prometheus once RemoteWrite CR is deleted
		//desired := currentPrometheus
		//desired.Spec.RemoteWrite = nil
		//updateMeta(currentPrometheus, desired)
		//_, err = r.prometheusClient.MonitoringV1().
		//	Prometheuses(currentPrometheus.GetNamespace()).
		//	Update(ctx, desired, metav1.UpdateOptions{})
		//if err != nil {
		//	return microerror.Mask(err)
		//}

	}
	r.logger.Debugf(ctx, "deleted")

	return nil
}
