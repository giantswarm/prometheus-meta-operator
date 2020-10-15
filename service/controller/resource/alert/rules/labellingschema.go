package rules

import (
	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func LabellingSchemaValidationRule(obj metav1.Object) promv1.RuleGroup {
	return promv1.RuleGroup{
		Name: "labelling-schema",
		Rules: []promv1.Rule{
			promv1.Rule{
				Alert: "InvalidLabellingSchema",
				Expr:  intstr.FromString(`up{cluster_type!~"control_plane|tenant_cluster"} or up{app=""} or up{installation=""} or up{cluster_id=""} or up{instance!~"(\\d+\\.\\d+\\.\\d+\\.\\d+:\\d+)||vault.*"} or up{provider!~"aws|azure|kvm"}`),
				For:   "10m",
				Labels: map[string]string{
					"cancel_if_cluster_status_creating":         "true",
					"cancel_if_cluster_status_updating":         "true",
					"cancel_if_cluster_status_deleting":         "true",
					"cancel_if_cluster_with_notready_nodepools": "true",
					"severity": "notify",
					"area":     "empowerment",
					"team":     "atlas",
					"topic":    "observability",
				},
				Annotations: map[string]string{
					"description": "Labelling schema is invalid.",
				},
			},
		},
	}
}
