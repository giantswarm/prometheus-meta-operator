package prometheus

import (
	"context"
	"math"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/nodecounter"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

func ComputePrometheusMinMemory(ctx context.Context, k8sclient k8sclient.Interface, cluster metav1.Object, installation string, provider string) (*resource.Quantity, error) {
	prometheusMemory := resource.NewQuantity(1024*1024*1024, resource.DecimalSI)

	if !key.IsManagementCluster(installation, cluster) || key.IsCAPIManagementCluster(provider) {
		nodeCount, err := nodecounter.CountClusterNodes(ctx, k8sclient, cluster)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		if nodeCount > 2 {
			// We consider that a node requires 500Mb of Prometheus RAM
			prometheusMemory = resource.NewQuantity(
				int64(math.Floor(
					1024*1024*512*float64(nodeCount),
				)),
				resource.DecimalSI,
			)
		}
	}

	return prometheusMemory, nil
}
