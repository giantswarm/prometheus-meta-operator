package remotewriteconfig

import (
	"context"
	"errors"
	"net"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/organization"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/prometheus/agent"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/prometheusquerier"
	remotewriteconfiguration "github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewrite/configuration"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "remotewriteconfig"
)

type Config struct {
	K8sClient          k8sclient.Interface
	Logger             micrologger.Logger
	OrganizationReader organization.Reader

	Customer     string
	Installation string
	Pipeline     string
	Provider     cluster.Provider
	Region       string
	Version      string

	ShardingStrategy agent.ShardingStrategy
}

type Resource struct {
	k8sClient          k8sclient.Interface
	logger             micrologger.Logger
	organizationReader organization.Reader

	customer     string
	installation string
	pipeline     string
	provider     cluster.Provider
	region       string
	version      string

	shardingStrategy agent.ShardingStrategy
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.OrganizationReader == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.OrganizationReader must not be empty")
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
	if reflect.ValueOf(config.Provider).IsZero() {
		return nil, microerror.Maskf(invalidConfigError, "config.Provider must not be empty")
	}
	if config.Version == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Version must not be empty")
	}

	r := &Resource{
		k8sClient:          config.K8sClient,
		logger:             config.Logger,
		organizationReader: config.OrganizationReader,

		customer:     config.Customer,
		installation: config.Installation,
		pipeline:     config.Pipeline,
		provider:     config.Provider,
		region:       config.Region,
		version:      config.Version,

		shardingStrategy: config.ShardingStrategy,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) desiredConfigMap(ctx context.Context, cluster metav1.Object, name string, namespace string, shards int) (*corev1.ConfigMap, error) {
	organization, err := r.organizationReader.Read(ctx, cluster)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	provider, err := key.ClusterProvider(cluster, r.provider)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	externalLabels := map[string]string{
		key.ClusterIDKey:       key.ClusterID(cluster),
		key.ClusterTypeKey:     key.ClusterType(r.installation, cluster),
		key.CustomerKey:        r.customer,
		key.InstallationKey:    r.installation,
		key.OrganizationKey:    organization,
		key.PipelineKey:        r.pipeline,
		key.ProviderKey:        provider,
		key.RegionKey:          r.region,
		key.ServicePriorityKey: key.GetServicePriority(cluster),
	}

	prometheusAgentConfig := remotewriteconfiguration.PrometheusAgentConfig{
		ExternalLabels: externalLabels,
		Image: remotewriteconfiguration.PrometheusAgentImage{
			Tag: r.version,
		},
		Shards:  shards,
		Version: r.version,
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
func (r *Resource) getShardsCountForCluster(cluster metav1.Object, currentShardCount int) (int, error) {
	clusterShardingStrategy, err := key.GetClusterShardingStrategy(cluster)
	if err != nil {
		return 0, microerror.Mask(err)
	}

	shardingStrategy := r.shardingStrategy.Merge(clusterShardingStrategy)
	headSeries, err := prometheusquerier.QueryTSDBHeadSeries(key.ClusterID(cluster))
	if err != nil {
		// If prometheus is not accessible (for instance, not running because this is a new cluster, we check if prometheus is accessible)
		var dnsError *net.DNSError
		if errors.As(err, &dnsError) {
			return shardingStrategy.ComputeShards(currentShardCount, 1), nil
		}

		return 0, microerror.Mask(err)
	}
	return shardingStrategy.ComputeShards(currentShardCount, headSeries), nil
}

func (r *Resource) createConfigMap(ctx context.Context, cluster metav1.Object, name string, namespace string, version string) error {
	shards, err := r.getShardsCountForCluster(cluster, 1)
	if err != nil {
		return microerror.Mask(err)
	}

	configMap, err := r.desiredConfigMap(ctx, cluster, name, namespace, shards)
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
