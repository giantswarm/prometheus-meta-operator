package namespace

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	desired, err := r.toNamespace(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.config.Logger.Debugf(ctx, "creating")
	_, err = r.config.K8sClient.K8sClient().CoreV1().Namespaces().Get(ctx, desired.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = r.config.K8sClient.K8sClient().CoreV1().Namespaces().Create(ctx, desired, metav1.CreateOptions{})
	}
	if err != nil {
		return microerror.Mask(err)
	}
	r.config.Logger.Debugf(ctx, "created")

	return nil
}
