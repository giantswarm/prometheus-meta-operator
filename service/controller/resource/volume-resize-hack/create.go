package volumeresizehack

import (
	"context"
	"fmt"
	"reflect"

	"github.com/giantswarm/microerror"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	stsName := fmt.Sprintf("prometheus-%s", cluster.GetName())
	r.logger.LogCtx(ctx, "level", "debug", "message", "creating")
	currentStS, err := r.k8sClient.K8sClient().AppsV1().StatefulSets(key.Namespace(cluster)).Get(ctx, stsName, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	// TODO find by name : prometheus-df3yn-db-prometheus-df3yn-0
	desiredPVC := currentStS.Spec.VolumeClaimTemplates[0]

	currentPVC, err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(desiredPVC.GetNamespace()).Get(ctx, desiredPVC.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "PVC NOT FOUND")
		if *currentStS.Spec.Replicas <= 0 {
			r.logger.LogCtx(ctx, "level", "debug", "message", "SCALING UP")

			scale := &autoscalingv1.Scale{
				Spec: autoscalingv1.ScaleSpec{
					Replicas: 1,
				},
			}
			_, err = r.k8sClient.K8sClient().AppsV1().StatefulSets(key.Namespace(cluster)).UpdateScale(ctx, stsName, scale, metav1.UpdateOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
			r.logger.LogCtx(ctx, "level", "debug", "message", "SCALED UP")
		} else {
			return microerror.Mask(err)
		}
	} else {
		if err != nil {
			return microerror.Mask(err)
		}

		if !reflect.DeepEqual(desiredPVC, currentPVC) {
			r.logger.LogCtx(ctx, "level", "debug", "message", "PVC DO NOT MATCH; DELETING")
			err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(desiredPVC.GetNamespace()).Delete(ctx, desiredPVC.GetName(), metav1.DeleteOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
			r.logger.LogCtx(ctx, "level", "debug", "message", "PVC DELETED; SCALING DOWN")

			scale := &autoscalingv1.Scale{
				Spec: autoscalingv1.ScaleSpec{
					Replicas: 0,
				},
			}
			_, err = r.k8sClient.K8sClient().AppsV1().StatefulSets(key.Namespace(cluster)).UpdateScale(ctx, stsName, scale, metav1.UpdateOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
			r.logger.LogCtx(ctx, "level", "debug", "message", "SCALED DOWN")
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", "PVC MATCH")
		}
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "created")

	return nil
}
