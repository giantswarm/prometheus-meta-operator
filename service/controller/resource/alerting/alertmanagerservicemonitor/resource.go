package alertmanagerservicemonitor

import (
	"context"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringv1client "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "alertmanagerservicemonitor"
)

type Config struct {
	Client       monitoringv1client.Interface
	Logger       micrologger.Logger
	Installation string
	Provider     string
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.Client.MonitoringV1().ServiceMonitors(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:    clientFunc,
		Logger:        config.Logger,
		Name:          Name,
		GetObjectMeta: getObjectMeta,
		GetDesiredObject: func(ctx context.Context, v interface{}) (metav1.Object, error) {
			return toServiceMonitor(ctx, v, config)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
	return metav1.ObjectMeta{
		Name:      "alertmanager",
		Namespace: key.NamespaceMonitoring(),
		Labels:    key.AlertmanagerLabels(),
	}, nil
}

func toServiceMonitor(ctx context.Context, v interface{}, config Config) (metav1.Object, error) {
	objectMeta, err := getObjectMeta(ctx, v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	enableHttp2 := true
	sm := &monitoringv1.ServiceMonitor{
		ObjectMeta: objectMeta,
		Spec: monitoringv1.ServiceMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"operated-alertmanager": "true",
				},
			},
			NamespaceSelector: monitoringv1.NamespaceSelector{
				MatchNames: []string{key.NamespaceMonitoring()},
			},
			Endpoints: []monitoringv1.Endpoint{
				monitoringv1.Endpoint{
					Port:        "web",
					Path:        "/metrics",
					EnableHttp2: &enableHttp2,
				},
			},
		},
	}

	return sm, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*monitoringv1.ServiceMonitor)
	d := desired.(*monitoringv1.ServiceMonitor)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
