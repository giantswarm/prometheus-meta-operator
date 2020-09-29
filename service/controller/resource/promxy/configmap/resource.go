package configmap

import (
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/promxy"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

type Resource struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger
}

func New(config Config) (*Resource, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	r := &Resource{
		logger:    config.Logger,
		k8sClient: config.K8sClient,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return "promxy-configmap"
}

func (r *Resource) getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.PromxyConfigMapName(),
		Namespace: key.PromxyConfigMapNamespace(),
		Labels:    key.Labels(cluster),
	}, nil
}

func (r *Resource) toConfigMap(objectMeta metav1.ObjectMeta) (*corev1.ConfigMap, error) {
	data := make(map[string]string)

	content, err := promxy.Serialize(promxy.Promxy{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	data["values"] = content
	return &corev1.ConfigMap{
		ObjectMeta: objectMeta,
		Data:       data,
	}, nil
}
