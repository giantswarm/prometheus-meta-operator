package key

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/project"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/prometheus/agent"
)

const (
	MonitoringNamespace = "monitoring"

	DefaultServicePriority string = "highest"

	ClusterLabel                 string = "giantswarm.io/cluster"
	MonitoringLabel              string = "giantswarm.io/monitoring"
	ServicePriorityLabel         string = "giantswarm.io/service-priority"
	TeamLabel                    string = "application.giantswarm.io/team"
	OpsGenieApiKey               string = "opsGenieApiKey" // #nosec G101
	AlertmanagerGlobalSecretName string = "alertmanager-global"

	// PrometheusCPULimitCoefficient is the number used to compute the CPU limit from the CPU request.
	// It is used when computing VPA settings, to set `max requests` so that `max limits` respects MaxCPU factor.
	PrometheusCPULimitCoefficient float64 = 1.5
	// PrometheusMemoryLimitCoefficient is the number used to compute the memory limit from the memory request.
	// It is used when computing VPA settings, to set `max request` so that `max limits` respect MaxMemory factor.
	PrometheusMemoryLimitCoefficient      float64 = 1
	PrometheusMetaOperatorRemoteWriteName string  = "prometheus-meta-operator"
	PrometheusServiceName                         = "prometheus-operated"
	// RemoteWriteAPIEndpointConfigSecretNameKey is the secret name used by a Prometheus client to access the Prometheus remote write endpoint. It is used at https://github.com/giantswarm/observability-bundle/blob/main/helm/observability-bundle/templates/apps.yaml
	RemoteWriteAPIEndpointConfigSecretNameKey string = "remote-write-api-endpoint-config"
	// RemoteWriteIngressAuthSecretName is the secret name referenced in the ingress to enable authentication against the Prometheus remote write endpoint.
	RemoteWriteIngressAuthSecretName string = "remote-write-ingress-auth"
	// PrometheusVolumeSizeAnnotation is the annotation referenced in the Cluster CR to define the size of Prometheus Volume.
	PrometheusVolumeSizeAnnotation string = "monitoring.giantswarm.io/prometheus-volume-size"
	// We apply a ratio to the volume storage size to compute the RetentionSize property (RetentionSize = 90% volume storage size)
	PrometheusVolumeStorageLimitRatio = 0.85

	ClusterIDKey       string = "cluster_id"
	ClusterTypeKey     string = "cluster_type"
	CustomerKey        string = "customer"
	InstallationKey    string = "installation"
	OrganizationKey    string = "organization"
	PipelineKey        string = "pipeline"
	ProviderKey        string = "provider"
	RegionKey          string = "region"
	ServicePriorityKey string = "service_priority"
	TypeKey            string = "type"

	IngressClassName string = "nginx"

	BearerTokenPath string = "/var/run/secrets/kubernetes.io/serviceaccount/token" // nolint:gosec
	CAFilePath      string = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"

	EtcdSecretSourceName      string = "etcd-certs"
	EtcdSecretSourceNamespace string = "giantswarm"

	APIServerCertificatesSecretName string = "cluster-certificates" // nolint:gosec
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

func ClusterNamespace(cluster metav1.Object) string {
	return cluster.GetNamespace()
}

func CAPICertificateName(cluster metav1.Object) string {
	return fmt.Sprintf("%s-kubeconfig", ClusterID(cluster))
}

func GetServicePriority(cluster metav1.Object) string {
	if servicePriority, ok := cluster.GetLabels()[ServicePriorityLabel]; ok && servicePriority != "" {
		return servicePriority
	}
	return DefaultServicePriority
}

func IsMonitoringDisabled(cluster metav1.Object) bool {
	ignored, ok := cluster.GetLabels()[MonitoringLabel]
	return ok && ignored == "false"
}

func EtcdSecret(installation string, obj interface{}) string {
	if IsManagementCluster(installation, obj) {
		return "etcd-certificates"
	}

	return APIServerCertificatesSecretName
}

func IsCAPIManagementCluster(provider cluster.Provider) bool {
	return provider.Flavor == "capi"
}

func ClusterProvider(obj metav1.Object, provider cluster.Provider) (string, error) {
	// TODO remove once all clusters are on CAPI
	// We keep the existing behavior for vintage management clusters
	if !IsCAPIManagementCluster(provider) {
		return provider.Kind, nil
	}

	if c, ok := obj.(*capi.Cluster); ok {
		switch c.Spec.InfrastructureRef.Kind {
		case cluster.AWSClusterKind:
			return cluster.AWSClusterKindProvider, nil
		case cluster.AWSManagedClusterKind:
			return cluster.AWSManagedClusterKindProvider, nil
		case cluster.AzureClusterKind:
			return cluster.AzureClusterKindProvider, nil
		case cluster.AzureManagedClusterKind:
			return cluster.AzureManagedClusterKindProvider, nil
		case cluster.VCDClusterKind:
			return cluster.VCDClusterKindProvider, nil
		case cluster.VSphereClusterKind:
			return cluster.VSphereClusterKindProvider, nil
		case cluster.GCPClusterKind:
			return cluster.GCPClusterKindProvider, nil
		case cluster.GCPManagedClusterKind:
			return cluster.GCPManagedClusterKindProvider, nil
		}
	}

	return "", infrastructureRefNotFoundError
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

func RemoteWriteAuthenticationAnnotations(baseDomain string, externalDNS bool) map[string]string {
	annotations := map[string]string{
		"nginx.ingress.kubernetes.io/auth-type":   "basic",
		"nginx.ingress.kubernetes.io/auth-secret": RemoteWriteIngressAuthSecretName,
		"nginx.ingress.kubernetes.io/auth-realm":  "Authentication Required",
		// Set this annotation to avoid using a temporary buffer file for remote write requests
		"nginx.ingress.kubernetes.io/client-body-buffer-size": "50m",
		// Remote write requests can be quite big. (default max body size: 1m)
		"nginx.ingress.kubernetes.io/proxy-body-size": "50m",
	}

	// create external-dns required annotations
	if externalDNS {
		annotations["external-dns.alpha.kubernetes.io/hostname"] = baseDomain
		annotations["giantswarm.io/external-dns"] = "managed"
	}

	return annotations
}

func RemoteWriteConfigName(cluster metav1.Object) string {
	return fmt.Sprintf("%s-remote-write-config", ClusterID(cluster))
}

func RemoteWriteSecretName(cluster metav1.Object) string {
	return fmt.Sprintf("%s-remote-write-secret", ClusterID(cluster))
}

func GetClusterAppsNamespace(cluster metav1.Object, installation string, provider cluster.Provider) string {
	if IsCAPIManagementCluster(provider) {
		return cluster.GetNamespace()
	} else if IsManagementCluster(installation, cluster) {
		return "giantswarm"
	}
	return ClusterID(cluster)
}

func RemoteWriteAPIEndpointConfigSecretName(cluster metav1.Object, provider cluster.Provider) string {
	// TODO remove once all clusters are on v19
	if IsCAPIManagementCluster(provider) {
		return fmt.Sprintf("%s-%s", ClusterID(cluster), RemoteWriteAPIEndpointConfigSecretNameKey)
	}
	return RemoteWriteAPIEndpointConfigSecretNameKey
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
	return resource.NewMilliQuantity(
		int64(math.Floor(
			100*PrometheusCPULimitCoefficient,
		)),
		resource.DecimalSI,
	)
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

func PrometheusPort() int32 {
	return 9090
}

func ClusterID(cluster metav1.Object) string {
	return cluster.GetName()
}

func GetClusterShardingStrategy(cluster metav1.Object) (*agent.ShardingStrategy, error) {
	var err error
	var scaleUpSeriesCount, scaleDownPercentage float64
	if value, ok := cluster.GetAnnotations()["monitoring.giantswarm.io/prometheus-agent-scale-up-series-count"]; ok {
		if scaleUpSeriesCount, err = strconv.ParseFloat(value, 64); err != nil {
			return nil, err
		}
	}
	if value, ok := cluster.GetAnnotations()["monitoring.giantswarm.io/prometheus-agent-scale-down-percentage"]; ok {
		if scaleDownPercentage, err = strconv.ParseFloat(value, 64); err != nil {
			return nil, err
		}
	}
	return &agent.ShardingStrategy{
		ScaleUpSeriesCount:  scaleUpSeriesCount,
		ScaleDownPercentage: scaleDownPercentage,
	}, nil
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

func AlertmanagerSecretName() string {
	return "alertmanager-config"
}

func AlertmanagerKey() string {
	return "alertmanager-additional.yaml"
}

func APIUrl(obj interface{}) string {
	switch v := obj.(type) {
	case *v1.Service:
		return "kubernetes.default:443"
	case *capi.Cluster:
		// We remove any https:// prefix from the api-server host due to a bug in CAPA Managed EKS clusters (cf. https://gigantic.slack.com/archives/C02HLSDH3DZ/p1695213116360889)
		return fmt.Sprintf("%s:%d", strings.TrimPrefix(v.Spec.ControlPlaneEndpoint.Host, "https://"), v.Spec.ControlPlaneEndpoint.Port)
	case metav1.Object: // TODO remove once all clusters are on v19
		return fmt.Sprintf("master.%s:443", v.GetName())
	}

	return ""
}

func IsManagementCluster(installation string, obj interface{}) bool {
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

func IsEKSCluster(obj interface{}) bool {
	if c, ok := obj.(*capi.Cluster); ok {
		return c.Spec.InfrastructureRef.Kind == cluster.AWSManagedClusterKind
	}
	return false
}

func ClusterType(installation string, obj interface{}) string {
	if IsManagementCluster(installation, obj) {
		return "management_cluster"
	}

	return "workload_cluster"
}

func ApiServerAuthenticationType(ctx context.Context, k8sClient k8sclient.Interface, clusterNamespace string) (string, error) {
	secret, err := k8sClient.K8sClient().CoreV1().Secrets(clusterNamespace).Get(ctx, APIServerCertificatesSecretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if secret.Data["token"] != nil && len(secret.Data["token"]) > 0 {
		return "token", nil
	} else if (secret.Data["crt"] != nil && len(secret.Data["crt"]) > 0) && (secret.Data["key"] != nil && len(secret.Data["key"]) > 0) {
		return "certificates", nil
	}
	return "", errors.New("no authentication found")
}
