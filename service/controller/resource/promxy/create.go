package promxy

import (
	"context"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	configMap, err := r.getConfigMap(ctx, obj)
	if err != nil {
		return microerror.Mask(err)
	} else if configMap == nil {
		return nil // Missing config map, we return immediately
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

	if !config.Promxy.contains(serverGroup) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "adding server group")
		config.Promxy.add(serverGroup)

		err = r.updateConfig(ctx, configMap, config)
		if err != nil {
			return microerror.Mask(err)
		}

		r.logger.LogCtx(ctx, "level", "debug", "message", "added server group")
	}
	return nil
}
