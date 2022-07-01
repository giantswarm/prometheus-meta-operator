package remotewriteutils

import (
	"context"
	"reflect"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/api/v1alpha1"
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

	kind := reflect.TypeOf(remotewrite).Elem().Name()
	remotewrite.TypeMeta = metav1.TypeMeta{
		Kind:       kind,
		APIVersion: pmov1alpha1.SchemeGroupVersion.String(),
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

	return prometheusList, nil
}