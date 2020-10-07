package promxy

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if promxy config map already exists")
	configMap, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.PromxyConfigMapNamespace()).Get(ctx, key.PromxyConfigMapName(), metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "promxy needs to be updated")
		return nil
	} else if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "reading promxy config")
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

	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if server group must be added")
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
