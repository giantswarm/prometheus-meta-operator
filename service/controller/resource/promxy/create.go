package promxy

import (
	"context"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if promxy configmap needs to be updated")
	configMap, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.PromxyConfigMapNamespace()).Get(ctx, key.PromxyConfigMapName(), metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	config, err := r.readFromConfig(configMap)
	if err != nil {
		return microerror.Mask(err)
	}

	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	serverGroup, err := r.toServerGroup(cluster)
	if err != nil {
		return microerror.Mask(err)
	}

	if !promxyContains(config.PromxyConfig, serverGroup) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "promxy configmap needs to be updated")
		r.logger.LogCtx(ctx, "level", "debug", "message", "adding server group")
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
