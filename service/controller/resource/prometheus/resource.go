package prometheus

import (
	"fmt"
	"reflect"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "prometheus"
)

type Config struct {
	PrometheusClient promclient.Interface
	Logger           micrologger.Logger
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.PrometheusClient.MonitoringV1().Prometheuses(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:     clientFunc,
		Logger:         config.Logger,
		Name:           Name,
		ToCR:           toPrometheus,
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func toPrometheus(v interface{}) (metav1.Object, error) {
	if v == nil {
		return nil, nil
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	name := cluster.GetName()
	var replicas int32 = 1

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
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					// cpu: 100m
					corev1.ResourceCPU: *resource.NewMilliQuantity(100, resource.DecimalSI),
					// memory: 100Mi
					corev1.ResourceMemory: *resource.NewQuantity(1*1024*1024*1024, resource.BinarySI),
				},
				Requests: corev1.ResourceList{
					// cpu: 100m
					corev1.ResourceCPU: *resource.NewMilliQuantity(100, resource.DecimalSI),
					// memory: 100Mi
					corev1.ResourceMemory: *resource.NewQuantity(100*1024*1024, resource.BinarySI),
				},
			},
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
			AdditionalScrapeConfigs: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: key.PrometheusAdditionalScrapeConfigsSecretName(),
				},
				Key: key.PrometheusAdditionalScrapeConfigsName(),
			},
		},
	}

	return prometheus, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*promv1.Prometheus)
	d := desired.(*promv1.Prometheus)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
