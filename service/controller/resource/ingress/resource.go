package ingress

import (
	"reflect"

	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Name = "ingress"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger

	BaseDomain string
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger

	baseDomain string
}

func New(config Config) (*generic.Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.BaseDomain == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.BaseDomain must not be empty", config)
	}

	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().ExtensionsV1beta1().Ingresses(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:       clientFunc,
		Logger:           config.Logger,
		Name:             Name,
		GetObjectMeta:    getObjectMeta,
		GetDesiredObject: toIngress,
		HasChangedFunc:   hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, err
	}

	return metav1.ObjectMeta{
		Name:      cluster.GetName(),
		Namespace: key.Namespace(cluster),
	}, nil
}

func toIngress(v interface{}) (metav1.Object, error) {
	return &v1beta1.Ingress{}, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*v1beta1.Ingress)
	d := desired.(*v1beta1.Ingress)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
