package generic

import (
	"context"
	"reflect"

	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	desired, err := r.toCR(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "creating")
	current, err := r.client.Create(desired)
	if apierrors.IsAlreadyExists(err) {
		resetMeta(current)
		if !reflect.DeepEqual(current, desired) {
			_, err = r.client.Update(desired)
		}
	}

	if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "level", "debug", "message", "created")

	return nil
}

func resetMeta(o metav1.Object) {
	var (
		uid  types.UID
		time metav1.Time
	)

	o.SetNamespace("")
	o.SetName("")
	o.SetGenerateName("")
	o.SetUID(uid)
	o.SetResourceVersion("")
	o.SetGeneration(0)
	o.SetSelfLink("")
	o.SetCreationTimestamp(time)
	o.SetDeletionTimestamp(nil)
	o.SetDeletionGracePeriodSeconds(nil)
	o.SetLabels(nil)
	o.SetAnnotations(nil)
	o.SetFinalizers(nil)
	o.SetOwnerReferences(nil)
	o.SetClusterName("")
	o.SetManagedFields(nil)
}
