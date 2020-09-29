package resource

import (
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v2/pkg/resource"
	"github.com/giantswarm/operatorkit/v2/pkg/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/v2/pkg/resource/wrapper/retryresource"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alert"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/certificates"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/ingress"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/namespace"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/prometheus"
	promxyApp "github.com/giantswarm/prometheus-meta-operator/service/controller/resource/promxy/app"
	promxyConfigmap "github.com/giantswarm/prometheus-meta-operator/service/controller/resource/promxy/configmap"
	promxyServerGroup "github.com/giantswarm/prometheus-meta-operator/service/controller/resource/promxy/servergroup"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/scrapeconfigs"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/servicemonitor"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

type Config struct {
	Address                 string
	BaseDomain              string
	Provider                string
	Installation            string
	CreatePVC               bool
	StorageSize             string
	RestrictedAccessEnabled bool
	WhitelistedSubnets      string
	K8sClient               k8sclient.Interface
	Logger                  micrologger.Logger
	PrometheusClient        promclient.Interface
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

	var apiCertificatesResource resource.Interface
	{
		c := certificates.Config{
			Name:                "api-certificates",
			K8sClient:           config.K8sClient,
			Logger:              config.Logger,
			SourceNameFunc:      key.Namespace,
			SourceNamespaceFunc: key.NamespaceDefault,
			TargetNameFunc:      key.SecretAPICertificates,
		}

		apiCertificatesResource, err = certificates.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var tlsCertificatesResource resource.Interface
	{
		c := certificates.Config{
			Name:                "tls-certificates",
			K8sClient:           config.K8sClient,
			Logger:              config.Logger,
			SourceNameFunc:      key.SecretTLSCertificates,
			SourceNamespaceFunc: key.NamespaceMonitoring,
			TargetNameFunc:      key.SecretTLSCertificates,
		}

		tlsCertificatesResource, err = certificates.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var prometheusResource resource.Interface
	{
		c := prometheus.Config{
			Address:          config.Address,
			PrometheusClient: config.PrometheusClient,
			Logger:           config.Logger,
			CreatePVC:        config.CreatePVC,
			StorageSize:      config.StorageSize,
		}

		prometheusResource, err = prometheus.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var serviceMonitorResource resource.Interface
	{
		c := servicemonitor.Config{
			PrometheusClient: config.PrometheusClient,
			Logger:           config.Logger,
			Provider:         config.Provider,
			Installation:     config.Installation,
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
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
			Provider:     config.Provider,
			Installation: config.Installation,
		}

		scrapeConfigResource, err = scrapeconfigs.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var ingressResource resource.Interface
	{
		c := ingress.Config{
			K8sClient:               config.K8sClient,
			Logger:                  config.Logger,
			BaseDomain:              config.BaseDomain,
			RestrictedAccessEnabled: config.RestrictedAccessEnabled,
			WhitelistedSubnets:      config.WhitelistedSubnets,
		}

		ingressResource, err = ingress.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var promxyConfigmapResource resource.Interface
	{
		c := promxyConfigmap.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		promxyConfigmapResource, err = promxyConfigmap.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	var promxyAppResource resource.Interface
	{
		c := promxyApp.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		promxyAppResource, err = promxyApp.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var promxyServerGroupResource resource.Interface
	{
		c := promxyServerGroup.Config{
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
			Installation: config.Installation,
		}

		promxyServerGroupResource, err = promxyServerGroup.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		namespaceResource,
		apiCertificatesResource,
		tlsCertificatesResource,
		prometheusResource,
		serviceMonitorResource,
		alertResource,
		scrapeConfigResource,
		ingressResource,
		promxyConfigmapResource,
		promxyAppResource,
		promxyServerGroupResource,
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
