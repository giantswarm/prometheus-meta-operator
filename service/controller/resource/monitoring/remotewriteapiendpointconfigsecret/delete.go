package remotewriteapiendpointconfigsecret

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "deleting prometheus remote write api endpoint secret")
	{
		cluster, err := key.ToCluster(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		name := key.RemoteWriteAPIEndpointConfigSecretName(cluster, r.Provider)
		namespace := key.GetClusterAppsNamespace(cluster, r.Installation, r.Provider)

		_, err = r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			// We ignore cases where the secret is not found (it it was manually deleted for instance)
			return nil
		} else if err != nil {
			return microerror.Mask(err)
		}

		// Delete the finalizer
		patch := []byte(`{"metadata":{"finalizers":null}}`)
		current, err := r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Patch(ctx, name, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		err = r.deleteSecret(ctx, current)
		if err != nil {
			return microerror.Mask(err)
		}

	}
	r.logger.Debugf(ctx, "deleted prometheus remote write api endpoint secret")

	return nil
}
