package certificates

import (
	"context"
	"fmt"
	"reflect"

	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

type Config struct {
	Name      string
	Sources   []CertificateSource
	Target    NameFunc
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

type NameFunc func(metav1.Object) string

type CertificateSource struct {
	NameFunc      NameFunc
	NamespaceFunc NameFunc
}

type Resource struct {
	name      string
	sources   []CertificateSource
	target    NameFunc
	k8sClient k8sclient.Interface
	logger    micrologger.Logger
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
	if len(config.Sources) == 0 {
		return nil, microerror.Maskf(invalidConfigError, "%T.Sources must not be empty", config)
	}
	if config.Target == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Target must not be empty", config)
	}

	r := &Resource{
		name:      config.Name,
		logger:    config.Logger,
		k8sClient: config.K8sClient,
		sources:   config.Sources,
		target:    config.Target,
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
		Name:      r.target(cluster),
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

	var secret *v1.Secret
	for _, source := range r.sources {
		secretName := source.NameFunc(cluster)
		secretNamespace := source.NamespaceFunc(cluster)

		secret, err = r.k8sClient.K8sClient().CoreV1().Secrets(secretNamespace).Get(ctx, secretName, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	if secret == nil {
		err := fmt.Errorf("No certificates found for %s", key.ClusterID(cluster))
		return nil, microerror.Mask(err)
	}

	return secret, nil
}

func (r *Resource) hasChanged(current, desired *corev1.Secret) bool {
	return !reflect.DeepEqual(current.Data, desired.Data)
}
