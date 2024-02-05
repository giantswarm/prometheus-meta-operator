package ciliumnetpol

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "creating")
	{
		desired, err := toCiliumNetworkPolicy(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		current, err := r.k8sClient.K8sClient().Get(ctx, desired.GetName(), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			current, err = r.k8sClient.K8sClient().Create(ctx, desired, metav1.CreateOptions{})
		}
		if err != nil {
			return microerror.Mask(err)
		}
	}
	r.logger.Debugf(ctx, "created")

	return nil
}
