package prometheusremotewrite

import (
	"context"
	"strings"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/google/go-cmp/cmp"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
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

	HTTPProxy  string
	HTTPSProxy string
	NoProxy    string
}

type Resource struct {
	k8sClient        k8sclient.Interface
	logger           micrologger.Logger
	prometheusClient promclient.Interface

	HTTPProxy  string
	HTTPSProxy string
	NoProxy    string
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		k8sClient:        config.K8sClient,
		logger:           config.Logger,
		prometheusClient: config.PrometheusClient,

		HTTPProxy:  config.HTTPProxy,
		HTTPSProxy: config.HTTPSProxy,
		NoProxy:    config.NoProxy,
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

func ensurePrometheusRemoteWrite(r *Resource, rw pmov1alpha1.RemoteWrite, p promv1.Prometheus) (*promv1.Prometheus, bool) {
	rw.Spec.RemoteWrite.Name = rw.GetName()

	if !strings.Contains(r.NoProxy, rw.Spec.RemoteWrite.URL) {
		if len(r.HTTPSProxy) > 0 {
			rw.Spec.RemoteWrite.ProxyURL = r.HTTPSProxy
		} else if len(r.HTTPProxy) > 0 {
			rw.Spec.RemoteWrite.ProxyURL = r.HTTPProxy
		}
	}

	if p.Spec.RemoteWrite != nil {
		if rwIndex, ok := remoteWriteExists(rw.GetName(), p.Spec.RemoteWrite); !ok { // item not found
			p.Spec.RemoteWrite = append(p.Spec.RemoteWrite, rw.Spec.RemoteWrite)
			return &p, true
		} else if !cmp.Equal(rw.Spec.RemoteWrite, p.Spec.RemoteWrite[rwIndex]) { //  item found
			p.Spec.RemoteWrite[rwIndex] = rw.Spec.RemoteWrite
			return &p, true
		} else {
			// no update!!
			return &p, false
		}
	} else {
		p.Spec.RemoteWrite = []promv1.RemoteWriteSpec{rw.Spec.RemoteWrite}
		return &p, true
	}
}

func removePrometheusRemoteWrite(r pmov1alpha1.RemoteWrite, p promv1.Prometheus) (*promv1.Prometheus, bool) {
	r.Spec.RemoteWrite.Name = r.GetName()
	if p.Spec.RemoteWrite != nil {
		if rwIndex, ok := remoteWriteExists(r.GetName(), p.Spec.RemoteWrite); ok { // item found
			p.Spec.RemoteWrite = remove(p.Spec.RemoteWrite, rwIndex)
			return &p, true
		}
	}
	return &p, false
}

// remoteWriteExists checks if the item exists and return the item index
func remoteWriteExists(name string, items []promv1.RemoteWriteSpec) (int, bool) {
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
