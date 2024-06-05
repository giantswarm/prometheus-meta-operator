package etcdcertificates

import (
	"context"
	"os"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

// getSource retrieves data for the desired Secret.
func (r *Resource) getSource(ctx context.Context) (map[string]string, error) {
	var data map[string]string
	var err error
	if key.IsCAPIManagementCluster(r.config.Provider) {
		// In CAPI clusters, etcd certificates are stored in a secret
		data, err = r.getSourceFromSecret(ctx, key.EtcdSecretSourceName, key.EtcdSecretSourceNamespace)
	} else {
		// In Vintage clusters, etcd certificates are mounted as files on the node filesystem
		data, err = r.getSourceFromDisk()
	}
	if err != nil {
		return nil, microerror.Mask(err)
	}
	return data, nil
}

// getSourceFromSecret retrieves etcd certificates from a kubernetes secret.
func (r *Resource) getSourceFromSecret(ctx context.Context, name, namespace string) (map[string]string, error) {
	secret, err := r.config.K8sClient.K8sClient().CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var ca, crt, key []byte
	var ok bool

	if ca, ok = secret.Data["ca.crt"]; !ok {
		return nil, microerror.Maskf(keyMissingError, "ca.crt key missing in secret %s/%s", namespace, name)
	}
	if crt, ok = secret.Data["tls.crt"]; !ok {
		return nil, microerror.Maskf(keyMissingError, "tls.crt key missing in secret %s/%s", namespace, name)
	}
	if key, ok = secret.Data["tls.key"]; !ok {
		return nil, microerror.Maskf(keyMissingError, "tls.key key missing in secret %s/%s", namespace, name)
	}

	data := map[string]string{
		"ca":  string(ca),
		"crt": string(crt),
		"key": string(key),
	}

	return data, nil
}

// getSourceFromDisk retrieves etcd certificates from the filesystem.
func (r *Resource) getSourceFromDisk() (map[string]string, error) {
	ca, err := os.ReadFile("/etcd-client-certs/ca.pem")
	if err != nil {
		return nil, microerror.Mask(err)
	}

	crt, err := os.ReadFile("/etcd-client-certs/crt.pem")
	if err != nil {
		return nil, microerror.Mask(err)
	}

	key, err := os.ReadFile("/etcd-client-certs/key.pem")
	if err != nil {
		return nil, microerror.Mask(err)
	}

	data := map[string]string{
		"ca":  string(ca),
		"crt": string(crt),
		"key": string(key),
	}

	return data, nil
}
