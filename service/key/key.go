package key

import (
	"fmt"

	"github.com/giantswarm/apiextensions/v2/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/api/v1alpha2"
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

func Secret() string {
	return "cluster-certificates"
}

func ClusterIDKey() string {
	return "cluster_id"
}

func ClusterID(cluster metav1.Object) string {
	return cluster.GetName()
}

func PrometheusAdditionalScrapeConfigsSecretName() string {
	return "additional-scrape-configs"
}

func PrometheusAdditionalScrapeConfigsName() string {
	return "prometheus-additional.yaml"
}

func APIUrl(obj interface{}) string {
	switch v := obj.(type) {
	case *v1alpha2.Cluster:
		return fmt.Sprintf("master.%s", v.GetName())
	case *v1alpha1.AWSConfig:
		return fmt.Sprintf("master.%s", v.GetName())
	case *v1alpha1.AzureConfig:
		return fmt.Sprintf("master.%s", v.GetName())
	case *v1alpha1.KVMConfig:
		return fmt.Sprintf("master.%s", v.GetName())
	case *v1.Service:
		return v.Spec.ClusterIP
	}

	return ""
}

func IsInCluster(obj interface{}) bool {
	switch obj.(type) {
	case *v1.Service:
		return true
	default:
		return false
	}
}

func ClusterType(obj interface{}) string {
	if IsInCluster(obj) {
		return "control_plane"
	}

	return "tenant_cluster"
}
