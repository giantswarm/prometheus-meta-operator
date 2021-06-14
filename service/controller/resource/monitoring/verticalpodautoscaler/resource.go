package verticalpodautoscaler

import (
	"context"

	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	autoscaling "k8s.io/api/autoscaling/v1"
	v1 "k8s.io/api/core/v1"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
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
		Labels:    key.PrometheusLabels(cluster),
	}, nil
}

func (r *Resource) getObject(ctx context.Context, v interface{}) (*vpa_types.VerticalPodAutoscaler, error) {
	objectMeta, err := r.getObjectMeta(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	minCpu := key.PrometheusDefaultCPU()
	minMemory := key.PrometheusDefaultMemory()

	maxCpu, err := r.getMaxCPU(ctx)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	maxMemory, err := r.getMaxMemory(ctx)
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
						MinAllowed: v1.ResourceList{
							v1.ResourceCPU:    *minCpu,
							v1.ResourceMemory: *minMemory,
						},
						MaxAllowed: v1.ResourceList{
							v1.ResourceCPU:    *maxCpu,
							v1.ResourceMemory: *maxMemory,
						},
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

func (r *Resource) getMaxCPU(ctx context.Context) (*resource.Quantity, error) {
	nodes, err := r.k8sClient.K8sClient().CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var nodeCpu *resource.Quantity
	if len(nodes.Items) > 0 {
		n := nodes.Items[0]
		s, ok := n.Status.Allocatable[v1.ResourceCPU]
		if ok {
			nodeCpu = &s
		}

		for _, n := range nodes.Items[1:] {
			s, ok := n.Status.Allocatable[v1.ResourceCPU]
			if ok && nodeCpu.Cmp(s) == -1 {
				nodeCpu = &s
			}
		}
	}
	if nodeCpu == nil || nodeCpu.IsZero() {
		return nil, microerror.Mask(nodeCpuNotFoundError)
	}

	q, err := quantityMultiply(nodeCpu, 0.9)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return q, nil
}

func (r *Resource) getMaxMemory(ctx context.Context) (*resource.Quantity, error) {
	nodes, err := r.k8sClient.K8sClient().CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var nodeMemory *resource.Quantity
	if len(nodes.Items) > 0 {
		n := nodes.Items[0]
		s, ok := n.Status.Allocatable[v1.ResourceMemory]
		if ok {
			nodeMemory = &s
		}

		for _, n := range nodes.Items[1:] {
			s, ok := n.Status.Allocatable[v1.ResourceMemory]
			if ok && nodeMemory.Cmp(s) == -1 {
				nodeMemory = &s
			}
		}
	}
	if nodeMemory == nil || nodeMemory.IsZero() {
		return nil, microerror.Mask(nodeMemoryNotFoundError)
	}

	q, err := quantityMultiply(nodeMemory, 0.9)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	q, err = quantityMultiply(q, 1/key.PrometheusMemoryLimitCoefficient)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return q, nil
}

func quantityMultiply(q *resource.Quantity, multiplier float64) (*resource.Quantity, error) {
	i, ok := q.AsInt64()
	if !ok {
		return nil, microerror.Maskf(quantityConversionError, q.String())
	}

	n := float64(i) * multiplier

	q.Set(int64(n))

	return q, nil
}

func hasChanged(current, desired metav1.Object) bool {
	c := current.(*vpa_types.VerticalPodAutoscaler)
	d := desired.(*vpa_types.VerticalPodAutoscaler)

	return !apiequality.Semantic.DeepEqual(c.Spec, d.Spec)
}
