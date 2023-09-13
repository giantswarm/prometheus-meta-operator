package remotewriteconfig

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "deleting prometheus remote write config")
	{
		cluster, err := key.ToCluster(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		name := key.RemoteWriteConfigName(cluster)
		namespace := key.GetClusterAppsNamespace(cluster, r.installation, r.provider)

		_, err = r.k8sClient.K8sClient().CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			// We ignore cases where the configmap is not found (it it was manually deleted for instance)
			return nil
		} else if err != nil {
			return microerror.Mask(err)
		}

		// Delete the finalizer
		patch := []byte(`{"metadata":{"finalizers":null}}`)
		current, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(namespace).Patch(ctx, name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		err = r.deleteConfigMap(ctx, current)
		if err != nil {
			return microerror.Mask(err)
		}

	}
	r.logger.Debugf(ctx, "deleted prometheus remote write config")

	return nil
}
