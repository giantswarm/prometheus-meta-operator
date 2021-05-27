package resource

import (
	"github.com/giantswarm/k8sclient/v5/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v4/pkg/resource"
	"github.com/giantswarm/operatorkit/v4/pkg/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/v4/pkg/resource/wrapper/retryresource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alertmanagerconfig"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/certificates"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/heartbeat"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/heartbeatrouting"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/ingress"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/namespace"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/remotewriteconfig"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/scrapeconfigs"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/tlscleanup"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/verticalpodautoscaler"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/volumeresizehack"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/wrapper/monitoringdisabledresource"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

type Config struct {
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
	VpaClient        vpa_clientset.Interface

	HTTPProxy  string
	HTTPSProxy string
	NoProxy    string

	Bastions     []string
	Customer     string
	Installation string
	Pipeline     string
	Provider     string
	Region       string
	Registry     string

	OpsgenieKey string

	PrometheusAddress             string
	PrometheusBaseDomain          string
	PrometheusCreatePVC           bool
	PrometheusStorageSize         string
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
					NameFunc:      key.CAPICertificateName,
					NamespaceFunc: key.CAPICertificateNamespace,
				},
			},
			Target: key.SecretAPICertificates,
		}

		apiCertificatesResource, err = certificates.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var tlsCleanupResource resource.Interface
	{
		c := tlscleanup.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		tlsCleanupResource, err = tlscleanup.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var alertmanagerConfig resource.Interface
	{
		c := alertmanagerconfig.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		alertmanagerConfig, err = alertmanagerconfig.New(c)
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
			PrometheusVersion: config.PrometheusVersion,
			Provider:          config.Provider,
			Region:            config.Region,
			Registry:          config.Registry,
			StorageSize:       config.PrometheusStorageSize,
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

	var volumeResizeHack resource.Interface
	{
		c := volumeresizehack.Config{
			Logger:           config.Logger,
			K8sClient:        config.K8sClient,
			PrometheusClient: config.PrometheusClient,
		}

		volumeResizeHack, err = volumeresizehack.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var scrapeConfigResource resource.Interface
	{
		c := scrapeconfigs.Config{
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
			Bastions:     config.Bastions,
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

	var heartbeatRoutingResource resource.Interface
	{
		c := heartbeatrouting.Config{
			Installation: config.Installation,
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
			OpsgenieKey:  config.OpsgenieKey,
		}

		heartbeatRoutingResource, err = heartbeatrouting.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	resources := []resource.Interface{
		namespaceResource,
		apiCertificatesResource,
		tlsCleanupResource,
		alertmanagerConfig,
		scrapeConfigResource,
		remoteWriteConfigResource,
		prometheusResource,
		verticalPodAutoScalerResource,
		volumeResizeHack,
		ingressResource,
		heartbeatResource,
		heartbeatRoutingResource,
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
