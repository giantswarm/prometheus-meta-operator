package volumeresizehack

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "checking if pvc need to be re-created")

	namespace := key.Namespace(cluster)

	// Get StS
	stsName := fmt.Sprintf("prometheus-%s", cluster.GetName())
	currentStS, err := r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).Get(ctx, stsName, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	if len(currentStS.Spec.VolumeClaimTemplates) < 1 {
		// No PVC template found in StS, nothing to resize. Skip this resource.
		r.logger.LogCtx(ctx, "level", "debug", "message", "skipping, no pvc found in sts")
		return nil
	}

	// Get PVC
	index := 0
	desiredPVC := currentStS.Spec.VolumeClaimTemplates[index]
	pvcName := fmt.Sprintf("%s-%s-%d", desiredPVC.GetName(), currentStS.GetName(), index)
	currentPVC, err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "pvc is missing, need to be re-created")
	} else {
		if err != nil {
			return microerror.Mask(err)
		}
		if !reflect.DeepEqual(desiredPVC.Spec.Resources.Requests, currentPVC.Spec.Resources.Requests) {
			r.logger.LogCtx(ctx, "level", "debug", "message", "pvc has wrong size, need to be re-created")

			// delete pvc
			{
				err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
				if err != nil {
					return microerror.Mask(err)
				}
				r.logger.LogCtx(ctx, "level", "debug", "message", "PVC DELETED")
			}

			// scale down
			{
				r.logger.LogCtx(ctx, "level", "debug", "message", "SCALING DOWN")
				*currentStS.Spec.Replicas = 0
				_, err := r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).Update(ctx, currentStS, metav1.UpdateOptions{})
				if err != nil {
					return microerror.Mask(err)
				}
				r.logger.LogCtx(ctx, "level", "debug", "message", "SCALED DOWN")
			}

			// wait 30s for pvc gone
			{
				r.logger.LogCtx(ctx, "level", "debug", "message", "WAITING PVC")
				o := func() error {
					_, err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
					if apierrors.IsNotFound(err) {
						return nil
					}

					return microerror.Mask(pvcExist)
				}
				b := backoff.NewMaxRetries(6, 5*time.Second)
				err := backoff.Retry(o, b)
				if err != nil {
					return microerror.Mask(err)
				}
				r.logger.LogCtx(ctx, "level", "debug", "message", "WAITED PVC")
			}
		} else {
			r.logger.LogCtx(ctx, "level", "debug", "message", "pvc do not need to be re-created")
			return nil
		}
	}

	// scale down again
	r.logger.LogCtx(ctx, "level", "debug", "message", "SCALING DOWN AGAIN")
	currentStS, err = r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).Get(ctx, stsName, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}
	*currentStS.Spec.Replicas = 0
	_, err = r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).Update(ctx, currentStS, metav1.UpdateOptions{})
	if err != nil {
		return microerror.Mask(err)
	}
	r.logger.LogCtx(ctx, "level", "debug", "message", "SCALED DOWN AGAIN")

	r.logger.LogCtx(ctx, "level", "debug", "message", "pvc re-created")

	return nil
}
