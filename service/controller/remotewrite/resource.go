package remotewrite

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/prometheusremotewrite"
)

func newResources(config ControllerConfig) ([]resource.Interface, error) {
	var err error

	var prometheusRemoteWrite resource.Interface
	{
		c := prometheusremotewrite.Config{
			K8sClient:        config.K8sClient,
			Logger:           config.Logger,
			PrometheusClient: config.PrometheusClient,
		}

		prometheusRemoteWrite, err = prometheusremotewrite.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		prometheusRemoteWrite,
	}

	return resources, nil
}
