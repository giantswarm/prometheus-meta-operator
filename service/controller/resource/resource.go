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
	ingressv1 "github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/ingress/v1"
	ingressv1beta1 "github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/ingress/v1beta1"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteapiendpointconfigsecret"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteapiendpointsecret"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteconfig"
	remotewriteingressv1 "github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteingress/v1"
	remotewriteingressv1beta1 "github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/remotewriteingress/v1beta1"
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
	Pipeline                string
	Provider                string
	Region                  string
	Registry                string
	IngressAPIVersion       string

	OpsgenieKey string

	PrometheusAddress             string
	PrometheusBaseDomain          string
	PrometheusCreatePVC           bool
	PrometheusLogLevel            string
	PrometheusRemoteWriteURL      string
	PrometheusRemoteWriteUsername string
	PrometheusRemoteWritePassword string
	PrometheusRetentionDuration   string
	PrometheusRetentionSize       string
	PrometheusVersion             string

	RestrictedAccessEnabled bool
	WhitelistedSubnets      string
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

	var remoteWriteIngressResource resource.Interface
	if config.IngressAPIVersion == "networking.k8s.io/v1beta1" {
		c := remotewriteingressv1beta1.Config{
			K8sClient:  config.K8sClient,
			Logger:     config.Logger,
			BaseDomain: config.PrometheusBaseDomain,
		}

		remoteWriteIngressResource, err = remotewriteingressv1beta1.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	} else {
		c := remotewriteingressv1.Config{
			K8sClient:  config.K8sClient,
			Logger:     config.Logger,
			BaseDomain: config.PrometheusBaseDomain,
		}

		remoteWriteIngressResource, err = remotewriteingressv1.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var remoteWriteAgentSecretResource resource.Interface
	{
		c := remotewriteapiendpointsecret.Config{
			K8sClient:       config.K8sClient,
			Logger:          config.Logger,
			PasswordManager: password.SimpleManager{},
		}

		remoteWriteAgentSecretResource, err = remotewriteapiendpointsecret.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var remoteWriteAgentConfigSecretResource resource.Interface
	{
		c := remotewriteapiendpointconfigsecret.Config{
			K8sClient:  config.K8sClient,
			Logger:     config.Logger,
			BaseDomain: config.PrometheusBaseDomain,
		}

		remoteWriteAgentConfigSecretResource, err = remotewriteapiendpointconfigsecret.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var remoteWriteConfigResource resource.Interface
	{
		c := remotewriteconfig.Config{
			K8sClient:           config.K8sClient,
			Logger:              config.Logger,
			RemoteWriteUsername: config.PrometheusRemoteWriteUsername,
			RemoteWritePassword: config.PrometheusRemoteWritePassword,
		}

		remoteWriteConfigResource, err = remotewriteconfig.New(c)
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
	if config.IngressAPIVersion == "networking.k8s.io/v1beta1" {
		c := ingressv1beta1.Config{
			K8sClient:               config.K8sClient,
			Logger:                  config.Logger,
			BaseDomain:              config.PrometheusBaseDomain,
			RestrictedAccessEnabled: config.RestrictedAccessEnabled,
			WhitelistedSubnets:      config.WhitelistedSubnets,
		}

		ingressResource, err = ingressv1beta1.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	} else {
		c := ingressv1.Config{
			K8sClient:               config.K8sClient,
			Logger:                  config.Logger,
			BaseDomain:              config.PrometheusBaseDomain,
			RestrictedAccessEnabled: config.RestrictedAccessEnabled,
			WhitelistedSubnets:      config.WhitelistedSubnets,
		}

		ingressResource, err = ingressv1.New(c)
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

	resources := []resource.Interface{
		namespaceResource,
		apiCertificatesResource,
		rbacResource,
		heartbeatWebhookConfigResource,
		scrapeConfigResource,
		remoteWriteIngressResource,
		remoteWriteAgentSecretResource,
		remoteWriteAgentConfigSecretResource,
		remoteWriteConfigResource,
		alertmanagerWiringResource,
		prometheusResource,
		verticalPodAutoScalerResource,
		ingressResource,
		heartbeatResource,
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
