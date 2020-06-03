package ingress

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
	Name = "ingress"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	BaseDomain            string
	LetsEncryptEnabled    bool
	WhitelistingEnabled   bool
	WhitelistingSourceIPs string
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger

	baseDomain            string
	letsEncryptEnabled    bool
	whitelistingEnabled   bool
	whitelistingSourceIPs string
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

		baseDomain:            config.BaseDomain,
		letsEncryptEnabled:    config.LetsEncryptEnabled,
		whitelistingEnabled:   config.WhitelistingEnabled,
		whitelistingSourceIPs: config.WhitelistingSourceIPs,
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

	host := key.TenantClusterHost(cluster, r.baseDomain)
	secretName := key.TenantClusterSecret(cluster)
	ingressName := fmt.Sprintf("%s-prometheus", cluster.GetName())

	annotations := map[string]string{
		"kubernetes.io/ingress.class":             "nginx",
		"nginx.ingress.kubernetes.io/auth-signin": fmt.Sprintf("https://%s/oauth2/start", host),
		"nginx.ingress.kubernetes.io/auth-url":    fmt.Sprintf("https://%s/oauth2/start", host),
	}

	if r.letsEncryptEnabled {
		annotations["kubernetes.io/tls-acme"] = "true"
	}
	if r.whitelistingEnabled {
		annotations["nginx.ingress.kubernetes.io/whitelist-source-range"] = r.whitelistingSourceIPs
	}

	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressName,
			Namespace: key.Namespace(cluster),
			Labels: map[string]string{
				key.ClusterIDKey(): key.ClusterID(cluster),
			},
			Annotations: annotations,
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
										ServiceName: cluster.GetName(),
										ServicePort: intstr.FromInt(9090),
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
