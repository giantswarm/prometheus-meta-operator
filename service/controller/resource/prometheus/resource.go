package prometheus

import (
	"fmt"
	"net/url"
	"reflect"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "prometheus"
)

type Config struct {
	Address          string
	PrometheusClient promclient.Interface
	Logger           micrologger.Logger

	CreatePVC         bool
	StorageSize       string
	RetentionDuration string
	RetentionSize     string
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.PrometheusClient.MonitoringV1().Prometheuses(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:    clientFunc,
		Logger:        config.Logger,
		Name:          Name,
		GetObjectMeta: getObjectMeta,
		GetDesiredObject: func(v interface{}) (metav1.Object, error) {
			return toPrometheus(v, config)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.ClusterID(cluster),
		Namespace: key.Namespace(cluster),
		Labels:    key.Labels(cluster),
	}, nil
}

func toPrometheus(v interface{}, config Config) (metav1.Object, error) {
	if v == nil {
		return nil, nil
	}

	objectMeta, err := getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	address, err := url.Parse(config.Address)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var replicas int32 = 1
	// Configured following https://github.com/coreos/prometheus-operator/issues/541#issuecomment-451884171
	// as the volume could not mount otherwise
	var uid int64 = 1000
	var fsGroup int64 = 2000
	var runAsNonRoot bool = true
	// Prometheus default image runs using the nobody user (65534)
	var gid int64 = 65534

	var walCompression bool = true

	var storage promv1.StorageSpec
	storageSize := resource.MustParse(config.StorageSize)

	if config.CreatePVC {
		storage = promv1.StorageSpec{
			VolumeClaimTemplate: promv1.EmbeddedPersistentVolumeClaim{
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: storageSize,
						},
					},
				},
			},
		}
	} else {
		storage = promv1.StorageSpec{
			EmptyDir: &v1.EmptyDirVolumeSource{
				SizeLimit: &storageSize,
			},
		}
	}

	externalURL, err := address.Parse("/" + key.ClusterID(cluster))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	labels := make(map[string]string)
	for k, v := range key.Labels(cluster) {
		labels[k] = v
	}

	labels["giantswarm.io/monitoring"] = "true"

	prometheus := &promv1.Prometheus{
		ObjectMeta: objectMeta,
		Spec: promv1.PrometheusSpec{
			ExternalLabels: map[string]string{
				key.ClusterIDKey(): key.ClusterID(cluster),
				"cluster_type":     key.ClusterType(cluster),
			},
			ExternalURL: externalURL.String(),
			RoutePrefix: fmt.Sprintf("/%s", key.ClusterID(cluster)),
			PodMetadata: &promv1.EmbeddedObjectMetadata{
				Labels: labels,
			},
			Replicas: &replicas,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					// cpu: 100m
					corev1.ResourceCPU: *resource.NewMilliQuantity(100, resource.DecimalSI),
					// memory: 100Mi
					corev1.ResourceMemory: *resource.NewQuantity(100*1024*1024, resource.BinarySI),
				},
			},
			Retention:      config.RetentionDuration,
			RetentionSize:  config.RetentionSize,
			WALCompression: &walCompression,
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
			SecurityContext: &v1.PodSecurityContext{
				RunAsUser:    &uid,
				RunAsGroup:   &gid,
				RunAsNonRoot: &runAsNonRoot,
				FSGroup:      &fsGroup,
			},
			Storage: &storage,
		},
	}

	if !key.IsInCluster(cluster) {
		prometheus.Spec.APIServerConfig = &promv1.APIServerConfig{
			Host: fmt.Sprintf("https://%s", key.APIUrl(cluster)),
			TLSConfig: &promv1.TLSConfig{
				CAFile:   fmt.Sprintf("/etc/prometheus/secrets/%s/ca", key.Secret()),
				CertFile: fmt.Sprintf("/etc/prometheus/secrets/%s/crt", key.Secret()),
				KeyFile:  fmt.Sprintf("/etc/prometheus/secrets/%s/key", key.Secret()),
			},
		}

		prometheus.Spec.Secrets = []string{
			key.Secret(),
		}
	} else {
		prometheus.Spec.APIServerConfig = &promv1.APIServerConfig{
			Host:            fmt.Sprintf("https://%s", key.APIUrl(cluster)),
			BearerTokenFile: key.ControlPlaneBearerToken(),
			TLSConfig: &promv1.TLSConfig{
				CAFile:             key.ControlPlaneCAFile(),
				InsecureSkipVerify: true,
			},
		}

		prometheus.Spec.Secrets = []string{
			key.EtcdSecret(cluster),
		}
	}

	prometheus.Spec.Alerting = alertManagerConfig()

	return prometheus, nil
}

func alertManagerConfig() *promv1.AlertingSpec {
	return &promv1.AlertingSpec{
		Alertmanagers: []promv1.AlertmanagerEndpoints{
			promv1.AlertmanagerEndpoints{
				Namespace: "monitoring",
				Name:      "alertmanager",
				Port:      intstr.FromInt(9093),
				Scheme:    "http",
			},
		},
	}
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*promv1.Prometheus)
	d := desired.(*promv1.Prometheus)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
