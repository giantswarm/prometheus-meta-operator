package service

import (
	"reflect"

	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
	toService := func(v interface{}) (metav1.Object, error) {
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

		current, err := clientFunc(cluster.GetNamespace()).Get("prometheus", metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return nil, microerror.Mask(err)
		} else {
			c := current.(*corev1.Service)
			service.Spec.ClusterIP = c.Spec.ClusterIP
		}

		return service, nil
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

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Service)
	d := desired.(*corev1.Service)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
