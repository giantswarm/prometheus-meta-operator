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
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
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

type RemoteWrite struct {
	Name        string             `json:"name"`
	Password    string             `json:"password"`
	Username    string             `json:"username"`
	URL         string             `json:"url"`
	QueueConfig promv1.QueueConfig `json:"queueConfig"`
	TLSConfig   promv1.TLSConfig   `json:"tlsConfig"`
}

type GlobalRemoteWriteValues struct {
	Global RemoteWriteValues `json:"global"`
}

type RemoteWriteValues struct {
	RemoteWrite    []RemoteWrite     `json:"remoteWrite"`
	ExternalLabels map[string]string `json:"externalLabels"`
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().CoreV1().Secrets(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc: clientFunc,
		Logger:     config.Logger,
		Name:       Name,
		GetObjectMeta: func(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
			return getObjectMeta(ctx, v, config.Installation, config.Provider)
		},
		GetDesiredObject: func(ctx context.Context, v interface{}) (metav1.Object, error) {
			return toSecret(ctx, v, config)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(ctx context.Context, v interface{}, installation string, provider string) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	name, namespace := key.RemoteWriteAPIEndpointConfigSecretNameAndNamespace(cluster, installation, provider)

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    key.PrometheusLabels(cluster),
	}, nil
}

func toSecret(ctx context.Context, v interface{}, config Config) (*corev1.Secret, error) {
	objectMeta, err := getObjectMeta(ctx, v, config.Installation, config.Provider)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	config.Logger.Debugf(ctx, "generating password for the prometheus agent")
	password, err := config.PasswordManager.GeneratePassword(32)
	if err != nil {
		config.Logger.Errorf(ctx, err, "failed to generate the prometheus agent password")
		return nil, microerror.Mask(err)
	}

	config.Logger.Debugf(ctx, "generate password for the prometheus agent")

	remoteWrites := []RemoteWrite{
		{
			Name:        key.PrometheusMetaOperatorRemoteWriteName,
			URL:         fmt.Sprintf(remoteWriteEndpointTemplateURL, config.BaseDomain, key.ClusterID(cluster)),
			Username:    key.ClusterID(cluster),
			Password:    password,
			QueueConfig: defaultQueueConfig(),
			TLSConfig: promv1.TLSConfig{
				SafeTLSConfig: promv1.SafeTLSConfig{
					InsecureSkipVerify: config.InsecureCA,
				},
			},
		},
	}

	externalLabels := map[string]string{
		key.ClusterIDKey:       key.ClusterID(cluster),
		key.ClusterTypeKey:     key.ClusterType(config.Installation, cluster),
		key.CustomerKey:        config.Customer,
		key.InstallationKey:    config.Installation,
		key.OrganizationKey:    key.GetOrganization(cluster),
		key.PipelineKey:        config.Pipeline,
		key.ProviderKey:        config.Provider,
		key.RegionKey:          config.Region,
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

	var immutable bool = true
	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		Data: map[string][]byte{
			"values": []byte(marshalledValues),
		},
		Type:      "Opaque",
		Immutable: &immutable,
	}
	return secret, nil
}

func defaultQueueConfig() promv1.QueueConfig {
	return promv1.QueueConfig{
		Capacity:          30000,
		MaxSamplesPerSend: 10000,
		MaxShards:         10,
	}
}

func hasChanged(current, desired metav1.Object) bool {
	return false
}
