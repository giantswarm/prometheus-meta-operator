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

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	configMapName := key.PromxyConfigMapName()
	r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("checking if %s configmap already exists ", configMapName))
	configmap, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.PromxyConfigMapNamespace()).Get(ctx, configMapName, metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("configmap %s does not exists", configMapName))
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	promxyConfiguration, err := promxy.Deserialize(configmap.Data["values.promxy"])
	if err != nil {
		return microerror.Mask(err)
	}

	url, err := url.Parse("https://kubernetes.default:443")
	if err != nil {
		return microerror.Mask(err)
	}

	serverGroup := r.toServerGroup(url, cluster)
	if !promxyConfiguration.Promxy.Contains(serverGroup) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "adding server group")
		promxyConfiguration.Promxy.Add(serverGroup)
		content, err := promxy.Serialize(promxyConfiguration)
		if err != nil {
			return microerror.Mask(err)
		}

		configmap.Data["values.promxy"] = content
		_, err = r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.PromxyConfigMapNamespace()).Update(ctx, configmap, metav1.UpdateOptions{})

		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.LogCtx(ctx, "level", "debug", "message", "added server group")
	}
	return nil
}
