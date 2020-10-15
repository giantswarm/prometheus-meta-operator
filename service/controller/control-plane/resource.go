package controlplane

import (
	promclient "github.com/coreos/prometheus-operator/pkg/client/versioned"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v2/pkg/resource"
	"github.com/giantswarm/operatorkit/v2/pkg/resource/wrapper/metricsresource"
	"github.com/giantswarm/operatorkit/v2/pkg/resource/wrapper/retryresource"

	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alert"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/alertmanagerconfig"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/certificates"
	etcdcertificates "github.com/giantswarm/prometheus-meta-operator/service/controller/resource/etcd-certificates"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/ingress"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/namespace"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/promxy"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/rbac"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/scrapeconfigs"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/servicemonitor"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/volumeresizehack"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/wrapper/monitoringdisabledresource"
	"github.com/giantswarm/prometheus-meta-operator/service/key"
)

type resourcesConfig struct {
	Address                 string
	BaseDomain              string
	Provider                string
	Installation            string
	CreatePVC               bool
	StorageSize             string
	Vault                   string
	RestrictedAccessEnabled bool
	WhitelistedSubnets      string
	RetentionDuration       string
	RetentionSize           string
	K8sClient               k8sclient.Interface
	Logger                  micrologger.Logger
	PrometheusClient        promclient.Interface
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

	var tlsCertificatesResource resource.Interface
	{
		c := certificates.Config{
			Name:                "tls-certificates",
			K8sClient:           config.K8sClient,
			Logger:              config.Logger,
			SourceNameFunc:      key.SecretTLSCertificates,
			SourceNamespaceFunc: key.NamespaceMonitoring,
			TargetNameFunc:      key.SecretTLSCertificates,
		}

		tlsCertificatesResource, err = certificates.New(c)
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

	var prometheusResource resource.Interface
	{
		c := prometheus.Config{
			Address:           config.Address,
			PrometheusClient:  config.PrometheusClient,
			Logger:            config.Logger,
			CreatePVC:         config.CreatePVC,
			StorageSize:       config.StorageSize,
			RetentionDuration: config.RetentionDuration,
			RetentionSize:     config.RetentionSize,
		}

		prometheusResource, err = prometheus.New(c)
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

	var alertResource resource.Interface
	{
		c := alert.Config{
			PrometheusClient: config.PrometheusClient,
			Logger:           config.Logger,
		}

		alertResource, err = alert.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var scrapeConfigResource resource.Interface
	{
		c := scrapeconfigs.Config{
			K8sClient:    config.K8sClient,
			Logger:       config.Logger,
			Provider:     config.Provider,
			Installation: config.Installation,
			Vault:        config.Vault,
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
	resources := []resource.Interface{
		namespaceResource,
		tlsCertificatesResource,
		etcdCertificatesResource,
		rbacResource,
		alertmanagerConfig,
		serviceMonitorResource,
		alertResource,
		scrapeConfigResource,
		prometheusResource,
		volumeResizeHack,
		ingressResource,
		promxyResource,
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
		resources, err = ControlPlaneWrap(resources, config)
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
