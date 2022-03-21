package alertmanagerwiring

import (
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "alertmanagerwiring"
)

var (
	alertmanagerConfig = `- static_configs:
  - targets:
    - alertmanager.monitoring:9093
  - targets:
    - alertmanager-operated.monitoring:9093	
  scheme: http
  timeout: 10s
  api_version: v2
`
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	Installation string
}

func New(config Config) (*generic.Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().CoreV1().Secrets(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:    clientFunc,
		Logger:        config.Logger,
		Name:          Name,
		GetObjectMeta: getObjectMeta,
		GetDesiredObject: func(v interface{}) (metav1.Object, error) {
			return toSecret(v, config)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.AlertManagerSecretName(),
		Namespace: key.Namespace(cluster),
	}, nil
}

func toData(v interface{}, config Config) []byte {
	return []byte(alertmanagerConfig)
}

func toSecret(v interface{}, config Config) (metav1.Object, error) {
	objectMeta, err := getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		Data: map[string][]byte{
			key.AlertManagerKey(): toData(v, config),
		},
	}

	return secret, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
