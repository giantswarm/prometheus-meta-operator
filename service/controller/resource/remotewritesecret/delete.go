package remotewritesecret

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/resourcecanceledcontext"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/giantswarm/prometheus-meta-operator/pkg/remotewriteutils"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "deleting prometheus remoteWrite secrets")
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
		if len(prometheusList.Items) == 0 {
			r.logger.Debugf(ctx, "no prometheus found, cancel reconciliation")
			resourcecanceledcontext.SetCanceled(ctx)
			return nil
		}

		for _, current := range prometheusList.Items {

			/*
			  Cleanup deleted secrets from RemoteWrite CR
			*/
			l := labels.SelectorFromSet(labels.Set(map[string]string{
				label:          Name,
				labelName:      remoteWrite.GetName(),
				labelNamespace: remoteWrite.GetNamespace(),
			}))
			listOptions := metav1.ListOptions{
				LabelSelector: l.String(),
			}
			err := r.k8sClient.K8sClient().CoreV1().Secrets(current.GetNamespace()).DeleteCollection(ctx, metav1.DeleteOptions{}, listOptions)
			if err != nil {
				return microerror.Mask(err)
			}

		}

	}
	r.logger.Debugf(ctx, "deleted prometheus remoteWrite secrets")

	return nil
}
