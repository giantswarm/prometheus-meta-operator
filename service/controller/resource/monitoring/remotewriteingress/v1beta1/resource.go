package v1beta1

import (
	"context"
	"fmt"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

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
		c := config.K8sClient.K8sClient().NetworkingV1beta1().Ingresses(namespace)
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
		Annotations: map[string]string{
			"nginx.ingress.kubernetes.io/auth-type":   "basic",
			"nginx.ingress.kubernetes.io/auth-secret": key.RemoteWriteAgentSecretName,
			"nginx.ingress.kubernetes.io/auth-realm":  "Authentication Required",
		},
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

	// Note we only configure a path that will cause a location/HTTP path to be
	// added to the proxy for the base domain here.
	//
	// The common server configuration like TLS is configured in the `Ingress`
	// installed by PMO chart itself. We want to avoid duplicating the TLS
	// configuration here as only one certificate can be used for a given
	// domain, and if multiple `Ingress` resources specify that configuration
	// the Ingress Controller just picks a random one (first one it finds). And
	// if it picks up a copied certificate, and that copy happens to become out
	// of date at some point, it would break HTTPS access to that domain.
	//
	// So we want TLS configuration to be controlled only by the `Ingress`
	// resource that also defines the source of the certificates (i.e. the
	// Let's Encrypt annotation or the static source for the installation)
	// so we know as soon as it's updated IC will be using it.
	ingress := &networkingv1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: networkingv1beta1.SchemeGroupVersion.Version,
			Kind:       "Ingress",
		},
		ObjectMeta: objectMeta,
		Spec: networkingv1beta1.IngressSpec{
			Rules: []networkingv1beta1.IngressRule{
				{
					Host: config.BaseDomain,
					IngressRuleValue: networkingv1beta1.IngressRuleValue{
						HTTP: &networkingv1beta1.HTTPIngressRuleValue{
							Paths: []networkingv1beta1.HTTPIngressPath{
								{
									Path: fmt.Sprintf("/%s/api/v1/write", key.ClusterID(cluster)),
									Backend: networkingv1beta1.IngressBackend{
										ServiceName: key.PrometheusServiceName,
										ServicePort: intstr.FromInt(int(key.PrometheusPort())),
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
	c := current.(*networkingv1beta1.Ingress)
	d := desired.(*networkingv1beta1.Ingress)

	return !reflect.DeepEqual(c.Spec, d.Spec) || !reflect.DeepEqual(c.GetAnnotations(), d.GetAnnotations())
}
