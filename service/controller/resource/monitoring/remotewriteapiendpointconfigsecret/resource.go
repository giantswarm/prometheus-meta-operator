package remotewriteapiendpointconfigsecret

import (
	"context"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/organization"
	remotewriteconfiguration "github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewrite/configuration"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "remotewriteapiendpointconfigsecret"
)

type Config struct {
	K8sClient          k8sclient.Interface
	Logger             micrologger.Logger
	OrganizationReader organization.Reader

	BaseDomain   string
	Customer     string
	Installation string
	InsecureCA   bool
	Pipeline     string
	Provider     cluster.Provider
	Region       string
	Version      string
}

type Resource struct {
	k8sClient          k8sclient.Interface
	logger             micrologger.Logger
	organizationReader organization.Reader

	baseDomain   string
	customer     string
	installation string
	insecureCA   bool
	pipeline     string
	provider     cluster.Provider
	region       string
	version      string
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
	if config.BaseDomain == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.BaseDomain must not be empty")
	}
	if config.Customer == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Customer must not be empty")
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Installation must not be empty")
	}
	if config.Pipeline == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Pipeline must not be empty")
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

		baseDomain:   config.BaseDomain,
		customer:     config.Customer,
		installation: config.Installation,
		insecureCA:   config.InsecureCA,
		pipeline:     config.Pipeline,
		provider:     config.Provider,
		region:       config.Region,
		version:      config.Version,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) desiredSecret(ctx context.Context, cluster metav1.Object, name string, namespace string, password string, version string) (*corev1.Secret, error) {
	organization, err := r.organizationReader.Read(ctx, cluster)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	provider, err := key.ClusterProvider(cluster, r.provider)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	globalConfig := remotewriteconfiguration.GlobalConfig{
		RemoteWrite: []remotewriteconfiguration.RemoteWrite{
			remotewriteconfiguration.DefaultRemoteWrite(key.ClusterID(cluster), r.baseDomain, password, r.insecureCA),
		},
		ExternalLabels: map[string]string{
			key.ClusterIDKey:       key.ClusterID(cluster),
			key.ClusterTypeKey:     key.ClusterType(r.installation, cluster),
			key.CustomerKey:        r.customer,
			key.InstallationKey:    r.installation,
			key.OrganizationKey:    organization,
			key.PipelineKey:        r.pipeline,
			key.ProviderKey:        provider,
			key.RegionKey:          r.region,
			key.ServicePriorityKey: key.GetServicePriority(cluster),
		},
	}

	remoteWriteConfig := remotewriteconfiguration.RemoteWriteConfig{
		GlobalConfig: globalConfig,
	}

	marshalledValues, err := yaml.Marshal(remoteWriteConfig)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    key.PrometheusLabels(cluster),
			Finalizers: []string{
				"monitoring.giantswarm.io/prometheus-remote-write",
			},
		},
		Data: map[string][]byte{
			"values": marshalledValues,
		},
		Type: "Opaque",
	}, nil
}

func (r *Resource) createSecret(ctx context.Context, cluster metav1.Object, name string, namespace string, password, version string) error {
	secret, err := r.desiredSecret(ctx, cluster, name, namespace, password, version)
	if err != nil {
		return microerror.Mask(err)
	}

	_, err = r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	return microerror.Mask(err)
}

func (r *Resource) deleteSecret(ctx context.Context, secret *corev1.Secret) error {
	err := r.k8sClient.K8sClient().CoreV1().Secrets(secret.Namespace).Delete(ctx, secret.Name, metav1.DeleteOptions{})
	return microerror.Mask(err)
}
