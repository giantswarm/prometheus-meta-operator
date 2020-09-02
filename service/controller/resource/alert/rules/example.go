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
		Name:      "example-rules",
		Namespace: key.Namespace(cluster),
		Labels: map[string]string{
			key.ClusterIDKey(): key.ClusterID(cluster),
		},
	}, nil
}

func ExampleRule(obj interface{}) (*promv1.PrometheusRule, error) {
	objectMeta, err := GetObjectMeta(obj)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return &promv1.PrometheusRule{
		ObjectMeta: objectMeta,
		Spec: promv1.PrometheusRuleSpec{
			Groups: []promv1.RuleGroup{
				promv1.RuleGroup{
					Name: "apiserver",
					Rules: []promv1.Rule{
						promv1.Rule{
							Alert: "APIServerLatencyTooHigh",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: `histogram_quantile(0.95, sum(rate(apiserver_request_latencies_bucket{subresource!~"log", verb=~"DELETE|GET|PATCH|POST|PUT"}[1h])) WITHOUT (instance, resource)) / 1e+06 > 3`,
							},
							For: "3h",
							Labels: map[string]string{
								"l_if_cluster_status_creating":              "true",
								"cancel_if_cluster_status_updating":         "true",
								"cancel_if_cluster_status_deleting":         "true",
								"cancel_if_cluster_with_notready_nodepools": "true",
								"severity": "notify",
								"area":     "kaas",
								"team":     "ludacris",
								"topic":    "kubernetes",
							},
							Annotations: map[string]string{
								"description": "Kubernetes API Server {{ $labels.verb }} request latency is too high.",
								"opsrecipe":   "apiserver-overloaded.md",
							},
						},
					},
				},
			},
		},
	}, nil
}
