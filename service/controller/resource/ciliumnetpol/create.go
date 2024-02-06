package ciliumnetpol

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "creating")
	{
		desired, err := toCiliumNetworkPolicy(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		var current unstructured.Unstructured
		err = r.k8sClient.CtrlClient().Get(ctx, client.ObjectKey{Name: desired.GetName(), Namespace: desired.GetNamespace()}, &current)
		if apierrors.IsNotFound(err) {
			err = r.k8sClient.CtrlClient().Create(ctx, desired)
		}
		if err != nil {
			return microerror.Mask(err)
		}
	}
	r.logger.Debugf(ctx, "created")

	return nil
}
