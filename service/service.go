// Package service implements business logic to create Kubernetes resources
// against the Kubernetes API.
package service

import (
	"context"
	"sync"

	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/k8sclient/v4/pkg/k8srestconfig"
	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/versionbundle"
	"github.com/spf13/viper"
	"k8s.io/client-go/rest"

	providerv1alpha1 "github.com/giantswarm/apiextensions/v2/pkg/apis/provider/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	capiv1alpha2 "sigs.k8s.io/cluster-api/api/v1alpha2"
	capiv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"

	"github.com/giantswarm/prometheus-meta-operator/flag"
	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/clusterapi"
	controlplane "github.com/giantswarm/prometheus-meta-operator/service/controller/control-plane"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/legacy"
)

// Config represents the configuration used to create a new service.
type Config struct {
	Logger micrologger.Logger

	Flag  *flag.Flag
	Viper *viper.Viper
}

type Service struct {
	Version *version.Service

	bootOnce               sync.Once
	legacyController       *legacy.Controller
	clusterapiController   *clusterapi.Controller
	controlplaneController *controlplane.Controller
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
				capiv1alpha2.AddToScheme,
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

	var clusterapiController *clusterapi.Controller
	{
		c := clusterapi.ControllerConfig{
			K8sClient:        k8sClient,
			Logger:           config.Logger,
			PrometheusClient: prometheusClient,
			BaseDomain:       config.Viper.GetString(config.Flag.Service.Prometheus.BaseDomain),
			Provider:         config.Viper.GetString(config.Flag.Service.Provider.Kind),
			CreatePVC:        config.Viper.GetBool(config.Flag.Service.Prometheus.Storage.CreatePVC),
			StorageSize:      config.Viper.GetString(config.Flag.Service.Prometheus.Storage.Size),
		}
		clusterapiController, err = clusterapi.NewController(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var legacyController *legacy.Controller
	{
		c := legacy.ControllerConfig{
			BaseDomain:       config.Viper.GetString(config.Flag.Service.Prometheus.BaseDomain),
			Provider:         config.Viper.GetString(config.Flag.Service.Provider.Kind),
			CreatePVC:        config.Viper.GetBool(config.Flag.Service.Prometheus.Storage.CreatePVC),
			StorageSize:      config.Viper.GetString(config.Flag.Service.Prometheus.Storage.Size),
			K8sClient:        k8sClient,
			Logger:           config.Logger,
			PrometheusClient: prometheusClient,
		}
		legacyController, err = legacy.NewController(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var controlplaneController *controlplane.Controller
	{
		c := controlplane.ControllerConfig{
			BaseDomain:       config.Viper.GetString(config.Flag.Service.Prometheus.BaseDomain),
			Provider:         config.Viper.GetString(config.Flag.Service.Provider.Kind),
			CreatePVC:        config.Viper.GetBool(config.Flag.Service.Prometheus.Storage.CreatePVC),
			StorageSize:      config.Viper.GetString(config.Flag.Service.Prometheus.Storage.Size),
			Vault:            config.Viper.GetString(config.Flag.Service.Vault.Host),
			K8sClient:        k8sClient,
			Logger:           config.Logger,
			PrometheusClient: prometheusClient,
		}
		controlplaneController, err = controlplane.NewController(c)
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

		bootOnce:               sync.Once{},
		legacyController:       legacyController,
		clusterapiController:   clusterapiController,
		controlplaneController: controlplaneController,
	}

	return s, nil
}

func (s *Service) Boot(ctx context.Context) {
	s.bootOnce.Do(func() {

		go s.legacyController.Boot(ctx)
		go s.clusterapiController.Boot(ctx)
		go s.controlplaneController.Boot(ctx)
	})
}
