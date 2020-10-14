package rules

import (
	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Heartbeat(obj interface{}) promv1.RuleGroup {
	return promv1.RuleGroup{
		Name: "heartbeat",
		Rules: []promv1.Rule{
			promv1.Rule{
				Alert: "Heartbeat",
				Expr:  intstr.FromString(`time() - max(container_start_time_seconds{cluster_type="host", namespace="monitoring", container="prometheus"}) > 60`),
				Labels: map[string]string{
					"severity": "heartbeat",
					"team":     "atlas",
				},
				Annotations: map[string]string{
					"description": "This alert is used to ensure the entire alerting pipeline is functionnal.",
				},
			},
		},
	}
}
