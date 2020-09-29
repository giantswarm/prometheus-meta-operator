package configmap

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	objectMeta, err := r.getObjectMeta(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	message := fmt.Sprintf("checking if %s configmap already exists ", objectMeta.Name)
	r.logger.LogCtx(ctx, "level", "debug", "message", message)
	_, err = r.k8sClient.K8sClient().CoreV1().ConfigMaps(objectMeta.Namespace).Get(ctx, objectMeta.Name, metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("configmap %s does not exists", objectMeta.Name))

		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("creating configmap %s", objectMeta.Name))
		configmap, err := r.toConfigMap(objectMeta)
		if err != nil {
			return microerror.Mask(err)
		}
		_, err = r.k8sClient.K8sClient().CoreV1().ConfigMaps(objectMeta.Namespace).Create(ctx, configmap, metav1.CreateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("created configmap %s", objectMeta.Name))

	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
