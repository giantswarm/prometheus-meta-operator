package heartbeatrouting

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/alertmanager/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/microerror"

	commoncfg "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func toRoute(cluster metav1.Object, installation string) config.Route {
	name := fmt.Sprintf("heartbeat_%s_%s", installation, key.ClusterID(cluster))

	return config.Route{
		Receiver: name,
		Match: map[string]string{
			"cluster":      key.ClusterID(cluster),
			"installation": installation,
			"type":         "heartbeat",
		},
		Continue:       false,
		GroupInterval:  &model.Duration(1 * time.Second),
		GroupWait:      &model.Duration(1 * time.Second),
		RepeatInterval: &model.Duration(15 * time.Second),
	}
}

func toReceiver(cluster metav1.Object, installation string, opsgenieKey string) config.Receiver {
	name := fmt.Sprintf("heartbeat_%s_%s", installation, key.ClusterID(cluster))

	return config.Receiver{
		Name: name,
		WebhookConfigs: []*config.WebhookConfigs{
			&config.WebhookConfigs{
				URL: nil,
				HTTPConfig: &commoncfg.HTTPClientConfig{
					BasicAuth: &commoncfg.BasicAuth{
						Password: commoncfg.Secret(opsgenieKey),
					},
				},
				NotifierConfig: config.NotifierConfig{
					VSendResolved: false,
				},
			},
		},
	}
}

func contains(cfg config.Config, receiver config.Receiver, route config.Route) bool {
	// TODO: implement me.
	return false
}

func hasChanged(cfg config.Config, receiver config.Receiver, route config.Route) bool {
	// TODO: implement me.
	return false
}

func add(cfg config.Config, receiver config.Receiver, route config.Route) config.Config {
	// TODO: implement me.
	return config.Config{}
}

func remove(cfg config.Config, receiver config.Receiver, route config.Route) config.Config {
	// TODO: implement me.
	return config.Config{}
}

func (r *Resource) readConfig() (*v1.ConfigMap, config.Config, error) {
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

	return configMap, cfg, nil
}

func (r *Resource) updateConfig(ctx context.Context, configMap *v1.ConfigMap, cfg config.Config) error {
	content, err := config.String()
	if err != nil {
		return microerror.Mask(err)
	}

	configMap.Data[key.AlertmanagerConfigMapKey()] = content
	_, err = r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.AlertmanagerConfigMapNamespace()).Update(ctx, configMap, metav1.UpdateOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
