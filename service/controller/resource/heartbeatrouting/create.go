package heartbeatrouting

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/heartbeatrouting/receiver"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/heartbeatrouting/route"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if alertmanager configmap needs to be updated")

	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	configMap, cfg, err := r.readConfig(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	cfg, receiverUpdated, err := receiver.EnsureCreated(cfg, cluster, r.installation, r.opsgenieKey)
	if err != nil {
		return microerror.Mask(err)
	}

	cfg, routeUpdated, err := route.EnsureCreated(cfg, cluster, r.installation)
	if err != nil {
		return microerror.Mask(err)
	}

	if receiverUpdated || routeUpdated {
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
