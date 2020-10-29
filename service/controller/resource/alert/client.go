package alert

import (
	"context"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monv1 "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type wrappedClient struct {
	client monv1.PrometheusRuleInterface
}

func (c wrappedClient) Create(ctx context.Context, o metav1.Object, options metav1.CreateOptions) (metav1.Object, error) {
	return c.client.Create(ctx, o.(*promv1.PrometheusRule), options)
}

func (c wrappedClient) Update(ctx context.Context, o metav1.Object, options metav1.UpdateOptions) (metav1.Object, error) {
	return c.client.Update(ctx, o.(*promv1.PrometheusRule), options)
}

func (c wrappedClient) Get(ctx context.Context, name string, options metav1.GetOptions) (metav1.Object, error) {
	return c.client.Get(ctx, name, options)
}

func (c wrappedClient) Delete(ctx context.Context, name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(ctx, name, *options)
}
