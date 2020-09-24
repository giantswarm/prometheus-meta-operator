package key

import (
	"fmt"

	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
)

func ToCluster(obj interface{}) (metav1.Object, error) {
	clusterMetaObject, ok := obj.(metav1.Object)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "'%T' does not implements 'metav1.Object'", obj)
	}

	return clusterMetaObject, nil
}

func Namespace(cluster metav1.Object) string {
	return fmt.Sprintf("%s-prometheus", cluster.GetName())
}

func Secret() string {
	return "cluster-certificates"
}

func EtcdSecret(obj interface{}) string {
	if IsInCluster(obj) {
		return "etcd-certificates"
	}

	return Secret()
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
	case *v1.Service:
		return v.Spec.ClusterIP
	case metav1.Object:
		return fmt.Sprintf("master.%s", v.GetName())
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

func ControlPlaneBearerToken() string {
	return "/var/run/secrets/kubernetes.io/serviceaccount/token"
}

func ControlPlaneCAFile() string {
	return "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
}

func Labels(cluster metav1.Object) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "prometheus",
		"app.kubernetes.io/managed-by": project.Name(),
		"app.kubernetes.io/instance":   cluster.GetName(),
	}
}

func Oauth2ServicePort() int {
	return 4180
}
func Oauth2ServiceName() string {
	return "oauth2-proxy"
}
