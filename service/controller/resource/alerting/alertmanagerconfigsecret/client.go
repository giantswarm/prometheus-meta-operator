package alertmanagerconfigsecret

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type wrappedClient struct {
	client v1.SecretInterface
}

func (c wrappedClient) Create(ctx context.Context, o metav1.Object, options metav1.CreateOptions) (metav1.Object, error) {
	return c.client.Create(ctx, o.(*corev1.Secret), options)
}
func (c wrappedClient) Update(ctx context.Context, o metav1.Object, options metav1.UpdateOptions) (metav1.Object, error) {
	return c.client.Update(ctx, o.(*corev1.Secret), options)
}
func (c wrappedClient) Get(ctx context.Context, name string, options metav1.GetOptions) (metav1.Object, error) {
	return c.client.Get(ctx, name, options)
}
func (c wrappedClient) Delete(ctx context.Context, name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(ctx, name, *options)
}
