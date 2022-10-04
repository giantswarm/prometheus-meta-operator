package remotewriteapiendpointconfigsecret

import (
	"context"
	"fmt"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	pov1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "remotewriteapiendpointconfigsecret"
)

type Config struct {
	K8sClient  k8sclient.Interface
	Logger     micrologger.Logger
	BaseDomain string
}

type remoteWriteValues struct {
	RemoteWrite []pov1.RemoteWriteSpec `json:"remoteWrite"`
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
		Name:      key.RemoteWriteAgentConfigSecretName,
		Namespace: key.ClusterID(cluster),
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

	remoteWrites := []pov1.RemoteWriteSpec{
		{
			URL: fmt.Sprintf("https://%s/%s/api/v1/write", config.BaseDomain, key.ClusterID(cluster)),
			BasicAuth: &pov1.BasicAuth{
				Username: corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: key.RemoteWriteAgentSecretName,
					},
					Key: "username",
				},
				Password: corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: key.RemoteWriteAgentSecretName,
					},
					Key: "password",
				},
			},
		},
	}

	values := remoteWriteValues{RemoteWrite: remoteWrites}
	marshalledValues, err := yaml.Marshal(values)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		Data: map[string][]byte{
			"values": []byte(marshalledValues),
		},
		Type: "Opaque",
	}
	return secret, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
