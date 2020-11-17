package prometheus

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
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

	Address           string
	Bastions          []string
	CreatePVC         bool
	Customer          string
	Installation      string
	Pipeline          string
	Provider          string
	Region            string
	StorageSize       string
	RetentionDuration string
	RetentionSize     string
	RemoteWriteURL    string
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
	// Configured following https://github.com/prometheus-operator/prometheus-operator/issues/541#issuecomment-451884171
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
				key.ClusterIDKey():    key.ClusterID(cluster),
				"cluster_type":        key.ClusterType(cluster),
				"customer":            config.Customer,
				key.InstallationKey(): config.Installation,
				"pipeline":            config.Pipeline,
				"provider":            config.Provider,
				"region":              config.Region,
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
					// memory: 5Gi
					corev1.ResourceMemory: *resource.NewQuantity(5000*1024*1024, resource.BinarySI),
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
			AdditionalScrapeConfigs: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: key.PrometheusAdditionalScrapeConfigsSecretName(),
				},
				Key: key.PrometheusAdditionalScrapeConfigsName(),
			},
			AdditionalAlertManagerConfigs: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: key.AlertManagerSecretName(),
				},
				Key: key.AlertManagerKey(),
			},
			SecurityContext: &v1.PodSecurityContext{
				RunAsUser:    &uid,
				RunAsGroup:   &gid,
				RunAsNonRoot: &runAsNonRoot,
				FSGroup:      &fsGroup,
			},
			Storage: &storage,
			Affinity: &v1.Affinity{
				NodeAffinity: &v1.NodeAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{
						v1.PreferredSchedulingTerm{
							Weight: 60,
							Preference: v1.NodeSelectorTerm{
								MatchExpressions: []v1.NodeSelectorRequirement{
									v1.NodeSelectorRequirement{
										Key:      "role",
										Operator: v1.NodeSelectorOpNotIn,
										Values: []string{
											"master",
										},
									},
								},
							},
						},
					},
				},
			},
			RuleNamespaceSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": key.NamespaceMonitoring(cluster),
				},
			},
		},
	}

	if config.RemoteWriteURL != "" {
		prometheus.Spec.RemoteWrite = []promv1.RemoteWriteSpec{
			promv1.RemoteWriteSpec{
				URL: config.RemoteWriteURL,
				BasicAuth: &promv1.BasicAuth{
					Username: v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: key.RemoteWriteSecretName(),
						},
						Key: key.RemoteWriteUsernameKey(),
					},
					Password: v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: key.RemoteWriteSecretName(),
						},
						Key: key.RemoteWritePasswordKey(),
					},
				},
				Name: key.ClusterID(cluster),
				WriteRelabelConfigs: []promv1.RelabelConfig{
					promv1.RelabelConfig{
						SourceLabels: []string{"__name__"},
						Regex:        "(^aggregation:.+|prometheus_tsdb_head_series|^slo_.+)",
						Action:       "keep",
					},
				},
			},
		}
	}

	if !key.IsInCluster(cluster) {
		// Tenant cluster
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

		prometheus.Spec.RuleSelector = &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				metav1.LabelSelectorRequirement{
					Key:      "cluster_type",
					Operator: metav1.LabelSelectorOpNotIn,
					Values:   []string{"control_plane"},
				},
			},
		}
	} else {
		// Control plane
		prometheus.Spec.APIServerConfig = &promv1.APIServerConfig{
			Host:            fmt.Sprintf("https://%s", key.APIUrl(cluster)),
			BearerTokenFile: key.ControlPlaneBearerToken(),
			TLSConfig: &promv1.TLSConfig{
				CAFile: key.ControlPlaneCAFile(),
				SafeTLSConfig: promv1.SafeTLSConfig{
					InsecureSkipVerify: true,
				},
			},
		}

		prometheus.Spec.Secrets = []string{
			key.EtcdSecret(cluster),
		}

		// empty LabelSelector matches all objects.
		prometheus.Spec.RuleSelector = &metav1.LabelSelector{}
	}

	return prometheus, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*promv1.Prometheus)
	d := desired.(*promv1.Prometheus)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
