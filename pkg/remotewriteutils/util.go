package remotewriteutils

import (
	"context"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/v2/api/v1alpha1"
)

type ResourceWrapper struct {
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
}

func ToRemoteWrite(obj interface{}) (*pmov1alpha1.RemoteWrite, error) {
	remotewrite, ok := obj.(*pmov1alpha1.RemoteWrite)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "'%T' is not a 'pmov1alpha1.RemoteWrite'", obj)
	}

	return remotewrite, nil
}

func FetchPrometheusList(ctx context.Context, r *ResourceWrapper, rw *pmov1alpha1.RemoteWrite) (*promv1.PrometheusList, error) {
	selector, err := metav1.LabelSelectorAsSelector(&rw.Spec.ClusterSelector)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	prometheusList, err := r.PrometheusClient.
		MonitoringV1().
		Prometheuses(metav1.NamespaceAll).
		List(ctx, metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, microerror.Maskf(errorFetchingPrometheus, "Could not fetch Prometheus with label selector %#q", rw.Spec.ClusterSelector.String())
	}

	for i, p := range prometheusList.Items {
		if p.GetNamespace() == metav1.NamespaceSystem {
			// Remove the prometheus-agent from the list.
			prometheusList.Items = append(prometheusList.Items[:i], prometheusList.Items[i+1:]...)
		}
	}

	return prometheusList, nil
}
