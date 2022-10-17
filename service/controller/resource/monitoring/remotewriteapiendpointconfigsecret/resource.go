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
	Name = "remotewriteapiendpointconfigsecret"
)

type Config struct {
	K8sClient       k8sclient.Interface
	Logger          micrologger.Logger
	PasswordManager password.Manager
	BaseDomain      string
	Provider        string
}

type RemoteWrite struct {
	Name        string             `json:"name"`
	Password    string             `json:"password"`
	Username    string             `json:"username"`
	URL         string             `json:"url"`
	QueueConfig promv1.QueueConfig `json:"queueConfig"`
}

type GlobalRemoteWriteValues struct {
	Global RemoteWriteValues `json:"global"`
}

type RemoteWriteValues struct {
	RemoteWrite []RemoteWrite `json:"remoteWrite"`
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
			return getObjectMeta(ctx, v, config.Provider)
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

func getObjectMeta(ctx context.Context, v interface{}, provider string) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	name, namespace := key.RemoteWriteAPIEndpointConfigSecretNameAndNamespace(cluster, provider)

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    key.PrometheusLabels(cluster),
	}, nil
}

func toSecret(ctx context.Context, v interface{}, config Config) (*corev1.Secret, error) {
	objectMeta, err := getObjectMeta(ctx, v, config.Provider)
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
			URL:         fmt.Sprintf("https://%s/%s/api/v1/write", config.BaseDomain, key.ClusterID(cluster)),
			Username:    key.ClusterID(cluster),
			Password:    password,
			QueueConfig: defaultQueueConfig(),
		},
	}

	values := RemoteWriteValues{RemoteWrite: remoteWrites}
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
		Capacity:          10000,
		MaxSamplesPerSend: 1000,
		MaxShards:         50,
	}
}

func hasChanged(current, desired metav1.Object) bool {
	return false
}
