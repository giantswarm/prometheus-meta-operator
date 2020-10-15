package heartbeatrouting

import (
	"fmt"
	"reflect"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

var (
	one, _     = model.ParseDuration("1s")
	fifteen, _ = model.ParseDuration("15s")
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
		GroupInterval:  &one,
		GroupWait:      &one,
		RepeatInterval: &fifteen,
	}
}

// ensureRoute ensure route exist in cfg.Route and is up to date. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func ensureRoute(cfg config.Config, route config.Route) (config.Config, bool, error) {
	r, _ := getRoute(&cfg, route)

	if r != nil {
		if !reflect.DeepEqual(*r, route) {
			*r = route
			return cfg, true, nil
		}
	} else {
		if cfg.Route == nil {
			return cfg, false, microerror.Mask(emptyRouteError)
		}
		cfg.Route.Routes = append(cfg.Route.Routes, &route)
		return cfg, true, nil
	}

	return cfg, false, nil
}

// removeRoute ensure route is removed from cfg.Receivers. Returns true when changes have been made to cfg.
// Return untouched cfg and false when no changes are made.
func removeRoute(cfg config.Config, route config.Route) (config.Config, bool) {
	r, index := getRoute(&cfg, route)

	if r != nil {
		cfg.Route.Routes = append(cfg.Route.Routes[:index], cfg.Route.Routes[index+1:]...)
		return cfg, true
	}

	return cfg, false
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
