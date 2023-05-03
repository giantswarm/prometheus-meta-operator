package remotewriteapiendpointconfigsecret

import (
	"context"
	"reflect"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "ensuring prometheus remote write api endpoint secret")
	{

		cluster, err := key.ToCluster(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		name, namespace := key.RemoteWriteAPIEndpointConfigSecretNameAndNamespace(cluster, r.Installation, r.Provider)

		// Get the current secret if it exists.
		current, err := r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			err = r.createSecret(ctx, cluster, name, namespace, r.Version)
			if err != nil {
				return microerror.Mask(err)
			}
		} else if err != nil {
			return microerror.Mask(err)
		}

		if current != nil {
			// As it takes a long time to apply the new password to the agent due to a built-in delay in the app-platform,
			// we keep the already generated remote write password.
			password, err := readRemoteWritePasswordFromSecret(*current)
			if err != nil {
				return microerror.Mask(err)
			}

			desired, err := r.desiredSecret(cluster, name, namespace, password, r.Version)
			if err != nil {
				return microerror.Mask(err)
			}
			if !reflect.DeepEqual(current.Data, desired.Data) {
				updateMeta(current, desired)
				_, err := r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Update(ctx, desired, metav1.UpdateOptions{})
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	r.logger.Debugf(ctx, "ensured prometheus remote write api endpoint secret")

	return nil
}

func readRemoteWritePasswordFromSecret(secret corev1.Secret) (string, error) {
	secretValues := GlobalRemoteWriteValues{}
	err := yaml.Unmarshal(secret.Data["values"], &secretValues)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return secretValues.Global.RemoteWrite[0].Password, nil
}

func updateMeta(c, d metav1.Object) {
	d.SetGenerateName(c.GetGenerateName())
	d.SetUID(c.GetUID())
	d.SetResourceVersion(c.GetResourceVersion())
	d.SetGeneration(c.GetGeneration())
	d.SetSelfLink(c.GetSelfLink())
	d.SetCreationTimestamp(c.GetCreationTimestamp())
	d.SetDeletionTimestamp(c.GetDeletionTimestamp())
	d.SetDeletionGracePeriodSeconds(c.GetDeletionGracePeriodSeconds())
	d.SetLabels(c.GetLabels())
	d.SetAnnotations(c.GetAnnotations())
	d.SetFinalizers(c.GetFinalizers())
	d.SetOwnerReferences(c.GetOwnerReferences())
	d.SetManagedFields(c.GetManagedFields())
}
