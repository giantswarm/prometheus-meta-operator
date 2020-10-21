package heartbeatrouting

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/heartbeatrouting/receiver"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/heartbeatrouting/route"
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

	cfg, receiverNeedsUpdate, err := receiver.EnsureDeleted(cfg, cluster, r.installation, r.opsgenieKey)
	if err != nil {
		return microerror.Mask(err)
	}

	cfg, routeNeedsUpdate, err := route.EnsureDeleted(cfg, cluster, r.installation)
	if err != nil {
		return microerror.Mask(err)
	}

	alertManagerConfigMapNeedsUpdate := receiverNeedsUpdate || routeNeedsUpdate

	if alertManagerConfigMapNeedsUpdate {
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
