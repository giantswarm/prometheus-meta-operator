package heartbeatrouting

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/prometheus/alertmanager/config"
	commoncfg "github.com/prometheus/common/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func toReceiver(cluster metav1.Object, installation string, opsgenieKey string) (config.Receiver, error) {
	u, err := url.Parse(fmt.Sprintf("https://api.opsgenie.com/v2/heartbeats/%s/ping", key.HeartbeatName(cluster, installation)))
	if err != nil {
		return config.Receiver{}, microerror.Mask(err)
	}

	r := config.Receiver{
		Name: key.HeartbeatReceiverName(cluster, installation),
		WebhookConfigs: []*config.WebhookConfig{
			&config.WebhookConfig{
				URL: &config.URL{
					URL: u,
				},
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

	return r, nil
}

// ensureReceiver ensure receiver exist in cfg.Receivers and is up to date. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func ensureReceiver(cfg config.Config, receiver config.Receiver) (config.Config, bool) {
	r, _ := getReceiver(cfg, receiver)

	if r != nil {
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
	r, index := getReceiver(cfg, receiver)

	if r != nil {
		cfg.Receivers = append(cfg.Receivers[:index], cfg.Receivers[index+1:]...)
		return cfg, true
	}

	return cfg, false
}

func getReceiver(cfg config.Config, receiver config.Receiver) (*config.Receiver, int) {
	for index, r := range cfg.Receivers {
		if r.Name == receiver.Name {
			return r, index
		}
	}

	return nil, -1
}
