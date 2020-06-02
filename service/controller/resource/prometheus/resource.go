package prometheus

import (
	"fmt"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "prometheus"
)

type Config struct {
	PrometheusClient promclient.Interface
	Logger           micrologger.Logger
}

type Resource struct {
	prometheusClient promclient.Interface
	logger           micrologger.Logger
}

func New(config Config) (*Resource, error) {
	if config.PrometheusClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.PrometheusClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &Resource{
		prometheusClient: config.PrometheusClient,
		logger:           config.Logger,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func toPrometheus(v interface{}) (*promv1.Prometheus, error) {
	if v == nil {
		return nil, nil
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	name := cluster.GetName()
	var replicas int32 = 2

	prometheus := &promv1.Prometheus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: key.Namespace(cluster),
		},
		Spec: promv1.PrometheusSpec{
			APIServerConfig: &promv1.APIServerConfig{
				Host: fmt.Sprintf("https://master.%s", name),
				TLSConfig: &promv1.TLSConfig{
					CAFile:   fmt.Sprintf("/etc/prometheus/secrets/%s/ca", key.Secret()),
					CertFile: fmt.Sprintf("/etc/prometheus/secrets/%s/crt", key.Secret()),
					KeyFile:  fmt.Sprintf("/etc/prometheus/secrets/%s/key", key.Secret()),
				},
			},
			Replicas: &replicas,
			Secrets: []string{
				key.Secret(),
			},
			ServiceMonitorSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					key.ClusterIDKey(): key.ClusterID(cluster),
				},
			},
			RuleSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					key.ClusterIDKey(): key.ClusterID(cluster),
				},
			},
		},
	}

	return prometheus, nil
}
