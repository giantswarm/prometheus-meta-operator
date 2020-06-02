package certificates

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	targetSecret, err := toTargetSecret(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	sourceSecret, err := toSourceSecret(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "creating certificates")
	_, err = r.k8sClient.K8sClient().CoreV1().Secrets(targetSecret.GetNamespace()).Get(targetSecret.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		sourceSecret, err := r.k8sClient.K8sClient().CoreV1().Secrets(sourceSecret.GetNamespace()).Get(sourceSecret.GetName(), metav1.GetOptions{})
		if err != nil {
			return microerror.Mask(err)
		}

		targetSecret.Data = sourceSecret.Data
		_, err = r.k8sClient.K8sClient().CoreV1().Secrets(targetSecret.GetNamespace()).Create(targetSecret)
		if err != nil {
			return microerror.Mask(err)
		}
	} else if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "created certificates")

	return nil
}
