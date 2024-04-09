package remotewriteconfig

import (
	"context"
	"reflect"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	remotewriteconfiguration "github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewrite/configuration"
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
		namespace := key.GetClusterAppsNamespace(cluster, r.installation, r.provider)

		// Get the current configmap if it exists.
		current, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			err = r.createConfigMap(ctx, cluster, name, namespace)
			if err != nil {
				return microerror.Mask(err)
			}
		} else if err != nil {
			return microerror.Mask(err)
		}

		if current != nil {
			currentShards, err := readCurrentShardsFromConfig(*current)
			if err != nil {
				return microerror.Mask(err)
			}

			shards, err := r.getShardsCountForCluster(ctx, cluster, currentShards)
			if err != nil {
				return microerror.Mask(err)
			}

			desired, err := r.desiredConfigMap(ctx, cluster, name, namespace, shards)
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

func readCurrentShardsFromConfig(configMap corev1.ConfigMap) (int, error) {
	remoteWriteConfig := remotewriteconfiguration.RemoteWriteConfig{}
	err := yaml.Unmarshal([]byte(configMap.Data["values"]), &remoteWriteConfig)
	if err != nil {
		return 0, microerror.Mask(err)
	}

	return remoteWriteConfig.PrometheusAgentConfig.Shards, nil
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
