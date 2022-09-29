package remotewriteapiendpointsecret

import (
	"context"
	"fmt"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/password"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "remotewriteapiendpointsecret"
)

type Config struct {
	K8sClient       k8sclient.Interface
	Logger          micrologger.Logger
	PasswordManager password.Manager
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
		Name:      key.RemoteWriteAgentSecretName,
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
	username := key.ClusterID(cluster)

	config.Logger.Debugf(ctx, "generating password for the prometheus agent")
	password, err := config.PasswordManager.GeneratePassword(32)
	if err != nil {
		config.Logger.Errorf(ctx, err, "failed to generate the prometheus agent password")
		return nil, microerror.Mask(err)
	}

	hashedPassword, err := config.PasswordManager.Hash(password)
	if err != nil {
		config.Logger.Errorf(ctx, err, "failed to hash the prometheus agent password")
		return nil, microerror.Mask(err)
	}
	config.Logger.Debugf(ctx, "generate password for the prometheus agent")

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		Data: map[string][]byte{
			// To be used by the ingress basic auth
			"auth":     []byte(fmt.Sprintf("%s:%s", username, hashedPassword)),
			"password": []byte(password),
			"username": []byte(username),
		},
		Type: "Opaque",
	}
	return secret, nil
}

func hasChanged(current, desired metav1.Object) bool {
	return false
}
