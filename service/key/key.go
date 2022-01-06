package key

import (
	"fmt"
	"math"
	"strings"

	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	capiv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"

	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
)

const (
	monitoring = "monitoring"

	PrometheusMemoryLimitCoefficient float64 = 1.2
)

func ToCluster(obj interface{}) (metav1.Object, error) {
	clusterMetaObject, ok := obj.(metav1.Object)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "'%T' does not implements 'metav1.Object'", obj)
	}

	return clusterMetaObject, nil
}

type MetaRuner interface {
	metav1.Object
	runtime.Object
}

func ToClusterMR(obj interface{}) (MetaRuner, error) {
	clusterMetaObject, ok := obj.(MetaRuner)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "'%T' does not implements 'MetaRuner'", obj)
	}

	return clusterMetaObject, nil
}

func Namespace(cluster metav1.Object) string {
	return fmt.Sprintf("%s-prometheus", ClusterID(cluster))
}

func NamespaceDefault(cluster metav1.Object) string {
	return v1.NamespaceDefault
}

func OrganizationNamespace(cluster metav1.Object) string {
	return cluster.GetNamespace()
}

func NamespaceMonitoring() string {
	return monitoring
}

func Secret() string {
	return SecretAPICertificates(nil)
}

func SecretAPICertificates(cluster metav1.Object) string {
	return "cluster-certificates"
}

func CAPICertificateName(cluster metav1.Object) string {
	return fmt.Sprintf("%s-kubeconfig", ClusterID(cluster))
}

func CAPICertificateNamespace(cluster metav1.Object) string {
	return cluster.GetNamespace()
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

func AlertmanagerLabels(cluster metav1.Object) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "alertmanager",
		"app.kubernetes.io/managed-by": project.Name(),
		"app.kubernetes.io/instance":   "alertmanager",
	}
}

func PrometheusLabels(cluster metav1.Object) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       "prometheus",
		"app.kubernetes.io/managed-by": project.Name(),
		"app.kubernetes.io/instance":   ClusterID(cluster),
		"giantswarm.io/cluster":        ClusterID(cluster),
	}
}

func AlertmanagerDefaultCPU() *resource.Quantity {
	return resource.NewMilliQuantity(100, resource.DecimalSI)
}

func AlertmanagerDefaultMemory() *resource.Quantity {
	return resource.NewQuantity(200*1024*1024, resource.DecimalSI)
}

func PrometheusDefaultCPU() *resource.Quantity {
	return resource.NewMilliQuantity(100, resource.DecimalSI)
}

func PrometheusDefaultCPULimit() *resource.Quantity {
	return resource.NewMilliQuantity(100*1.5, resource.DecimalSI)
}

func PrometheusDefaultMemory() *resource.Quantity {
	return resource.NewQuantity(1024*1024*1024, resource.DecimalSI)
}

func PrometheusDefaultMemoryLimit() *resource.Quantity {
	return resource.NewQuantity(
		int64(math.Floor(
			1024*1024*1024*PrometheusMemoryLimitCoefficient,
		)),
		resource.DecimalSI,
	)
}

func AlertmanagerKey() string {
	return "alertmanager.yaml"
}

func AlertmanagerPort() int32 {
	return 9093
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

func PrometheusContainerName() string {
	return "prometheus"
}

func PrometheusSTSName(cluster metav1.Object) string {
	return fmt.Sprintf("prometheus-%s", ClusterID(cluster))
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
		return fmt.Sprintf("%s:443", v.Spec.ClusterIP)
	case *capiv1alpha3.Cluster: // Support CAPI Clusters
		host := strings.TrimPrefix(v.Spec.ControlPlaneEndpoint.Host, "https://")
		host = strings.TrimPrefix(host, "http://")
		return fmt.Sprintf("%s:%d", host, v.Spec.ControlPlaneEndpoint.Port)
	case metav1.Object:
		return fmt.Sprintf("master.%s:443", v.GetName())
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
		return "management_cluster"
	}

	return "workload_cluster"
}

func BearerTokenPath() string {
	return "/var/run/secrets/kubernetes.io/serviceaccount/token"
}

func CAFilePath() string {
	return "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
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

// IsCAPICluster returns true if the cluster is in v1alpha3 and does not have the "azure-operator.giantswarm.io/version" label added by the azure operator.
// We do not have a provider agnostic label like "giantswarm.io/version" to define this.
func IsCAPICluster(obj metav1.Object) bool {
	// TODO once we have migrated all clusters to CAPI, we can remove this
	switch v := obj.(type) {
	case *capiv1alpha3.Cluster:
		if _, ok := v.Labels["azure-operator.giantswarm.io/version"]; ok {
			return false
		}
		return true
	default:
		return false
	}
}
