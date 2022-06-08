package promremotewrite

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/google/go-cmp/cmp"
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

func toPrometheusRemoteWrite(r pmov1alpha1.RemoteWrite, p promv1.Prometheus) (*promv1.Prometheus, bool) {

	return ensurePrometheusRemoteWrite(r, p)
}

func ToRemoteWrite(obj interface{}) (*pmov1alpha1.RemoteWrite, error) {
	remotewrite, ok := obj.(*pmov1alpha1.RemoteWrite)
	if !ok {
		return nil, microerror.Maskf(wrongTypeError, "'%T' is not a 'pmov1alpha1.RemoteWrite'", obj)
	}

	return remotewrite, nil
}

func ensurePrometheusRemoteWrite(r pmov1alpha1.RemoteWrite, p promv1.Prometheus) (*promv1.Prometheus, bool) {
	r.Spec.RemotWrite.Name = r.GetName()
	if p.Spec.RemoteWrite != nil {
		if rwIndex, ok := isRemoteWriteExists(r.GetName(), p.Spec.RemoteWrite); !ok { // item not found
			p.Spec.RemoteWrite = append(p.Spec.RemoteWrite, r.Spec.RemotWrite)
			return &p, true
		} else if !cmp.Equal(r.Spec.RemotWrite, p.Spec.RemoteWrite[rwIndex]) { //  item found
			p.Spec.RemoteWrite[rwIndex] = r.Spec.RemotWrite
			return &p, true
		} else {
			// no update!!
			return &p, false
		}
	} else {
		p.Spec.RemoteWrite = []promv1.RemoteWriteSpec{r.Spec.RemotWrite}
		return &p, true
	}
}

func omitPrometheusRemoteWrite(r pmov1alpha1.RemoteWrite, p promv1.Prometheus) (*promv1.Prometheus, bool) {
	r.Spec.RemotWrite.Name = r.GetName()
	if p.Spec.RemoteWrite != nil {
		if rwIndex, ok := isRemoteWriteExists(r.GetName(), p.Spec.RemoteWrite); ok { // item found
			p.Spec.RemoteWrite = remove(p.Spec.RemoteWrite, rwIndex)
			return &p, true
		}
	}
	return &p, false
}

// isRemoteWriteExists checks if the item exists and return the item index
func isRemoteWriteExists(name string, items []promv1.RemoteWriteSpec) (int, bool) {
	for i, item := range items {
		if name == item.Name {
			return i, true
		}
	}
	return -1, false
}

func remove(slice []promv1.RemoteWriteSpec, s int) []promv1.RemoteWriteSpec {
	return append(slice[:s], slice[s+1:]...)
}
