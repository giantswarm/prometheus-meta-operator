package prometheus

import (
	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
)

type wrappedClient struct {
	client monv1.PrometheusInterface
}

func (c wrappedClient) Create(o metav1.Object) (metav1.Object, error) {
	return c.client.Create(o.(*promv1.Prometheus))
}
func (c wrappedClient) Update(o metav1.Object) (metav1.Object, error) {
	return c.client.Update(o.(*promv1.Prometheus))
}
func (c wrappedClient) Get(name string, options metav1.GetOptions) (metav1.Object, error) {
	return c.client.Get(name, options)
}
func (c wrappedClient) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(name, options)
}
