package etcdcertificates

import (
	"context"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "etcd-certificates"
)

type Config struct {
	K8sClient    k8sclient.Interface
	Logger       micrologger.Logger
	Installation string
	Provider     cluster.Provider
}

type Resource struct {
	config Config
}

func (r *Resource) Name() string {
	return Name
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Installation == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Installation must not be empty", config)
	}

	return &Resource{config}, nil
}

// getObjectMeta returns the target secret metadata.
func (r *Resource) getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.EtcdSecret(r.config.Installation, v),
		Namespace: key.Namespace(cluster),
	}, nil
}

// ToSecret returns the target secret by combining results of getObjectMeta and getSource.
func (r *Resource) toSecret(ctx context.Context, v interface{}) (*corev1.Secret, error) {
	objectMeta, err := r.getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	data, err := r.getSource(ctx)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		StringData: data,
	}

	return secret, nil
}

// hasChanged determines if secret data have changed.
func (r *Resource) hasChanged(current, desired metav1.Object) bool {
	c := current.(*corev1.Secret)
	d := desired.(*corev1.Secret)

	return !reflect.DeepEqual(c.Data, d.Data)
}
