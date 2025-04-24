package namespace

import (
	"context"

	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "namespace"
)

type Config struct {
	K8sClient    k8sclient.Interface
	Logger       micrologger.Logger
	MimirEnabled bool
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().CoreV1().Namespaces()
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:           clientFunc,
		Logger:               config.Logger,
		Name:                 Name,
		GetObjectMeta:        getObjectMeta,
		GetDesiredObject:     toNamespace,
		HasChangedFunc:       hasChanged,
		DeleteIfMimirEnabled: config.MimirEnabled,
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
		Name:   key.Namespace(cluster),
		Labels: key.PrometheusLabels(cluster),
	}, nil
}

func toNamespace(ctx context.Context, v interface{}) (metav1.Object, error) {
	objectMeta, err := getObjectMeta(ctx, v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	namespace := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.Version,
			Kind:       "Namespace",
		},
		ObjectMeta: objectMeta,
	}

	return namespace, nil
}

func hasChanged(current, desired metav1.Object) bool {
	return false
}
