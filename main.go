package main

import (
	"context"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/microkit/command"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/versionbundle"
	"github.com/spf13/viper"

	"github.com/giantswarm/prometheus-meta-operator/flag"
	"github.com/giantswarm/prometheus-meta-operator/pkg/project"
	"github.com/giantswarm/prometheus-meta-operator/server"
	"github.com/giantswarm/prometheus-meta-operator/service"
)

var (
	f *flag.Flag = flag.New()
)

func main() {
	err := mainE(context.Background())
	if err != nil {
		panic(microerror.JSON(err))
	}
}

func mainE(ctx context.Context) error {
	var err error

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
	serverFactory := func(v *viper.Viper) microserver.Server {
		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			c := service.Config{
				Logger: logger,

				Flag:  f,
				Viper: v,
			}

			newService, err = service.New(c)
			if err != nil {
				panic(microerror.JSON(err))
			}

			go newService.Boot(ctx)
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			c := server.Config{
				Logger:  logger,
				Service: newService,

				Viper: v,
			}

			newServer, err = server.New(c)
			if err != nil {
				panic(microerror.JSON(err))
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		c := command.Config{
			Logger:        logger,
			ServerFactory: serverFactory,

			Description:    project.Description(),
			GitCommit:      project.GitSHA(),
			Name:           project.Name(),
			Source:         project.Source(),
			Version:        project.Version(),
			VersionBundles: []versionbundle.Bundle{project.NewVersionBundle()},
		}

		newCommand, err = command.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	daemonCommand := newCommand.DaemonCommand().CobraCommand()

	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.Address, "http://127.0.0.1:6443", "Address used to connect to Kubernetes. When empty in-cluster config is created.")
	daemonCommand.PersistentFlags().Bool(f.Service.Kubernetes.InCluster, false, "Whether to use the in-cluster config to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.KubeConfig, "", "KubeConfig used to connect to Kubernetes. When empty other settings are used.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.CAFile, "", "Certificate authority file path to use to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.CrtFile, "", "Certificate file path to use to authenticate with Kubernetes.")
	daemonCommand.PersistentFlags().String(f.Service.Kubernetes.TLS.KeyFile, "", "Key file path to use to authenticate with Kubernetes.")

	daemonCommand.PersistentFlags().String(f.Service.Installation.Name, "", "Name of the installation.")
	daemonCommand.PersistentFlags().String(f.Service.Installation.Pipeline, "", "Pipeline of the installation (stable or testing).")
	daemonCommand.PersistentFlags().String(f.Service.Installation.Region, "", "Region where the installation is located.")
	daemonCommand.PersistentFlags().String(f.Service.Installation.Customer, "", "Customer the installation belongs to.")
	daemonCommand.PersistentFlags().String(f.Service.Installation.Registry, "", "Container image registry.")
	daemonCommand.PersistentFlags().String(f.Service.Installation.IngressApiVersion, "networking.k8s.io/v1", "Installation supported ingress api version")
	daemonCommand.PersistentFlags().String(f.Service.Provider.Kind, "", "Provider of the installation. One of aws, azure, kvm.")

	daemonCommand.PersistentFlags().String(f.Service.Opsgenie.Key, "", "Opsgenie Key used for API authentication.")

	daemonCommand.PersistentFlags().String(f.Service.Alertmanager.Address, "", "Address to access Alertmanager UI.")
	daemonCommand.PersistentFlags().String(f.Service.Alertmanager.LogLevel, "info", "Prometheus log level.")
	daemonCommand.PersistentFlags().Bool(f.Service.Alertmanager.Storage.CreatePVC, false, "Should the operator create a PVC for storage.")
	daemonCommand.PersistentFlags().String(f.Service.Alertmanager.Storage.Size, "1Gi", "Storage size for alertmanagers.")
	daemonCommand.PersistentFlags().String(f.Service.Alertmanager.Version, "v0.22.1", "Alertmanager container image version.")

	daemonCommand.PersistentFlags().String(f.Service.Prometheus.AdditionalScrapeConfigs, "", "Additional installation-specific scrape configs.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.Address, "", "Address to access Prometheus UI.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.BaseDomain, "", "Base domain to create Prometheus Ingress resources under.")
	daemonCommand.PersistentFlags().StringSlice(f.Service.Prometheus.Bastions, make([]string, 0), "Address of the bastions.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.LogLevel, "info", "Prometheus log level.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.Mayu, "", "Mayu host.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.RemoteWrite.URL, "", "URL to send prometheus data to.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.RemoteWrite.BasicAuth.Username, "", "Remote write username.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.RemoteWrite.BasicAuth.Password, "", "Remote write password.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.Retention.Duration, "2w", "Retention duration for prometheus.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.Retention.Size, "90Gi", "Retention size for prometheus.")
	daemonCommand.PersistentFlags().Bool(f.Service.Prometheus.Storage.CreatePVC, false, "Should the operator create a PVC for storage.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.Storage.Size, "100Gi", "Storage size for prometheus.")
	daemonCommand.PersistentFlags().String(f.Service.Prometheus.Version, "v2.28.1", "Prometheus container image version.")

	daemonCommand.PersistentFlags().String(f.Service.Grafana.Address, "", "Grafana url.")
	daemonCommand.PersistentFlags().String(f.Service.Slack.ApiURL, "", "Slack api url.")
	daemonCommand.PersistentFlags().String(f.Service.Slack.ProjectName, "", "Slack project name.")
	daemonCommand.PersistentFlags().Bool(f.Service.Security.RestrictedAccess.Enabled, false, "Is the access to the prometheus restricted to certain subnets?")
	daemonCommand.PersistentFlags().String(f.Service.Security.RestrictedAccess.Subnets, "", "List of subnets to restrict the access to.")

	daemonCommand.PersistentFlags().String(f.Service.Vault.Host, "", "Host used to connect to Vault.")

	err = newCommand.CobraCommand().Execute()
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
