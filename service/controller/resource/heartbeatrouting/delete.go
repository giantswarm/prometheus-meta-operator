package heartbeatrouting

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if alertmanager configmap needs to be updated")

	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	configMap, cfg, err := r.readConfig(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	receiver, err := toReceiver(cluster, r.installation, r.opsgenieKey)
	if err != nil {
		return microerror.Mask(err)
	}
	cfg, receiverUpdate := removeReceiver(cfg, receiver)

	route, err := toRoute(cluster, r.installation)
	if err != nil {
		return microerror.Mask(err)
	}
	cfg, routeUpdate := removeRoute(cfg, route)

	if receiverUpdate || routeUpdate {
		r.logger.LogCtx(ctx, "level", "debug", "message", "alertmanager configmap needs to be updated")
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
