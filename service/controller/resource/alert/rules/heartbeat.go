package rules

import (
	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func Heartbeat(obj metav1.Object, installation string) promv1.RuleGroup {
	return promv1.RuleGroup{
		Name: "heartbeat",
		Rules: []promv1.Rule{
			promv1.Rule{
				Alert: "Heartbeat",
				Expr:  intstr.FromString(`time() - max(container_start_time_seconds{cluster_type="host", namespace="monitoring", container="prometheus"}) > 20 * 60`),
				Labels: map[string]string{
					"cluster":      key.ClusterID(obj),
					"installation": installation,
					"type":         "heartbeat",
					"team":         "atlas",
				},
				Annotations: map[string]string{
					"description": "This alert is used to ensure the entire alerting pipeline is functionnal.",
				},
			},
		},
	}
}
