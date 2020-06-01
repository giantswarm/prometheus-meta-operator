package frontend

import (
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

const (
	Name = "frontend"
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
		k8sClient: config.K8sClient,
		logger:    config.Logger,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func toFrontend(v interface{}) (*v1beta1.Deployment, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var replicas int32 = 2

	deployment := &v1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "frontend",
			Namespace: key.Namespace(cluster),
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "frontend",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "frontend",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "frontend",
							Image: "quay.io/cortexproject/cortex:v0.4.0",
							Args: []string{
								"-auth.enabled=false",
								"-fifocache.duration=24h",
								"-frontend.cache.enable-fifocache=true",
								"-frontend.default-validity=24h",
								"-frontend.downstream-url=http://prometheus-operated:9090",
								"-frontend.fifocache.duration=24h",
								"-frontend.fifocache.size=1024",
								"-frontend.log-queries-longer-than=1s",
								"-http.prefix=",
								"-log.level=debug",
								"-querier.align-querier-with-step=true",
								"-querier.cache-results=true",
								"-querier.compress-http-responses=true",
								"-querier.split-queries-by-interval=24h",
								"-server.http-listen-port=9091",
								"-target=query-frontend",
							},
						},
					},
				},
			},
		},
	}

	return deployment, nil
}
