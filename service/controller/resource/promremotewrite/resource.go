package promremotewrite

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"

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

func toPrometheusRemoteWrite(r pmov1alpha1.RemoteWrite, p promv1.Prometheus) (*promv1.Prometheus, error) {

	if p.Spec.RemoteWrite != nil {
		p.Spec.RemoteWrite = append(p.Spec.RemoteWrite, r.Spec.RemotWrite)
	} else {
		p.Spec.RemoteWrite = []promv1.RemoteWriteSpec{r.Spec.RemotWrite}
	}

	return &p, nil
}

func ToRemoteWrite(obj interface{}) (pmov1alpha1.RemoteWrite, error) {
	remotewrite, ok := obj.(pmov1alpha1.RemoteWrite)
	if !ok {
		return pmov1alpha1.RemoteWrite{}, microerror.Maskf(wrongTypeError, "'%T' is not a 'pmov1alpha1.RemoteWrite'", obj)
	}

	return remotewrite, nil
}
