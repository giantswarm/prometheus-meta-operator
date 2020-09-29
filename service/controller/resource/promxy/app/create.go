package app

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	objectMeta, err := r.getObjectMeta(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("checking if %s app exists", objectMeta.Name))
	app, err := r.k8sClient.G8sClient().ApplicationV1alpha1().Apps(objectMeta.Namespace).Get(ctx, objectMeta.Name, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	app.Spec.Config.ConfigMap.Name = key.PromxyConfigMapName()
	app.Spec.Config.ConfigMap.Namespace = key.PromxyConfigMapNamespace()

	r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("configuring the server groups configmap for app %s", objectMeta.Name))
	_, err = r.k8sClient.G8sClient().ApplicationV1alpha1().Apps(objectMeta.Namespace).Update(ctx, app, metav1.UpdateOptions{})
	if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("Configured the server groups configmap for app %s", objectMeta.Name))

	return nil
}
