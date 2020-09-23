package rules

import (
	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

func GetObjectMeta(obj interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(obj)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      "labelling-schema-rules",
		Namespace: key.Namespace(cluster),
		Labels: map[string]string{
			key.ClusterIDKey(): key.ClusterID(cluster),
		},
	}, nil
}

func LabellingSchemaValidationRule(obj interface{}) (metav1.Object, error) {
	objectMeta, err := GetObjectMeta(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &promv1.PrometheusRule{
		ObjectMeta: objectMeta,
		Spec: promv1.PrometheusRuleSpec{
			Groups: []promv1.RuleGroup{
				promv1.RuleGroup{
					Name: "labelling-schema",
					Rules: []promv1.Rule{
						promv1.Rule{
							Alert: "InvalidLabellingSchema",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: `up{cluster_type!~"control_plane|tenant_cluster"} or up{app=""} or up{installation=""} or up{cluster_id=""} or up{instance!~"(\\d+\\.\\d+\\.\\d+\\.\\d+:\\d+)||vault.*"} or up{provider!~"aws|azure|kvm"}`,
							},
							For: "10m",
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
				},
			},
		},
	}, nil
}
