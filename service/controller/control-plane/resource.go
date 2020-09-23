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
	etcdcertificates "github.com/giantswarm/prometheus-meta-operator/service/controller/resource/etcd-certificates"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/namespace"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/prometheus"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/rbac"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/scrapeconfigs"
	"github.com/giantswarm/prometheus-meta-operator/service/controller/resource/servicemonitor"
)

type resourcesConfig struct {
	BaseDomain       string
	Provider         string
	Installation     string
	CreatePVC        bool
	StorageSize      string
	Vault            string
	K8sClient        k8sclient.Interface
	Logger           micrologger.Logger
	PrometheusClient promclient.Interface
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

	var prometheusResource resource.Interface
	{
		c := prometheus.Config{
			PrometheusClient: config.PrometheusClient,
			Logger:           config.Logger,
			CreatePVC:        config.CreatePVC,
			StorageSize:      config.StorageSize,
		}

		prometheusResource, err = prometheus.New(c)
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

	resources := []resource.Interface{
		namespaceResource,
		etcdCertificatesResource,
		rbacResource,
		prometheusResource,
		serviceMonitorResource,
		alertResource,
		scrapeConfigResource,
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

	// TODO: remove this block in the next release.
	var namespaceDeleterResource resource.Interface
	{
		c := namespace.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
		}

		namespaceDeleterResource, err = namespace.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		cd := deleteResourceConfig{
			Resource: namespaceDeleterResource,
		}
		namespaceDeleterResource, err = newDeleteResource(cd)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	r := []resource.Interface{
		namespaceDeleterResource,
	}
	r = append(r, resources...)

	return r, nil
}
