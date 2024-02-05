package ciliumnetpol

import (
	"context"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "creating")
	{
		desired, err := toCiliumNetworkPolicy(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		ciliumClient, err := ciliumv2.client.NewForConfig(config)
		if err != nil {
			return microerror.Mask(err)
		}

		current, err := ciliumClient.CiliumV2().CiliumNetworkPolicies(key.Namespace(cluster)).Get(ctx, desired.GetName(), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			current, err = r.k8sClient.K8sClient().RbacV1().ClusterRoleBindings().Create(ctx, desired, metav1.CreateOptions{})
		}
		if err != nil {
			return microerror.Mask(err)
		}
	}
	r.logger.Debugf(ctx, "created")

	return nil
}
