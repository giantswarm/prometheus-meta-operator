package service

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	service, err := toService(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "creating service")
	_, err = r.k8sClient.K8sClient().CoreV1().Services(service.GetNamespace()).Create(service)
	if apierrors.IsAlreadyExists(err) {
		// fall through
	} else if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "created service")

	return nil
}
