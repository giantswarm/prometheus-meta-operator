package ciliumnetpol

import (
	"net/url"
	"os"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
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
	Proxy            func(reqURL *url.URL) (*url.URL, error)
	Logger           micrologger.Logger
}

type Resource struct {
	dynamicK8sClient dynamic.Interface
	logger           micrologger.Logger
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		dynamicK8sClient: config.DynamicK8sClient,
		logger:           config.Logger,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func toCiliumNetworkPolicy(v interface{}) (*unstructured.Unstructured, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	ports := []map[string]string{
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
	proxyURL := os.Getenv("HTTP_PROXY")
	if proxyURL != "" {
		proxyURL, err := url.Parse(proxyURL)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		proxyPort := proxyURL.Port()
		if proxyPort == "" {
			proxyPort = "80"
		}
		ports = append(ports, map[string]string{"port": proxyPort})
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
								"ports": ports,
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
