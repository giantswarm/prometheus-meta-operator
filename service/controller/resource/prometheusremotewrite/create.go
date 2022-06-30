package prometheusremotewrite

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/resourcecanceledcontext"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

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

			data := fmt.Sprint(`[{ "op": "add", "path": "/status/configuredPrometheuses/-", "value": { "name": "%s", "namespace": "%s" } }]`, current.GetName(), current.GetNamespace())
			patch := patch{
				patchType: types.JSONPatchType,
				data:      []byte(data),
			}
			err = r.k8sClient.CtrlClient().Status().Patch(ctx, remoteWrite, patch)
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

type patch struct {
	patchType types.PatchType
	data      []byte
}

func (p patch) Type() types.PatchType {
	return p.patchType
}

func (p patch) Data(obj runtime.Object) ([]byte, error) {
	return p.data, nil
}
