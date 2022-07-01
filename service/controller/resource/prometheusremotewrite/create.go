package prometheusremotewrite

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/resourcecanceledcontext"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	"github.com/giantswarm/prometheus-meta-operator/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/pkg/remotewriteutils"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "ensuring prometheus remoteWrite config")
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
		if prometheusList == nil || len(prometheusList.Items) == 0 {
			r.logger.Debugf(ctx, "no prometheus found, cancel reconciliation")
			resourcecanceledcontext.SetCanceled(ctx)
			return nil
		}

		for _, current := range prometheusList.Items {

			desired, ok := r.ensurePrometheusRemoteWrite(*remoteWrite, *current)
			if !ok {
				r.logger.Debugf(ctx, fmt.Sprintf("no update required for Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))
				continue
			}

			r.logger.Debugf(ctx, fmt.Sprintf("updating Prometheus CR %#q in namespace %#q", desired.Name, desired.Namespace))

			updateMeta(current, desired)
			_, err = r.prometheusClient.MonitoringV1().
				Prometheuses(current.GetNamespace()).
				Update(ctx, desired, metav1.UpdateOptions{})
			if err != nil {
				return microerror.Mask(err)
			}

			newStatus := v1alpha1.RemoteWriteStatusConfiguredPrometheus{
				Name:      current.GetName(),
				Namespace: current.GetNamespace(),
			}

			gvk, err := apiutil.GVKForObject(remoteWrite, r.k8sClient.Scheme())
			if err != nil {
				return microerror.Mask(err)
			}
			gvk.Kind = remoteWrite.Kind
			remoteWrite.SetGroupVersionKind(gvk)

			r.logger.Debugf(ctx, fmt.Sprintf("remotewrite kind %v", remoteWrite.GetObjectKind().GroupVersionKind()))
			remoteWrite.Status.ConfiguredPrometheuses = append(remoteWrite.Status.ConfiguredPrometheuses, newStatus)
			err = r.k8sClient.CtrlClient().Status().Update(ctx, remoteWrite)
			if err != nil {
				return microerror.Mask(err)
			}
		}

	}

	r.logger.Debugf(ctx, "ensured prometheus remoteWrite config")

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
