package namespace

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureDeleted(ctx context.Context, obj interface{}) error {
	namespace, err := toNamespace(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "deleting namespace")
	err = r.k8sClient.K8sClient().CoreV1().Namespaces().Delete(namespace.GetName(), &metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		// fall through
	} else if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "level", "debug", "message", "deleted namespace")

	return nil
}
