package key

import (
	"fmt"

	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
)

const monitoring = "monitoring"

func ToCluster(obj interface{}) (metav1.Object, error) {
	clusterMetaObject, ok := obj.(metav1.Object)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "'%T' does not implements 'metav1.Object'", obj)
	}

	return clusterMetaObject, nil
}

func Namespace(cluster metav1.Object) string {
	return fmt.Sprintf("%s-prometheus", ClusterID(cluster))
}

func NamespaceDefault(cluster metav1.Object) string {
	return v1.NamespaceDefault
}

func NamespaceMonitoring(cluster metav1.Object) string {
	return monitoring
}

func Secret() string {
	return SecretAPICertificates(nil)
}

func SecretAPICertificates(cluster metav1.Object) string {
	return "cluster-certificates"
}

func SecretTLSCertificates(cluster metav1.Object) string {
	return "prometheus-tls"
}

func IsMonitoringDisabled(cluster metav1.Object) bool {
	ignored, ok := cluster.GetLabels()["giantswarm.io/monitoring"]
	return ok && ignored == "false"
}

func EtcdSecret(obj interface{}) string {
	if IsInCluster(obj) {
		return "etcd-certificates"
	}

	return Secret()
}

func Labels(cluster metav1.Object) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "prometheus",
		"app.kubernetes.io/managed-by": project.Name(),
		"app.kubernetes.io/instance":   ClusterID(cluster),
	}
}

func PrometheusPort() int32 {
	return 9090
}

func ClusterIDKey() string {
	return "cluster_id"
}

func ClusterID(cluster metav1.Object) string {
	return cluster.GetName()
}

func InstallationKey() string {
	return "installation"
}

func TypeKey() string {
	return "type"
}

func Heartbeat() string {
	return "heartbeat"
}

func HeartbeatName(cluster metav1.Object, installation string) string {
	return fmt.Sprintf("%s-%s", installation, ClusterID(cluster))
}

func HeartbeatReceiverName(cluster metav1.Object, installation string) string {
	return fmt.Sprintf("heartbeat_%s_%s", installation, ClusterID(cluster))
}

func PrometheusAdditionalScrapeConfigsSecretName() string {
	return "additional-scrape-configs"
}

func PrometheusAdditionalScrapeConfigsName() string {
	return "prometheus-additional.yaml"
}

func AlertManagerSecretName() string {
	return "alertmanager-config"
}

func AlertManagerKey() string {
	return "alertmanager-additional.yaml"
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

func PromxyConfigMapName() string {
	return "promxy-app-unique"
}

func PromxyConfigMapNamespace() string {
	return monitoring
}

func PromxyConfigFileName() string {
	return "config.yaml"
}

func AlertmanagerConfigMapName() string {
	return "alertmanager"
}

func AlertmanagerConfigMapNamespace() string {
	return monitoring
}

func AlertmanagerConfigMapKey() string {
	return "config.yml"
}

func RemoteWriteSecretName() string {
	return "remote-write"
}

func RemoteWriteUsernameKey() string {
	return "username"
}

func RemoteWritePasswordKey() string {
	return "password"
}
