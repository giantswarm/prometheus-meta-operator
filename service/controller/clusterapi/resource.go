package clusterapi

import (
	"net/url"

	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/wrapper/retryresource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	"k8s.io/client-go/dynamic"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/cluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/organization"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/password"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/alerting/alertmanagerwiring"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/alerting/heartbeat"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/alerting/heartbeatwebhookconfig"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/certificates"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/ciliumnetpol"
	ingress "github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/ingress"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/pvcresizingresource"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteapiendpointconfigsecret"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteconfig"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteingress"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteingressauth"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewritesecret"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/scrapeconfigs"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/verticalpodautoscaler"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/namespace"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/noop"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/rbac"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/wrapper/monitoringdisabledresource"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

type Config struct {
	K8sClient        k8sclient.Interface
	DynamicK8sClient dynamic.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
	VpaClient        vpa_clientset.Interface
	Proxy            func(reqURL *url.URL) (*url.URL, error)

	AdditionalScrapeConfigs string
	Bastions                []string
	Customer                string
	Installation            string
	InsecureCA              bool
	Pipeline                string
	Provider                cluster.Provider
	Region                  string
	Registry                string

	OpsgenieKey string

	MimirEnabled bool

	PrometheusAddress            string
	PrometheusBaseDomain         string
	PrometheusEvaluationInterval string
	PrometheusLogLevel           string
	PrometheusScrapeInterval     string
	PrometheusImageRepository    string
	PrometheusVersion            string

	RestrictedAccessEnabled bool
	WhitelistedSubnets      string

	ExternalDNS bool
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
			Provider:  config.Provider,
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Sources: []certificates.CertificateSource{
				{
					NameFunc:      key.Namespace,
					NamespaceFunc: key.NamespaceDefault,
				},
				{
					NameFunc:      key.Namespace,
					NamespaceFunc: key.ClusterNamespace,
				},
				{
					NameFunc:      key.CAPICertificateName,
					NamespaceFunc: key.ClusterNamespace,
				},
			},
			Target: key.APIServerCertificatesSecretName,
		}

		apiCertificatesResource, err = certificates.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var ciliumnetpolResource resource.Interface
	{
		c := ciliumnetpol.Config{
			DynamicK8sClient: config.DynamicK8sClient,
			Logger:           config.Logger,
		}

		ciliumnetpolResource, err = ciliumnetpol.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var heartbeatWebhookConfigResource resource.Interface
	{
		c := heartbeatwebhookconfig.Config{
			Client: config.PrometheusClient,
			Logger: config.Logger,

			Installation: config.Installation,
			Proxy:        config.Proxy,
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
			K8sClient:   config.K8sClient,
			Logger:      config.Logger,
			BaseDomain:  config.PrometheusBaseDomain,
			ExternalDNS: config.ExternalDNS,
		}

		remoteWriteIngressResource, err = remotewriteingress.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	organizationReader := organization.NewNamespaceReader(config.K8sClient.K8sClient(), config.Installation, config.Provider)
	// This resource creates a the prometheus agent remote write configuration.
	// This is now managed by the observability-operator when mimir is enabled.
	var remoteWriteConfigResource resource.Interface
	if config.MimirEnabled {
		remoteWriteConfigResource = noop.New(noop.Config{Logger: config.Logger})
	} else {
		c := remotewriteconfig.Config{
			K8sClient:          config.K8sClient,
			Logger:             config.Logger,
			OrganizationReader: organizationReader,

			Customer:     config.Customer,
			Installation: config.Installation,
			Pipeline:     config.Pipeline,
			Provider:     config.Provider,
			Region:       config.Region,
			Version:      config.PrometheusVersion,
		}

		remoteWriteConfigResource, err = remotewriteconfig.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// This resource creates a the prometheus agent remote write secret.
	// This is now managed by the observability-operator when mimir is enabled.
	var remoteWriteSecretResource resource.Interface
	if config.MimirEnabled {
		remoteWriteSecretResource = noop.New(noop.Config{Logger: config.Logger})
	} else {
		c := remotewritesecret.Config{
			K8sClient:       config.K8sClient,
			Logger:          config.Logger,
			PasswordManager: passwordManager,
			BaseDomain:      config.PrometheusBaseDomain,
			Installation:    config.Installation,
			InsecureCA:      config.InsecureCA,
			Provider:        config.Provider,
		}

		remoteWriteSecretResource, err = remotewritesecret.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	// This resource is not used in latest observability bundle versions.
	var remoteWriteAPIEndpointConfigSecretResource resource.Interface
	if config.MimirEnabled {
		remoteWriteAPIEndpointConfigSecretResource = noop.New(noop.Config{Logger: config.Logger})
	} else {
		c := remotewriteapiendpointconfigsecret.Config{
			K8sClient:          config.K8sClient,
			Logger:             config.Logger,
			OrganizationReader: organizationReader,

			BaseDomain:   config.PrometheusBaseDomain,
			Customer:     config.Customer,
			Installation: config.Installation,
			InsecureCA:   config.InsecureCA,
			Pipeline:     config.Pipeline,
			Provider:     config.Provider,
			Region:       config.Region,
			Version:      config.PrometheusVersion,
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
			Address:            config.PrometheusAddress,
			K8sClient:          config.K8sClient,
			PrometheusClient:   config.PrometheusClient,
			Logger:             config.Logger,
			Customer:           config.Customer,
			Installation:       config.Installation,
			Pipeline:           config.Pipeline,
			Version:            config.PrometheusVersion,
			Provider:           config.Provider,
			Region:             config.Region,
			Registry:           config.Registry,
			LogLevel:           config.PrometheusLogLevel,
			EvaluationInterval: config.PrometheusEvaluationInterval,
			ImageRepository:    config.PrometheusImageRepository,
			ScrapeInterval:     config.PrometheusScrapeInterval,

			MimirEnabled: config.MimirEnabled,
		}

		prometheusResource, err = prometheus.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var verticalPodAutoScalerResource resource.Interface
	{
		c := verticalpodautoscaler.Config{
			Logger:       config.Logger,
			K8sClient:    config.K8sClient,
			VpaClient:    config.VpaClient,
			Installation: config.Installation,
			Provider:     config.Provider,
		}

		verticalPodAutoScalerResource, err = verticalpodautoscaler.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var scrapeConfigResource resource.Interface
	{
		c := scrapeconfigs.Config{
			K8sClient:          config.K8sClient,
			Logger:             config.Logger,
			OrganizationReader: organizationReader,

			AdditionalScrapeConfigs: config.AdditionalScrapeConfigs,
			Bastions:                config.Bastions,
			Customer:                config.Customer,
			Pipeline:                config.Pipeline,
			Provider:                config.Provider,
			Region:                  config.Region,
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
			ExternalDNS:             config.ExternalDNS,
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
	// This resource creates a static secret to connect Prometheus to Alertmanager. When using mimir, this is not needed anymore
	if config.MimirEnabled {
		alertmanagerWiringResource = noop.New(noop.Config{Logger: config.Logger})
	} else {
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
		ciliumnetpolResource,
		rbacResource,
		heartbeatWebhookConfigResource,
		scrapeConfigResource,
		alertmanagerWiringResource,
		prometheusResource,
		remoteWriteConfigResource,
		remoteWriteSecretResource,
		remoteWriteAPIEndpointConfigSecretResource,
		remoteWriteIngressAuthResource,
		remoteWriteIngressResource,
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
