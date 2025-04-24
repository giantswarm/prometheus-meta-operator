package etcdcertificates

import (
	"context"
	"reflect"

	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
)

const (
	Name = "etcd-certificates"
)

type Config struct {
	Installation string

	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
	Provider  cluster.Provider

	MimirEnabled bool
}

// secretCopier provides a way to create a new secret from different data source.
type secretCopier struct {
	logger       micrologger.Logger
	clientFunc   func(string) generic.Interface
	k8sClient    k8sclient.Interface
	installation string
}

func New(config Config) (*generic.Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}

	// Wrapping the secret client into a generic interface.
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().CoreV1().Secrets(namespace)
		return wrappedClient{client: c}
	}

	sc := secretCopier{
		logger:       config.Logger,
		clientFunc:   clientFunc,
		k8sClient:    config.K8sClient,
		installation: config.Installation,
	}

	c := generic.Config{
		ClientFunc: clientFunc,
		Logger:     config.Logger,
		Name:       Name,
		GetObjectMeta: func(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
			return getObjectMeta(v, config.Installation)
		},
		GetDesiredObject: func(ctx context.Context, v interface{}) (metav1.Object, error) {
			return sc.ToSecret(ctx, v, config)
		},
		HasChangedFunc:       hasChanged,
		DeleteIfMimirEnabled: config.MimirEnabled,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

// hasChanged determines if secret data have changed.
func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
