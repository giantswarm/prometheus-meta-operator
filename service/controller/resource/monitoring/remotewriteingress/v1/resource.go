package v1

import (
	"context"
	"fmt"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "monitoringremotewriteingress"
)

type Config struct {
	K8sClient  k8sclient.Interface
	Logger     micrologger.Logger
	BaseDomain string
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().NetworkingV1().Ingresses(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc: clientFunc,
		Logger:     config.Logger,
		Name:       Name,
		GetObjectMeta: func(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
			return getObjectMeta(v, config)
		},
		GetDesiredObject: func(ctx context.Context, v interface{}) (metav1.Object, error) {
			return toIngress(v, config)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(v interface{}, config Config) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      fmt.Sprintf("prometheus-%s-remote-write", key.ClusterID(cluster)),
		Namespace: key.Namespace(cluster),
		Labels:    key.PrometheusLabels(cluster),
	}, nil
}

func toIngress(v interface{}, config Config) (metav1.Object, error) {
	if v == nil {
		return nil, nil
	}

	return &networkingv1.Ingress{}, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*networkingv1.Ingress)
	d := desired.(*networkingv1.Ingress)

	return !reflect.DeepEqual(c.Spec, d.Spec) || !reflect.DeepEqual(c.GetAnnotations(), d.GetAnnotations())
}
