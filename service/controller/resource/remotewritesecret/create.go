package remotewritesecret

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/controller/context/resourcecanceledcontext"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/giantswarm/prometheus-meta-operator/pkg/remotewriteutils"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "ensuring prometheus remotewrite secret")
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

			installedSecrets := make([]corev1.Secret, 0)
			// Loop over remote write secrets
			for _, item := range remoteWrite.Spec.Secrets {
				desired := r.ensureRemoteWriteSecret(item, remoteWrite.ObjectMeta, current.GetNamespace())
				fmt.Println("Debug:", desired)
				secret, err := r.k8sClient.K8sClient().CoreV1().Secrets(current.GetNamespace()).Get(ctx, item.Name, metav1.GetOptions{})
				if apierrors.IsNotFound(err) {
					r.logger.Debugf(ctx, fmt.Sprintf("creating Secret %#q in namespace %#q", desired.Name, desired.Namespace))
					secret, err = r.k8sClient.K8sClient().CoreV1().Secrets(current.GetNamespace()).Create(ctx, &desired, metav1.CreateOptions{})
				}
				if err != nil {
					return microerror.Mask(err)
				}
				installedSecrets = append(installedSecrets, *secret)
				r.logger.Debugf(ctx, fmt.Sprintf("Secret %#q in namespace %#q created", secret.Name, secret.Namespace))

			}

			/*
			  Cleanup deleted secrets from RemoteWrite CR
			*/
			l := labels.SelectorFromSet(labels.Set(map[string]string{
				label:          Name,
				labelName:      remoteWrite.GetName(),
				labelNamespace: remoteWrite.GetNamespace(),
			}))
			options := metav1.ListOptions{
				LabelSelector: l.String(),
			}
			secrets, err := r.k8sClient.K8sClient().CoreV1().Secrets(current.GetNamespace()).List(ctx, options)
			if err != nil {
				return microerror.Mask(err)
			}
			for _, secret := range secrets.Items {
				// delete secret if it doesn't exist in the remotewrite secrets field
				if !secretInstalled(secret, installedSecrets) {
					err := r.k8sClient.K8sClient().CoreV1().Secrets(secret.GetNamespace()).Delete(ctx, secret.GetName(), metav1.DeleteOptions{})
					if err != nil {
						return microerror.Mask(err)
					}
				}
			}

		}

	}

	r.logger.Debugf(ctx, "ensured prometheus remotewrite secret")

	return nil
}
