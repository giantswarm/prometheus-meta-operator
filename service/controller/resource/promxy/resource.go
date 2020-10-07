package promxy

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/ghodss/yaml"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/kubernetes"
	"github.com/prometheus/prometheus/pkg/relabel"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	proxyconfig "github.com/giantswarm/prometheus-meta-operator/service/controller/resource/promxy/config"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

type Config struct {
	K8sClient    k8sclient.Interface
	Logger       micrologger.Logger
	Installation string
	Provider     string
}

type Resource struct {
	k8sClient    k8sclient.Interface
	logger       micrologger.Logger
	installation string
	provider     string
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}

	r := &Resource{
		logger:       config.Logger,
		k8sClient:    config.K8sClient,
		installation: config.Installation,
		provider:     config.Provider,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return "promxy"
}

func (r *Resource) toServerGroup(cluster metav1.Object) (*ServerGroup, error) {
	httpClient := config.HTTPClientConfig{
		TLSConfig: config.TLSConfig{
			CAFile:             "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
			InsecureSkipVerify: true,
		},
		BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
	}

	apiServerHost := r.k8sClient.RESTConfig().Host
	apiServerURL, err := url.Parse(apiServerHost)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &ServerGroup{
		Scheme:         "http",
		RemoteReadPath: "api/v1/read",
		RemoteRead:     true,
		HTTPConfig: HTTPClientConfig{
			DialTimeout: time.Millisecond * 200, // Default dial timeout of 200ms
		},
		Labels: model.LabelSet{
			model.LabelName("installation"):     model.LabelValue(r.installation),
			model.LabelName(key.ClusterIDKey()): model.LabelValue(key.ClusterID(cluster)),
			model.LabelName("provider"):         model.LabelValue(r.provider),
		},
		PathPrefix: fmt.Sprintf("/%s", key.ClusterID(cluster)),
		RelabelConfigs: []*relabel.Config{
			{
				SourceLabels: []model.LabelName{
					"__meta_kubernetes_pod_label_app_kubernetes_io_managed_by",
					"__meta_kubernetes_pod_label_app_kubernetes_io_name",
					"__meta_kubernetes_pod_label_app_kubernetes_io_instance",
				},
				Action:    "keep",
				Separator: ";",
				Regex:     relabel.MustNewRegexp(fmt.Sprintf("prometheus-meta-operator;prometheus;%s", key.ClusterID(cluster))),
			},
		},
		KubernetesSDConfigs: []*kubernetes.SDConfig{
			&kubernetes.SDConfig{
				APIServer: config.URL{
					URL: apiServerURL,
				},
				Role:             kubernetes.RolePod,
				HTTPClientConfig: httpClient,
				NamespaceDiscovery: kubernetes.NamespaceDiscovery{
					Names: []string{
						key.Namespace(cluster),
					},
				},
			},
		},
	}, nil
}

func (r *Resource) readFromConfig(configMap *v1.ConfigMap) (*proxyconfig.Config, error) {
	content, ok := configMap.Data[key.PromxyConfigFileName()]
	if !ok {
		return nil, microerror.Mask(invalidConfigError)
	}

	config := proxyconfig.Config{}
	err := yaml.Unmarshal([]byte(content), &config)
	return &config, microerror.Mask(err)

}

func (r *Resource) updateConfig(ctx context.Context, configMap *v1.ConfigMap, config *proxyconfig.Config) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	configMap.Data[key.PromxyConfigFileName()] = string(bytes)
	_, err = r.k8sClient.K8sClient().CoreV1().ConfigMaps(key.PromxyConfigMapNamespace()).Update(ctx, configMap, metav1.UpdateOptions{})

	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
func promxyAdd(p proxyconfig.PromxyConfig, group *ServerGroup) proxyconfig.PromxyConfig {
	p.ServerGroups = append(p.ServerGroups, group)

	return p
}

func promxyRemove(p proxyconfig.PromxyConfig, group *ServerGroup) proxyconfig.PromxyConfig {
	var index int
	for key, val := range p.ServerGroups {
		if val.PathPrefix == group.PathPrefix {
			index = key
		}
	}
	p.ServerGroups = append(p.ServerGroups[:index], p.ServerGroups[index+1:]...)

	return p
}
