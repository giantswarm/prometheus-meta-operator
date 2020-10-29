package service

import (
	"fmt"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func APIServer(cluster metav1.Object, provider string, installation string) *promv1.ServiceMonitor {
	serviceMonitor := &promv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("kubernetes-apiserver-%s", cluster.GetName()),
			Namespace: key.Namespace(cluster),
			Labels: map[string]string{
				key.ClusterIDKey(): key.ClusterID(cluster),
			},
		},
		Spec: promv1.ServiceMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"component": "apiserver",
				},
			},
			NamespaceSelector: promv1.NamespaceSelector{
				Any: true,
			},
			Endpoints: []promv1.Endpoint{
				{
					Port:        "https",
					Scheme:      "https",
					HonorLabels: true,
					RelabelConfigs: []*promv1.RelabelConfig{
						{
							SourceLabels: []string{"__meta_kubernetes_service_name"},
							TargetLabel:  "app",
						},
						{
							SourceLabels: []string{"__meta_kubernetes_pod_node_name"},
							TargetLabel:  "node",
						},
						{
							TargetLabel: "cluster_id",
							Replacement: cluster.GetName(),
						},
						{
							TargetLabel: "cluster_type",
							Replacement: key.ClusterType(cluster),
						},
						{
							TargetLabel: "installation",
							Replacement: installation,
						},
						{
							TargetLabel: "provider",
							Replacement: provider,
						},
						{
							SourceLabels: []string{"__meta_kubernetes_service_label_giantswarm_io_monitoring"},
							Regex:        "true",
							Action:       "drop",
						},
					},
				},
			},
		},
	}

	if !key.IsInCluster(cluster) {
		serviceMonitor.Spec.Endpoints[0].TLSConfig = &promv1.TLSConfig{
			CAFile:             fmt.Sprintf("/etc/prometheus/secrets/%s/ca", key.Secret()),
			CertFile:           fmt.Sprintf("/etc/prometheus/secrets/%s/crt", key.Secret()),
			KeyFile:            fmt.Sprintf("/etc/prometheus/secrets/%s/key", key.Secret()),
			InsecureSkipVerify: true,
		}
	} else {
		serviceMonitor.Spec.Endpoints[0].TLSConfig = &promv1.TLSConfig{
			CAFile:             key.ControlPlaneCAFile(),
			InsecureSkipVerify: true,
		}
		serviceMonitor.Spec.Endpoints[0].BearerTokenFile = key.ControlPlaneBearerToken()
	}

	return serviceMonitor
}
