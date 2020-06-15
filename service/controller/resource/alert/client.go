package alert

import (
	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	monv1 "github.com/coreos/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type wrappedClient struct {
	client monv1.PrometheusRuleInterface
}

func (c wrappedClient) Create(o metav1.Object) (metav1.Object, error) {
	return c.client.Create(o.(*promv1.PrometheusRule))
}

func (c wrappedClient) Update(o metav1.Object) (metav1.Object, error) {
	return c.client.Update(o.(*promv1.PrometheusRule))
}

func (c wrappedClient) Get(name string, options metav1.GetOptions) (metav1.Object, error) {
	return c.client.Get(name, options)
}

func (c wrappedClient) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(name, options)
}
