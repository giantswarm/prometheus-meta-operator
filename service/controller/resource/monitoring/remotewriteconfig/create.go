package remotewriteconfig

import (
	"context"
	"reflect"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "ensuring prometheus remote write config")
	{

		cluster, err := key.ToCluster(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		name := key.RemoteWriteConfigName(cluster)
		namespace := key.GetClusterAppsNamespace(cluster, r.Installation, r.Provider)

		// Get the current configmap if it exists.
		current, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			err = r.createConfigMap(ctx, cluster, name, namespace, r.Version)
			if err != nil {
				return microerror.Mask(err)
			}
		} else if err != nil {
			return microerror.Mask(err)
		}

		if current != nil {
			desired, err := r.desiredConfigMap(cluster, name, namespace, r.Version)
			if err != nil {
				return microerror.Mask(err)
			}
			if !reflect.DeepEqual(current.Data, desired.Data) {
				updateMeta(current, desired)
				_, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(namespace).Update(ctx, desired, metav1.UpdateOptions{})
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}

	r.logger.Debugf(ctx, "ensured prometheus remote write config")

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
