package remotewriteapiendpointconfigsecret

import (
	"context"
	"reflect"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	remotewriteconfiguration "github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewrite/configuration"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

// /////////////////////////////////////////////////////////////
// TODO: Remove this resource when all WC are migrated to V19
// /////////////////////////////////////////////////////////////
func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	if r.mimirEnabled {
		r.logger.Debugf(ctx, "mimir is enabled, deleting")
		return r.EnsureDeleted(ctx, obj)
	}

	r.logger.Debugf(ctx, "ensuring prometheus remote write api endpoint secret")
	{
		cluster, err := key.ToCluster(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		// Get password from remote-write-secret
		r.logger.Debugf(ctx, "looking up for secret remote write secret")
		_, password, err := remotewriteconfiguration.GetUsernameAndPassword(r.k8sClient.K8sClient(), ctx, cluster, r.installation, r.provider)
		if err != nil {
			r.logger.Errorf(ctx, err, "lookup for remote write secret failed")
			return microerror.Mask(err)
		}

		name := key.RemoteWriteAPIEndpointConfigSecretName(cluster, r.provider)
		namespace := key.GetClusterAppsNamespace(cluster, r.installation, r.provider)
		// Get the current secret if it exists.
		current, err := r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			err = r.createSecret(ctx, cluster, name, namespace, password)
			if err != nil {
				return microerror.Mask(err)
			}
		} else if err != nil {
			return microerror.Mask(err)
		}

		if current != nil {
			desired, err := r.desiredSecret(ctx, cluster, name, namespace, password)
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
