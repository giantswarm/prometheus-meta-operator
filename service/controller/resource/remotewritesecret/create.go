package remotewritesecret

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/resourcecanceledcontext"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "ensuring prometheus remotewrite secret")
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
		if prometheusList == nil || len(prometheusList.Items) == 0 {
			r.logger.Debugf(ctx, "no prometheus found, cancel reconciliation")
			resourcecanceledcontext.SetCanceled(ctx)
			return nil
		}

		for _, current := range prometheusList.Items {

			// Loop over remote write secrets
			for _, item := range remoteWrite.Spec.Secrets {
				desired := r.ensureRemoteWriteSecret(item, remoteWrite.ObjectMeta, current.GetNamespace())
				secret, err := r.k8sClient.K8sClient().CoreV1().Secrets(current.GetNamespace()).Get(ctx, item.Name, metav1.GetOptions{})
				if apierrors.IsNotFound(err) {
					r.logger.Debugf(ctx, fmt.Sprintf("creating Secret %#q in namespace %#q", desired.Name, desired.Namespace))

					//if err := r.k8sClient.K8sClient().SetControllerReference(p, &desiredSa, r.Scheme); err != nil {
					//	return err
					//}
					secret, err = r.k8sClient.K8sClient().CoreV1().Secrets(current.GetNamespace()).Create(ctx, &desired, metav1.CreateOptions{})
				}
				if err != nil {
					return microerror.Mask(err)
				}
				r.logger.Debugf(ctx, fmt.Sprintf("Secret %#q in namespace %#q created", secret.Name, secret.Namespace))

			}
		}

	}

	r.logger.Debugf(ctx, "ensured prometheus remotewrite secret")

	return nil
}
