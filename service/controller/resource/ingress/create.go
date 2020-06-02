package ingress

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	ingress, err := r.toIngress(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "creating ingress")
	_, err = r.k8sClient.K8sClient().ExtensionsV1beta1().Ingresses(ingress.GetNamespace()).Create(ingress)
	if apierrors.IsAlreadyExists(err) {
		// fall through
	} else if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "created ingress")

	return nil
}
