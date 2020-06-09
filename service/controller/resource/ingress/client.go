package ingress

import (
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	extv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

type wrappedClient struct {
	client extv1beta1.IngressInterface
}

func (c wrappedClient) Create(o metav1.Object) (metav1.Object, error) {
	return c.client.Create(o.(*v1beta1.Ingress))
}
func (c wrappedClient) Update(o metav1.Object) (metav1.Object, error) {
	return c.client.Update(o.(*v1beta1.Ingress))
}
func (c wrappedClient) Get(name string, options metav1.GetOptions) (metav1.Object, error) {
	return c.client.Get(name, options)
}
func (c wrappedClient) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(name, options)
}
