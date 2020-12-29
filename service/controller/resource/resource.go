package resource

import (
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v2/pkg/resource"
	"github.com/giantswarm/operatorkit/v2/pkg/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/v2/pkg/resource/wrapper/retryresource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alertmanagerconfig"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/certificates"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/heartbeat"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/heartbeatrouting"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/ingress"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/namespace"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/promxy"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/remotewriteconfig"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/scrapeconfigs"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/servicemonitor"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/verticalpodautoscaler"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/volumeresizehack"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/wrapper/monitoringdisabledresource"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

type Config struct {
	Address           string
	BaseDomain        string
	Bastions          []string
	Provider          string
	Installation      string
	Pipeline          string
	Region            string
	Registry          string
	PrometheusVersion string
	Customer          string

	CreatePVC               bool
	StorageSize             string
	RestrictedAccessEnabled bool
	WhitelistedSubnets      string
	RetentionDuration       string
	RetentionSize           string
	OpsgenieKey             string
	RemoteWriteURL          string
	RemoteWriteUsername     string
	RemoteWritePassword     string
	K8sClient               k8sclient.Interface
	Logger                  micrologger.Logger
	PrometheusClient        promclient.Interface
	VpaClient               vpa_clientset.Interface
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
			Name:                "api-certificates",
			K8sClient:           config.K8sClient,
			Logger:              config.Logger,
			SourceNameFunc:      key.Namespace,
			SourceNamespaceFunc: key.NamespaceDefault,
			TargetNameFunc:      key.SecretAPICertificates,
		}

		apiCertificatesResource, err = certificates.New(c)
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
			RemoteWriteUsername: config.RemoteWriteUsername,
			RemoteWritePassword: config.RemoteWritePassword,
		}

		remoteWriteConfigResource, err = remotewriteconfig.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var prometheusResource resource.Interface
	{
		c := prometheus.Config{
			Address:           config.Address,
			PrometheusClient:  config.PrometheusClient,
			Logger:            config.Logger,
			CreatePVC:         config.CreatePVC,
			Customer:          config.Customer,
			Installation:      config.Installation,
			Pipeline:          config.Pipeline,
			PrometheusVersion: config.PrometheusVersion,
			Provider:          config.Provider,
			Region:            config.Region,
			Registry:          config.Registry,
			StorageSize:       config.StorageSize,
			RetentionDuration: config.RetentionDuration,
			RetentionSize:     config.RetentionSize,
			RemoteWriteURL:    config.RemoteWriteURL,
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

	var serviceMonitorResource resource.Interface
	{
		c := servicemonitor.Config{
			PrometheusClient: config.PrometheusClient,
			Logger:           config.Logger,
			Provider:         config.Provider,
			Installation:     config.Installation,
		}

		serviceMonitorResource, err = servicemonitor.New(c)
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
			BaseDomain:              config.BaseDomain,
			RestrictedAccessEnabled: config.RestrictedAccessEnabled,
			WhitelistedSubnets:      config.WhitelistedSubnets,
		}

		ingressResource, err = ingress.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var promxyResource resource.Interface
	{
		c := promxy.Config{
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
			Installation: config.Installation,
			Provider:     config.Provider,
		}

		promxyResource, err = promxy.New(c)
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
		alertmanagerConfig,
		serviceMonitorResource,
		scrapeConfigResource,
		remoteWriteConfigResource,
		prometheusResource,
		verticalPodAutoScalerResource,
		volumeResizeHack,
		ingressResource,
		promxyResource,
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
