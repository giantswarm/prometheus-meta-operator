package alert

import (
	"path"
	"reflect"
	"runtime"
	"testing"

	"github.com/giantswarm/micrologger"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cluster = &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster",
			Namespace: "cluster-namespace",
		},
	}

	installation = "installation"
)

func TestGetRules(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot get current filename")
	}

	path := path.Join(path.Dir(filename), "../../../..", ruleFilesPath)

	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		t.Fatal(err)
	}

	c := Config{
		Installation:     installation,
		PrometheusClient: &promclient.Clientset{},
		Logger:           logger,
		TemplatePath:     path,
	}

	r, err := New(c)
	if err != nil {
		t.Fatal(err)
	}

	rules, err := r.GetRules(cluster)
	if err != nil {
		t.Fatal(err)
	}

	if len(rules) <= 0 {
		t.Fatalf("no rules")
	}

	for _, r := range rules {
		valid := validateRule(r)
		if !valid {
			t.Errorf("rule %#q is invalid\n", r.GetName())
		}
	}
}

// validateRule validate a promv1.PrometheusRule object.
// This is to avoid adding invalid rule into the codebase.
// Rule are as follow:
// * object must not be empty
// * .spec.groups > 0
// * .spec.groups[].name != ""
// * .spec.groups[].rules > 0
// * .spec.groups[].rules.expr != ""
// taken from: https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#prometheusrule
func validateRule(rule *promv1.PrometheusRule) bool {
	if reflect.DeepEqual(rule, &promv1.PrometheusRule{}) {
		return false
	}

	if len(rule.Spec.Groups) <= 0 {
		return false
	}

	for _, g := range rule.Spec.Groups {
		if g.Name == "" {
			return false
		}
		if len(g.Rules) <= 0 {
			return false
		}

		for _, r := range g.Rules {
			if r.Expr.String() == "" {
				return false
			}
		}
	}

	return true
}
