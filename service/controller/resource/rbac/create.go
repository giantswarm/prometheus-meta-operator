package rbac

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/resourceutils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "creating")
	{
		desired, err := toClusterRoleBinding(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		current, err := r.k8sClient.K8sClient().RbacV1().ClusterRoleBindings().Get(ctx, desired.GetName(), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			current, err = r.k8sClient.K8sClient().RbacV1().ClusterRoleBindings().Create(ctx, desired, metav1.CreateOptions{})
		}
		if err != nil {
			return microerror.Mask(err)
		}

		if hasClusterRoleBindingChanged(current, desired) {
			resourceutils.UpdateMeta(current, desired)
			_, err = r.k8sClient.K8sClient().RbacV1().ClusterRoleBindings().Update(ctx, desired, metav1.UpdateOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}
	r.logger.Debugf(ctx, "created")

	return nil
}
