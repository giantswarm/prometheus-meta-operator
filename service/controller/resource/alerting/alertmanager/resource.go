package alertmanager

import (
	"net/url"
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringv1client "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

var replicas int32 = 1

const (
	Name = "alertmanager"
)

type Config struct {
	Client monitoringv1client.Interface
	Logger micrologger.Logger

	Address     string
	CreatePVC   bool
	LogLevel    string
	StorageSize string
	Version     string
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.Client.MonitoringV1().Alertmanagers(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:    clientFunc,
		Logger:        config.Logger,
		Name:          Name,
		GetObjectMeta: getObjectMeta,
		GetDesiredObject: func(v interface{}) (metav1.Object, error) {
			return toAlertmanager(v, config)
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
		Name:      "alertmanager",
		Namespace: key.NamespaceMonitoring(),
		Labels:    key.AlertmanagerLabels(cluster),
	}, nil
}

func toAlertmanager(v interface{}, config Config) (metav1.Object, error) {
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

	labels := make(map[string]string)
	for k, v := range key.AlertmanagerLabels(cluster) {
		labels[k] = v
	}

	var storage monitoringv1.StorageSpec
	storageSize := resource.MustParse(config.StorageSize)
	if config.CreatePVC {
		storage = monitoringv1.StorageSpec{
			VolumeClaimTemplate: monitoringv1.EmbeddedPersistentVolumeClaim{
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
	} else {
		storage = monitoringv1.StorageSpec{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				SizeLimit: &storageSize,
			},
		}
	}

	address, err := url.Parse(config.Address)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Alertmanager default image runs using the nobody user (65534)
	var uid int64 = 65534
	var fsGroup int64 = 65534
	var runAsNonRoot bool = true
	var gid int64 = 65534
	alertmanager := &monitoringv1.Alertmanager{
		ObjectMeta: objectMeta,
		Spec: monitoringv1.AlertmanagerSpec{
			PodMetadata: &monitoringv1.EmbeddedObjectMetadata{
				Labels: labels,
			},
			Version:     config.Version,
			LogLevel:    config.LogLevel,
			ExternalURL: address.String(),
			Replicas:    &replicas,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					// cpu: 100m
					corev1.ResourceCPU: *key.AlertmanagerDefaultCPU(),
					// memory: 200m
					corev1.ResourceMemory: *key.AlertmanagerDefaultMemory(),
				},
				Limits: corev1.ResourceList{
					// cpu: 100m
					corev1.ResourceCPU: *key.AlertmanagerDefaultCPU(),
					// memory: 200m
					corev1.ResourceMemory: *key.AlertmanagerDefaultMemory(),
				},
			},
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    &uid,
				RunAsGroup:   &gid,
				RunAsNonRoot: &runAsNonRoot,
				FSGroup:      &fsGroup,
			},
			Storage: &storage,
			Affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "role",
										Operator: corev1.NodeSelectorOpNotIn,
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
			TopologySpreadConstraints: []corev1.TopologySpreadConstraint{
				{
					MaxSkew:           1,
					TopologyKey:       "kubernetes.io/hostname",
					WhenUnsatisfiable: corev1.ScheduleAnyway,
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name": "alertmanager",
						},
					},
				},
			},
			PriorityClassName: "prometheus",
		},
	}

	return alertmanager, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*monitoringv1.Alertmanager)
	d := desired.(*monitoringv1.Alertmanager)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
