package namespace

import (
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "namespace"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

type wrappedClient struct {
	client v1.NamespaceInterface
}

func (c wrappedClient) Create(o metav1.Object) (metav1.Object, error) {
	return c.client.Create(o.(*corev1.Namespace))
}
func (c wrappedClient) Update(o metav1.Object) (metav1.Object, error) {
	return c.client.Update(o.(*corev1.Namespace))
}
func (c wrappedClient) Get(name string, options metav1.GetOptions) (metav1.Object, error) {
	return c.client.Get(name, options)
}
func (c wrappedClient) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(name, options)
}

func New(config Config) (*generic.Resource, error) {
	c := generic.Config{
		Client: wrappedClient{client: config.K8sClient.K8sClient().CoreV1().Namespaces()},
		Logger: config.Logger,
		Name:   Name,
		ToCR:   toNamespace,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func toNamespace(v interface{}) (metav1.Object, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: key.Namespace(cluster),
		},
	}

	return namespace, nil
}
