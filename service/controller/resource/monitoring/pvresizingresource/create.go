package pvresizingresource

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/pvresizing"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "ensuring pv resizing")
	{
		cluster, err := key.ToCluster(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		namespace := fmt.Sprintf("%s-%s", cluster.GetName()+"-prometheus")
		pvc, err := r.fetchPVC(ctx, cluster.GetName(), namespace)
		if err != nil {
			return microerror.Mask(err)
		}

		currentValue := pvc.Spec.Resources.Requests.Storage().String()
		anoVal := cluster.GetAnnotations()[key.PrometheusDiskSizeAnnotation]
		desireDiskSize := pvresizing.PrometheusDiskSize(anoVal)

		// Check the value of annotation with the current value in PVC.
		if currentValue != desireDiskSize {
			// if different we can trigger an update
			// https://github.com/prometheus-operator/prometheus-operator/issues/4079#issuecomment-1211989005
			err = r.updatePrometheusDiskSize(ctx, cluster.GetName(), namespace, desireDiskSize, pvc)
			if err != nil {
				return microerror.Mask(err)
			}
		}

	}
	r.logger.Debugf(ctx, "ensured pv resizing")

	return nil
}

func (r *Resource) fetchPVC(ctx context.Context, clusterID, namespace string) (*corev1.PersistentVolumeClaim, error) {
	l := labels.SelectorFromSet(labels.Set(map[string]string{
		"app.kubernetes.io/instance":   clusterID,
		"app.kubernetes.io/managed-by": "prometheus-operator",
		"app.kubernetes.io/name":       "prometheus",
		"operator.prometheus.io/name":  clusterID,
		"prometheus":                   clusterID,
	}))
	options := metav1.ListOptions{
		LabelSelector: l.String(),
	}
	list, err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).List(ctx, options)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if len(list.Items) == 0 {
		return nil, fmt.Errorf("couldn't retrieve PVC")
	}
	return &list.Items[0], err
}

func (r *Resource) updatePrometheusDiskSize(ctx context.Context, clusterID, namespace, desiredValue string, pvc *corev1.PersistentVolumeClaim) error {

	prometheus, err := r.prometheusClient.MonitoringV1().Prometheuses(namespace).Get(ctx, clusterID, metav1.GetOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	newSize := resource.MustParse(desiredValue)
	// Updating Prometheus CR Storage size if it doesn't match with the desiredValue
	if !(prometheus.Spec.Storage.VolumeClaimTemplate.Spec.Resources.Requests.Storage().String() == desiredValue) {
		prometheus.Spec.Storage.VolumeClaimTemplate.Spec.Resources.Requests["Storage"] = newSize
		_, err = r.prometheusClient.MonitoringV1().Prometheuses(namespace).Update(ctx, prometheus, metav1.UpdateOptions{})
		if err != nil {
			return microerror.Mask(err)
		}
	}

	// Patch PVC with the new size
	pvc.Spec.Resources.Requests["Storage"] = newSize
	patch := []byte(fmt.Sprintf(`{"spec": { "resources": { "requests": { "storage": "%v" } } } }`, newSize))
	_, err = r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).
		Patch(ctx, pvc.GetName(), types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return microerror.Mask(err)
	}

	// Delete the STS to update its definition, the recreation is done right away by the Prometheus operator
	orphan := metav1.DeletePropagationOrphan
	err = r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).
		Delete(ctx, fmt.Sprintf("prometheus-%v", clusterID), metav1.DeleteOptions{PropagationPolicy: &orphan})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
