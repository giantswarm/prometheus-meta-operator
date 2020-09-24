package ingress

import (
	"golang.org/x/net/context"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientextv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

type wrappedClient struct {
	client clientextv1beta1.IngressInterface
}

func (c wrappedClient) Create(ctx context.Context, object metav1.Object, options metav1.CreateOptions) (metav1.Object, error) {
	return c.client.Create(ctx, object.(*v1beta1.Ingress), options)
}
func (c wrappedClient) Update(ctx context.Context, object metav1.Object, options metav1.UpdateOptions) (metav1.Object, error) {
	return c.client.Update(ctx, object.(*v1beta1.Ingress), options)
}
func (c wrappedClient) Get(ctx context.Context, name string, options metav1.GetOptions) (metav1.Object, error) {
	return c.client.Get(ctx, name, options)
}
func (c wrappedClient) Delete(ctx context.Context, name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(ctx, name, *options)
}
