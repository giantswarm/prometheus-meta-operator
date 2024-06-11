package remotewriteingress

import (
	"fmt"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "monitoringremotewriteingress"
)

type Config struct {
	K8sClient    k8sclient.Interface
	Logger       micrologger.Logger
	BaseDomain   string
	ExternalDNS  bool
	MimirEnabled bool
}

type Resource struct {
	config Config
}

func New(config Config) (*Resource, error) {
	return &Resource{config}, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:        fmt.Sprintf("prometheus-%s-remote-write", key.ClusterID(cluster)),
		Namespace:   key.Namespace(cluster),
		Labels:      key.PrometheusLabels(cluster),
		Annotations: key.RemoteWriteAuthenticationAnnotations(r.config.BaseDomain, r.config.ExternalDNS),
	}, nil
}

func (r *Resource) toIngress(v interface{}) (*networkingv1.Ingress, error) {
	if v == nil {
		return nil, nil
	}

	objectMeta, err := r.getObjectMeta(v)
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
	pathType := networkingv1.PathTypeImplementationSpecific
	ingressClassName := key.IngressClassName
	ingress := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: networkingv1.SchemeGroupVersion.Version,
			Kind:       "Ingress",
		},
		ObjectMeta: objectMeta,
		Spec: networkingv1.IngressSpec{
			IngressClassName: &ingressClassName,
			Rules: []networkingv1.IngressRule{
				{
					Host: r.config.BaseDomain,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: fmt.Sprintf("/%s/api/v1/write", key.ClusterID(cluster)),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: key.PrometheusServiceName,
											Port: networkingv1.ServiceBackendPort{
												Number: key.PrometheusPort(),
											},
										},
									},
									PathType: &pathType,
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

func (r *Resource) hasChanged(current, desired metav1.Object) bool {
	c := current.(*networkingv1.Ingress)
	d := desired.(*networkingv1.Ingress)

	return !reflect.DeepEqual(c.Spec, d.Spec) || !reflect.DeepEqual(c.GetAnnotations(), d.GetAnnotations())
}
