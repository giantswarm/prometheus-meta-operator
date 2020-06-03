package key

import (
	"fmt"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ToCluster(obj interface{}) (metav1.Object, error) {
	clusterMetaObject, ok := obj.(metav1.Object)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "'%T' does not implements '%T'", obj, clusterMetaObject)
	}

	return clusterMetaObject, nil
}

func Namespace(cluster metav1.Object) string {
	return fmt.Sprintf("%s-prometheus", cluster.GetName())
}

func TenantClusterSecret(cluster metav1.Object) string {
	return fmt.Sprintf("%s-prometheus-secret", cluster.GetName())
}

func TenantClusterHost(cluster metav1.Object, baseDomain string) string {
	return fmt.Sprintf("%s.prometheus.%s", cluster.GetName(), baseDomain)
}

func Secret() string {
	return "cluster-certificates"
}

func ClusterIDKey() string {
	return "cluster_id"
}

func ClusterID(cluster metav1.Object) string {
	return cluster.GetName()
}
