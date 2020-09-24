package oauth

import (
	"fmt"
	"reflect"

	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "ingress"
)

type Config struct {
	K8sClient      k8sclient.Interface
	Logger         micrologger.Logger
	PrometheusHost string
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().ExtensionsV1beta1().Ingresses(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc: clientFunc,
		Logger:     config.Logger,
		Name:       Name,
		GetObjectMeta: func(v interface{}) (metav1.ObjectMeta, error) {
			return getObjectMeta(v, config)
		},
		GetDesiredObject: func(v interface{}) (metav1.Object, error) {
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
		Name:      key.ClusterID(cluster),
		Namespace: key.Namespace(cluster),
		Labels:    key.Labels(cluster),
	}, nil
}

func toIngress(v interface{}, config Config) (metav1.Object, error) {
	if v == nil {
		return nil, nil
	}

	objectMeta, err := getObjectMeta(v, config)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	ingress := &extensionsv1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: extensionsv1beta1.SchemeGroupVersion.Version,
			Kind:       "Ingress",
		},
		ObjectMeta: objectMeta,
		Spec: extensionsv1beta1.IngressSpec{
			TLS: []extensionsv1beta1.IngressTLS{
				{
					Hosts:      []string{config.PrometheusHost},
					SecretName: "prometheus-certificates",
				},
			},
			Rules: []extensionsv1beta1.IngressRule{
				{
					Host: config.PrometheusHost,
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Path: fmt.Sprintf("/%s/oauth2", key.ClusterID(cluster)),
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: key.Oauth2ServiceName(),
										ServicePort: intstr.FromInt(key.Oauth2ServicePort()),
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
	c := current.(*extensionsv1beta1.Ingress)
	d := desired.(*extensionsv1beta1.Ingress)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
