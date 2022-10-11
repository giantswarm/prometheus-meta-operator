package remotewriteingressauth

import (
	"context"
	"fmt"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/password"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteapiendpointconfigsecret"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "remotewriteingressauth"
)

type Config struct {
	K8sClient       k8sclient.Interface
	Logger          micrologger.Logger
	PasswordManager password.Manager
	Provider        string
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().CoreV1().Secrets(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:    clientFunc,
		Logger:        config.Logger,
		Name:          Name,
		GetObjectMeta: getObjectMeta,
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

func getObjectMeta(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.RemoteWriteIngressAuthSecretName,
		Namespace: key.Namespace(cluster),
		Labels:    key.PrometheusLabels(cluster),
	}, nil
}

func toSecret(ctx context.Context, v interface{}, config Config) (*corev1.Secret, error) {
	objectMeta, err := getObjectMeta(ctx, v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secretName, secretNamespace := key.RemoteWriteAPIEndpointConfigSecretNameAndNamespace(cluster, config.Provider)

	apiEndpointSecret, err := config.K8sClient.K8sClient().CoreV1().Secrets(secretNamespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		config.Logger.Errorf(ctx, err, "failed to get api endpoint secret")
		return nil, microerror.Mask(err)
	}

	username, password, err := extractUsernameAndPasswordFromSecret(apiEndpointSecret)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	config.Logger.Debugf(ctx, "hashing password for the prometheus agent")
	hashedPassword, err := config.PasswordManager.Hash([]byte(password))
	if err != nil {
		config.Logger.Errorf(ctx, err, "failed to hash the prometheus agent password")
		return nil, microerror.Mask(err)
	}

	config.Logger.Debugf(ctx, "hashed password for the prometheus agent")

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		Data: map[string][]byte{
			"auth": []byte(fmt.Sprintf("%s:%s", username, hashedPassword)),
		},
		Type: "Opaque",
	}
	return secret, nil
}

func extractUsernameAndPasswordFromSecret(secret *corev1.Secret) (string, string, error) {
	data := secret.Data["values"]
	remoteWriteValues := remotewriteapiendpointconfigsecret.GlobalRemoteWriteValues{}
	err := yaml.Unmarshal(data, &remoteWriteValues)
	if err != nil {
		return "", "", microerror.Mask(err)
	}
	if len(remoteWriteValues.Global.RemoteWrite) == 0 {
		// skipping
		return "", "", secretNotFound
	}

	username := remoteWriteValues.Global.RemoteWrite[0].Username
	password := remoteWriteValues.Global.RemoteWrite[0].Password

	return username, password, nil
}

func hasChanged(current, desired metav1.Object) bool {
	return false
}
