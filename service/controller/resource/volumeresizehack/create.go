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

	r.logger.Debugf(ctx, "checking if pvc need to be re-created")

	namespace := key.Namespace(cluster)

	// get sts
	stsName := key.PrometheusSTSName(cluster)
	currentStS, err := r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).Get(ctx, stsName, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	if len(currentStS.Spec.VolumeClaimTemplates) < 1 {
		// No pvc template found in sts, nothing to resize. Skip this resource.
		r.logger.Debugf(ctx, "skipping, no pvc found in sts volumeclaimtemplates")
		return nil
	}

	// get pvc
	index := 0
	desiredPVC := currentStS.Spec.VolumeClaimTemplates[index]
	pvcName := fmt.Sprintf("%s-%s-%d", desiredPVC.GetName(), currentStS.GetName(), index)
	currentPVC, err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		r.logger.Debugf(ctx, "pvc is missing, need to be re-created")
	} else {
		if err != nil {
			return microerror.Mask(err)
		}
		if !reflect.DeepEqual(desiredPVC.Spec.Resources.Requests, currentPVC.Spec.Resources.Requests) {
			r.logger.Debugf(ctx, "pvc has wrong size, need to be re-created")

			// delete pvc
			{
				err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
				if err != nil {
					return microerror.Mask(err)
				}
			}

			// scale down sts
			{
				*currentStS.Spec.Replicas = 0
				_, err := r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).Update(ctx, currentStS, metav1.UpdateOptions{})
				if err != nil {
					return microerror.Mask(err)
				}
			}

			// wait 30s for pvc gone
			{
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
			}
		} else {
			r.logger.Debugf(ctx, "pvc do not need to be re-created")
			return nil
		}
	}

	// scale down again sts
	currentStS, err = r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).Get(ctx, stsName, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}
	*currentStS.Spec.Replicas = 0
	_, err = r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).Update(ctx, currentStS, metav1.UpdateOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	r.logger.Debugf(ctx, "pvc re-created")

	return nil
}
