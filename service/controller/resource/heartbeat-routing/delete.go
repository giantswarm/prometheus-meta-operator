package heartbeatrouting

import (
	"context"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if alertmanager configmap needs to be updated")

	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	configMap, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.AlertmanagerConfigMapNamespace()).Get(ctx, key.AlertmanagerConfigMapName(), metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	cfg, err := r.readFromConfig(configMap)
	if err != nil {
		return microerror.Mask(err)
	}

	receiver := toReceiver(cluster, r.installation)

	route := toRoute(cluster, r.installation)

	if contains(cfg, receiver, route) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "alertmanager configmap needs to be updated")
		r.logger.LogCtx(ctx, "level", "debug", "message", "removing receiver and route")
		cfg = remove(cfg, receiver, route)

		err = r.updateConfig(ctx, configMap, cfg)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "removed receiver and route")
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "alertmanager configmap does not need to be updated")
	}

	return nil
}
