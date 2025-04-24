package remotewritesecret

import (
	"github.com/giantswarm/k8sclient/v8/pkg/k8sclient"
	"github.com/giantswarm/micrologger"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/v2/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/remotewriteutils"
)

const (
	Name           = "remotewrite"
	labelName      = "remotewrite/name"
	labelNamespace = "remotewrite/namespace"
)

type Config struct {
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
}

type Resource struct {
	k8sClient        k8sclient.Interface
	logger           micrologger.Logger
	prometheusClient promclient.Interface
}

func New(config Config) (*Resource, error) {
	r := &Resource{
		k8sClient:        config.K8sClient,
		logger:           config.Logger,
		prometheusClient: config.PrometheusClient,
	}

	return r, nil
}

func (r *Resource) Name() string {
	return Name
}

func (r *Resource) ensureRemoteWriteSecret(scSpec pmov1alpha1.RemoteWriteSecretSpec, meta metav1.ObjectMeta, ns string) corev1.Secret {

	labels := meta.GetLabels()
	if labels != nil {
		labels[labelName] = meta.GetName()
		labels[labelNamespace] = meta.GetNamespace()
	} else {
		labels = map[string]string{
			labelName:      meta.GetName(),
			labelNamespace: meta.GetNamespace(),
		}
	}

	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        scSpec.Name,
			Namespace:   ns,
			Labels:      labels,
			Annotations: meta.GetAnnotations(),
		},
		Data: scSpec.Data,
	}
}

func toResourceWrapper(r *Resource) *remotewriteutils.ResourceWrapper {
	return &remotewriteutils.ResourceWrapper{
		K8sClient:        r.k8sClient,
		Logger:           r.logger,
		PrometheusClient: r.prometheusClient,
	}
}

func secretInstalled(sc corev1.Secret, list []corev1.Secret) bool {
	for _, item := range list {
		if item.GetNamespace() == sc.GetNamespace() &&
			item.GetName() == sc.GetName() {
			return true
		}
	}
	return false
}
