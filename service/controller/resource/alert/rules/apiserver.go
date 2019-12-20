package rules

import (
	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/api/v1alpha2"
)

func APIServer(cluster *v1alpha2.Cluster) *promv1.PrometheusRule {
	return &promv1.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "alerts",
			Namespace: key.Namespace(cluster),
			Labels: map[string]string{
				key.ClusterIDKey(): key.ClusterID(cluster),
			},
		},
		Spec: promv1.PrometheusRuleSpec{
			Groups: []promv1.RuleGroup{
				promv1.RuleGroup{
					Name:     "",
					Interval: "",
					Rules: []promv1.Rules{
						promv1.Rule{
							Record:      "",
							Alert:       "",
							Expr:        "",
							For:         "",
							Labels:      "",
							Annotations: "",
						},
					},
				},
			},
		},
	}
}
