package remotewrite

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v8/pkg/resource"

	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/prometheusremotewrite"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/remotewritesecret"
)

func newResources(config ControllerConfig) ([]resource.Interface, error) {
	var err error

	var prometheusRemoteWrite resource.Interface
	{
		c := prometheusremotewrite.Config{
			K8sClient:        config.K8sClient,
			Logger:           config.Logger,
			PrometheusClient: config.PrometheusClient,

			ProxyConfiguration: config.ProxyConfiguration,
		}

		prometheusRemoteWrite, err = prometheusremotewrite.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var rwSecretResource resource.Interface
	{
		c := remotewritesecret.Config{
			K8sClient:        config.K8sClient,
			Logger:           config.Logger,
			PrometheusClient: config.PrometheusClient,
		}

		rwSecretResource, err = remotewritesecret.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		rwSecretResource,
		prometheusRemoteWrite,
	}

	return resources, nil
}
