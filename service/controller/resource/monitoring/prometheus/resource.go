package prometheus

import (
	"fmt"
	"net/url"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/pvcresizing"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "prometheus"
)

type Config struct {
	PrometheusClient promclient.Interface
	Logger           micrologger.Logger

	Address            string
	Bastions           []string
	Customer           string
	EvaluationInterval string
	ImageRepository    string
	Installation       string
	Pipeline           string
	Provider           string
	Region             string
	Registry           string
	LogLevel           string
	RetentionDuration  string
	ScrapeInterval     string
	Version            string
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
		GetDesiredObject: func(ctx context.Context, v interface{}) (metav1.Object, error) {
			return toPrometheus(ctx, v, config)
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(ctx context.Context, v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.ClusterID(cluster),
		Namespace: key.Namespace(cluster),
		Labels:    key.PrometheusLabels(cluster),
	}, nil
}

func toPrometheus(ctx context.Context, v interface{}, config Config) (metav1.Object, error) {
	if v == nil {
		return nil, nil
	}

	objectMeta, err := getObjectMeta(ctx, v)
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

	annotationValue := cluster.GetAnnotations()[key.PrometheusVolumeSizeAnnotation]
	volumeSize := pvcresizing.PrometheusVolumeSize(annotationValue)
	storageSize := resource.MustParse(volumeSize)

	storage := promv1.StorageSpec{
		VolumeClaimTemplate: promv1.EmbeddedPersistentVolumeClaim{
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: storageSize,
					},
				},
			},
		},
	}

	externalURL, err := address.Parse("/" + key.ClusterID(cluster))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	labels := make(map[string]string)
	for k, v := range key.PrometheusLabels(cluster) {
		labels[k] = v
	}

	labels[key.MonitoringLabel] = "true"

	image := fmt.Sprintf("%s/%s:%s", config.Registry, config.ImageRepository, config.Version)
	pageTitle := fmt.Sprintf("%s/%s Prometheus", config.Installation, key.ClusterID(cluster))
	prometheus := &promv1.Prometheus{
		ObjectMeta: objectMeta,
		Spec: promv1.PrometheusSpec{
			// We need to use this to connect each WC prometheus with the central alertmanager instead of the alerting section of the Prometheus CR
			// because the alerting section tries to find the alertmanager service in the workload cluster and not in the management cluster
			// as it is using the secrets defined under prometheus.Spec.APIServerConfig.
			//
			// This forces us to use the static config defined in resource/alerting/alertmanagerwiring.
			AdditionalAlertManagerConfigs: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: key.AlertmanagerSecretName(),
				},
				Key: key.AlertmanagerKey(),
			},

			CommonPrometheusFields: promv1.CommonPrometheusFields{
				AdditionalScrapeConfigs: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: key.PrometheusAdditionalScrapeConfigsSecretName(),
					},
					Key: key.PrometheusAdditionalScrapeConfigsName(),
				},
				Affinity: &corev1.Affinity{
					NodeAffinity: &corev1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
							NodeSelectorTerms: []corev1.NodeSelectorTerm{
								{
									MatchExpressions: []corev1.NodeSelectorRequirement{
										{
											Key:      "node-role.kubernetes.io/control-plane",
											Operator: corev1.NodeSelectorOpDoesNotExist,
										},
									},
								},
							},
						},
					},
				},
				EnableFeatures: []string{"remote-write-receiver"},
				ExternalLabels: map[string]string{
					key.ClusterIDKey:    key.ClusterID(cluster),
					key.ClusterTypeKey:  key.ClusterType(config.Installation, cluster),
					key.CustomerKey:     config.Customer,
					key.InstallationKey: config.Installation,
					key.PipelineKey:     config.Pipeline,
					key.ProviderKey:     config.Provider,
					key.RegionKey:       config.Region,
				},
				ExternalURL: externalURL.String(),
				Image:       &image,
				LogLevel:    config.LogLevel,
				PodMetadata: &promv1.EmbeddedObjectMetadata{
					Labels: labels,
				},
				PriorityClassName: "prometheus",
				Replicas:          &replicas,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						// cpu: 100m
						corev1.ResourceCPU: *key.PrometheusDefaultCPU(),
						// memory: 1Gi
						corev1.ResourceMemory: *key.PrometheusDefaultMemory(),
					},
					Limits: corev1.ResourceList{
						// cpu: 150m
						corev1.ResourceCPU: *key.PrometheusDefaultCPULimit(),
						// memory: 1Gi
						corev1.ResourceMemory: *key.PrometheusDefaultMemoryLimit(),
					},
				},
				RoutePrefix:    fmt.Sprintf("/%s", key.ClusterID(cluster)),
				ScrapeInterval: promv1.Duration(config.ScrapeInterval),
				SecurityContext: &corev1.PodSecurityContext{
					RunAsUser:    &uid,
					RunAsGroup:   &gid,
					RunAsNonRoot: &runAsNonRoot,
					FSGroup:      &fsGroup,
				},
				Storage: &storage,
				TopologySpreadConstraints: []corev1.TopologySpreadConstraint{
					{
						MaxSkew:           1,
						TopologyKey:       "kubernetes.io/hostname",
						WhenUnsatisfiable: corev1.ScheduleAnyway,
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app.kubernetes.io/name": "prometheus",
							},
						},
					},
				},
				Version:        config.Version,
				WALCompression: &walCompression,
				Web: &promv1.PrometheusWebSpec{
					PageTitle: &pageTitle,
				},
			},

			EvaluationInterval: promv1.Duration(config.EvaluationInterval),
			Retention:          promv1.Duration(config.RetentionDuration),
			RetentionSize:      promv1.ByteSize(pvcresizing.GetRetentionSize(storageSize)),
		},
	}

	if !key.IsManagementCluster(config.Installation, cluster) {
		// Workload cluster
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

		prometheus.Spec.ServiceMonitorSelector = &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "nonexistentkey",
					Operator: metav1.LabelSelectorOpExists,
				},
			},
		}

		prometheus.Spec.ServiceMonitorNamespaceSelector = &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "nonexistentkey",
					Operator: metav1.LabelSelectorOpExists,
				},
			},
		}
	} else {
		// Management cluster
		prometheus.Spec.APIServerConfig = &promv1.APIServerConfig{
			Host:            fmt.Sprintf("https://%s", key.APIUrl(cluster)),
			BearerTokenFile: key.BearerTokenPath(),
			TLSConfig: &promv1.TLSConfig{
				CAFile: key.CAFilePath(),
				SafeTLSConfig: promv1.SafeTLSConfig{
					InsecureSkipVerify: true,
				},
			},
		}

		prometheus.Spec.Secrets = []string{
			key.EtcdSecret(config.Installation, cluster),
		}

		// We do not discover the service monitors discovered by the agent running on the management cluster
		allMonitorSelector := []metav1.LabelSelectorRequirement{
			{
				Key:      key.TeamLabel,
				Operator: metav1.LabelSelectorOpDoesNotExist,
			},
		}

		// An empty label selector matches all objects.
		prometheus.Spec.ServiceMonitorSelector = &metav1.LabelSelector{
			MatchExpressions: allMonitorSelector,
		}

		// An empty label selector matches all objects.
		prometheus.Spec.ServiceMonitorNamespaceSelector = &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{},
		}

		// An empty label selector matches all objects.
		prometheus.Spec.PodMonitorSelector = &metav1.LabelSelector{
			MatchExpressions: allMonitorSelector,
		}

		// An empty label selector matches all objects.
		prometheus.Spec.PodMonitorNamespaceSelector = &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{},
		}
	}

	if config.PrometheusClient != nil {
		err = currentRemoteWrite(ctx, config, prometheus)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return prometheus, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*promv1.Prometheus)
	d := desired.(*promv1.Prometheus)

	return !cmp.Equal(c.Spec, d.Spec, cmpopts.IgnoreFields(promv1.PrometheusSpec{}, "RemoteWrite"))
}

// Fetch current Prometheus CR and update RemoteWrite field
func currentRemoteWrite(ctx context.Context, config Config, p *promv1.Prometheus) error {
	current, err := config.PrometheusClient.MonitoringV1().Prometheuses(p.GetNamespace()).Get(ctx, p.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return microerror.Mask(err)
	}
	p.Spec.RemoteWrite = current.Spec.RemoteWrite
	return nil
}
