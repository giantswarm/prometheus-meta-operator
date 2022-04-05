package etcdcertificates

import (
	"context"
	"io/ioutil"

	"github.com/giantswarm/microerror"
)

// getSource retrieve data for the desired Secret.
func (sc *secretCopier) getSource(ctx context.Context, v interface{}) (map[string]string, error) {
	data, err := sc.getSourceFromDisk()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return data, nil
}

// getSourceFromDisk retrieves etcd certificates from the filesystem.
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
