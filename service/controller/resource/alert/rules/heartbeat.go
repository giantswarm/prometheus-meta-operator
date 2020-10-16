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
				Expr:  intstr.FromString(`vector(1)`),
				Labels: map[string]string{
					key.ClusterIDKey():    key.ClusterID(obj),
					key.InstallationKey(): installation,
					key.TypeKey():         key.Heartbeat(),
					"team":                "atlas",
				},
				Annotations: map[string]string{
					"description": "This alert is used to ensure the entire alerting pipeline is functionnal.",
				},
			},
		},
	}
}
