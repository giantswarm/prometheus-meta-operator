package route

import (
	"reflect"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func toRoute(cluster metav1.Object, installation string) (config.Route, error) {
	one, err := model.ParseDuration("1s")
	if err != nil {
		return config.Route{}, microerror.Mask(err)
	}

	fifteen, err := model.ParseDuration("15s")
	if err != nil {
		return config.Route{}, microerror.Mask(err)
	}

	r := config.Route{
		Receiver: key.HeartbeatReceiverName(cluster, installation),
		Match: map[string]string{
			key.ClusterIDKey():    key.ClusterID(cluster),
			key.InstallationKey(): installation,
			key.TypeKey():         key.Heartbeat(),
		},
		Continue:       false,
		GroupInterval:  &one,
		GroupWait:      &one,
		RepeatInterval: &fifteen,
	}

	return r, nil
}

// EnsureCreated ensure route exist in cfg.Route and is up to date. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func EnsureCreated(cfg config.Config, cluster metav1.Object, installation string) (config.Config, bool, error) {
	desired, err := toRoute(cluster, installation)
	if err != nil {
		return cfg, false, microerror.Mask(err)
	}

	current, _ := getRoute(&cfg, desired)

	if current != nil {
		if !reflect.DeepEqual(*current, desired) {
			*current = desired
			return cfg, true, nil
		}
	} else {
		if cfg.Route == nil {
			return cfg, false, microerror.Mask(emptyRouteError)
		}
		cfg.Route.Routes = append(cfg.Route.Routes, &desired)
		return cfg, true, nil
	}

	return cfg, false, nil
}

// EnsureDeleted ensure route is removed from cfg.Receivers. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func EnsureDeleted(cfg config.Config, cluster metav1.Object, installation string) (config.Config, bool, error) {
	desired, err := toRoute(cluster, installation)
	if err != nil {
		return cfg, false, microerror.Mask(err)
	}

	current, index := getRoute(&cfg, desired)

	if current != nil {
		cfg.Route.Routes = append(cfg.Route.Routes[:index], cfg.Route.Routes[index+1:]...)
		return cfg, true, nil
	}

	return cfg, false, nil
}

func getRoute(cfg *config.Config, route config.Route) (*config.Route, int) {
	if cfg.Route != nil {
		for index, r := range cfg.Route.Routes {
			if r.Receiver == route.Receiver {
				return r, index
			}
		}
	}

	return nil, -1
}
