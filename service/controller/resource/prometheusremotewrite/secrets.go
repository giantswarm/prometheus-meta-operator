package prometheusremotewrite

import (
	"context"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/api/v1alpha1"
)

func (r *Resource) copyBasicAuthSecret(ctx context.Context, rw *v1alpha1.RemoteWrite, ns string) error {

	if rw.Spec.RemoteWrite.BasicAuth != nil {

		if rw.Spec.RemoteWrite.BasicAuth.Username.Name == rw.Spec.RemoteWrite.BasicAuth.Password.Name {

			sc, err := r.fetchSecret(ctx, rw.Spec.RemoteWrite.BasicAuth.Username.Name, rw.GetNamespace())
			if err != nil {
				return microerror.Mask(err)
			}

			// creating the new secret in the prometheus namespace
			_, err = r.createSecret(ctx, sc, ns)
			if err != nil {
				return microerror.Mask(err)
			}

		} else {

			scUsername, err := r.fetchSecret(ctx, rw.Spec.RemoteWrite.BasicAuth.Username.Name, rw.GetNamespace())
			if err != nil {
				return microerror.Mask(err)
			}
			_, err = r.createSecret(ctx, scUsername, ns)
			if err != nil {
				return microerror.Mask(err)
			}

			scPassword, err := r.fetchSecret(ctx, rw.Spec.RemoteWrite.BasicAuth.Password.Name, rw.GetNamespace())
			if err != nil {
				return microerror.Mask(err)
			}
			_, err = r.createSecret(ctx, scPassword, ns)
			if err != nil {
				return microerror.Mask(err)
			}

		}

	}
	return nil
}

func (r *Resource) fetchSecret(ctx context.Context, name, namespace string) (*corev1.Secret, error) {

	sc, err := r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Maskf(errorRetrievingSecret, "Could not fetch Secret with name %#s in namespace %#s", name, namespace)
	}
	return sc, nil
}

func (r *Resource) createSecret(ctx context.Context, sc *corev1.Secret, namespace string) (*corev1.Secret, error) {
	sc.SetNamespace(namespace)

	newSc, err := r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Create(ctx, sc, metav1.CreateOptions{})
	if err != nil {
		return nil, microerror.Maskf(errorCreatingSecret, "Could not create Secret with name %#s in namespace %#s", newSc.GetName(), namespace)
	}
	return newSc, nil
}
