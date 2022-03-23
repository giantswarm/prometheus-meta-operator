package heartbeatwebhookconfig

import (
	"golang.org/x/net/context"

	monitoringv1alpha1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1alpha1"
	monitoringv1alpha1client "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned/typed/monitoring/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type wrappedClient struct {
	client monitoringv1alpha1client.AlertmanagerConfigInterface
}

func (c wrappedClient) Create(ctx context.Context, object metav1.Object, options metav1.CreateOptions) (metav1.Object, error) {
	return c.client.Create(ctx, object.(*monitoringv1alpha1.AlertmanagerConfig), options)
}
func (c wrappedClient) Update(ctx context.Context, object metav1.Object, options metav1.UpdateOptions) (metav1.Object, error) {
	return c.client.Update(ctx, object.(*monitoringv1alpha1.AlertmanagerConfig), options)
}
func (c wrappedClient) Get(ctx context.Context, name string, options metav1.GetOptions) (metav1.Object, error) {
	return c.client.Get(ctx, name, options)
}
func (c wrappedClient) Delete(ctx context.Context, name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(ctx, name, *options)
}
