package ciliumnetpol

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.config.Logger.Debugf(ctx, "deleting")
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

		err = r.config.DynamicK8sClient.Resource(resource).Namespace(desired.GetNamespace()).Delete(ctx, desired.GetName(), metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}
	r.config.Logger.Debugf(ctx, "deleted")

	return nil
}
