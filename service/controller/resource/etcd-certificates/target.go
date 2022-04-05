package etcdcertificates

import (
	"context"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

// getObjectMeta returns the target secret metadata.
func getObjectMeta(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.EtcdSecret(v),
		Namespace: key.Namespace(cluster),
	}, nil
}

// ToSecret returns the target secret by combining results of getObjectMeta and getSource.
func (sc *secretCopier) ToSecret(ctx context.Context, v interface{}) (metav1.Object, error) {
	objectMeta, err := getObjectMeta(ctx, v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	data, err := sc.getSource(ctx, v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		StringData: data,
	}

	return secret, nil
}