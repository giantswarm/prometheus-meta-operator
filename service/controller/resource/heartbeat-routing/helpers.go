package heartbeatrouting

import (
	"context"

	"github.com/prometheus/alertmanager/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func toRoute(cluster metav1.Object, installation string) config.Route {
	// TODO: implement me.
	return config.Route{}
}

func toReceiver(cluster metav1.Object, installation string) config.Receiver {
	// TODO: implement me.
	return config.Receiver{}
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

func (r *Resource) readFromConfig(configMap *v1.ConfigMap) (config.Config, error) {
	// TODO: implement me.
	return config.Config{}, nil
}

func (r *Resource) updateConfig(ctx context.Context, configMap *v1.ConfigMap, cfg config.Config) error {
	// TODO: implement me.
	return nil
}
