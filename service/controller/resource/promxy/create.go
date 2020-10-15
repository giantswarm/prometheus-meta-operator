package promxy

import (
	"context"
	"net/url"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

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

	serverGroup, err := toServerGroup(cluster, apiServerURL, r.installation, r.provider)
	if err != nil {
		return microerror.Mask(err)
	}

	containsServerGroup := config.Promxy.contains(serverGroup)
	if !containsServerGroup || config.Promxy.needsUpdate(serverGroup) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "promxy configmap needs to be updated")
		// We remove the server group if it needs to be updated
		if containsServerGroup {
			config.Promxy.remove(serverGroup)
		}
		config.Promxy.add(serverGroup)

		err = r.updateConfig(ctx, configMap, config)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "added server group")
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "promxy configmap does not need to be updated")
	}
	return nil
}
