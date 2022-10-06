package pvcresizingresource

import (
	"context"
	"fmt"

	"github.com/giantswarm/microerror"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/pvcresizing"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

func (r *Resource) EnsureCreated(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "ensuring pvc resizing")
	{
		cluster, err := key.ToCluster(obj)
		if err != nil {
			return microerror.Mask(err)
		}

		namespace := key.Namespace(cluster)
		pvcList, err := r.listPVC(ctx, cluster.GetName(), namespace)
		if err != nil {
			return microerror.Mask(err)
		}

		for _, pvc := range pvcList {
			fmt.Println("PVC Name", pvc.GetName())
			currentPVCSize := pvc.Spec.Resources.Requests.Storage().String()
			desiredPVCSize := cluster.GetAnnotations()[key.PrometheusDiskSizeAnnotation]
			desiredVolumeSize := pvcresizing.PrometheusVolumeSize(desiredPVCSize)

			fmt.Println("currentPVCSize", currentPVCSize)
			fmt.Println("desiredPVCSize", desiredPVCSize)
			fmt.Println("desiredVolumeSize", desiredVolumeSize)
			// Check the value of annotation with the current value in PVC.
			if currentPVCSize < desiredVolumeSize {
				// Resizing requested. Following the procedure described here:
				// https://github.com/prometheus-operator/prometheus-operator/issues/4079#issuecomment-1211989005
				// until stateful set resizing made it into kubernetes:
				// https://github.com/kubernetes/enhancements/issues/661
				fmt.Println("resize..........", pvc.GetName())
				err = r.resize(ctx, desiredVolumeSize, pvc)
				if err != nil {
					return microerror.Mask(err)
				}
			} else if currentPVCSize > desiredPVCSize {
				// Since Resizing to lower storage is forbidden
				// Therefore, we replace the PVC and STS
				// But this will cause data loss
				fmt.Println("replacePVC..........", pvc.GetName())
				err = r.replacePVC(ctx, pvc)
				if err != nil {
					return microerror.Mask(err)
				}
			}
		}
	}
	r.logger.Debugf(ctx, "ensured pvc resizing")

	return nil
}

func pvcSelector(clusterID string) labels.Selector {
	return labels.SelectorFromSet(labels.Set(map[string]string{
		"app.kubernetes.io/instance":   clusterID,
		"app.kubernetes.io/managed-by": "prometheus-operator",
		"app.kubernetes.io/name":       "prometheus",
		"operator.prometheus.io/name":  clusterID,
		"prometheus":                   clusterID,
	}))
}

func (r *Resource) listPVC(ctx context.Context, clusterID, namespace string) ([]corev1.PersistentVolumeClaim, error) {
	options := metav1.ListOptions{
		LabelSelector: pvcSelector(clusterID).String(),
	}
	list, err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).List(ctx, options)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if len(list.Items) == 0 {
		return nil, fmt.Errorf("could not find PVCs")
	}
	return list.Items, err
}

func (r *Resource) resize(ctx context.Context, desiredVolumeSize string, pvc corev1.PersistentVolumeClaim) error {

	namespace := pvc.GetNamespace()
	clusterID := pvc.GetLabels()["prometheus"]

	newVolumeSize := resource.MustParse(desiredVolumeSize)

	// Patch PVC with the new size
	pvc.Spec.Resources.Requests["Storage"] = newVolumeSize
	patch := []byte(fmt.Sprintf(`{"spec": { "resources": { "requests": { "storage": "%v" } } } }`, newVolumeSize.String()))
	_, err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).
		Patch(ctx, pvc.GetName(), types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return microerror.Mask(err)
	}
	fmt.Println("PVC patched.....")

	// Delete the sts without the PVC (using orphan)
	orphan := metav1.DeletePropagationOrphan
	err = r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).
		Delete(ctx, fmt.Sprintf("prometheus-%v", clusterID), metav1.DeleteOptions{PropagationPolicy: &orphan})
	if err != nil {
		return microerror.Mask(err)
	}
	fmt.Println("Sts deleted.....")

	return nil
}

func (r *Resource) replacePVC(ctx context.Context, pvc corev1.PersistentVolumeClaim) error {

	namespace := pvc.GetNamespace()
	clusterID := pvc.GetLabels()["prometheus"]

	patch := []byte(`{"metadata":{"finalizers":null}}`)
	_, err := r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).
		Patch(ctx, pvc.GetName(), types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return microerror.Mask(err)
	}
	err = r.k8sClient.K8sClient().CoreV1().PersistentVolumeClaims(namespace).
		Delete(ctx, pvc.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return microerror.Mask(err)
	}
	fmt.Println("PVC deleted.....")

	// Delete the sts without the PVC (using orphan)
	orphan := metav1.DeletePropagationOrphan
	err = r.k8sClient.K8sClient().AppsV1().StatefulSets(namespace).
		Delete(ctx, fmt.Sprintf("prometheus-%v", clusterID), metav1.DeleteOptions{PropagationPolicy: &orphan})
	if err != nil {
		return microerror.Mask(err)
	}
	fmt.Println("Sts deleted.....")

	return nil
}
