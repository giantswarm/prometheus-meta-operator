package verticalpodautoscaler

import (
	"reflect"

	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	autoscaling "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "verticalpodautoscaler"
)

type Config struct {
	K8sClient k8sclient.Interface
	VpaClient vpa_clientset.Interface
	Logger    micrologger.Logger
}

type Resource struct {
	k8sClient k8sclient.Interface
	vpaClient vpa_clientset.Interface
	logger    micrologger.Logger
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.VpaClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.VpaClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &Resource{
		k8sClient: config.K8sClient,
		vpaClient: config.VpaClient,
		logger:    config.Logger,
	}

	return r, nil
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
		Name:      "prometheus",
		Namespace: key.Namespace(cluster),
		Labels:    key.Labels(cluster),
	}, nil
}

func (r *Resource) getObject(v interface{}) (*vpa_types.VerticalPodAutoscaler, error) {
	objectMeta, err := r.getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	updateModeAuto := vpa_types.UpdateModeAuto
	containerScalingModeAuto := vpa_types.ContainerScalingModeAuto
	containerScalingModeOff := vpa_types.ContainerScalingModeOff
	containerControlledValuesRequestsAndLimits := vpa_types.ContainerControlledValuesRequestsAndLimits
	vpa := &vpa_types.VerticalPodAutoscaler{
		ObjectMeta: objectMeta,
		Spec: vpa_types.VerticalPodAutoscalerSpec{
			TargetRef: &autoscaling.CrossVersionObjectReference{
				Kind:       "StatefulSet",
				Name:       key.PrometheusSTSName(cluster),
				APIVersion: "apps/v1",
			},
			UpdatePolicy: &vpa_types.PodUpdatePolicy{
				UpdateMode: &updateModeAuto,
			},
			ResourcePolicy: &vpa_types.PodResourcePolicy{
				ContainerPolicies: []vpa_types.ContainerResourcePolicy{
					{
						ContainerName:    key.PrometheusContainerName(),
						Mode:             &containerScalingModeAuto,
						ControlledValues: &containerControlledValuesRequestsAndLimits,
					},
					{
						ContainerName: "prometheus-config-reloader",
						Mode:          &containerScalingModeOff,
					},
					{
						ContainerName: "rules-configmap-reloader",
						Mode:          &containerScalingModeOff,
					},
				},
			},
		},
	}

	return vpa, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*vpa_types.VerticalPodAutoscaler)
	d := desired.(*vpa_types.VerticalPodAutoscaler)

	return !reflect.DeepEqual(c.Spec, d.Spec)
}
