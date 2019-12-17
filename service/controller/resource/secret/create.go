package secret

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	namespace, err := toNamespace(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	_, err = r.k8sClient.K8sClient().CoreV1().Namespaces().Create(namespace)
	if apierrors.IsAlreadyExists(err) {
		// fall through
	} else if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
