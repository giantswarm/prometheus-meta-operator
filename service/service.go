// Package service implements business logic to create Kubernetes resources
// against the Kubernetes API.
package service

import (
	"context"
	"os"
	"sync"

	providerv1alpha1 "github.com/giantswarm/apiextensions/v3/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/k8sclient/v5/pkg/k8srestconfig"
	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/versionbundle"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	"github.com/spf13/viper"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	"k8s.io/client-go/rest"
	capiv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"

	"github.com/giantswarm/prometheus-meta-operator/flag"
	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/clusterapi"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/legacy"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/managementcluster"
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
				capiv1alpha3.AddToScheme,
				providerv1alpha1.AddToScheme,
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

	var clusterapiController *clusterapi.Controller
	{
		if shouldCreateCAPIController(provider) {
			c := clusterapi.ControllerConfig{
				K8sClient:        k8sClient,
				Logger:           config.Logger,
				PrometheusClient: prometheusClient,
				VpaClient:        vpaClient,

				HTTPProxy:  os.Getenv("HTTP_PROXY"),
				HTTPSProxy: os.Getenv("HTTPS_PROXY"),
				NoProxy:    os.Getenv("NO_PROXY"),

				AdditionalScrapeConfigs: config.Viper.GetString(config.Flag.Service.Prometheus.AdditionalScrapeConfigs),
				Bastions:                config.Viper.GetStringSlice(config.Flag.Service.Prometheus.Bastions),
				Customer:                config.Viper.GetString(config.Flag.Service.Installation.Customer),
				Installation:            config.Viper.GetString(config.Flag.Service.Installation.Name),
				Pipeline:                config.Viper.GetString(config.Flag.Service.Installation.Pipeline),
				Provider:                provider,
				Region:                  config.Viper.GetString(config.Flag.Service.Installation.Region),
				Registry:                config.Viper.GetString(config.Flag.Service.Installation.Registry),

				OpsgenieKey: config.Viper.GetString(config.Flag.Service.Opsgenie.Key),

				PrometheusAddress:             config.Viper.GetString(config.Flag.Service.Prometheus.Address),
				PrometheusBaseDomain:          config.Viper.GetString(config.Flag.Service.Prometheus.BaseDomain),
				PrometheusCreatePVC:           config.Viper.GetBool(config.Flag.Service.Prometheus.Storage.CreatePVC),
				PrometheusStorageSize:         config.Viper.GetString(config.Flag.Service.Prometheus.Storage.Size),
				PrometheusLogLevel:            config.Viper.GetString(config.Flag.Service.Prometheus.LogLevel),
				PrometheusRemoteWriteURL:      config.Viper.GetString(config.Flag.Service.Prometheus.RemoteWrite.URL),
				PrometheusRemoteWriteUsername: config.Viper.GetString(config.Flag.Service.Prometheus.RemoteWrite.BasicAuth.Username),
				PrometheusRemoteWritePassword: config.Viper.GetString(config.Flag.Service.Prometheus.RemoteWrite.BasicAuth.Password),
				PrometheusRetentionDuration:   config.Viper.GetString(config.Flag.Service.Prometheus.Retention.Duration),
				PrometheusRetentionSize:       config.Viper.GetString(config.Flag.Service.Prometheus.Retention.Size),
				PrometheusVersion:             config.Viper.GetString(config.Flag.Service.Prometheus.Version),

				RestrictedAccessEnabled: config.Viper.GetBool(config.Flag.Service.Security.RestrictedAccess.Enabled),
				WhitelistedSubnets:      config.Viper.GetString(config.Flag.Service.Security.RestrictedAccess.Subnets),
			}

			clusterapiController, err = clusterapi.NewController(c)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}
	}

	var legacyController *legacy.Controller = nil
	createLegacyController, err := shouldCreateLegacyController(k8sClient, provider)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if createLegacyController {
		c := legacy.ControllerConfig{
			K8sClient:        k8sClient,
			Logger:           config.Logger,
			PrometheusClient: prometheusClient,
			VpaClient:        vpaClient,

			HTTPProxy:  os.Getenv("HTTP_PROXY"),
			HTTPSProxy: os.Getenv("HTTPS_PROXY"),
			NoProxy:    os.Getenv("NO_PROXY"),

			AdditionalScrapeConfigs: config.Viper.GetString(config.Flag.Service.Prometheus.AdditionalScrapeConfigs),
			Bastions:                config.Viper.GetStringSlice(config.Flag.Service.Prometheus.Bastions),
			Customer:                config.Viper.GetString(config.Flag.Service.Installation.Customer),
			Installation:            config.Viper.GetString(config.Flag.Service.Installation.Name),
			Pipeline:                config.Viper.GetString(config.Flag.Service.Installation.Pipeline),
			Provider:                provider,
			Region:                  config.Viper.GetString(config.Flag.Service.Installation.Region),
			Registry:                config.Viper.GetString(config.Flag.Service.Installation.Registry),

			OpsgenieKey: config.Viper.GetString(config.Flag.Service.Opsgenie.Key),

			PrometheusAddress:             config.Viper.GetString(config.Flag.Service.Prometheus.Address),
			PrometheusBaseDomain:          config.Viper.GetString(config.Flag.Service.Prometheus.BaseDomain),
			PrometheusCreatePVC:           config.Viper.GetBool(config.Flag.Service.Prometheus.Storage.CreatePVC),
			PrometheusStorageSize:         config.Viper.GetString(config.Flag.Service.Prometheus.Storage.Size),
			PrometheusLogLevel:            config.Viper.GetString(config.Flag.Service.Prometheus.LogLevel),
			PrometheusRemoteWriteURL:      config.Viper.GetString(config.Flag.Service.Prometheus.RemoteWrite.URL),
			PrometheusRemoteWriteUsername: config.Viper.GetString(config.Flag.Service.Prometheus.RemoteWrite.BasicAuth.Username),
			PrometheusRemoteWritePassword: config.Viper.GetString(config.Flag.Service.Prometheus.RemoteWrite.BasicAuth.Password),
			PrometheusRetentionDuration:   config.Viper.GetString(config.Flag.Service.Prometheus.Retention.Duration),
			PrometheusRetentionSize:       config.Viper.GetString(config.Flag.Service.Prometheus.Retention.Size),
			PrometheusVersion:             config.Viper.GetString(config.Flag.Service.Prometheus.Version),

			RestrictedAccessEnabled: config.Viper.GetBool(config.Flag.Service.Security.RestrictedAccess.Enabled),
			WhitelistedSubnets:      config.Viper.GetString(config.Flag.Service.Security.RestrictedAccess.Subnets),
		}
		legacyController, err = legacy.NewController(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var managementclusterController *managementcluster.Controller
	{
		c := managementcluster.ControllerConfig{
			K8sClient:        k8sClient,
			Logger:           config.Logger,
			PrometheusClient: prometheusClient,
			VpaClient:        vpaClient,

			HTTPProxy:  os.Getenv("HTTP_PROXY"),
			HTTPSProxy: os.Getenv("HTTPS_PROXY"),
			NoProxy:    os.Getenv("NO_PROXY"),

			AdditionalScrapeConfigs: config.Viper.GetString(config.Flag.Service.Prometheus.AdditionalScrapeConfigs),
			Bastions:                config.Viper.GetStringSlice(config.Flag.Service.Prometheus.Bastions),
			Customer:                config.Viper.GetString(config.Flag.Service.Installation.Customer),
			Installation:            config.Viper.GetString(config.Flag.Service.Installation.Name),
			Pipeline:                config.Viper.GetString(config.Flag.Service.Installation.Pipeline),
			Provider:                provider,
			Region:                  config.Viper.GetString(config.Flag.Service.Installation.Region),
			Registry:                config.Viper.GetString(config.Flag.Service.Installation.Registry),

			AlertmanagerAddress:     config.Viper.GetString(config.Flag.Service.Alertmanager.Address),
			AlertmanagerCreatePVC:   config.Viper.GetBool(config.Flag.Service.Alertmanager.Storage.CreatePVC),
			AlertmanagerLogLevel:    config.Viper.GetString(config.Flag.Service.Alertmanager.LogLevel),
			AlertmanagerStorageSize: config.Viper.GetString(config.Flag.Service.Alertmanager.Storage.Size),
			AlertmanagerVersion:     config.Viper.GetString(config.Flag.Service.Alertmanager.Version),
			GrafanaAddress:          config.Viper.GetString(config.Flag.Service.Grafana.Address),
			SlackApiURL:             config.Viper.GetString(config.Flag.Service.Slack.ApiURL),
			SlackProjectName:        config.Viper.GetString(config.Flag.Service.Slack.ProjectName),

			OpsgenieKey: config.Viper.GetString(config.Flag.Service.Opsgenie.Key),

			PrometheusAddress:             config.Viper.GetString(config.Flag.Service.Prometheus.Address),
			PrometheusBaseDomain:          config.Viper.GetString(config.Flag.Service.Prometheus.BaseDomain),
			PrometheusCreatePVC:           config.Viper.GetBool(config.Flag.Service.Prometheus.Storage.CreatePVC),
			PrometheusStorageSize:         config.Viper.GetString(config.Flag.Service.Prometheus.Storage.Size),
			PrometheusLogLevel:            config.Viper.GetString(config.Flag.Service.Prometheus.LogLevel),
			PrometheusRemoteWriteURL:      config.Viper.GetString(config.Flag.Service.Prometheus.RemoteWrite.URL),
			PrometheusRemoteWriteUsername: config.Viper.GetString(config.Flag.Service.Prometheus.RemoteWrite.BasicAuth.Username),
			PrometheusRemoteWritePassword: config.Viper.GetString(config.Flag.Service.Prometheus.RemoteWrite.BasicAuth.Password),
			PrometheusRetentionDuration:   config.Viper.GetString(config.Flag.Service.Prometheus.Retention.Duration),
			PrometheusRetentionSize:       config.Viper.GetString(config.Flag.Service.Prometheus.Retention.Size),
			PrometheusVersion:             config.Viper.GetString(config.Flag.Service.Prometheus.Version),

			RestrictedAccessEnabled: config.Viper.GetBool(config.Flag.Service.Security.RestrictedAccess.Enabled),
			WhitelistedSubnets:      config.Viper.GetString(config.Flag.Service.Security.RestrictedAccess.Subnets),

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

	s := &Service{
		Version: versionService,

		bootOnce:                    sync.Once{},
		legacyController:            legacyController,
		clusterapiController:        clusterapiController,
		managementclusterController: managementclusterController,
	}

	return s, nil
}

func shouldCreateLegacyController(clients k8sclient.Interface, provider string) (bool, error) {
	switch provider {
	case "aws":
		_, err := clients.G8sClient().ProviderV1alpha1().AWSConfigs("default").List(context.Background(), v1.ListOptions{})
		if apierrors.IsNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, microerror.Mask(err)
		}
	case "azure":
		_, err := clients.G8sClient().ProviderV1alpha1().AzureConfigs("default").List(context.Background(), v1.ListOptions{})
		if apierrors.IsNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, microerror.Mask(err)
		}
	case "kvm":
		_, err := clients.G8sClient().ProviderV1alpha1().KVMConfigs("default").List(context.Background(), v1.ListOptions{})
		if apierrors.IsNotFound(err) {
			return false, nil
		} else if err != nil {
			return false, microerror.Mask(err)
		}
	case "vmware":
		return false, nil
	}
	return true, nil
}

func shouldCreateCAPIController(provider string) bool {
	// KVM installations do not currently support cluster-api clusters.
	// This is being tracked here: https://github.com/giantswarm/roadmap/issues/319
	if provider == "kvm" {
		return false
	}

	return true
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
	})
}
