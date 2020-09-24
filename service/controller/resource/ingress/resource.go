package ingress

import (
	"fmt"
	"reflect"

	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
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

	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().ExtensionsV1beta1().Ingresses(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:    clientFunc,
		Logger:        config.Logger,
		Name:          Name,
		GetObjectMeta: getObjectMeta,
		GetDesiredObject: func(v interface{}) (metav1.Object, error) {
			return toIngress(v, config.BaseDomain)
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
		return metav1.ObjectMeta{}, err
	}

	return metav1.ObjectMeta{
		Name:      cluster.GetName(),
		Namespace: key.Namespace(cluster),
		Labels: map[string]string{
			"app": "prometheus",
		},
		Annotations: map[string]string{
			"kubernetes.io/ingress.class":                        "nginx",
			"kubernetes.io/tls-acme":                             "true",
			"nginx.ingress.kubernetes.io/auth-signin":            "https://$host/oauth2/start?rd=$escaped_request_uri",
			"nginx.ingress.kubernetes.io/auth-url":               "https://$host/oauth2/auth",
			"nginx.ingress.kubernetes.io/whitelist-source-range": "185.102.95.187/32,95.179.153.65/32",
		},
	}, nil
}

func toIngress(v interface{}, baseDomain string) (metav1.Object, error) {
	objectMeta, err := getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	path := fmt.Sprintf("/%s", cluster.GetName())

	return &v1beta1.Ingress{
		ObjectMeta: objectMeta,
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					Hosts:      []string{baseDomain},
					SecretName: "monitoring-prometheus", // TODO: is this new?
				},
			},
			Rules: []v1beta1.IngressRule{
				{
					Host: baseDomain,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Backend: v1beta1.IngressBackend{
										ServiceName: "prometheus-operated",
										ServicePort: intstr.FromInt(9090),
									},
									Path: path,
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*v1beta1.Ingress)
	d := desired.(*v1beta1.Ingress)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
