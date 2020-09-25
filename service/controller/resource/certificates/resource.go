package certificates

import (
	"context"
	"reflect"

	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

type Config struct {
	Name                string
	SourceNameFunc      NameFunc
	SourceNamespaceFunc NameFunc
	TargetNameFunc      NameFunc
	K8sClient           k8sclient.Interface
	Logger              micrologger.Logger
}

type NameFunc func(metav1.Object) string

type Resource struct {
	name                string
	sourceNameFunc      NameFunc
	sourceNamespaceFunc NameFunc
	targetNameFunc      NameFunc
	k8sClient           k8sclient.Interface
	logger              micrologger.Logger
}

func New(config Config) (*Resource, error) {
	if config.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Name must not be empty", config)
	}
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.SourceNameFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.SourceNameFunc must not be empty", config)
	}
	if config.SourceNamespaceFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.SourceNamespaceFunc must not be empty", config)
	}
	if config.TargetNameFunc == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.TargetNameFunc must not be empty", config)
	}

	r := &Resource{
		name:                config.Name,
		logger:              config.Logger,
		k8sClient:           config.K8sClient,
		sourceNameFunc:      config.SourceNameFunc,
		sourceNamespaceFunc: config.SourceNamespaceFunc,
		targetNameFunc:      config.TargetNameFunc,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return r.name
}

func (r *Resource) getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      r.targetNameFunc(cluster),
		Namespace: key.Namespace(cluster),
	}, nil
}

func (r *Resource) getDesiredObject(v interface{}) (*corev1.Secret, error) {
	objectMeta, err := r.getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	sourceSecret, err := r.getSource(context.TODO(), v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		Data:       sourceSecret.Data,
		Type:       sourceSecret.Type,
	}

	return secret, nil
}

// getSource returns the Secret to be copied, i.e. default/$CLUSTER_ID-prometheus
func (r *Resource) getSource(ctx context.Context, v interface{}) (*corev1.Secret, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secretName := r.sourceNameFunc(cluster)
	secretNamespace := r.sourceNamespaceFunc(cluster)

	s, err := r.k8sClient.K8sClient().CoreV1().Secrets(secretNamespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return s, nil
}

func (r *Resource) hasChanged(current, desired *corev1.Secret) bool {
	return !reflect.DeepEqual(current.Data, desired.Data)
}
