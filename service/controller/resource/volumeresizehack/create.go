package volumeresizehack

import (
	"context"
	"fmt"
	"reflect"
	"time"

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

	r.logger.LogCtx(ctx, "level", "debug", "message", "creating")

	// 1. pause prometheus
	// 2. delete pvc
	// 3. scale down sts
	// 4. wait for pvc gone
	// 5. scale up sts
	// 6. unpause prometheus

	namespace := key.Namespace(cluster)

	// Get prometheus
	prometheusName := key.ClusterID(cluster)
	currentPrometheus, err := r.prometheusClient.MonitoringV1().Prometheuses(namespace).Get(ctx, prometheusName, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	// Get StS
	stsName := fmt.Sprintf("prometheus-%s", cluster.GetName())
	currentStS, err := r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).Get(ctx, stsName, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	noPVCnoReplicas := false
	resize := false
	// Get PVC
	// TODO find by name : prometheus-df3yn-db-prometheus-df3yn-0
	index := 0
	desiredPVC := currentStS.Spec.VolumeClaimTemplates[index]
	pvcName := fmt.Sprintf("%s-%s-%d", desiredPVC.GetName(), currentStS.GetName(), index)
	currentPVC, err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		r.logger.LogCtx(ctx, "level", "debug", "message", "PVC NOT FOUND")
		if *currentStS.Spec.Replicas <= 0 {
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("StS SCALED DOWN replicas=%d", *currentStS.Spec.Replicas))
			noPVCnoReplicas = true
		} else {
			// TODO: return nil, when pvc is not present but replicas is set to 1, need time or there's another issue. This resource is not responsible here.
			return microerror.Mask(err)
		}
	} else {
		if err != nil {
			return microerror.Mask(err)
		}
		if !reflect.DeepEqual(desiredPVC.Spec, currentPVC.Spec) {
			r.logger.LogCtx(ctx, "level", "debug", "message", "PVC DO NOT MATCH")
			resize = true
		}
	}

	if noPVCnoReplicas || resize {
		if !noPVCnoReplicas {
			// delete pvc
			r.logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("DELETING PVC %s/%s", namespace, pvcName))
			err = r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
			r.logger.LogCtx(ctx, "level", "debug", "message", "DELETED PVC")
			time.Sleep(5 * time.Second)
		}
	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "PVC MATCH")
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "created")

	return nil
}
