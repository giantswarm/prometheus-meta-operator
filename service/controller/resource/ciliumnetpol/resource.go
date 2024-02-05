package ciliumnetpol

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	Name           = "ciliumnetpol"
	labelName      = "ciliumnetpol/name"
	labelNamespace = "ciliumnetpol/namespace"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

type Resource struct {
	k8sClient k8sclient.Interface
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
			"metadata": map[string]interface{}{
				"name":      key.Namespace(cluster) + "-prometheus",
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
					map[string]interface{}{
						"toEntities": []string{
							"world",
						},
						"toPorts": []map[string]interface{}{
							map[string]interface{}{
								"ports": []map[string]interface{}{
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
								"ports": []map[string]interface{}{
									map[string]string{
										"port": "9090",
									},
								},
							},
						},
					},
				},
			},
		}
	}

	return ciliumNetworkPolicy, nil
}
