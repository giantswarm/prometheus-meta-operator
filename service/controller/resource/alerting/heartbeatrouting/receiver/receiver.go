package receiver

import (
	"net/url"
	"reflect"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/microerror"

	alertmanagerconfig "github.com/giantswarm/prometheus-meta-operator/pkg/alertmanager/config"
	promcommonconfig "github.com/giantswarm/prometheus-meta-operator/pkg/prometheus/common/config"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func toReceiver(cfg alertmanagerconfig.Config, cluster metav1.Object, installation string, opsgenieKey string) (alertmanagerconfig.Receiver, error) {
	u, err := url.Parse(key.HeartbeatAPI(cluster, installation))
	if err != nil {
		return alertmanagerconfig.Receiver{}, microerror.Mask(err)
	}

	httpConfig := &promcommonconfig.HTTPClientConfig{}
	if cfg.Global != nil && cfg.Global.HTTPConfig != nil {
		httpConfigCp, err := cloneHttpConfig(cfg.Global.HTTPConfig)
		if err != nil {
			return alertmanagerconfig.Receiver{}, microerror.Mask(err)
		}
		httpConfig = httpConfigCp
	}
	httpConfig.BasicAuth = &promcommonconfig.BasicAuth{
		Password: promcommonconfig.Secret(opsgenieKey),
	}

	r := alertmanagerconfig.Receiver{
		Name: key.HeartbeatReceiverName(cluster, installation),
		WebhookConfigs: []*alertmanagerconfig.WebhookConfig{
			&alertmanagerconfig.WebhookConfig{
				URL: &alertmanagerconfig.URL{
					URL: u,
				},
				HTTPConfig: httpConfig,
				NotifierConfig: alertmanagerconfig.NotifierConfig{
					VSendResolved: false,
				},
			},
		},
	}

	return r, nil
}

// EnsureCreated ensure receiver exist in cfg.Receivers and is up to date. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func EnsureCreated(cfg alertmanagerconfig.Config, cluster metav1.Object, installation, opsgenieKey string) (alertmanagerconfig.Config, bool, error) {
	desired, err := toReceiver(cfg, cluster, installation, opsgenieKey)
	if err != nil {
		return cfg, false, microerror.Mask(err)
	}

	current, _ := getReceiver(cfg, desired)

	if current != nil {
		if !reflect.DeepEqual(*current, desired) {
			*current = desired
			return cfg, true, nil
		}
	} else {
		cfg.Receivers = append(cfg.Receivers, &desired)
		return cfg, true, nil
	}

	return cfg, false, nil
}

// EnsureDeleted ensure receiver is removed from cfg.Receivers. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func EnsureDeleted(cfg alertmanagerconfig.Config, cluster metav1.Object, installation, opsgenieKey string) (alertmanagerconfig.Config, bool, error) {
	desired, err := toReceiver(cfg, cluster, installation, opsgenieKey)
	if err != nil {
		return cfg, false, microerror.Mask(err)
	}

	current, index := getReceiver(cfg, desired)

	if current != nil {
		cfg.Receivers = append(cfg.Receivers[:index], cfg.Receivers[index+1:]...)
		return cfg, true, nil
	}

	return cfg, false, nil
}

func getReceiver(cfg alertmanagerconfig.Config, receiver alertmanagerconfig.Receiver) (*alertmanagerconfig.Receiver, int) {
	for index, r := range cfg.Receivers {
		if r.Name == receiver.Name {
			return r, index
		}
	}

	return nil, -1
}

func cloneHttpConfig(hc *promcommonconfig.HTTPClientConfig) (*promcommonconfig.HTTPClientConfig, error) {
	s, err := yaml.Marshal(hc)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	r := &promcommonconfig.HTTPClientConfig{}
	err = yaml.Unmarshal(s, r)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return r, nil
}
