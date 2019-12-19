package service

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	service, err := toService(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.k8sClient.K8sClient().CoreV1().Services(service.GetNamespace()).Delete(service.GetName(), &metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		// fall through
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
