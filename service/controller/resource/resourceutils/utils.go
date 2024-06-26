package resourceutils

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func UpdateMeta(c, d metav1.Object) {
	d.SetGenerateName(c.GetGenerateName())
	d.SetUID(c.GetUID())
	d.SetResourceVersion(c.GetResourceVersion())
	d.SetGeneration(c.GetGeneration())
	d.SetSelfLink(c.GetSelfLink())
	d.SetCreationTimestamp(c.GetCreationTimestamp())
	d.SetDeletionTimestamp(c.GetDeletionTimestamp())
	d.SetDeletionGracePeriodSeconds(c.GetDeletionGracePeriodSeconds())
	// without this, it's impossible to change labels on resources
	if len(d.GetLabels()) == 0 {
		d.SetLabels(c.GetLabels())
	}
	// without this, it's impossible to change annotations on resources
	if len(d.GetAnnotations()) == 0 {
		d.SetAnnotations(c.GetAnnotations())
	}
	d.SetFinalizers(c.GetFinalizers())
	d.SetOwnerReferences(c.GetOwnerReferences())
	d.SetManagedFields(c.GetManagedFields())
}
