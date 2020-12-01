package verticalpodautoscaler

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
	autoscalingv1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned/typed/autoscaling.k8s.io/v1"
)

type wrappedClient struct {
	client autoscalingv1.VerticalPodAutoscalerInterface
}

func (c wrappedClient) Create(ctx context.Context, o metav1.Object, options metav1.CreateOptions) (metav1.Object, error) {
	return c.client.Create(ctx, o.(*vpa_types.VerticalPodAutoscaler), options)
}

func (c wrappedClient) Update(ctx context.Context, o metav1.Object, options metav1.UpdateOptions) (metav1.Object, error) {
	return c.client.Update(ctx, o.(*vpa_types.VerticalPodAutoscaler), options)
}

func (c wrappedClient) Get(ctx context.Context, name string, options metav1.GetOptions) (metav1.Object, error) {
	return c.client.Get(ctx, name, options)
}

func (c wrappedClient) Delete(ctx context.Context, name string, options *metav1.DeleteOptions) error {
	return c.client.Delete(ctx, name, *options)
}
