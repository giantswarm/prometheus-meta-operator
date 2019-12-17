package key

import (
	"fmt"

	"github.com/giantswarm/microerror"
	"sigs.k8s.io/cluster-api/api/v1alpha2"
)

func ToCluster(obj interface{}) (*v1alpha2.Cluster, error) {
	cluster, ok := obj.(*v1alpha2.Cluster)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "expected '%T', got '%T'", &v1alpha2.Cluster{}, obj)
	}

	return cluster, nil
}

func Namespace(cluster *v1alpha2.Cluster) string {
	return fmt.Sprintf("%s-prometheus", cluster.GetName())
}

func Secret() string {
	return "cluster-certificates"
}
