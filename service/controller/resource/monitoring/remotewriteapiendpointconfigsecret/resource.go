package remotewriteapiendpointconfigsecret

import (
	"context"
	"fmt"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/password"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name                           = "remotewriteapiendpointconfigsecret"
	remoteWriteEndpointTemplateURL = "https://%s/%s/api/v1/write"
)

type Config struct {
	K8sClient       k8sclient.Interface
	Logger          micrologger.Logger
	PasswordManager password.Manager
	BaseDomain      string
	Customer        string
	Installation    string
	InsecureCA      bool
	Pipeline        string
	Provider        string
	Region          string
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger

	PasswordManager password.Manager
	BaseDomain      string
	Customer        string
	Installation    string
	InsecureCA      bool
	Pipeline        string
	Provider        string
	Region          string
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		PasswordManager: config.PasswordManager,
		BaseDomain:      config.BaseDomain,
		Customer:        config.Customer,
		Installation:    config.Installation,
		InsecureCA:      config.InsecureCA,
		Pipeline:        config.Pipeline,
		Provider:        config.Provider,
		Region:          config.Region,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

type RemoteWrite struct {
	Name        string             `yaml:"name" json:"name"`
	Password    string             `yaml:"password" json:"password"`
	Username    string             `yaml:"username" json:"username"`
	URL         string             `yaml:"url" json:"url"`
	QueueConfig promv1.QueueConfig `yaml:"queueConfig" json:"queueConfig"`
	TLSConfig   promv1.TLSConfig   `yaml:"tlsConfig" json:"tlsConfig"`
}

type GlobalRemoteWriteValues struct {
	Global RemoteWriteValues `yaml:"global" json:"global"`
}

type RemoteWriteValues struct {
	RemoteWrite    []RemoteWrite     `yaml:"remoteWrite" json:"remoteWrite"`
	ExternalLabels map[string]string `yaml:"externalLabels" json:"externalLabels"`
}

func (r *Resource) desiredSecret(cluster metav1.Object, name string, namespace string, password string) (*corev1.Secret, error) {
	url := fmt.Sprintf(remoteWriteEndpointTemplateURL, r.BaseDomain, key.ClusterID(cluster))
	remoteWrites := []RemoteWrite{
		{
			Name:        key.PrometheusMetaOperatorRemoteWriteName,
			URL:         url,
			Username:    key.ClusterID(cluster),
			Password:    password,
			QueueConfig: defaultQueueConfig(),
			TLSConfig: promv1.TLSConfig{
				SafeTLSConfig: promv1.SafeTLSConfig{
					InsecureSkipVerify: r.InsecureCA,
				},
			},
		},
	}

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

	values := RemoteWriteValues{
		RemoteWrite:    remoteWrites,
		ExternalLabels: externalLabels,
	}

	marshalledValues, err := yaml.Marshal(GlobalRemoteWriteValues{values})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    key.PrometheusLabels(cluster),
		},
		Data: map[string][]byte{
			"values": []byte(marshalledValues),
		},
		Type: "Opaque",
	}, nil
}

func (r *Resource) createSecret(ctx context.Context, cluster metav1.Object, name string, namespace string) error {
	r.logger.Debugf(ctx, "generating password for the prometheus agent")
	password, err := r.PasswordManager.GeneratePassword(32)
	if err != nil {
		r.logger.Errorf(ctx, err, "failed to generate the prometheus agent password")
		return microerror.Mask(err)
	}

	secret, err := r.desiredSecret(cluster, name, namespace, password)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Debugf(ctx, "generated password for the prometheus agent")

	_, err = r.k8sClient.K8sClient().CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{})
	return microerror.Mask(err)
}

func (r *Resource) deleteSecret(ctx context.Context, secret *corev1.Secret) error {
	err := r.k8sClient.K8sClient().CoreV1().Secrets(secret.Namespace).Delete(ctx, secret.Name, metav1.DeleteOptions{})
	return microerror.Mask(err)
}

func defaultQueueConfig() promv1.QueueConfig {
	return promv1.QueueConfig{
		Capacity:          30000,
		MaxSamplesPerSend: 150000,
		MaxShards:         10,
	}
}
