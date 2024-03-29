package nodecounter

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expv1beta1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

func CountClusterNodes(ctx context.Context, k8sClient k8sclient.Interface, cluster metav1.Object) (int, error) {
	mdCount, err := countMachineDeploymentClusterNodes(ctx, k8sClient, cluster)
	if err != nil {
		return 0, microerror.Mask(err)
	}

	mpCount, err := countMachinePoolClusterNodes(ctx, k8sClient, cluster)
	if err != nil {
		return 0, microerror.Mask(err)
	}

	return int(mdCount + mpCount), nil
}

func countMachineDeploymentClusterNodes(ctx context.Context, k8sClient k8sclient.Interface, cluster metav1.Object) (int32, error) {
	var machinedeployments apiv1beta1.MachineDeploymentList = apiv1beta1.MachineDeploymentList{}

	opts := []client.ListOption{
		client.MatchingLabels{
			key.ClusterLabel: key.ClusterID(cluster),
		},
		client.InNamespace(cluster.GetNamespace()),
	}

	err := k8sClient.CtrlClient().List(ctx, &machinedeployments, opts...)
	if err != nil {
		return 0, microerror.Mask(err)
	}

	var nodeCount int32 = 0
	for _, machinedeployment := range machinedeployments.Items {
		// If the machinedeployment has no role label or one without the bastion value then we keep it
		if role, ok := machinedeployment.Labels["cluster.x-k8s.io/role"]; !ok || role != "bastion" {
			nodeCount += machinedeployment.Status.Replicas
		}
	}

	return nodeCount, nil
}

func countMachinePoolClusterNodes(ctx context.Context, k8sClient k8sclient.Interface, cluster metav1.Object) (int32, error) {
	var machinepools expv1beta1.MachinePoolList = expv1beta1.MachinePoolList{}

	opts := []client.ListOption{
		client.MatchingLabels{
			key.ClusterLabel: key.ClusterID(cluster),
		},
		client.InNamespace(cluster.GetNamespace()),
	}

	err := k8sClient.CtrlClient().List(ctx, &machinepools, opts...)
	if err != nil {
		return 0, microerror.Mask(err)
	}

	var nodeCount int32 = 0
	for _, machinepool := range machinepools.Items {
		// If the machinepool has no role label or one without the bastion value then we keep it
		if role, ok := machinepool.Labels["cluster.x-k8s.io/role"]; !ok || role != "bastion" {
			nodeCount += machinepool.Status.Replicas
		}
	}

	return nodeCount, nil
}
