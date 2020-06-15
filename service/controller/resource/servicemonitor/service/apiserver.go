package service

import (
	"fmt"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func APIServer(cluster metav1.Object) *promv1.ServiceMonitor {
	return &promv1.ServiceMonitor{
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
				promv1.Endpoint{
					Port:   "https",
					Scheme: "https",
					TLSConfig: &promv1.TLSConfig{
						CAFile:             fmt.Sprintf("/etc/prometheus/secrets/%s/ca", key.Secret()),
						CertFile:           fmt.Sprintf("/etc/prometheus/secrets/%s/crt", key.Secret()),
						KeyFile:            fmt.Sprintf("/etc/prometheus/secrets/%s/key", key.Secret()),
						InsecureSkipVerify: true,
					},
				},
			},
		},
	}
}
