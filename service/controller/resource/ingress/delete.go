package ingress

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	ingress, err := r.toIngress(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.k8sClient.K8sClient().ExtensionsV1beta1().Ingresses(ingress.GetNamespace()).Delete(ingress.GetName(), &metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		// fall through
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
