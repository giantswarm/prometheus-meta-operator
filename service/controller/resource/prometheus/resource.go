package prometheus

import (
	"fmt"
	"reflect"

	promv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
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

	CreatePVC   bool
	StorageSize string
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.PrometheusClient.MonitoringV1().Prometheuses(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc: clientFunc,
		Logger:     config.Logger,
		Name:       Name,
		ToCR: func(v interface{}) (metav1.Object, error) {
			return toPrometheus(v, config.CreatePVC, resource.MustParse(config.StorageSize))
		},
		HasChangedFunc: hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func toPrometheus(v interface{}, createPVC bool, storageSize resource.Quantity) (metav1.Object, error) {
	if v == nil {
		return nil, nil
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	name := cluster.GetName()
	var replicas int32 = 1
	var uid int64 = 1000
	var gid int64 = 65534
	var fsGroup int64 = 2000
	var runAsNonRoot bool = true

	var storage promv1.StorageSpec
	if createPVC {
		storage = promv1.StorageSpec{
			VolumeClaimTemplate: v1.PersistentVolumeClaim{
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
			ExternalLabels: map[string]string{
				key.ClusterIDKey(): key.ClusterID(cluster),
				"cluster_type":     "tenant_cluster",
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
			SecurityContext: &v1.PodSecurityContext{
				RunAsUser:    &uid,
				RunAsGroup:   &gid,
				RunAsNonRoot: &runAsNonRoot,
				FSGroup:      &fsGroup,
			},
			Storage: &storage,
		},
	}

	return prometheus, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*promv1.Prometheus)
	d := desired.(*promv1.Prometheus)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
