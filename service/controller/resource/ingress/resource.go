package ingress

import (
	"fmt"
	"reflect"

	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "ingress"
)

var (
	baseDomain string
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	BaseDomain string
}

func New(config Config) (*generic.Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.BaseDomain == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.BaseDomain must not be empty", config)
	}

	baseDomain = config.BaseDomain

	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().ExtensionsV1beta1().Ingresses(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:     clientFunc,
		Logger:         config.Logger,
		Name:           Name,
		ToCR:           toIngress,
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func toIngress(v interface{}) (metav1.Object, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	host := fmt.Sprintf("%s-%s", cluster.GetName(), baseDomain)

	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.GetName(),
			Namespace: key.Namespace(cluster),
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: host,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: v1beta1.IngressBackend{
										ServiceName: "prometheus",
										ServicePort: intstr.FromInt(9091),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return ingress, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*v1beta1.Ingress)
	d := desired.(*v1beta1.Ingress)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
