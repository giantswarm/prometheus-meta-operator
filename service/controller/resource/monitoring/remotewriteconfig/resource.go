package remotewriteconfig

import (
	"context"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/prometheusquerier"
	remotewriteconfiguration "github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewrite/configuration"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "remotewriteconfig"
)

type Config struct {
	K8sClient    k8sclient.Interface
	Logger       micrologger.Logger
	Customer     string
	Installation string
	Pipeline     string
	Provider     string
	Region       string
	Version      string
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger

	Customer     string
	Installation string
	Pipeline     string
	Provider     string
	Region       string
	Version      string
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Customer == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Customer must not be empty")
	}
	if config.Pipeline == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Pipeline must not be empty")
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Installation must not be empty")
	}
	if config.Provider == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Provider must not be empty")
	}
	if config.Version == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Version must not be empty")
	}

	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		Customer:     config.Customer,
		Installation: config.Installation,
		Pipeline:     config.Pipeline,
		Provider:     config.Provider,
		Region:       config.Region,
		Version:      config.Version,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) desiredConfigMap(cluster metav1.Object, name string, namespace string, version string, shards int) (*corev1.ConfigMap, error) {
	externalLabels := map[string]string{
		key.ClusterIDKey:       key.ClusterID(cluster),
		key.ClusterTypeKey:     key.ClusterType(r.Installation, cluster),
		key.CustomerKey:        r.Customer,
		key.InstallationKey:    r.Installation,
		key.OrganizationKey:    key.GetOrganization(cluster),
		key.PipelineKey:        r.Pipeline,
		key.ProviderKey:        r.Provider,
		key.RegionKey:          r.Region,
		key.ServicePriorityKey: key.GetServicePriority(cluster),
	}

	prometheusAgentConfig := remotewriteconfiguration.PrometheusAgentConfig{
		ExternalLabels: externalLabels,
		Image: remotewriteconfiguration.PrometheusAgentImage{
			Tag: r.Version,
		},
		Shards:  shards,
		Version: r.Version,
	}

	marshalledValues, err := yaml.Marshal(remotewriteconfiguration.RemoteWriteConfig{
		PrometheusAgentConfig: prometheusAgentConfig,
	})

	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    key.PrometheusLabels(cluster),
			Finalizers: []string{
				"monitoring.giantswarm.io/prometheus-remote-write",
			},
		},
		Data: map[string]string{
			"values": string(marshalledValues),
		},
	}, nil
}

// We want to compute the number of shards based on the number of nodes.
func (r *Resource) getShardsCountForCluster(ctx context.Context, cluster metav1.Object, currentShardCount int) (int, error) {
	headSeries, err := prometheusquerier.QueryTSDBHeadSeries(key.ClusterID(cluster))
	if err != nil {
		return 0, microerror.Mask(err)
	}
	return computeShards(currentShardCount, headSeries), nil
}

func (r *Resource) createConfigMap(ctx context.Context, cluster metav1.Object, name string, namespace string, version string) error {
	shards, err := r.getShardsCountForCluster(ctx, cluster, 1)
	if err != nil {
		return microerror.Mask(err)
	}

	configMap, err := r.desiredConfigMap(cluster, name, namespace, version, shards)
	if err != nil {
		return microerror.Mask(err)
	}

	_, err = r.k8sClient.K8sClient().CoreV1().ConfigMaps(namespace).Create(ctx, configMap, metav1.CreateOptions{})
	return microerror.Mask(err)
}

func (r *Resource) deleteConfigMap(ctx context.Context, configmap *corev1.ConfigMap) error {
	err := r.k8sClient.K8sClient().CoreV1().ConfigMaps(configmap.Namespace).Delete(ctx, configmap.Name, metav1.DeleteOptions{})
	return microerror.Mask(err)
}
