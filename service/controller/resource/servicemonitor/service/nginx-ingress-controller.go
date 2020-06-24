package service

import (
	"fmt"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func NginxIngressController(cluster metav1.Object) *promv1.ServiceMonitor {
	return &promv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("nginx-ingress-controller-%s", cluster.GetName()),
			Namespace: key.Namespace(cluster),
			Labels: map[string]string{
				key.ClusterIDKey(): key.ClusterID(cluster),
			},
		},
		Spec: promv1.ServiceMonitorSpec{
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"k8s-app": "nginx-ingress-controller",
				},
			},
			NamespaceSelector: promv1.NamespaceSelector{
				MatchNames: []string{"kube-system"},
			},
			Endpoints: []promv1.Endpoint{
				promv1.Endpoint{
					Port:   "http",
					Scheme: "http",
				},
			},
		},
	}
}
