package prometheusremotewrite

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/resourcecanceledcontext"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/pkg/remotewriteutils"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "deleting prometheus remoteWrite config")
	{
		remoteWrite, err := remotewriteutils.ToRemoteWrite(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		// fetch current prometheus using the selector provided in remoteWrite resource.
		prometheusList, err := remotewriteutils.FetchPrometheusList(ctx, toResourceWrapper(r), remoteWrite)
		if err != nil {
			return microerror.Mask(err)
		}
		if prometheusList == nil && len(prometheusList.Items) == 0 {
			r.logger.Debugf(ctx, "no prometheus found, cancel reconciliation")
			resourcecanceledcontext.SetCanceled(ctx)
			return nil
		}

		for _, current := range prometheusList.Items {
			err = r.unsetRemoteWrite(ctx, remoteWrite, current)
			if err != nil {
				return microerror.Mask(err)
			}
		}

		err = r.ensureCleanUp(ctx, remoteWrite, prometheusList.Items)
		if err != nil {
			return microerror.Mask(err)
		}

	}
	r.logger.Debugf(ctx, "deleted prometheus remoteWrite config")

	return nil
}

func (r *Resource) unsetRemoteWrite(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, prometheus *promv1.Prometheus) error {
	err := r.ensureStatusDeleted(ctx, remoteWrite, prometheus)
	if err != nil {
		return microerror.Mask(err)
	}

	// remove remotewrite config from Prometheus once RemoteWrite CR is deleted
	if desired, ok := removePrometheusRemoteWrite(*remoteWrite, *prometheus); ok {
		if !ok {
			r.logger.Debugf(ctx, fmt.Sprintf("no update required for Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))
			return nil
		}
		r.logger.Debugf(ctx, fmt.Sprintf("updating Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))
		updateMeta(prometheus, desired)
		_, err := r.prometheusClient.MonitoringV1().
			Prometheuses(prometheus.GetNamespace()).
			Update(ctx, desired, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (r *Resource) ensureStatusDeleted(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, prometheus *promv1.Prometheus) error {
	for index, ref := range remoteWrite.Status.ConfiguredPrometheuses {
		if ref.Name == prometheus.GetName() && ref.Namespace == prometheus.GetNamespace() {
			remoteWrite.Status.ConfiguredPrometheuses = append(remoteWrite.Status.ConfiguredPrometheuses[:index], remoteWrite.Status.ConfiguredPrometheuses[index+1:]...)
			err := r.k8sClient.CtrlClient().Status().Update(ctx, remoteWrite)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}
