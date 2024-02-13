package ciliumnetpol

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "creating")
	{
		resource := schema.GroupVersionResource{
			Group:    "cilium.io",
			Version:  "v2",
			Resource: "ciliumnetworkpolicies",
		}

		desired, err := toCiliumNetworkPolicy(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		ciliumnetpol, err := r.dynamicK8sClient.Resource(resource).Get(ctx, desired.GetName(), metav1.GetOptions{})
		r.logger.Debugf(ctx, "ciliumnetpol: %#v", ciliumnetpol)
		if apierrors.IsNotFound(err) {
			_, err = r.dynamicK8sClient.Resource(resource).Namespace(desired.GetNamespace()).Create(ctx, desired, metav1.CreateOptions{})
		}
		if err != nil {
			return microerror.Mask(err)
		}
	}
	r.logger.Debugf(ctx, "created")

	return nil
}
