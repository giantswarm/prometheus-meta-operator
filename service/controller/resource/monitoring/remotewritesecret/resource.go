package remotewritesecret

import (
	"context"
	"reflect"

	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/password"
	remotewriteconfiguration "github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewrite/configuration"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "remotewritesecret"
)

type Config struct {
	K8sClient       k8sclient.Interface
	Logger          micrologger.Logger
	PasswordManager password.Manager
	BaseDomain      string
	InsecureCA      bool
	Installation    string
	Provider        cluster.Provider
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger

	PasswordManager password.Manager
	BaseDomain      string
	InsecureCA      bool
	Installation    string
	Provider        cluster.Provider
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.K8sClient must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.PasswordManager == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.PasswordManager must not be empty")
	}
	if config.BaseDomain == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.BaseDomain must not be empty")
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.Installation must not be empty")
	}
	if reflect.ValueOf(config.Provider).IsZero() {
		return nil, microerror.Maskf(invalidConfigError, "config.Provider must not be empty")
	}
	r := &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,

		PasswordManager: config.PasswordManager,
		BaseDomain:      config.BaseDomain,
		InsecureCA:      config.InsecureCA,
		Installation:    config.Installation,
		Provider:        config.Provider,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) desiredSecret(cluster metav1.Object, name string, namespace string, password string) (*corev1.Secret, error) {
	remoteWriteConfig := remotewriteconfiguration.RemoteWriteConfig{
		PrometheusAgentConfig: remotewriteconfiguration.PrometheusAgentConfig{
			RemoteWrite: []remotewriteconfiguration.RemoteWrite{
				remotewriteconfiguration.DefaultRemoteWrite(key.ClusterID(cluster), r.BaseDomain, password, r.InsecureCA),
			},
		},
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
