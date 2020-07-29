package resource

import (
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/k8sclient/v3/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/resource"
	"github.com/giantswarm/operatorkit/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/resource/wrapper/retryresource"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alert"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/certificates"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/frontend"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/ingress"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/namespace"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/scrapeconfigs"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/service"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/servicemonitor"
)

type Config struct {
	BaseDomain                string
	SupportsPersistentStorage bool
	StorageSize               string
	K8sClient                 k8sclient.Interface
	Logger                    micrologger.Logger
	PrometheusClient          promclient.Interface
}

func New(config Config) ([]resource.Interface, error) {
	var err error

	var namespaceResource resource.Interface
	{
		c := namespace.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		namespaceResource, err = namespace.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var certificatesResource resource.Interface
	{
		c := certificates.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		certificatesResource, err = certificates.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var prometheusResource resource.Interface
	{
		c := prometheus.Config{
			PrometheusClient:          config.PrometheusClient,
			Logger:                    config.Logger,
			SupportsPersistentStorage: config.SupportsPersistentStorage,
			StorageSize:               config.StorageSize,
		}

		prometheusResource, err = prometheus.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var frontendResource resource.Interface
	{
		c := frontend.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		frontendResource, err = frontend.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var serviceResource resource.Interface
	{
		c := service.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		serviceResource, err = service.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var ingressResource resource.Interface
	{
		c := ingress.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,

			BaseDomain: config.BaseDomain,
		}

		ingressResource, err = ingress.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var serviceMonitorResource resource.Interface
	{
		c := servicemonitor.Config{
			PrometheusClient: config.PrometheusClient,
			Logger:           config.Logger,
		}

		serviceMonitorResource, err = servicemonitor.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var alertResource resource.Interface
	{
		c := alert.Config{
			PrometheusClient: config.PrometheusClient,
			Logger:           config.Logger,
		}

		alertResource, err = alert.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var scrapeConfigResource resource.Interface
	{
		c := scrapeconfigs.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		scrapeConfigResource, err = scrapeconfigs.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		namespaceResource,
		certificatesResource,
		prometheusResource,
		frontendResource,
		serviceResource,
		ingressResource,
		serviceMonitorResource,
		alertResource,
		scrapeConfigResource,
	}

	{
		c := retryresource.WrapConfig{
			Logger: config.Logger,
		}

		resources, err = retryresource.Wrap(resources, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	{
		c := metricsresource.WrapConfig{}

		resources, err = metricsresource.Wrap(resources, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return resources, nil
}
