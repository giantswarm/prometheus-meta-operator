package servergroup

import (
	"context"
	"fmt"
	"net/url"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/promxy"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	configmapName := key.PromxyConfigMapName()
	r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("checking if %s configmap already exists ", configmapName))
	configmap, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.PromxyConfigMapNamespace()).Get(ctx, configmapName, metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("configmap %s does not exists", configmapName))
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}
	promxyConfiguration, err := promxy.Deserialize(configmap.Data["values"])
	if err != nil {
		return microerror.Mask(err)
	}

	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	url, err := url.Parse("https://kubernetes.default:443")
	if err != nil {
		return microerror.Mask(err)
	}
	serverGroup := r.toServerGroup(url, cluster)

	if promxyConfiguration.Promxy.Contains(serverGroup) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "removing server group")
		promxyConfiguration.Promxy.Remove(serverGroup)
		content, err := promxy.Serialize(promxyConfiguration)
		if err != nil {
			return microerror.Mask(err)
		}

		configmap.Data["values"] = content
		_, err = r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.PromxyConfigMapNamespace()).Update(ctx, configmap, metav1.UpdateOptions{})

		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.LogCtx(ctx, "level", "debug", "message", "removed server group")
	}
	return nil
}
