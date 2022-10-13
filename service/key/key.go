package key

import (
	"fmt"
	"math"

	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/project"
)

var capiProviders = []string{"capa", "cloud-director", "gcp", "openstack", "vsphere"}

const (
	monitoring = "monitoring"

	PrometheusMemoryLimitCoefficient       float64 = 1.2
	PrometheusMetaOperatorRemoteWriteName  string  = "prometheus-meta-operator"
	PrometheusServiceName                          = "prometheus-operated"
	RemoteWriteAPIEndpointConfigSecretName string  = "remote-write-api-endpoint-config"
	RemoteWriteIngressAuthSecretName       string  = "remote-write-ingress-auth"
)

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

func OrganizationNamespace(cluster metav1.Object) string {
	return cluster.GetNamespace()
}

func NamespaceMonitoring() string {
	return monitoring
}

func Secret() string {
	return "cluster-certificates"
}

func CAPICertificateName(cluster metav1.Object) string {
	return fmt.Sprintf("%s-kubeconfig", ClusterID(cluster))
}

func IsMonitoringDisabled(cluster metav1.Object) bool {
	ignored, ok := cluster.GetLabels()["giantswarm.io/monitoring"]
	return ok && ignored == "false"
}

func EtcdSecret(installation string, obj interface{}) string {
	if IsInCluster(installation, obj) {
		return "etcd-certificates"
	}

	return Secret()
}

func IsCAPIManagementCluster(provider string) bool {
	for _, v := range capiProviders {
		if v == provider {
			return true
		}
	}

	return false
}

func EtcdSecretSourceName() string {
	return "etcd-certs"
}

func EtcdSecretSourceNamespace() string {
	return "giantswarm"
}

func AlertmanagerLabels() map[string]string {
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

func RemoteWriteAuthenticationAnnotations() map[string]string {
	return map[string]string{
		"nginx.ingress.kubernetes.io/auth-type":   "basic",
		"nginx.ingress.kubernetes.io/auth-secret": RemoteWriteIngressAuthSecretName,
		"nginx.ingress.kubernetes.io/auth-realm":  "Authentication Required",
	}
}

func RemoteWriteAPIEndpointConfigSecretNameAndNamespace(cluster metav1.Object, provider string) (string, string) {
	name := RemoteWriteAPIEndpointConfigSecretName
	namespace := ClusterID(cluster)

	if IsCAPIManagementCluster(provider) {
		name = ClusterID(cluster) + "-" + name
		namespace = cluster.GetNamespace()
	}
	return name, namespace
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

func HeartbeatAPI(cluster metav1.Object, installation string) string {
	return fmt.Sprintf("https://api.opsgenie.com/v2/heartbeats/%s/ping", HeartbeatName(cluster, installation))
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
		return "kubernetes.default:443"
	case *capi.Cluster:
		return fmt.Sprintf("%s:%d", v.Spec.ControlPlaneEndpoint.Host, v.Spec.ControlPlaneEndpoint.Port)
	case metav1.Object:
		return fmt.Sprintf("master.%s:443", v.GetName())
	}

	return ""
}

func IsInCluster(installation string, obj interface{}) bool {
	switch v := obj.(type) {
	case *v1.Service:
		return true
	case *capi.Cluster:
		if val, ok := v.Labels["cluster.x-k8s.io/cluster-name"]; ok && val == installation {
			return true
		}
		return false
	default:
		return false
	}
}

func ClusterType(installation string, obj interface{}) string {
	if IsInCluster(installation, obj) {
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

// IsCAPICluster returns false if the cluster has any of the legacy labels such as azure-operator.giantswarm.io/version.
func IsCAPICluster(obj metav1.Object) bool {
	// TODO once we have migrated all clusters to CAPI, we can remove this

	checker := func(labels map[string]string) bool {
		if _, ok := labels["azure-operator.giantswarm.io/version"]; ok {
			return false
		}
		if _, ok := labels["cluster-operator.giantswarm.io/version"]; ok {
			return false
		}
		return true
	}

	switch v := obj.(type) {
	case *capi.Cluster:
		return checker(v.Labels)
	case *v1.Service:
		// Legacy Management Clusters.
		return false
	}

	// We didn't recognize the type, we assume CAPI
	return true
}

func IngressClassName() string {
	return "nginx"
}

func OpsgenieKey() string {
	return "opsgenie.key"
}
