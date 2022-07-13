package remotewritesecret

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/google/go-cmp/cmp"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/v2/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewriteutils"
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
		}

		for _, current := range prometheusList.Items {

			installedSecrets := make([]corev1.Secret, 0)
			// Loop over remote write secrets
			for _, item := range remoteWrite.Spec.Secrets {
				secret, err := r.syncSecret(ctx, remoteWrite, item, current.GetNamespace())
				if err != nil {
					return microerror.Mask(err)
				}
				installedSecrets = append(installedSecrets, *secret)
				r.logger.Debugf(ctx, fmt.Sprintf("Secret %#q in namespace %#q created", secret.Name, secret.Namespace))
			}

			/*
			 * Cleanup deleted secrets from RemoteWrite CR
			 *  Once RemoteWrite.spec.secrets changed
			 */
			err := r.ensureCleanupSecrets(ctx, remoteWrite, current.GetNamespace(), installedSecrets)
			if err != nil {
				return microerror.Mask(err)
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

	r.logger.Debugf(ctx, "ensured prometheus remotewrite secret")

	return nil
}

func (r *Resource) syncSecret(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, item pmov1alpha1.RemoteWriteSecretSpec, ns string) (*corev1.Secret, error) {
	desired := r.ensureRemoteWriteSecret(item, remoteWrite.ObjectMeta, ns)
	secret, err := r.k8sClient.K8sClient().CoreV1().Secrets(ns).Get(ctx, item.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		r.logger.Debugf(ctx, fmt.Sprintf("creating Secret %#q in namespace %#q", desired.Name, desired.Namespace))
		secret, err = r.k8sClient.K8sClient().CoreV1().Secrets(ns).Create(ctx, &desired, metav1.CreateOptions{})
	} else if err == nil && !cmp.Equal(desired.Data, secret.Data) {
		r.logger.Debugf(ctx, fmt.Sprintf("updating Secret %#q in namespace %#q", desired.Name, desired.Namespace))
		secret, err = r.k8sClient.K8sClient().CoreV1().Secrets(ns).Update(ctx, &desired, metav1.UpdateOptions{})
	}
	if err != nil {
		return nil, microerror.Mask(err)
	}
	err = r.ensureStatusCreated(ctx, remoteWrite, secret)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return secret, err
}

func (r *Resource) retrieveSecrets(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, ns string) (*corev1.SecretList, error) {
	l := labels.SelectorFromSet(labels.Set(map[string]string{
		labelName:      remoteWrite.GetName(),
		labelNamespace: remoteWrite.GetNamespace(),
	}))
	options := metav1.ListOptions{
		LabelSelector: l.String(),
	}
	secrets, err := r.k8sClient.K8sClient().CoreV1().Secrets(ns).List(ctx, options)

	return secrets, err
}

func (r *Resource) ensureStatusCreated(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, sc *corev1.Secret) error {
	for _, ref := range remoteWrite.Status.SyncedSecrets {
		if ref.Name == sc.GetName() && ref.Namespace == sc.GetNamespace() {
			return nil
		}
	}

	newStatus := corev1.ObjectReference{
		Name:      sc.GetName(),
		Namespace: sc.GetNamespace(),
	}
	remoteWrite.Status.SyncedSecrets = append(remoteWrite.Status.SyncedSecrets, newStatus)

	err := r.k8sClient.CtrlClient().Status().Update(ctx, remoteWrite)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *Resource) ensureCleanupSecrets(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, ns string, installedSecrets []corev1.Secret) error {
	secrets, err := r.retrieveSecrets(ctx, remoteWrite, ns)
	if err != nil {
		return microerror.Mask(err)
	}
	for _, secret := range secrets.Items {
		// delete secret if it doesn't exist in the remotewrite secrets field
		if !secretInstalled(secret, installedSecrets) {
			err = r.deleteSecret(ctx, remoteWrite, corev1.ObjectReference{
				Namespace: secret.GetNamespace(),
				Name:      secret.GetName(),
			})
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}

func (r *Resource) ensureCleanUp(ctx context.Context, remoteWrite *pmov1alpha1.RemoteWrite, prometheuses []*promv1.Prometheus) error {
	// Copy the statuses, because it will be overwritten later on.
	statuses := remoteWrite.Status.SyncedSecrets

	for _, statusRef := range statuses {
		if !inList(statusRef, prometheuses) {
			err := r.deleteSecret(ctx, remoteWrite, statusRef)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}

func inList(o corev1.ObjectReference, list []*promv1.Prometheus) bool {
	for _, item := range list {
		if o.Namespace == item.GetNamespace() {
			return true
		}
	}

	return false
}
