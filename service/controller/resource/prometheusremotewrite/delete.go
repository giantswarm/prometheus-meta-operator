package prometheusremotewrite

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v8/pkg/controller/context/resourcecanceledcontext"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/v2/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewriteutils"
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
			err = r.unsetRemoteWrite(ctx, remoteWrite, prometheusAndMetadata{
				prometheus: current,
				name:       current.GetName(),
				namespace:  current.GetNamespace(),
			})
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

func (r *Resource) unsetRemoteWrite(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, p prometheusAndMetadata) error {
	// remove remotewrite config from Prometheus once RemoteWrite CR is deleted
	// Check if prometheus pointer is not nil
	// at some cases the cluster is deleted, and prometheus as well.
	if p.prometheus != nil {
		if desired, ok := removePrometheusRemoteWrite(*remoteWrite, *p.prometheus); ok {
			if !ok {
				r.logger.Debugf(ctx, fmt.Sprintf("no update required for Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))
				return nil
			}
			r.logger.Debugf(ctx, fmt.Sprintf("updating Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))
			updateMeta(p.prometheus, desired)
			_, err := r.prometheusClient.MonitoringV1().
				Prometheuses(p.namespace).
				Update(ctx, desired, metav1.UpdateOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}
	// Delete the status ref from remotewrite
	err := r.ensureStatusDeleted(ctx, remoteWrite, corev1.ObjectReference{
		Name:      p.name,
		Namespace: p.namespace})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) ensureStatusDeleted(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, objRef corev1.ObjectReference) error {
	for index, ref := range remoteWrite.Status.ConfiguredPrometheuses {
		if ref.Name == objRef.Name && ref.Namespace == objRef.Namespace {
			remoteWrite.Status.ConfiguredPrometheuses = append(remoteWrite.Status.ConfiguredPrometheuses[:index], remoteWrite.Status.ConfiguredPrometheuses[index+1:]...)
			err := r.k8sClient.CtrlClient().Status().Update(ctx, remoteWrite)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}
