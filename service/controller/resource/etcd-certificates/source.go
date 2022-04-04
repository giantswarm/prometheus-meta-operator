package etcdcertificates

import (
	"context"
	"io/ioutil"

	"github.com/giantswarm/microerror"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

// getSource retrieve data for the desired Secret.
func (sc *secretCopier) getSource(ctx context.Context, v interface{}) (map[string]string, error) {
	data, err := sc.getSourceFromSecret(ctx, key.EtcdSecretSourceName(), key.EtcdSecretSourceNamespace())
	if err != nil {
		sc.logger.Debugf(ctx, "could not get certificates from secret : %v", err)

		data, err = sc.getSourceFromDisk()
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return data, nil
}

func (sc *secretCopier) getSourceFromSecret(ctx context.Context, name, namespace string) (map[string]string, error) {
	secret, err := sc.k8sClient.K8sClient().CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
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

func (sc *secretCopier) getSourceFromDisk() (map[string]string, error) {
	ca, err := ioutil.ReadFile("/etcd-client-certs/ca.pem")
	if err != nil {
		return nil, microerror.Mask(err)
	}

	crt, err := ioutil.ReadFile("/etcd-client-certs/crt.pem")
	if err != nil {
		return nil, microerror.Mask(err)
	}

	key, err := ioutil.ReadFile("/etcd-client-certs/key.pem")
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