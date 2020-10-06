package volumeresizehack

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/giantswarm/backoff"
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
		// pause
		r.logger.LogCtx(ctx, "level", "debug", "message", "PAUSING PROMETHEUS")
		currentPrometheus.Spec.Paused = true
		pausedPrometheus, err := r.prometheusClient.MonitoringV1().Prometheuses(namespace).Update(ctx, currentPrometheus, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.LogCtx(ctx, "level", "debug", "message", "PAUSED PROMETHEUS")

		if !noPVCnoReplicas {
			// delete pvc
			r.logger.LogCtx(ctx, "level", "debug", "message", "DELETING PVC")
			err = r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
			r.logger.LogCtx(ctx, "level", "debug", "message", "DELETED PVC")

			// scale down
			r.logger.LogCtx(ctx, "level", "debug", "message", "SCALING DOWN")
			scale := &autoscalingv1.Scale{
				ObjectMeta: metav1.ObjectMeta{
					Name:      stsName,
					Namespace: namespace,
				},
				Spec: autoscalingv1.ScaleSpec{
					Replicas: 0,
				},
			}
			_, err = r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).UpdateScale(ctx, stsName, scale, metav1.UpdateOptions{})
			if err != nil {
				return microerror.Mask(err)
			}
			r.logger.LogCtx(ctx, "level", "debug", "message", "SCALED DOWN")

			// wait 30s for pvc gone
			r.logger.LogCtx(ctx, "level", "debug", "message", "WAITING PVC")
			o := func() error {
				_, err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
				if apierrors.IsNotFound(err) {
					return nil
				}

				return microerror.Mask(pvcExist)
			}
			b := backoff.NewMaxRetries(6, 5*time.Second)
			err = backoff.Retry(o, b)
			if err != nil {
				return microerror.Mask(err)
			}
			r.logger.LogCtx(ctx, "level", "debug", "message", "WAITED PVC")
		}

		// scale back up
		r.logger.LogCtx(ctx, "level", "debug", "message", "SCALING UP")
		scale := &autoscalingv1.Scale{
			ObjectMeta: metav1.ObjectMeta{
				Name:      stsName,
				Namespace: namespace,
			},
			Spec: autoscalingv1.ScaleSpec{
				Replicas: 1,
			},
		}
		_, err = r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).UpdateScale(ctx, stsName, scale, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.LogCtx(ctx, "level", "debug", "message", "SCALED UP")

		// unpause
		r.logger.LogCtx(ctx, "level", "debug", "message", "UNPAUSING PROMETHEUS")
		pausedPrometheus.Spec.Paused = false
		_, err = r.prometheusClient.MonitoringV1().Prometheuses(namespace).Update(ctx, pausedPrometheus, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
		r.logger.LogCtx(ctx, "level", "debug", "message", "UNPAUSED PROMETHEUS")

	} else {
		r.logger.LogCtx(ctx, "level", "debug", "message", "PVC MATCH")
	}

	r.logger.LogCtx(ctx, "level", "debug", "message", "created")

	return nil
}
