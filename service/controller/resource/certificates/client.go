package certificates

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type wrappedClient struct {
	client clientcorev1.SecretInterface
}

func (c wrappedClient) Create(o metav1.Object) (metav1.Object, error) {
	return c.client.Create(o.(*corev1.Secret))
}

func (c wrappedClient) Update(o metav1.Object) (metav1.Object, error) {
	return c.client.Update(o.(*corev1.Secret))
}

func (c wrappedClient) Get(name string, options metav1.GetOptions) (metav1.Object, error) {
	return c.client.Get(name, options)
}

func (c wrappedClient) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(name, options)
}
