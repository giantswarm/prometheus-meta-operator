package prometheusautoscaler

import (
	"reflect"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	autoscaling "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/generic"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "prometheusautoscaler"
)

type Config struct {
	VpaClient vpa_clientset.Interface
	Logger    micrologger.Logger
}

func New(config Config) (*generic.Resource, error) {
	clientFunc := func(namespace string) generic.Interface {
		c := config.VpaClient.AutoscalingV1().VerticalPodAutoscalers(namespace)
		return wrappedClient{client: c}
	}

	c := generic.Config{
		ClientFunc:       clientFunc,
		Logger:           config.Logger,
		Name:             Name,
		GetObjectMeta:    getObjectMeta,
		GetDesiredObject: getObject,
		HasChangedFunc:   hasChanged,
	}
	r, err := generic.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return r, nil
}

func getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	return metav1.ObjectMeta{
		Name: "prometheus",
	}, nil
}

func getObject(v interface{}) (metav1.Object, error) {
	objectMeta, err := getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// updateModeAuto := vpa_types.UpdateModeAuto
	// containerScalingModeAuto := vpa_types.ContainerScalingModeAuto
	// containerControlledValuesRequestsAndLimits := vpa_types.ContainerControlledValuesRequestsAndLimits
	vpa := &vpa_types.VerticalPodAutoscaler{
		ObjectMeta: objectMeta,
		Spec: vpa_types.VerticalPodAutoscalerSpec{
			TargetRef: &autoscaling.CrossVersionObjectReference{
				Kind:       "StatefulSet",
				Name:       key.PrometheusSTSName(cluster),
				APIVersion: "apps/v1",
			},
			// UpdatePolicy: &vpa_types.PodUpdatePolicy{
			// 	UpdateMode: &updateModeAuto,
			// },
			ResourcePolicy: &vpa_types.PodResourcePolicy{
				ContainerPolicies: []vpa_types.ContainerResourcePolicy{
					{
						ContainerName: key.PrometheusContainerName(),
						// Mode:             &containerScalingModeAuto,
						// ControlledValues: &containerControlledValuesRequestsAndLimits,
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
