package ciliumnetpol

import (
	"context"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "creating")
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

		current, err := r.dynamicK8sClient.Resource(resource).Namespace(desired.GetNamespace()).Get(ctx, desired.GetName(), metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			current, err = r.dynamicK8sClient.Resource(resource).Namespace(desired.GetNamespace()).Create(ctx, desired, metav1.CreateOptions{})
		}
		if err != nil {
			return microerror.Mask(err)
		}

		if hasCiliumNetworkPolicyChanged(current, desired) {
			updateMeta(current, desired)
			_, err = r.dynamicK8sClient.Resource(resource).Namespace(desired.GetNamespace()).Update(ctx, desired, metav1.UpdateOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}
	r.logger.Debugf(ctx, "created")

	return nil
}

func updateMeta(c, d metav1.Object) {
	d.SetGenerateName(c.GetGenerateName())
	d.SetUID(c.GetUID())
	d.SetResourceVersion(c.GetResourceVersion())
	d.SetGeneration(c.GetGeneration())
	d.SetSelfLink(c.GetSelfLink())
	d.SetCreationTimestamp(c.GetCreationTimestamp())
	d.SetDeletionTimestamp(c.GetDeletionTimestamp())
	d.SetDeletionGracePeriodSeconds(c.GetDeletionGracePeriodSeconds())
	d.SetLabels(c.GetLabels())
	d.SetAnnotations(c.GetAnnotations())
	d.SetFinalizers(c.GetFinalizers())
	d.SetOwnerReferences(c.GetOwnerReferences())
	d.SetManagedFields(c.GetManagedFields())
}
