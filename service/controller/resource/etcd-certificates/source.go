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
		sc.logger.Errorf(ctx, err, "could not get certificates from secret")

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

	return secret.Data, nil
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
