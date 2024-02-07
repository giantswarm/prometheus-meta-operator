package ciliumnetpol

import (
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
	K8sClient dynamic.Interface
	Logger    micrologger.Logger
}

type Resource struct {
	k8sClient dynamic.Interface
	logger    micrologger.Logger
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,
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
					map[string]interface{}{
						"toEntities": []string{
							"kube-apiserver",
							"cluster",
						},
					},
					map[string]interface{}{
						"toEntities": []string{
							"world",
						},
						"toPorts": []map[string]interface{}{
							map[string]interface{}{
								"ports": []map[string]string{
									map[string]string{
										"port": "443",
									},
									map[string]string{
										"port": "6443",
									},
								},
							},
						},
					},
				},
				"ingress": []map[string]interface{}{
					map[string]interface{}{
						"fromEntities": []string{
							"cluster",
						},
					},
					map[string]interface{}{
						"fromPorts": []map[string]interface{}{
							map[string]interface{}{
								"ports": []map[string]string{
									map[string]string{
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
