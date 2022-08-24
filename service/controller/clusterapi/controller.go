package clusterapi

import (
	"github.com/giantswarm/k8sclient/v7/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/operatorkit/v7/pkg/controller"
	"github.com/giantswarm/operatorkit/v7/pkg/resource"
	promclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	vpa_clientset "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/client/clientset/versioned"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/project"
	controllerresource "github.com/giantswarm/prometheus-meta-operator/v2/service/controller/resource"
)

type ControllerConfig struct {
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

	OpsgenieKey string

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
}

type Controller struct {
	*controller.Controller
}

func NewController(config ControllerConfig) (*Controller, error) {
	var err error

	var resources []resource.Interface
	{
		c := controllerresource.Config(config)

		resources, err = controllerresource.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var operatorkitController *controller.Controller
	{
		c := controller.Config{
			K8sClient: config.K8sClient,
			Logger:    config.Logger,
			Name:      project.Name() + "-cluster-api-controller",
			NewRuntimeObjectFunc: func() client.Object {
				return new(capi.Cluster)
			},
			Resources: resources,
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
