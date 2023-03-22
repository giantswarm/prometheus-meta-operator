package remotewriteapiendpointconfigsecret

import (
	"context"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "deleting prometheus remote write api endpoint secret")
	{
		cluster, err := key.ToCluster(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		name, namespace := key.RemoteWriteAPIEndpointConfigSecretNameAndNamespace(cluster, r.Installation, r.Provider)

		current, err := r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		err = r.deleteSecret(ctx, current)
		if err != nil {
			return microerror.Mask(err)
		}

		// Delete duplicate secret until the new observability bundle is upgraded everywhere
		if !key.IsCAPIManagementCluster(r.Provider) {
			duplicateName := getConfigMapCopyName(r.Installation, cluster, name)
			duplicate, err := r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Get(ctx, duplicateName, metav1.GetOptions{})
			if err != nil {
				return microerror.Mask(err)
			}

			err = r.deleteSecret(ctx, duplicate)
			if err != nil {
				return microerror.Mask(err)
			}
		}

	}
	r.logger.Debugf(ctx, "deleted prometheus remote write api endpoint secret")

	return nil
}
