package oauth2proxy

import (
	"fmt"

	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "oauth2proxy"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	BaseDomain string
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger

	baseDomain string
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.BaseDomain == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.BaseDomain must not be empty", config)
	}

	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		baseDomain: config.BaseDomain,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) toIngress(v interface{}) (*v1beta1.Ingress, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secretName := key.TenantClusterSecret(cluster)
	host := key.TenantClusterHost(cluster, r.baseDomain)
	ingressName := fmt.Sprintf("%s-prometheus-oauth2-proxy", cluster.GetName())

	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: key.Namespace(cluster),
			Labels: map[string]string{
				key.ClusterIDKey(): key.ClusterID(cluster),
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: host,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/oauth2",
									Backend: v1beta1.IngressBackend{
										ServiceName: "oauth2-proxy",
										ServicePort: intstr.FromInt(4180),
									},
								},
							},
						},
					},
				},
			},
			TLS: []v1beta1.IngressTLS{
				{
					Hosts: []string{
						host,
					},
					SecretName: secretName,
				},
			},
		},
	}

	return ingress, nil
}
