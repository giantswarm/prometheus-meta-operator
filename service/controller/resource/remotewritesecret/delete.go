package remotewritesecret

import (
	"context"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/giantswarm/prometheus-meta-operator/v2/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewriteutils"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	if r.mimirEnabled {
		r.logger.Debugf(ctx, "mimir is enabled, deleting")
		return r.EnsureDeleted(ctx, obj)
	}

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
			r.logger.Debugf(ctx, "no prometheus found, stop reconciliation")
			return nil
		}

		for _, current := range prometheusList.Items {

			/*
			 * Cleanup deleted secrets from RemoteWrite CR
			 *  Once RemoteWrite.spec.secrets changed
			 */

			l := labels.SelectorFromSet(labels.Set(map[string]string{
				labelName:      remoteWrite.GetName(),
				labelNamespace: remoteWrite.GetNamespace(),
			}))
			options := metav1.ListOptions{
				LabelSelector: l.String(),
			}
			err := r.k8sClient.K8sClient().CoreV1().Secrets(current.GetNamespace()).DeleteCollection(ctx, metav1.DeleteOptions{}, options)
			if err != nil {
				return microerror.Mask(err)
			}

			for _, sRef := range remoteWrite.Status.SyncedSecrets {
				if sRef.Namespace == current.GetNamespace() {
					err = r.ensureStatusDeleted(ctx, remoteWrite, sRef)
					if err != nil {
						return microerror.Mask(err)
					}
				}
			}

		}

		/*
		 * Cleanup secrets from RemoteWrite CR
		 *  Once RemoteWrite.spec.clusterSelector changed
		 */
		err = r.ensureCleanUp(ctx, remoteWrite, prometheusList.Items)
		if err != nil {
			return microerror.Mask(err)
		}

	}
	r.logger.Debugf(ctx, "deleted prometheus remoteWrite secrets")

	return nil
}

func (r *Resource) ensureStatusDeleted(ctx context.Context, remoteWrite *v1alpha1.RemoteWrite, objRef corev1.ObjectReference) error {
	for index, ref := range remoteWrite.Status.SyncedSecrets {
		if ref.Name == objRef.Name && ref.Namespace == objRef.Namespace {
			remoteWrite.Status.SyncedSecrets = append(remoteWrite.Status.SyncedSecrets[:index], remoteWrite.Status.SyncedSecrets[index+1:]...)
			err := r.k8sClient.CtrlClient().Status().Update(ctx, remoteWrite)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}

func (r *Resource) deleteSecret(ctx context.Context, remoteWrite *v1alpha1.RemoteWrite, ref corev1.ObjectReference) error {
	err := r.k8sClient.K8sClient().CoreV1().Secrets(ref.Namespace).Delete(ctx, ref.Name, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return microerror.Mask(err)
	}
	err = r.ensureStatusDeleted(ctx, remoteWrite, ref)
	if err != nil {
		return microerror.Mask(err)
	}

	return err
}
