package service

import (
	"reflect"

	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "service"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

func New(config Config) (*generic.Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().CoreV1().Services(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:     clientFunc,
		Logger:         config.Logger,
		Name:           Name,
		ToCR:           toService,
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func toService(v interface{}) (metav1.Object, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prometheus",
			Namespace: key.Namespace(cluster),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "prometheus",
					Port:     int32(9091),
					Protocol: "TCP",
				},
			},
			Selector: map[string]string{
				"app": "frontend",
			},
		},
	}

	return service, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Service)
	d := desired.(*corev1.Service)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
