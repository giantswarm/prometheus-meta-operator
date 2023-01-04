package resource

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/wrapper/retryresource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/domain"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/password"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/alerting/alertmanagerwiring"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/alerting/heartbeat"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/alerting/heartbeatwebhookconfig"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/certificates"
	ingress "github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/ingress"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/pvcresizingresource"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteapiendpointconfigsecret"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteingress"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteingressauth"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/scrapeconfigs"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/verticalpodautoscaler"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/namespace"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/rbac"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/wrapper/monitoringdisabledresource"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

type Config struct {
	K8sClient          k8sclient.Interface
	Logger             micrologger.Logger
	PrometheusClient   promclient.Interface
	VpaClient          vpa_clientset.Interface
	ProxyConfiguration domain.ProxyConfiguration

	AdditionalScrapeConfigs string
	Bastions                []string
	Customer                string
	Installation            string
	InsecureCA              bool
	Pipeline                string
	Provider                string
	Region                  string
	Registry                string

	OpsgenieKey string

	PrometheusAddress           string
	PrometheusBaseDomain        string
	PrometheusCreatePVC         bool
	PrometheusLogLevel          string
	PrometheusRemoteWriteURL    string
	PrometheusRetentionDuration string
	PrometheusRetentionSize     string
	PrometheusVersion           string

	RestrictedAccessEnabled bool
	WhitelistedSubnets      string
}

func New(config Config) ([]resource.Interface, error) {
	var err error

	passwordManager := password.SimpleManager{}

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
			Name:      "api-certificates",
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Sources: []certificates.CertificateSource{
				{
					NameFunc:      key.Namespace,
					NamespaceFunc: key.NamespaceDefault,
				},
				{
					NameFunc:      key.Namespace,
					NamespaceFunc: key.OrganizationNamespace,
				},
				{
					NameFunc:      key.CAPICertificateName,
					NamespaceFunc: key.OrganizationNamespace,
				},
			},
			Target: key.Secret,
		}

		apiCertificatesResource, err = certificates.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var heartbeatWebhookConfigResource resource.Interface
	{
		c := heartbeatwebhookconfig.Config{
			Client: config.PrometheusClient,
			Logger: config.Logger,

			Installation:       config.Installation,
			ProxyConfiguration: config.ProxyConfiguration,
		}

		heartbeatWebhookConfigResource, err = heartbeatwebhookconfig.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var remoteWriteIngressAuthResource resource.Interface
	{
		c := remotewriteingressauth.Config{
			K8sClient:       config.K8sClient,
			Logger:          config.Logger,
			PasswordManager: passwordManager,
			Installation:    config.Installation,
			Provider:        config.Provider,
		}

		remoteWriteIngressAuthResource, err = remotewriteingressauth.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var remoteWriteIngressResource resource.Interface
	{
		c := remotewriteingress.Config{
			K8sClient:  config.K8sClient,
			Logger:     config.Logger,
			BaseDomain: config.PrometheusBaseDomain,
		}

		remoteWriteIngressResource, err = remotewriteingress.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var remoteWriteAPIEndpointConfigSecretResource resource.Interface
	{
		c := remotewriteapiendpointconfigsecret.Config{
			K8sClient:          config.K8sClient,
			Logger:             config.Logger,
			PasswordManager:    passwordManager,
			ProxyConfiguration: config.ProxyConfiguration,
			BaseDomain:         config.PrometheusBaseDomain,
			Customer:           config.Customer,
			Installation:       config.Installation,
			InsecureCA:         config.InsecureCA,
			Pipeline:           config.Pipeline,
			Provider:           config.Provider,
			Region:             config.Region,
		}

		remoteWriteAPIEndpointConfigSecretResource, err = remotewriteapiendpointconfigsecret.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var rbacResource resource.Interface
	{
		c := rbac.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		rbacResource, err = rbac.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var prometheusResource resource.Interface
	{
		c := prometheus.Config{
			Address:           config.PrometheusAddress,
			PrometheusClient:  config.PrometheusClient,
			Logger:            config.Logger,
			CreatePVC:         config.PrometheusCreatePVC,
			Customer:          config.Customer,
			Installation:      config.Installation,
			Pipeline:          config.Pipeline,
			Version:           config.PrometheusVersion,
			Provider:          config.Provider,
			Region:            config.Region,
			Registry:          config.Registry,
			LogLevel:          config.PrometheusLogLevel,
			RetentionDuration: config.PrometheusRetentionDuration,
			RetentionSize:     config.PrometheusRetentionSize,
		}

		prometheusResource, err = prometheus.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var verticalPodAutoScalerResource resource.Interface
	{
		c := verticalpodautoscaler.Config{
			Logger:    config.Logger,
			K8sClient: config.K8sClient,
			VpaClient: config.VpaClient,
		}

		verticalPodAutoScalerResource, err = verticalpodautoscaler.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var scrapeConfigResource resource.Interface
	{
		c := scrapeconfigs.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,

			AdditionalScrapeConfigs: config.AdditionalScrapeConfigs,
			Bastions:                config.Bastions,
			Customer:                config.Customer,
			Provider:                config.Provider,
			Installation:            config.Installation,
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
			BaseDomain:              config.PrometheusBaseDomain,
			RestrictedAccessEnabled: config.RestrictedAccessEnabled,
			WhitelistedSubnets:      config.WhitelistedSubnets,
		}

		ingressResource, err = ingress.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var heartbeatResource resource.Interface
	{
		c := heartbeat.Config{
			Logger:       config.Logger,
			Installation: config.Installation,
			OpsgenieKey:  config.OpsgenieKey,
			Pipeline:     config.Pipeline,
		}

		heartbeatResource, err = heartbeat.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var alertmanagerWiringResource resource.Interface
	{
		c := alertmanagerwiring.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		alertmanagerWiringResource, err = alertmanagerwiring.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var pvcResizeResource resource.Interface
	{
		c := pvcresizingresource.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		pvcResizeResource, err = pvcresizingresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		namespaceResource,
		apiCertificatesResource,
		rbacResource,
		heartbeatWebhookConfigResource,
		scrapeConfigResource,
		remoteWriteAPIEndpointConfigSecretResource,
		remoteWriteIngressAuthResource,
		remoteWriteIngressResource,
		alertmanagerWiringResource,
		prometheusResource,
		verticalPodAutoScalerResource,
		ingressResource,
		heartbeatResource,
		pvcResizeResource,
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

	{
		c := monitoringdisabledresource.WrapConfig{
			Logger: config.Logger,
		}

		resources, err = monitoringdisabledresource.Wrap(resources, c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return resources, nil
}
