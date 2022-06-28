package remotewritesecret

import (
	"context"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/api/v1alpha1"
)

const (
	Name = "remotewrite"
)

type Config struct {
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
}

type Resource struct {
	k8sClient        k8sclient.Interface
	logger           micrologger.Logger
	prometheusClient promclient.Interface
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		k8sClient:        config.K8sClient,
		logger:           config.Logger,
		prometheusClient: config.PrometheusClient,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func ToRemoteWrite(obj interface{}) (*pmov1alpha1.RemoteWrite, error) {
	remotewrite, ok := obj.(*pmov1alpha1.RemoteWrite)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "'%T' is not a 'pmov1alpha1.RemoteWrite'", obj)
	}

	return remotewrite, nil
}

func fetchPrometheusList(ctx context.Context, r *Resource, rw *pmov1alpha1.RemoteWrite) (*promv1.PrometheusList, error) {
	selector, err := metav1.LabelSelectorAsSelector(&rw.Spec.ClusterSelector)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	prometheusList, err := r.prometheusClient.
		MonitoringV1().
		Prometheuses(metav1.NamespaceAll).
		List(ctx, metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, microerror.Maskf(errorFetchingPrometheus, "Could not fetch Prometheus with label selector %#q", rw.Spec.ClusterSelector.String())
	}

	return prometheusList, nil
}

func (r *Resource) ensureRemoteWriteSecret(scSpec pmov1alpha1.RemoteWriteSecretSpec, meta metav1.ObjectMeta, ns string) corev1.Secret {

	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        scSpec.Name,
			Namespace:   ns,
			Labels:      meta.GetLabels(),
			Annotations: meta.GetAnnotations(),
		},
		Data: scSpec.Data,
	}
}
