package certificates

import (
	"fmt"
	"reflect"

	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "certificates"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

// secretCopier provides a `ToCR` method which copies data from the source
// cluster secret CR
type secretCopier struct {
	clientFunc func(string) generic.Interface
}

func New(config Config) (*generic.Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	clientFunc := func(namespace string) generic.Interface {
		c := config.K8sClient.K8sClient().CoreV1().Secrets(namespace)
		return wrappedClient{client: c}
	}

	sc := secretCopier{clientFunc: clientFunc}

	c := generic.Config{
		ClientFunc:     clientFunc,
		Logger:         config.Logger,
		Name:           Name,
		ToCR:           sc.ToCR,
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func (sc *secretCopier) ToCR(v interface{}) (metav1.Object, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	sourceSecret, err := sc.getSource(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Secret(),
			Namespace: key.Namespace(cluster),
		},
		Data: sourceSecret.Data,
	}

	return secret, nil
}

// getSource returns the Secret to be copied, i.e. default/$CLUSTER_ID-prometheus
func (sc *secretCopier) getSource(v interface{}) (*corev1.Secret, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	secretName := fmt.Sprintf("%s-prometheus", cluster.GetName())

	s, err := sc.clientFunc("default").Get(secretName, metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return s.(*corev1.Secret), nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
