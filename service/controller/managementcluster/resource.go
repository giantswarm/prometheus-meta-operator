package managementcluster

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/v7/pkg/resource/wrapper/retryresource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alerting/alertmanager"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alerting/alertmanagerconfig"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alerting/alertmanagerrouting"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alerting/alertmanagerwiring"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alerting/heartbeat"
	etcdcertificates "github.com/giantswarm/prometheus-meta-operator/service/controller/resource/etcd-certificates"
	ingressv1 "github.com/giantswarm/prometheus-meta-operator/service/controller/resource/monitoring/ingress/v1"
	ingressv1beta1 "github.com/giantswarm/prometheus-meta-operator/service/controller/resource/monitoring/ingress/v1beta1"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/monitoring/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/monitoring/remotewriteconfig"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/monitoring/scrapeconfigs"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/monitoring/verticalpodautoscaler"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/namespace"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/rbac"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/wrapper/monitoringdisabledresource"
)

type resourcesConfig struct {
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
	VpaClient        vpa_clientset.Interface

	HTTPProxy  string
	HTTPSProxy string
	NoProxy    string

	AdditionalScrapeConfigs string
	Bastions                []string
	Customer                string
	Installation            string
	Pipeline                string
	Provider                string
	Region                  string
	Registry                string
	IngressAPIVersion       string

	AlertmanagerAddress     string
	AlertmanagerBaseDomain  string
	AlertmanagerCreatePVC   bool
	AlertmanagerLogLevel    string
	AlertmanagerStorageSize string
	AlertmanagerVersion     string
	GrafanaAddress          string
	OpsgenieKey             string
	SlackApiURL             string
	SlackProjectName        string

	PrometheusAddress             string
	PrometheusBaseDomain          string
	PrometheusCreatePVC           bool
	PrometheusStorageSize         string
	PrometheusLogLevel            string
	PrometheusRemoteWriteURL      string
	PrometheusRemoteWriteUsername string
	PrometheusRemoteWritePassword string
	PrometheusRetentionDuration   string
	PrometheusRetentionSize       string
	PrometheusVersion             string

	RestrictedAccessEnabled bool
	WhitelistedSubnets      string

	Mayu  string
	Vault string
}

func newResources(config resourcesConfig) ([]resource.Interface, error) {
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

	var etcdCertificatesResource resource.Interface
	{
		c := etcdcertificates.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
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

	var alertmanagerResource resource.Interface
	{
		c := alertmanager.Config{
			Address:     config.AlertmanagerAddress,
			Client:      config.PrometheusClient,
			Logger:      config.Logger,
			CreatePVC:   config.AlertmanagerCreatePVC,
			LogLevel:    config.AlertmanagerLogLevel,
			StorageSize: config.AlertmanagerStorageSize,
			Version:     config.AlertmanagerVersion,
		}

		alertmanagerResource, err = alertmanager.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var alertmanagerRoutingResource resource.Interface
	{
		c := alertmanagerrouting.Config{
			Client: config.PrometheusClient,
			Logger: config.Logger,

			Installation: config.Installation,
			HTTPProxy:    config.HTTPProxy,
			HTTPSProxy:   config.HTTPSProxy,
			NoProxy:      config.NoProxy,
		}

		alertmanagerRoutingResource, err = alertmanagerrouting.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var alertmanagerWiringResource resource.Interface
	{
		c := alertmanagerwiring.Config{
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
			Installation: config.Installation,
		}

		alertmanagerWiringResource, err = alertmanagerwiring.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	var alertmanagerConfigSecretResource resource.Interface
	{
		c := alertmanagerconfig.Config{
			K8sClient:        config.K8sClient,
			Logger:           config.Logger,
			Installation:     config.Installation,
			Provider:         config.Provider,
			HTTPProxy:        config.HTTPProxy,
			HTTPSProxy:       config.HTTPSProxy,
			NoProxy:          config.NoProxy,
			OpsgenieKey:      config.OpsgenieKey,
			GrafanaAddress:   config.GrafanaAddress,
			SlackApiURL:      config.SlackApiURL,
			SlackProjectName: config.SlackProjectName,
			Pipeline:         config.Pipeline,
		}

		alertmanagerConfigSecretResource, err = alertmanagerconfig.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	/*
		var alertmanagerIngressResource resource.Interface
		if config.IngressAPIVersion == "networking.k8s.io/v1beta1" {
			c := alertingingressv1beta1.Config{
				K8sClient:               config.K8sClient,
				Logger:                  config.Logger,
				BaseDomain:              config.AlertmanagerBaseDomain,
				RestrictedAccessEnabled: config.RestrictedAccessEnabled,
				WhitelistedSubnets:      config.WhitelistedSubnets,
			}
			alertmanagerIngressResource, err = alertingingressv1beta1.New(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		} else {
			c := alertingingressv1.Config{
				K8sClient:               config.K8sClient,
				Logger:                  config.Logger,
				BaseDomain:              config.AlertmanagerBaseDomain,
				RestrictedAccessEnabled: config.RestrictedAccessEnabled,
				WhitelistedSubnets:      config.WhitelistedSubnets,
			}
			alertmanagerIngressResource, err = alertingingressv1.New(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}
	*/

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
			StorageSize:       config.PrometheusStorageSize,
			LogLevel:          config.PrometheusLogLevel,
			RetentionDuration: config.PrometheusRetentionDuration,
			RetentionSize:     config.PrometheusRetentionSize,
			RemoteWriteURL:    config.PrometheusRemoteWriteURL,
			HTTPProxy:         config.HTTPProxy,
			HTTPSProxy:        config.HTTPSProxy,
			NoProxy:           config.NoProxy,
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
	if config.IngressAPIVersion == "networking.k8s.io/v1beta1" {
		c := ingressv1beta1.Config{
			K8sClient:               config.K8sClient,
			Logger:                  config.Logger,
			BaseDomain:              config.PrometheusBaseDomain,
			RestrictedAccessEnabled: config.RestrictedAccessEnabled,
			WhitelistedSubnets:      config.WhitelistedSubnets,
		}

		monitoringIngressResource, err = ingressv1beta1.New(c)
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

		monitoringIngressResource, err = ingressv1.New(c)
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

	resources := []resource.Interface{
		namespaceResource,
		etcdCertificatesResource,
		rbacResource,
		alertmanagerResource,
		alertmanagerRoutingResource,
		//alertmanagerIngressResource,
		alertmanagerWiringResource,
		alertmanagerConfigSecretResource,
		scrapeConfigResource,
		remoteWriteConfigResource,
		prometheusResource,
		verticalPodAutoScalerResource,
		monitoringIngressResource,
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
