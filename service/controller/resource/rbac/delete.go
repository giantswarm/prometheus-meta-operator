package rbac

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.config.Logger.Debugf(ctx, "deleting")
	{
		desired, err := toClusterRoleBinding(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		err = r.config.K8sClient.K8sClient().RbacV1().ClusterRoleBindings().Delete(ctx, desired.GetName(), metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// fall through
		} else if err != nil {
			return microerror.Mask(err)
		}
	}
	r.config.Logger.Debugf(ctx, "deleted")

	return nil
}
