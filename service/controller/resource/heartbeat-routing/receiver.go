package heartbeatrouting

import (
	"fmt"
	"reflect"

	"github.com/prometheus/alertmanager/config"
	commoncfg "github.com/prometheus/common/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func toReceiver(cluster metav1.Object, installation string, opsgenieKey string) config.Receiver {
	name := fmt.Sprintf("heartbeat_%s_%s", installation, key.ClusterID(cluster))

	return config.Receiver{
		Name: name,
		WebhookConfigs: []*config.WebhookConfig{
			&config.WebhookConfig{
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

// ensureReceiver ensure receiver exist in cfg.Receivers and is up to date. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func ensureReceiver(cfg config.Config, receiver config.Receiver) (config.Config, bool) {
	r, _, exist := getReceiver(cfg, receiver)

	if exist {
		if !reflect.DeepEqual(*r, receiver) {
			*r = receiver
			return cfg, true
		}
	} else {
		cfg.Receivers = append(cfg.Receivers, &receiver)
		return cfg, true
	}

	return cfg, false
}

// removeReceiver ensure receiver is removed from cfg.Receivers. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func removeReceiver(cfg config.Config, receiver config.Receiver) (config.Config, bool) {
	_, index, exist := getReceiver(cfg, receiver)

	if exist {
		cfg.Receivers = append(cfg.Receivers[:index], cfg.Receivers[index+1:]...)
		return cfg, true
	}

	return cfg, false
}

func getReceiver(cfg config.Config, receiver config.Receiver) (*config.Receiver, int, bool) {
	for index, r := range cfg.Receivers {
		if r.Name == receiver.Name {
			return r, index, true
		}
	}

	return nil, -1, false
}
