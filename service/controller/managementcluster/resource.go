package managementcluster

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

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/password"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/alerting/alertmanagerconfig"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/alerting/alertmanagerwiring"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/alerting/heartbeat"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/alerting/heartbeatwebhookconfig"
	etcdcertificates "github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/etcd-certificates"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/monitoring/ingress"
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
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/rbac"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource/wrapper/monitoringdisabledresource"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

type resourcesConfig struct {
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
	VpaClient        vpa_clientset.Interface

	Proxy func(reqURL *url.URL) (*url.URL, error)

	AdditionalScrapeConfigs string
	Bastions                []string
	Customer                string
	Installation            string
	InsecureCA              bool
	Pipeline                string
	Provider                string
	Region                  string
	Registry                string

	GrafanaAddress string
	OpsgenieKey    string
	SlackApiURL    string

	PrometheusAddress            string
	PrometheusBaseDomain         string
	PrometheusEvaluationInterval string
	PrometheusLogLevel           string
	PrometheusRetentionDuration  string
	PrometheusScrapeInterval     string
	PrometheusImageRepository    string
	PrometheusVersion            string

	RestrictedAccessEnabled bool
	WhitelistedSubnets      string

	ExternalDNS bool

	Mayu  string
	Vault string
}

func newResources(config resourcesConfig) ([]resource.Interface, error) {
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

	var etcdCertificatesResource resource.Interface
	{
		c := etcdcertificates.Config{
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
			Installation: config.Installation,
			Provider:     config.Provider,
		}

		etcdCertificatesResource, err = etcdcertificates.New(c)
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

	var alertmanagerConfigResource resource.Interface
	{
		c := alertmanagerconfig.Config{
			K8sClient:      config.K8sClient,
			Logger:         config.Logger,
			Installation:   config.Installation,
			Provider:       config.Provider,
			Proxy:          config.Proxy,
			OpsgenieKey:    config.OpsgenieKey,
			GrafanaAddress: config.GrafanaAddress,
			SlackApiURL:    config.SlackApiURL,
			Pipeline:       config.Pipeline,
		}

		alertmanagerConfigResource, err = alertmanagerconfig.New(c)
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

	var prometheusResource resource.Interface
	{
		c := prometheus.Config{
			Address:            config.PrometheusAddress,
			PrometheusClient:   config.PrometheusClient,
			K8sClient:          config.K8sClient,
			Logger:             config.Logger,
			Customer:           config.Customer,
			Installation:       config.Installation,
			Pipeline:           config.Pipeline,
			Version:            config.PrometheusVersion,
			Provider:           config.Provider,
			Region:             config.Region,
			Registry:           config.Registry,
			EvaluationInterval: config.PrometheusEvaluationInterval,
			LogLevel:           config.PrometheusLogLevel,
			RetentionDuration:  config.PrometheusRetentionDuration,
			ImageRepository:    config.PrometheusImageRepository,
			ScrapeInterval:     config.PrometheusScrapeInterval,
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
			K8sClient: config.K8sClient,
			Logger:    config.Logger,

			AdditionalScrapeConfigs: config.AdditionalScrapeConfigs,
			Bastions:                config.Bastions,
			Customer:                config.Customer,
			Provider:                config.Provider,
			Installation:            config.Installation,
			Mayu:                    config.Mayu,
			Vault:                   config.Vault,
		}

		scrapeConfigResource, err = scrapeconfigs.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var monitoringIngressResource resource.Interface
	{
		c := ingress.Config{
			K8sClient:               config.K8sClient,
			Logger:                  config.Logger,
			BaseDomain:              config.PrometheusBaseDomain,
			RestrictedAccessEnabled: config.RestrictedAccessEnabled,
			WhitelistedSubnets:      config.WhitelistedSubnets,
			ExternalDNS:             config.ExternalDNS,
		}

		monitoringIngressResource, err = ingress.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var heartbeatResource resource.Interface
	{
		c := heartbeat.Config{
			Installation: config.Installation,
			Logger:       config.Logger,
			OpsgenieKey:  config.OpsgenieKey,
			Pipeline:     config.Pipeline,
		}

		heartbeatResource, err = heartbeat.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var pvcResizeResource resource.Interface
	{
		c := pvcresizingresource.Config{
			Logger:    config.Logger,
			K8sClient: config.K8sClient,
		}

		pvcResizeResource, err = pvcresizingresource.New(c)
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

	var remoteWriteConfigResource resource.Interface
	{
		c := remotewriteconfig.Config{
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
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

	var remoteWriteSecretResource resource.Interface
	{
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

	var remoteWriteAPIEndpointConfigSecretResource resource.Interface
	{
		c := remotewriteapiendpointconfigsecret.Config{
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
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

	resources := []resource.Interface{
		namespaceResource,
		etcdCertificatesResource,
		rbacResource,
		alertmanagerConfigResource,
		heartbeatWebhookConfigResource,
		alertmanagerWiringResource,
		remoteWriteConfigResource,
		remoteWriteSecretResource,
		remoteWriteAPIEndpointConfigSecretResource,
		remoteWriteIngressAuthResource,
		remoteWriteIngressResource,
		scrapeConfigResource,
		prometheusResource,
		verticalPodAutoScalerResource,
		monitoringIngressResource,
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

	if !key.IsCAPIManagementCluster(config.Provider) {
		resources, err = Wrap(resources, config)
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
