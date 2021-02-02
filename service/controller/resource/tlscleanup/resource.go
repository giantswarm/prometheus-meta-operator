// tlscleanup provides a resource that cleans up unnecessary TLS certificates
// from the shard namespace. These were created by earlier versions but are now
// unused so we ensure to remove them to not leave them dangling.
//
// TODO: this resource can be removed after a couple of releases
// TODO: when removing, check and remove `key.SecretTLSCertificates` function
package tlscleanup

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

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

	return &Resource{
		k8sClient: config.K8sClient,
		logger:    config.Logger,
	}, nil
}

func (r *Resource) Name() string {
	return "tlscleanup"
}

func (r *Resource) getObjectMeta(v interface{}) (metav1.ObjectMeta, error) {
	cluster, err := key.ToCluster(v)
	if err != nil {
		return metav1.ObjectMeta{}, microerror.Mask(err)
	}

	return metav1.ObjectMeta{
		Name:      key.SecretTLSCertificates(cluster),
		Namespace: key.Namespace(cluster),
	}, nil
}
