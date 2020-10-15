package heartbeatrouting

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if alertmanager configmap needs to be updated")

	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	configMap, cfg, err := r.readConfig()
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
