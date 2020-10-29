package alert

import (
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alert/rules"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "alert"
)

type Config struct {
	Installation     string
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
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}

	clientFunc := func(namespace string) generic.Interface {
		c := config.PrometheusClient.MonitoringV1().PrometheusRules(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:    clientFunc,
		Logger:        config.Logger,
		Name:          Name,
		GetObjectMeta: getObjectMeta,
		GetDesiredObject: func(obj interface{}) (metav1.Object, error) {
			return getRules(obj, config.Installation)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(obj interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      "prometheus-meta-operator",
		Namespace: key.Namespace(cluster),
		Labels: map[string]string{
			key.ClusterIDKey(): key.ClusterID(cluster),
		},
	}, nil
}

func getRules(obj interface{}, installation string) (metav1.Object, error) {
	objectMeta, err := getObjectMeta(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	r := &promv1.PrometheusRule{
		ObjectMeta: objectMeta,
		Spec: promv1.PrometheusRuleSpec{
			Groups: []promv1.RuleGroup{
				rules.LabellingSchemaValidationRule(cluster, installation),
				rules.Heartbeat(cluster, installation),
			},
		},
	}

	return r, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*promv1.PrometheusRule)
	d := desired.(*promv1.PrometheusRule)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
