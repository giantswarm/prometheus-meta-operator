package managementcluster

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/controller"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/domain"
	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/project"
	"github.com/giantswarm/prometheus-meta-operator/v2/service/key"
)

type ControllerConfig struct {
	K8sClient          k8sclient.Interface
	Logger             micrologger.Logger
	PrometheusClient   promclient.Interface
	VpaClient          vpa_clientset.Interface
	ProxyConfiguration domain.ProxyConfiguration

	AdditionalScrapeConfigs string
	Bastions                []string
	Customer                string
	Installation            string
	InsecureCA              bool
	Pipeline                string
	Provider                string
	Region                  string
	Registry                string

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

	PrometheusAddress           string
	PrometheusBaseDomain        string
	PrometheusCreatePVC         bool
	PrometheusLogLevel          string
	PrometheusRetentionDuration string
	PrometheusRetentionSize     string
	PrometheusVersion           string

	RestrictedAccessEnabled bool
	WhitelistedSubnets      string

	Mayu  string
	Vault string
}

type Controller struct {
	*controller.Controller
}

func NewController(config ControllerConfig) (*Controller, error) {
	var err error

	var resources []resource.Interface
	{
		c := resourcesConfig(config)

		resources, err = newResources(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var operatorkitController *controller.Controller
	{
		selectors := labels.Set{
			"component": "apiserver",
		}
		if key.IsCAPIManagementCluster(config.Provider) {
			selectors = labels.Set{
				"cluster.x-k8s.io/cluster-name": config.Installation,
			}
		}
		c := controller.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Name:      project.Name() + "-management-cluster-controller",
			NewRuntimeObjectFunc: func() client.Object {
				if key.IsCAPIManagementCluster(config.Provider) {
					return new(capi.Cluster)
				}
				return new(v1.Service)
			},
			Resources: resources,
			Selector:  labels.SelectorFromSet(selectors),
		}

		operatorkitController, err = controller.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	c := &Controller{
		Controller: operatorkitController,
	}

	return c, nil
}
