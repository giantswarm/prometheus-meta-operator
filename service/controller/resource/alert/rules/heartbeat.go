package rules

import (
	"fmt"

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
				Alert: fmt.Sprintf("Heartbeat-%s-%s", installation, key.ClusterID(obj)),
				Expr:  intstr.FromString(`vector(1)`),
				Labels: map[string]string{
					"name": fmt.Sprintf("%s-%s", installation, key.ClusterID(obj)),
					"team": "atlas",
					"type": "heartbeat",
				},
				Annotations: map[string]string{
					"description": "This alert is used to ensure the entire alerting pipeline is functionnal.",
				},
			},
		},
	}
}
