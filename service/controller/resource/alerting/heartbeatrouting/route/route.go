package route

import (
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	alertmanagerconfig "github.com/giantswarm/prometheus-meta-operator/pkg/alertmanager/config"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func toRoute(cluster metav1.Object, installation string) (alertmanagerconfig.Route, error) {
	// We ping OpsGenie every minute
	repeatInterval, err := model.ParseDuration("1m")
	if err != nil {
		return alertmanagerconfig.Route{}, microerror.Mask(err)
	}

	r := alertmanagerconfig.Route{
		Receiver: key.HeartbeatReceiverName(cluster, installation),
		Match: map[string]string{
			key.ClusterIDKey():    key.ClusterID(cluster),
			key.InstallationKey(): installation,
			key.TypeKey():         key.Heartbeat(),
		},
		Continue:       false,
		GroupBy:        []model.LabelName{"..."},
		RepeatInterval: &repeatInterval,
	}

	return r, nil
}

// EnsureCreated ensure route exist in cfg.Route and is up to date. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func EnsureCreated(cfg alertmanagerconfig.Config, cluster metav1.Object, installation string) (alertmanagerconfig.Config, bool, error) {
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
func EnsureDeleted(cfg alertmanagerconfig.Config, cluster metav1.Object, installation string) (alertmanagerconfig.Config, bool, error) {
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

func getRoute(cfg *alertmanagerconfig.Config, route alertmanagerconfig.Route) (*alertmanagerconfig.Route, int) {
	if cfg.Route != nil {
		for index, r := range cfg.Route.Routes {
			if r.Receiver == route.Receiver {
				return r, index
			}
		}
	}

	return nil, -1
}
