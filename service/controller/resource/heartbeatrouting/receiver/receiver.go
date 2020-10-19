package receiver

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

// EnsureCreated ensure receiver exist in cfg.Receivers and is up to date. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func EnsureCreated(cfg config.Config, cluster metav1.Object, installation, opsgenieKey string) (config.Config, bool, error) {
	desired, err := toReceiver(cluster, installation, opsgenieKey)
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
func EnsureDeleted(cfg config.Config, cluster metav1.Object, installation, opsgenieKey string) (config.Config, bool, error) {
	desired, err := toReceiver(cluster, installation, opsgenieKey)
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

func getReceiver(cfg config.Config, receiver config.Receiver) (*config.Receiver, int) {
	for index, r := range cfg.Receivers {
		if r.Name == receiver.Name {
			return r, index
		}
	}

	return nil, -1
}
