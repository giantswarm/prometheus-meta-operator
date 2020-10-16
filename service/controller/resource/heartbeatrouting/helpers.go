package heartbeatrouting

import (
	"context"

	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	alertmanagerconfig "github.com/giantswarm/prometheus-meta-operator/pkg/alertmanager/config"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) readConfig(ctx context.Context) (*v1.ConfigMap, alertmanagerconfig.Config, error) {
	configMap, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.AlertmanagerConfigMapNamespace()).Get(ctx, key.AlertmanagerConfigMapName(), metav1.GetOptions{})
	if err != nil {
		return nil, alertmanagerconfig.Config{}, microerror.Mask(err)
	}

	content, ok := configMap.Data[key.AlertmanagerConfigMapKey()]
	if !ok {
		return nil, alertmanagerconfig.Config{}, microerror.Mask(invalidConfigError)
	}

	cfg, err := alertmanagerconfig.Load(content)
	if err != nil {
		return nil, alertmanagerconfig.Config{}, microerror.Mask(err)
	}

	return configMap, *cfg, nil
}

func (r *Resource) updateConfig(ctx context.Context, configMap *v1.ConfigMap, cfg alertmanagerconfig.Config) error {
	configMap.Data[key.AlertmanagerConfigMapKey()] = cfg.String()
	_, err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.AlertmanagerConfigMapNamespace()).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
