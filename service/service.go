// Package service implements business logic to create Kubernetes resources
// against the Kubernetes API.
package service

import (
	"context"
	"sync"

	"golang.org/x/net/http/httpproxy"

	appsv1alpha1 "github.com/giantswarm/apiextensions-application/api/v1alpha1"
	providerv1alpha1 "github.com/giantswarm/apiextensions/v6/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/k8sclient/v7/pkg/k8srestconfig"
	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/versionbundle"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"github.com/spf13/viper"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	capiexp "sigs.k8s.io/cluster-api/exp/api/v1beta1"

	pmov1alpha1 "github.com/giantswarm/prometheus-meta-operator/v2/api/v1alpha1"
	"github.com/giantswarm/prometheus-meta-operator/v2/flag"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/project"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/clusterapi"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/legacy"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/managementcluster"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/controller/remotewrite"
)

// Config represents the configuration used to create a new service.
type Config struct {
	Logger micrologger.Logger

	Flag  *flag.Flag
	Viper *viper.Viper
}

type Service struct {
	Version *version.Service

	bootOnce                    sync.Once
	legacyController            *legacy.Controller
	clusterapiController        *clusterapi.Controller
	managementclusterController *managementcluster.Controller
	remotewriteController       *remotewrite.Controller
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Flag must not be empty")
	}
	if config.Viper == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Viper must not be empty")
	}

	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}

	var err error

	var restConfig *rest.Config
	{
		c := k8srestconfig.Config{
			Logger: config.Logger,

			Address:    config.Viper.GetString(config.Flag.Service.Kubernetes.Address),
			InCluster:  config.Viper.GetBool(config.Flag.Service.Kubernetes.InCluster),
			KubeConfig: config.Viper.GetString(config.Flag.Service.Kubernetes.KubeConfig),
			TLS: k8srestconfig.ConfigTLS{
				CAFile:  config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CAFile),
				CrtFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.CrtFile),
				KeyFile: config.Viper.GetString(config.Flag.Service.Kubernetes.TLS.KeyFile),
			},
		}

		restConfig, err = k8srestconfig.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var k8sClient k8sclient.Interface
	{
		c := k8sclient.ClientsConfig{
			Logger: config.Logger,

			RestConfig: restConfig,
			SchemeBuilder: k8sclient.SchemeBuilder{
				apiextensionsv1.AddToScheme,
				capi.AddToScheme,
				capiexp.AddToScheme,
				appsv1alpha1.AddToScheme,
				providerv1alpha1.AddToScheme,
				pmov1alpha1.AddToScheme,
			},
		}

		k8sClient, err = k8sclient.NewClients(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var prometheusClient promclient.Interface
	{
		prometheusClient, err = promclient.NewForConfig(restConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var vpaClient vpa_clientset.Interface
	{
		vpaClient, err = vpa_clientset.NewForConfig(restConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var provider = config.Viper.GetString(config.Flag.Service.Provider.Kind)

	var proxyConfig = httpproxy.FromEnvironment()
	var clusterapiController *clusterapi.Controller
	{
		if shouldCreateCAPIController(provider) {
			c := clusterapi.ControllerConfig{
				K8sClient:        k8sClient,
				Logger:           config.Logger,
				PrometheusClient: prometheusClient,
				VpaClient:        vpaClient,
				Proxy:            proxyConfig.ProxyFunc(),

				AdditionalScrapeConfigs: config.Viper.GetString(config.Flag.Service.Prometheus.AdditionalScrapeConfigs),
				Bastions:                config.Viper.GetStringSlice(config.Flag.Service.Prometheus.Bastions),
				Customer:                config.Viper.GetString(config.Flag.Service.Installation.Customer),
				Installation:            config.Viper.GetString(config.Flag.Service.Installation.Name),
				InsecureCA:              config.Viper.GetBool(config.Flag.Service.Installation.InsecureCA),
				Pipeline:                config.Viper.GetString(config.Flag.Service.Installation.Pipeline),
				Provider:                provider,
				Region:                  config.Viper.GetString(config.Flag.Service.Installation.Region),
				Registry:                config.Viper.GetString(config.Flag.Service.Installation.Registry),

				OpsgenieKey: config.Viper.GetString(config.Flag.Service.Opsgenie.Key),

				PrometheusAddress:            config.Viper.GetString(config.Flag.Service.Prometheus.Address),
				PrometheusBaseDomain:         config.Viper.GetString(config.Flag.Service.Prometheus.BaseDomain),
				PrometheusEvaluationInterval: config.Viper.GetString(config.Flag.Service.Prometheus.EvaluationInterval),
				PrometheusLogLevel:           config.Viper.GetString(config.Flag.Service.Prometheus.LogLevel),
				PrometheusRetentionDuration:  config.Viper.GetString(config.Flag.Service.Prometheus.Retention.Duration),
				PrometheusScrapeInterval:     config.Viper.GetString(config.Flag.Service.Prometheus.ScrapeInterval),
				PrometheusImageRepository:    config.Viper.GetString(config.Flag.Service.Prometheus.ImageRepository),
				PrometheusVersion:            config.Viper.GetString(config.Flag.Service.Prometheus.Version),

				RestrictedAccessEnabled: config.Viper.GetBool(config.Flag.Service.Security.RestrictedAccess.Enabled),
				WhitelistedSubnets:      config.Viper.GetString(config.Flag.Service.Security.RestrictedAccess.Subnets),

				ExternalDNS: config.Viper.GetBool(config.Flag.Service.Ingress.ExternalDNS.Enabled),
			}

			clusterapiController, err = clusterapi.NewController(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}
	}

	var legacyController *legacy.Controller
	{
		if shouldCreateLegacyController(provider) {
			c := legacy.ControllerConfig{
				K8sClient:        k8sClient,
				Logger:           config.Logger,
				PrometheusClient: prometheusClient,
				VpaClient:        vpaClient,
				Proxy:            proxyConfig.ProxyFunc(),

				AdditionalScrapeConfigs: config.Viper.GetString(config.Flag.Service.Prometheus.AdditionalScrapeConfigs),
				Bastions:                config.Viper.GetStringSlice(config.Flag.Service.Prometheus.Bastions),
				Customer:                config.Viper.GetString(config.Flag.Service.Installation.Customer),
				Installation:            config.Viper.GetString(config.Flag.Service.Installation.Name),
				InsecureCA:              config.Viper.GetBool(config.Flag.Service.Installation.InsecureCA),
				Pipeline:                config.Viper.GetString(config.Flag.Service.Installation.Pipeline),
				Provider:                provider,
				Region:                  config.Viper.GetString(config.Flag.Service.Installation.Region),
				Registry:                config.Viper.GetString(config.Flag.Service.Installation.Registry),

				OpsgenieKey: config.Viper.GetString(config.Flag.Service.Opsgenie.Key),

				PrometheusAddress:            config.Viper.GetString(config.Flag.Service.Prometheus.Address),
				PrometheusBaseDomain:         config.Viper.GetString(config.Flag.Service.Prometheus.BaseDomain),
				PrometheusEvaluationInterval: config.Viper.GetString(config.Flag.Service.Prometheus.EvaluationInterval),
				PrometheusLogLevel:           config.Viper.GetString(config.Flag.Service.Prometheus.LogLevel),
				PrometheusRetentionDuration:  config.Viper.GetString(config.Flag.Service.Prometheus.Retention.Duration),
				PrometheusScrapeInterval:     config.Viper.GetString(config.Flag.Service.Prometheus.ScrapeInterval),
				PrometheusImageRepository:    config.Viper.GetString(config.Flag.Service.Prometheus.ImageRepository),
				PrometheusVersion:            config.Viper.GetString(config.Flag.Service.Prometheus.Version),

				RestrictedAccessEnabled: config.Viper.GetBool(config.Flag.Service.Security.RestrictedAccess.Enabled),
				WhitelistedSubnets:      config.Viper.GetString(config.Flag.Service.Security.RestrictedAccess.Subnets),

				ExternalDNS: config.Viper.GetBool(config.Flag.Service.Ingress.ExternalDNS.Enabled),
			}
			legacyController, err = legacy.NewController(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}
	}

	var managementclusterController *managementcluster.Controller
	{
		c := managementcluster.ControllerConfig{
			K8sClient:        k8sClient,
			Logger:           config.Logger,
			PrometheusClient: prometheusClient,
			VpaClient:        vpaClient,

			Proxy: proxyConfig.ProxyFunc(),

			AdditionalScrapeConfigs: config.Viper.GetString(config.Flag.Service.Prometheus.AdditionalScrapeConfigs),
			Bastions:                config.Viper.GetStringSlice(config.Flag.Service.Prometheus.Bastions),
			Customer:                config.Viper.GetString(config.Flag.Service.Installation.Customer),
			Installation:            config.Viper.GetString(config.Flag.Service.Installation.Name),
			InsecureCA:              config.Viper.GetBool(config.Flag.Service.Installation.InsecureCA),
			Pipeline:                config.Viper.GetString(config.Flag.Service.Installation.Pipeline),
			Provider:                provider,
			Region:                  config.Viper.GetString(config.Flag.Service.Installation.Region),
			Registry:                config.Viper.GetString(config.Flag.Service.Installation.Registry),

			OpsgenieKey: config.Viper.GetString(config.Flag.Service.Opsgenie.Key),

			PrometheusAddress:            config.Viper.GetString(config.Flag.Service.Prometheus.Address),
			PrometheusBaseDomain:         config.Viper.GetString(config.Flag.Service.Prometheus.BaseDomain),
			PrometheusEvaluationInterval: config.Viper.GetString(config.Flag.Service.Prometheus.EvaluationInterval),
			PrometheusLogLevel:           config.Viper.GetString(config.Flag.Service.Prometheus.LogLevel),
			PrometheusRetentionDuration:  config.Viper.GetString(config.Flag.Service.Prometheus.Retention.Duration),
			PrometheusImageRepository:    config.Viper.GetString(config.Flag.Service.Prometheus.ImageRepository),
			PrometheusVersion:            config.Viper.GetString(config.Flag.Service.Prometheus.Version),

			RestrictedAccessEnabled:  config.Viper.GetBool(config.Flag.Service.Security.RestrictedAccess.Enabled),
			PrometheusScrapeInterval: config.Viper.GetString(config.Flag.Service.Prometheus.ScrapeInterval),
			WhitelistedSubnets:       config.Viper.GetString(config.Flag.Service.Security.RestrictedAccess.Subnets),

			ExternalDNS: config.Viper.GetBool(config.Flag.Service.Ingress.ExternalDNS.Enabled),

			Mayu:  config.Viper.GetString(config.Flag.Service.Prometheus.Mayu),
			Vault: config.Viper.GetString(config.Flag.Service.Vault.Host),
		}
		managementclusterController, err = managementcluster.NewController(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		c := version.Config{
			Description:    project.Description(),
			GitCommit:      project.GitSHA(),
			Name:           project.Name(),
			Source:         project.Source(),
			Version:        project.Version(),
			VersionBundles: []versionbundle.Bundle{project.NewVersionBundle()},
		}

		versionService, err = version.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var remotewriteController *remotewrite.Controller
	{
		c := remotewrite.ControllerConfig{
			K8sClient:        k8sClient,
			Logger:           config.Logger,
			PrometheusClient: prometheusClient,
			Proxy:            proxyConfig.ProxyFunc(),
		}
		remotewriteController, err = remotewrite.NewController(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	s := &Service{
		Version: versionService,

		bootOnce:                    sync.Once{},
		legacyController:            legacyController,
		clusterapiController:        clusterapiController,
		managementclusterController: managementclusterController,
		remotewriteController:       remotewriteController,
	}

	return s, nil
}

func shouldCreateLegacyController(provider string) bool {
	// Only KVM require the legacy controller.
	// AWS and Azure do work with CAPI/Cluster CRs.
	return provider == "kvm"
}

func shouldCreateCAPIController(provider string) bool {
	// KVM installations do not currently support cluster-api clusters.
	// This is being tracked here: https://github.com/giantswarm/roadmap/issues/441
	return provider != "kvm"
}

func (s *Service) Boot(ctx context.Context) {
	s.bootOnce.Do(func() {
		if s.legacyController != nil {
			go s.legacyController.Boot(ctx)
		}
		if s.clusterapiController != nil {
			go s.clusterapiController.Boot(ctx)
		}
		go s.managementclusterController.Boot(ctx)
		go s.remotewriteController.Boot(ctx)
	})
}
