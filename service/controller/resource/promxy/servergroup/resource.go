package servergroup

import (
	"fmt"
	"net/url"
	"time"

	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/kubernetes"
	"github.com/prometheus/prometheus/pkg/relabel"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/promxy"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

type Config struct {
	K8sClient    k8sclient.Interface
	Logger       micrologger.Logger
	Installation string
}

type Resource struct {
	k8sClient    k8sclient.Interface
	logger       micrologger.Logger
	installation string
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
	}

	return r, nil
}

func (r *Resource) Name() string {
	return "promxy-server-group"
}

func (r *Resource) toServerGroup(apiServerURL *url.URL, cluster metav1.Object) promxy.ServerGroup {
	httpClient := config.HTTPClientConfig{
		TLSConfig: config.TLSConfig{
			CAFile:             "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
			InsecureSkipVerify: true,
		},
		BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
	}

	return promxy.ServerGroup{
		AntiAffinity:   time.Second * 10,
		Scheme:         "http",
		RemoteReadPath: "api/v1/read",
		RemoteRead:     true,
		HTTPConfig: promxy.HTTPClientConfig{
			DialTimeout: time.Millisecond * 200, // Default dial timeout of 200ms
		},
		Labels: model.LabelSet{
			model.LabelName("installation"):     model.LabelValue(r.installation),
			model.LabelName(key.ClusterIDKey()): model.LabelValue(key.ClusterID(cluster)),
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
	}
}
