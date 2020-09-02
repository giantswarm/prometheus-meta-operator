package alert

import (
	"reflect"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alert/rules"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
)

const (
	Name = "alert"
)

type Config struct {
	PrometheusClient promclient.Interface
	Logger           micrologger.Logger
}

func New(config Config) (*generic.Resource, error) {
	if config.PrometheusClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.PrometheusClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	clientFunc := func(namespace string) generic.Interface {
		c := config.PrometheusClient.MonitoringV1().PrometheusRules(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:     clientFunc,
		Logger:         config.Logger,
		Name:           Name,
		GetObjectMeta:  rules.GetObjectMeta,
		ToCR:           toPrometheusRule,
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func toPrometheusRule(obj interface{}) (metav1.Object, error) {
	return rules.ExampleRule(obj)
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*promv1.PrometheusRule)
	d := desired.(*promv1.PrometheusRule)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
