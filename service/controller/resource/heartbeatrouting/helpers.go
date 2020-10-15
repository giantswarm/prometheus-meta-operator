package heartbeatrouting

import (
	"context"

	"github.com/prometheus/alertmanager/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) readConfig(ctx context.Context) (*v1.ConfigMap, config.Config, error) {
	configMap, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.AlertmanagerConfigMapNamespace()).Get(ctx, key.AlertmanagerConfigMapName(), metav1.GetOptions{})
	if err != nil {
		return nil, config.Config{}, microerror.Mask(err)
	}

	content, ok := configMap.Data[key.AlertmanagerConfigMapKey()]
	if !ok {
		return nil, config.Config{}, microerror.Mask(invalidConfigError)
	}

	cfg, err := config.Load(content)
	if err != nil {
		return nil, config.Config{}, microerror.Mask(err)
	}

	return configMap, *cfg, nil
}

func (r *Resource) updateConfig(ctx context.Context, configMap *v1.ConfigMap, cfg config.Config) error {
	configMap.Data[key.AlertmanagerConfigMapKey()] = cfg.String()
	_, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.AlertmanagerConfigMapNamespace()).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
