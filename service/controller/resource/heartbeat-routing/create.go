package heartbeatrouting

import (
	"context"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
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

	exists := contains(cfg, receiver, route)
	if !exists || hasChanged(cfg, receiver, route) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "alertmanager configmap needs to be updated")
		if exists {
			cfg = remove(cfg, receiver, route)
		}

		cfg = add(cfg, receiver, route)

		err = r.updateConfig(ctx, configMap, cfg)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "alertmanager configmap updated")
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "alertmanager configmap does not need to be updated")
	}

	return nil
}
