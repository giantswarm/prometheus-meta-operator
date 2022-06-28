package remotewritesecret

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/resourcecanceledcontext"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "deleting prometheus remoteWrite secret")
	{
		remoteWrite, err := ToRemoteWrite(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		// fetch current prometheus using the selector provided in remoteWrite resource.
		prometheusList, err := fetchPrometheusList(ctx, r, remoteWrite)
		if err != nil {
			return microerror.Mask(err)
		}
		if prometheusList == nil && len(prometheusList.Items) == 0 {
			r.logger.Debugf(ctx, "no prometheus found, cancel reconciliation")
			resourcecanceledcontext.SetCanceled(ctx)
			return nil
		}

		for _, current := range prometheusList.Items {

			// Loop over remote write secrets
			for _, item := range remoteWrite.Spec.Secrets {
				err := r.k8sClient.K8sClient().CoreV1().Secrets(current.GetNamespace()).Delete(ctx, item.Name, metav1.DeleteOptions{})
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}

	}
	r.logger.Debugf(ctx, "deleted prometheus remoteWrite secret")

	return nil
}
