package prometheus

import (
	"fmt"
	"net/url"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
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

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/pvcresizing"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

const (
	Name = "prometheus"
)

type Config struct {
	PrometheusClient promclient.Interface
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger

	Address            string
	Bastions           []string
	Customer           string
	EvaluationInterval string
	ImageRepository    string
	Installation       string
	Pipeline           string
	Provider           cluster.Provider
	Region             string
	Registry           string
	LogLevel           string
	ScrapeInterval     string
	Version            string

	MimirEnabled bool
}

type Resource struct {
	prometheusClient promclient.Interface
	k8sClient        k8sclient.Interface
	logger           micrologger.Logger

	address            string
	bastions           []string
	customer           string
	evaluationInterval string
	imageRepository    string
	installation       string
	pipeline           string
	provider           cluster.Provider
	region             string
	registry           string
	logLevel           string
	scrapeInterval     string
	version            string

	mimirEnabled bool
}

func New(config Config) (*Resource, error) {
	return &Resource{
		prometheusClient: config.PrometheusClient,
		k8sClient:        config.K8sClient,
		logger:           config.Logger,

		address:            config.Address,
		bastions:           config.Bastions,
		customer:           config.Customer,
		evaluationInterval: config.EvaluationInterval,
		imageRepository:    config.ImageRepository,
		installation:       config.Installation,
		pipeline:           config.Pipeline,
		provider:           config.Provider,
		region:             config.Region,
		registry:           config.Registry,
		logLevel:           config.LogLevel,
		scrapeInterval:     config.ScrapeInterval,
		version:            config.Version,

		mimirEnabled: config.MimirEnabled,
	}, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
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

func (r *Resource) toPrometheus(ctx context.Context, v interface{}) (*promv1.Prometheus, error) {
	if v == nil {
		return nil, nil
	}

	objectMeta, err := r.getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	address, err := url.Parse(r.address)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var replicas int32 = 1
	var keepDroppedTargets uint64 = 5
	// Configured following https://github.com/prometheus-operator/prometheus-operator/issues/541#issuecomment-451884171
	// as the volume could not mount otherwise
	var uid int64 = 1000
	var fsGroup int64 = 2000
	var runAsNonRoot bool = true
	// Prometheus default image runs using the nobody user (65534)
	var gid int64 = 65534
	var walCompression bool = true
	var prometheusShards int32 = 1

	annotationValue := cluster.GetAnnotations()[key.PrometheusVolumeSizeAnnotation]
	volumeSize := pvcresizing.PrometheusVolumeSize(annotationValue)
	storageSize := resource.MustParse(volumeSize)

	storage := promv1.StorageSpec{
		VolumeClaimTemplate: promv1.EmbeddedPersistentVolumeClaim{
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
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

	image := fmt.Sprintf("%s/%s:%s", r.registry, r.imageRepository, r.version)
	pageTitle := fmt.Sprintf("%s/%s Prometheus", r.installation, key.ClusterID(cluster))
	provider, err := key.ClusterProvider(cluster, r.provider)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	prometheus := &promv1.Prometheus{
		ObjectMeta: objectMeta,
		Spec: promv1.PrometheusSpec{
			CommonPrometheusFields: promv1.CommonPrometheusFields{
				AdditionalScrapeConfigs: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: key.PrometheusAdditionalScrapeConfigsSecretName(),
					},
					Key: key.PrometheusAdditionalScrapeConfigsName(),
				},
				EnableRemoteWriteReceiver: true,
				ExternalLabels: map[string]string{
					key.ClusterIDKey:    key.ClusterID(cluster),
					key.ClusterTypeKey:  key.ClusterType(r.installation, cluster),
					key.CustomerKey:     r.customer,
					key.InstallationKey: r.installation,
					key.PipelineKey:     r.pipeline,
					key.ProviderKey:     provider,
					key.RegionKey:       r.region,
				},
				ExternalURL:        externalURL.String(),
				Image:              &image,
				KeepDroppedTargets: &keepDroppedTargets,
				LogLevel:           r.logLevel,
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
				ScrapeInterval: promv1.Duration(r.scrapeInterval),
				SecurityContext: &corev1.PodSecurityContext{
					RunAsUser:    &uid,
					RunAsGroup:   &gid,
					RunAsNonRoot: &runAsNonRoot,
					FSGroup:      &fsGroup,
					SeccompProfile: &corev1.SeccompProfile{
						Type: corev1.SeccompProfileTypeRuntimeDefault,
					},
				},
				Shards:  &prometheusShards,
				Storage: &storage,
				TopologySpreadConstraints: []promv1.TopologySpreadConstraint{
					{
						CoreV1TopologySpreadConstraint: promv1.CoreV1TopologySpreadConstraint{
							MaxSkew:           1,
							TopologyKey:       "kubernetes.io/hostname",
							WhenUnsatisfiable: corev1.ScheduleAnyway,
							// We want to spread the pods across the nodes as much as possible
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app.kubernetes.io/name": "prometheus",
								},
							},
						},
					},
				},
				Version:        r.version,
				WALCompression: &walCompression,
				Web: &promv1.PrometheusWebSpec{
					PageTitle: &pageTitle,
				},
			},

			EvaluationInterval: promv1.Duration(r.evaluationInterval),
			RetentionSize:      promv1.ByteSize(pvcresizing.GetRetentionSize(storageSize)),
			// Fetches Prometheus rules from any namespace on the Management Cluster
			// using https://v1-22.docs.kubernetes.io/docs/reference/labels-annotations-taints/#kubernetes-io-metadata-name
			RuleNamespaceSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "kubernetes.io/metadata.name",
						Operator: metav1.LabelSelectorOpExists,
					},
				},
			},
		},
	}

	if !key.IsManagementCluster(r.installation, cluster) {
		// Workload cluster
		prometheus.Spec.APIServerConfig = &promv1.APIServerConfig{
			Host: fmt.Sprintf("https://%s", key.APIUrl(cluster)),
			TLSConfig: &promv1.TLSConfig{
				CAFile: fmt.Sprintf("/etc/prometheus/secrets/%s/ca", key.APIServerCertificatesSecretName),
			},
		}

		authenticationType, err := key.ApiServerAuthenticationType(ctx, r.k8sClient, key.Namespace(cluster))
		if err != nil {
			return nil, microerror.Mask(err)
		}
		if authenticationType == "token" {
			prometheus.Spec.APIServerConfig.Authorization = &promv1.Authorization{
				CredentialsFile: fmt.Sprintf("/etc/prometheus/secrets/%s/token", key.APIServerCertificatesSecretName),
			}
		} else if authenticationType == "certificates" {
			prometheus.Spec.APIServerConfig.TLSConfig.CertFile = fmt.Sprintf("/etc/prometheus/secrets/%s/crt", key.APIServerCertificatesSecretName)
			prometheus.Spec.APIServerConfig.TLSConfig.KeyFile = fmt.Sprintf("/etc/prometheus/secrets/%s/key", key.APIServerCertificatesSecretName)
		}

		prometheus.Spec.Secrets = []string{
			key.APIServerCertificatesSecretName,
		}

		prometheus.Spec.RuleSelector = &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      key.ClusterTypeKey,
					Operator: metav1.LabelSelectorOpNotIn,
					Values:   []string{"management_cluster"},
				},
				{
					Key:      key.TeamLabel,
					Operator: metav1.LabelSelectorOpExists,
				},
			},
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
		insecureSkipVerify := true
		prometheus.Spec.APIServerConfig = &promv1.APIServerConfig{
			Host: fmt.Sprintf("https://%s", key.APIUrl(cluster)),
			Authorization: &promv1.Authorization{
				CredentialsFile: key.BearerTokenPath,
			},
			TLSConfig: &promv1.TLSConfig{
				CAFile: key.CAFilePath,
				SafeTLSConfig: promv1.SafeTLSConfig{
					InsecureSkipVerify: &insecureSkipVerify,
				},
			},
		}

		prometheus.Spec.Secrets = []string{
			key.EtcdSecret(r.installation, cluster),
		}

		prometheus.Spec.RuleSelector = &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      key.ClusterTypeKey,
					Operator: metav1.LabelSelectorOpNotIn,
					Values:   []string{"workload_cluster"},
				},
				{
					Key:      key.TeamLabel,
					Operator: metav1.LabelSelectorOpExists,
				},
			},
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

	if r.mimirEnabled {
		emptyExternalLabels := ""
		// Remove prometheus and prometheus_replica external labels to avoid conflicts with our existing rules.
		prometheus.Spec.PrometheusExternalLabelName = &emptyExternalLabels
		prometheus.Spec.ReplicaExternalLabelName = &emptyExternalLabels
		prometheus.Spec.RuleNamespaceSelector = nil
		prometheus.Spec.RuleSelector = nil
	} else {
		// We need to use this to connect each WC prometheus with the central alertmanager instead of the alerting section of the Prometheus CR
		// because the alerting section tries to find the alertmanager service in the workload cluster and not in the management cluster
		// as it is using the secrets defined under prometheus.Spec.APIServerConfig.
		//
		// This forces us to use the static config defined in resource/alerting/alertmanagerwiring.

		// We enable alertmanager on Prometheus only if mimir is not enabled
		prometheus.Spec.AdditionalAlertManagerConfigs = &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: key.AlertmanagerSecretName(),
			},
			Key: key.AlertmanagerKey(),
		}
	}

	if r.prometheusClient != nil {
		err = r.currentRemoteWrite(ctx, prometheus)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return prometheus, nil
}

func (r *Resource) hasChanged(current, desired metav1.Object) bool {
	c := current.(*promv1.Prometheus)
	d := desired.(*promv1.Prometheus)

	return !cmp.Equal(c.Spec, d.Spec, cmpopts.IgnoreFields(promv1.PrometheusSpec{}, "RemoteWrite"))
}

// Fetch current Prometheus CR and update RemoteWrite field
func (r *Resource) currentRemoteWrite(ctx context.Context, p *promv1.Prometheus) error {
	current, err := r.prometheusClient.MonitoringV1().Prometheuses(p.GetNamespace()).Get(ctx, p.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return microerror.Mask(err)
	}
	p.Spec.RemoteWrite = current.Spec.RemoteWrite
	return nil
}
