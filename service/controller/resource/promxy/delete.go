package promxy

import (
	"context"
	"net/url"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if promxy configmap needs to be updated")
	configMap, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.PromxyConfigMapNamespace()).Get(ctx, key.PromxyConfigMapName(), metav1.GetOptions{})

	if err != nil {
		return microerror.Mask(err)
	}

	config, err := r.readFromConfig(configMap)
	if err != nil {
		return microerror.Mask(err)
	}

	apiServerHost := r.k8sClient.RESTConfig().Host
	apiServerURL, err := url.Parse(apiServerHost)
	if err != nil {
		return microerror.Mask(err)
	}

	serverGroup, err := toServerGroup(obj, apiServerURL, r.installation, r.provider)
	if err != nil {
		return microerror.Mask(err)
	}

	if config.Promxy.contains(serverGroup) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "promxy configmap needs to be updated")
		r.logger.LogCtx(ctx, "level", "debug", "message", "removing server group")
		config.Promxy.remove(serverGroup)

		err = r.updateConfig(ctx, configMap, config)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "removed server group")
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "promxy configmap does not need to be updated")
	}
	return nil
}
