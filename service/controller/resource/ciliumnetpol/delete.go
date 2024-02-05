package ciliumnetpol

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "deleting")
	{
		desired, err := toCiliumNetworkPolicy(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		ciliumClient, err := ciliumv2.client.NewForConfig(config)
		if err != nil {
			return microerror.Mask(err)
		}

		err = ciliumClient.CiliumV2().CiliumNetworkPolicies(key.Namespace(cluster)).Delete(ctx, desired.GetName(), metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}
	r.logger.Debugf(ctx, "deleted")

	return nil
}
