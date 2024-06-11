package ciliumnetpol

import (
	"net/url"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"golang.org/x/net/http/httpproxy"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name           = "ciliumnetpol"
	labelName      = "ciliumnetpol/name"
	labelNamespace = "ciliumnetpol/namespace"
)

type Config struct {
	DynamicK8sClient dynamic.Interface
	Logger           micrologger.Logger

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

func toCiliumNetworkPolicy(v interface{}) (*unstructured.Unstructured, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	worldPorts := []map[string]string{
		{
			"port": "443",
		},
		// Grafana cloud mimir port
		{
			"port": "6443",
		},
		// Grafana cloud squid proxy port
		{
			"port": "3128",
		},
	}
	// We need to retrieve the proxy port from the environment variables
	// and add it to the CiliumNetworkPolicy.
	proxyConfig := httpproxy.FromEnvironment()
	proxyUrl := proxyConfig.HTTPProxy
	proxyDefaultPort := "80"
	if proxyUrl == "" {
		proxyUrl = proxyConfig.HTTPSProxy
		proxyDefaultPort = "443"
	}
	if proxyUrl != "" {
		proxyURL, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		proxyPort := proxyURL.Port()
		if proxyPort == "" {
			proxyPort = proxyDefaultPort
		}
		worldPorts = append(worldPorts, map[string]string{"port": proxyPort})
	}

	ciliumNetworkPolicy := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cilium.io/v2",
			"kind":       "CiliumNetworkPolicy",
			"metadata": map[string]interface{}{
				"name":      key.ClusterID(cluster) + "-prometheus",
				"namespace": key.Namespace(cluster),
				"labels": map[string]string{
					"app.kubernetes.io/name": "prometheus",
				},
			},
			"spec": map[string]interface{}{
				"endpointSelector": map[string]interface{}{
					"matchLabels": map[string]string{
						"app.kubernetes.io/name": "prometheus",
					},
				},
				"egress": []map[string]interface{}{
					{
						"toEntities": []string{
							"kube-apiserver",
							"cluster",
						},
					},
					{
						"toEntities": []string{
							"world",
						},
						"toPorts": []map[string]interface{}{
							{
								"ports": worldPorts,
							},
						},
					},
				},
				"ingress": []map[string]interface{}{
					{
						"fromEntities": []string{
							"cluster",
						},
						"toPorts": []map[string]interface{}{
							{
								"ports": []map[string]string{
									{
										"port": "9090",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return ciliumNetworkPolicy, nil
}

func hasCiliumNetworkPolicyChanged(current *unstructured.Unstructured, desired *unstructured.Unstructured) bool {
	return !reflect.DeepEqual(current.Object["spec"], desired.Object["spec"])
}
