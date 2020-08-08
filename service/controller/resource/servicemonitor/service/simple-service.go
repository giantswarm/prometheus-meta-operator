package service

import (
	"fmt"
	"strings"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

// SimpleService returns a ServiceMonitor that selects services labelled with
// the "giantswarm.io/monitoring" label.
func SimpleService(cluster metav1.Object) *promv1.ServiceMonitor {
	return &promv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "simple-service",
			Namespace: key.Namespace(cluster),
			Labels: map[string]string{
				key.ClusterIDKey(): key.ClusterID(cluster),
			},
		},
		Spec: promv1.ServiceMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"giantswarm.io/monitoring": "true",
				},
			},
			NamespaceSelector: promv1.NamespaceSelector{
				Any: true,
			},
			Endpoints: []promv1.Endpoint{
				{
					Port:   "https",
					Scheme: "https",
					RelabelConfigs: []*promv1.RelabelConfig{
						{
							SourceLabels: []string{"__meta_kubernetes_namespace"},
							TargetLabel:  "namespace",
						},
						{
							SourceLabels: []string{"__meta_kubernetes_pod_name"},
							TargetLabel:  "pod", // TODO or "pod_name" ??
						},
						// Preserve selected service discovery metadata
						// TODO if this works, add "namespace" here and remove above
						{
							Action: "labelmap",
							Regex: fmt.Sprintf(
								"__meta_kubernetes_(%s)",
								strings.Join([]string{
									"node_name", "service_name", "pod_name", "pod_ip",
									"pod_container_name", "pod_ready", "pod_phase",
									"pod_node_name", "endpoints_name",
								}, "|"),
							),
						},
						// Default metrics path
						{
							Replacement: "/metrics",
							TargetLabel: "__tmp_metrics_path",
						},
						// Extract default port from __address__ before rewriting it.
						{
							Regex:        "(.+):(\\d+)",
							Replacement:  "$2",
							SourceLabels: []string{"__address__"},
							TargetLabel:  "__tmp_metrics_port",
						},
						// Override metrics path if specified
						{
							Regex: "(.+)",
							SourceLabels: []string{
								"__meta_kubernetes_service_annotation_giantswarm_io_monitoring_path",
							},
							TargetLabel: "__tmp_metrics_path",
						},
						// Override metrics port if specified
						{
							Regex: "(.+)",
							SourceLabels: []string{
								"__meta_kubernetes_service_annotation_giantswarm_io_monitoring_port",
							},
							TargetLabel: "__tmp_metrics_port",
						},
						// Rewrite address to API server to use its proxy
						{
							Replacement:  fmt.Sprintf("master.%s", key.ClusterID(cluster)),
							SourceLabels: []string{"__address__"},
							TargetLabel:  "__address__",
						},
						// Rewrite metrics path to go through API server proxy
						{
							Regex:       "(.*);(.*);(.*);(.*)",
							Replacement: "/api/v1/namespaces/${1}/pods/${2}:${3}/proxy/${4}",
							SourceLabels: []string{
								"namespace", "pod",
								"__tmp_metrics_port", "__tmp_metrics_path",
							},
							TargetLabel: "__metrics_path__",
						},
						// Expose all labels on the service as normalised
						// labels on the metric (i.e. with `_` replacements)
						{
							Action: "labelmap",
							Regex:  "__meta_kubernetes_service_label_(.+)",
						},
					},
					TLSConfig: &promv1.TLSConfig{
						CAFile:   fmt.Sprintf("/etc/prometheus/secrets/%s/ca", key.Secret()),
						CertFile: fmt.Sprintf("/etc/prometheus/secrets/%s/crt", key.Secret()),
						KeyFile:  fmt.Sprintf("/etc/prometheus/secrets/%s/key", key.Secret()),
					},
				},
			},
		},
	}
}
